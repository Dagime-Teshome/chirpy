package main

type Req_body struct {
	Body string `json:"body"`
}

type Res_body struct {
	Valid bool `json:"valid"`
}

type clean_body struct {
	Cleaned_Body string `json:"cleaned_body"`
}

type err_resp struct {
	Error string `json:"error"`
}
