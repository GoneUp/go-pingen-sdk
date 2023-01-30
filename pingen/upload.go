package pingen

type UploadData struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			URL          string `json:"url"`
			URLSignature string `json:"url_signature"`
			ExpiresAt    string `json:"expires_at"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

type CreateData struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			FileOriginalName string `json:"file_original_name"`
			FileURL          string `json:"file_url"`
			FileURLSignature string `json:"file_url_signature"`
			AddressPosition  string `json:"address_position"`
			AutoSend         bool   `json:"auto_send"`
		} `json:"attributes"`
	} `json:"data"`
}
