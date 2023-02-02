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
