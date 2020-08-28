package responses

import (
	"encoding/json"
	"net/http"
)

type GenericResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type StorageResizeReq struct {
	Namespace string `json:"namespace"`
}

func ResponseWithFailedMessage(statusCode int, message interface{}, writer *http.ResponseWriter) {
	resp := GenericResponse{Success: false, Data: message}
	(*writer).Header().Add("Content-Type", "application/json")
	(*writer).WriteHeader(statusCode)
	json.NewEncoder(*writer).Encode(resp)
}
