package common

type RespostaComErro struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Resposta struct {
	Message string `json:"message"`
}