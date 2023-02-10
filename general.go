package pingen

type AuthSuccess struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

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
