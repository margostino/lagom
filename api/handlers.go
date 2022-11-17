package api

import (
	"encoding/json"
	"github.com/margostino/lagom/configuration"
	"log"
	"net/http"
)

type ApiHandler struct {
	configChannel chan *configuration.Params
}

func NewApiHandler(configChannel chan *configuration.Params) *ApiHandler {
	return &ApiHandler{
		configChannel: configChannel,
	}
}

func (a *ApiHandler) UpdateConfiguration(writer http.ResponseWriter, request *http.Request) {
	var params *configuration.Params
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewDecoder(request.Body).Decode(&params)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}
	a.configChannel <- params
}
