package pingen

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	clientID     string
	clientSecret string

	baseUrl        string
	identiyUrl     string
	organisationId string

	bearer       *AuthSuccess
	bearerExpiry time.Time
	httpClient   *resty.Client
	authClient   *resty.Client
}

type AuthSuccess struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

//Creates a new client struct
func NewClient(clientID string, clientSecret string, useProd bool, organisationId string, ctx context.Context) (*Client, error) {
	c := &Client{}

	c.clientID = clientID
	c.clientSecret = clientSecret
	c.organisationId = organisationId

	if useProd {
		c.baseUrl = "https://api.v2.pingen.com"
		c.identiyUrl = "https://identity.pingen.com"
	} else {
		c.baseUrl = "https://api-staging.v2.pingen.com"
		c.identiyUrl = "https://identity-staging.pingen.com"
	}

	//order is important, auth needs to be initalized first
	c.authClient = resty.New()
	c.authClient.SetBaseURL(c.identiyUrl)

	c.httpClient = resty.New()
	c.httpClient.SetBaseURL(c.baseUrl)
	c.httpClient.SetHeader("Accept", "application/json")
	bearerToken, err := c.getBearer()
	if err == nil {
		c.httpClient.SetAuthToken(bearerToken)
	} else {
		return nil, err
	}

	return c, nil
}

