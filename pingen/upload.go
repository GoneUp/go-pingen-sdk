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

			DeliveryProduct string `json:"delivery_product"`
			PrintMode       string `json:"print_mode"`
			PrintSpectrum   string `json:"print_spectrum"`
		} `json:"attributes"`
	} `json:"data"`
}

type SendData struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			DeliveryProduct string `json:"delivery_product"`
			PrintMode       string `json:"print_mode"`
			PrintSpectrum   string `json:"print_spectrum"`
		} `json:"attributes"`
	} `json:"data"`
}
