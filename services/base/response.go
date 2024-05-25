package servicebase

type Pagination struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	Records    int `json:"records"`
}

type ResponseBody struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
