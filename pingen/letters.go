package pingen

type ApiError struct {
	Errors []struct {
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Source struct {
			Pointer   string `json:"pointer"`
			Parameter string `json:"parameter"`
		} `json:"source"`
	} `json:"errors"`
}

type LetterList struct {
	Data []LetterData `json:"data"`
}

type Letter struct {
	Data LetterData `json:"data"`
}

type LetterData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Status           string   `json:"status"`
		FileOriginalName string   `json:"file_original_name"`
		FilePages        int      `json:"file_pages"`
		Address          string   `json:"address"`
		AddressPosition  string   `json:"address_position"`
		Country          string   `json:"country"`
		DeliveryProduct  string   `json:"delivery_product"`
		PrintMode        string   `json:"print_mode"`
		PrintSpectrum    string   `json:"print_spectrum"`
		PriceCurrency    string   `json:"price_currency"`
		PriceValue       float64  `json:"price_value"`
		PaperTypes       []string `json:"paper_types"`
		Fonts            []struct {
			Name       string `json:"name"`
			IsEmbedded bool   `json:"is_embedded"`
		} `json:"fonts"`
		TrackingNumber string `json:"tracking_number"`
		SubmittedAt    string `json:"submitted_at"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	} `json:"attributes"`
}