func (c *Client) getBearer() (string, error) {
	//Check if we have cached values
	if !c.bearerExpiry.IsZero() && c.bearer != nil {
		if c.bearerExpiry.Before(time.Now()) && c.bearer.AccessToken != "" {
			//bearer not expired
			return c.bearer.AccessToken, nil
		}
	}

	authResult := &AuthSuccess{}
	resp, err := c.authClient.R().
		SetBody(map[string]interface{}{
			"grant_type":    "client_credentials",
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		SetResult(authResult).
		Post("/auth/access-tokens")

	if err != nil {
		log.Errorf("Auth error: %v", err)
		return "", fmt.Errorf("failed to get Bearer token, check your client credentials and if you use the correct environment. Error: %w", err)
	}
	if resp.IsError() {
		log.Errorf("Auth error: %v", resp)
		return "", fmt.Errorf("failed to get Bearer token, check your client credentials and if you use the correct environment. Error: %v", resp)
	}

	if err != nil || resp.StatusCode() != 200 {
		log.Fatalf("auth code wrong: %v", err)

	}

	c.bearer = authResult
	c.bearerExpiry = resp.ReceivedAt().Add(time.Second * time.Duration(c.bearer.ExpiresIn))
	log.Debugf("got auth! " + c.bearerExpiry.String())
	return c.bearer.AccessToken, nil
}

func (c *Client) ListLetters() (result *LetterList, err error) {
	result = &LetterList{}

	resp, err := c.httpClient.R().
		SetResult(&result).
		SetError(&ApiError{}).
		Get(fmt.Sprintf("/organisations/%s/letters", c.organisationId))

	if err != nil || resp.StatusCode() != 200 {
		log.Errorf("ListLetters error: %v", err)
		return nil, fmt.Errorf("listerror failed %w", err)

	}

	for _, v := range result.Data {
		log.Debugf("Letter %s: State %s, FileName %s, to Addr [%s]",
			v.ID, v.Attributes.Status, v.Attributes.FileOriginalName, v.Attributes.Address)
	}

	return result, nil
}

//Letter CRUD
func (c *Client) GetLetter(letterID string) (result *Letter, err error) {
	result = &Letter{}

	resp, err := c.httpClient.R().
		SetResult(&result).
		SetError(&ApiError{}).
		Get(fmt.Sprintf("/organisations/%s/letters/%s", c.organisationId, letterID))

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, fmt.Errorf("GetLetter failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return nil, fmt.Errorf("GetLetter with api message: %v", errData.Errors)
	}

	log.Debugf("Letter %s: State %s, FileName %s, to Addr [%s]",
		result.Data.ID, result.Data.Attributes.Status, result.Data.Attributes.FileOriginalName, result.Data.Attributes.Address)

	return result, nil
}

func (c *Client) CancelLetter(letterID string) (err error) {
	resp, err := c.httpClient.R().
		SetError(&ApiError{}).
		Patch(fmt.Sprintf("/organisations/%s/letters/%s/cancel", c.organisationId, letterID))

	if err != nil {
		log.Errorf("Error: %v", err)
		return fmt.Errorf("CancelLetter failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return fmt.Errorf("CancelLetter with api message: %v", errData.Errors)
	}

	log.Debugf("Letter %s canceled", letterID)
	return nil
}

func (c *Client) DeleteLetter(letterID string) (err error) {
	resp, err := c.httpClient.R().
		SetError(&ApiError{}).
		Delete(fmt.Sprintf("/organisations/%s/letters/%s", c.organisationId, letterID))

	if err != nil {
		log.Errorf("Error: %v", err)
		return fmt.Errorf("DeleteLetter failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return fmt.Errorf("DeleteLetter with api message: %v", errData.Errors)
	}

	log.Debugf("Letter %s deleted", letterID)
	return nil
}

//Upload process!
/*
Function that uploads the pdf file and creates a letter. Required letter data needs to be provided in createData.

--

createData.Data.Attributes.AddressPosition must be right or left

createData.Data.Attributes.DeliveryProduct can be empty, "fast", "cheap", "bulk", "premium", "registered", defaults to cheap

createData.Data.Attributes.PrintMode can be "simplex", "duplex", defaults to simplex

createData.Data.Attributes.PrintSpectrum can be "grayscale", "color", defaults to grayscale

createData.Data.Attributes.FileURL & createData.Data.Attributes.FileURLSignature will be autofilled.
*/
func (c *Client) CreateLetter(uploadPDF []byte, createData *CreateData) (result *Letter, err error) {
	result = &Letter{}
	validateCreateData(createData)

	//get upload url
	uploadData := &UploadData{}

	resp, err := c.httpClient.R().
		SetResult(uploadData).
		SetError(&ApiError{}).
		Get("file-upload")

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, fmt.Errorf("UploadFile failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return nil, fmt.Errorf("UploadFile with api message: %v", errData.Errors)
	}

	log.Debugf("Got file upload endpoint: url %s, sig %s", uploadData.Data.Attributes.URL, uploadData.Data.Attributes.URLSignature)

	//upload file
	//Important: NO AUTH HEADER
	uploadClient := resty.New()
	resp, err = uploadClient.R().
		SetHeader("Content-Type", "application/pdf").
		SetBody(uploadPDF).
		Put(uploadData.Data.Attributes.URL)

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, fmt.Errorf("UploadFile failed %w", err)
	}
	if resp.IsError() {
		log.Errorf("Error: %v", resp)
		return nil, fmt.Errorf("UploadFile failed %w", err)
	}
	log.Debugf("Uploaded file, resp code %d", resp.StatusCode())

	//create the letter
	createData.Data.Type = "letters"
	createData.Data.Attributes.FileURL = uploadData.Data.Attributes.URL
	createData.Data.Attributes.FileURLSignature = uploadData.Data.Attributes.URLSignature

	//name
	//addrees pos
	//auto_send
	//delivery_product
	//print_mode
	//print_spectrum

	resp, err = c.httpClient.R().
		SetHeader("Content-Type", "application/vnd.api+json").
		SetBody(createData).
		SetResult(result).
		SetError(&ApiError{}).
		Post(fmt.Sprintf("/organisations/%s/letters", c.organisationId))

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, fmt.Errorf("UploadFile failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return nil, fmt.Errorf("UploadFile with api message: %v", errData.Errors)
	}

	log.Debugf("Created letter: %s", result.Data.ID)
	return result, nil
}

func validateCreateData(createData *CreateData) error {
	if createData.Data.Attributes.PrintMode == "" {
		createData.Data.Attributes.PrintMode = "simplex"
	}
	if createData.Data.Attributes.PrintSpectrum == "" {
		createData.Data.Attributes.PrintSpectrum = "grayscale"
	}
	if createData.Data.Attributes.DeliveryProduct == "" {
		createData.Data.Attributes.DeliveryProduct = "cheap"
	}
	return nil
}

/*
Function that sends a letter.

--

data.Data.Attributes.DeliveryProduct can be empty, "fast", "cheap", "bulk", "premium", "registered"

data.Data.Attributes.PrintMode can be empty, "simplex", "duplex"

data.Data.Attributes.PrintSpectrum can be empty, "grayscale", "color"

Data.Type & Data.ID will be autofilled
*/
func (c *Client) SendLetter(letterID string, data *SendData) (result *Letter, err error) {
	result = &Letter{}
	//todo validate senddata
	data.Data.Type = "letters"
	data.Data.ID = letterID

	resp, err := c.httpClient.R().
		SetHeader("Content-Type", "application/vnd.api+json").
		SetBody(data).
		SetError(&ApiError{}).
		SetResult(result).
		Patch(fmt.Sprintf("/organisations/%s/letters/%s/send", c.organisationId, letterID))

	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, fmt.Errorf("SendLetter failed %w", err)
	}
	if resp.IsError() {
		errData := resp.Error().(*ApiError)
		log.Errorf("Api error: %v", errData)
		return nil, fmt.Errorf("SendLetter with api message: %v", errData.Errors)
	}

	log.Debugf("Letter %s sent", letterID)
	return result, nil
}

/*
status:
valid
action_required
proccesing
printing
sent
undeliverable
unprintable
*/

//file upload
//create a letter
//get letter details
//get file of letter
//get letter collection
