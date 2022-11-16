package loadgen

import (
	"bytes"
	"fmt"
	"github.com/margostino/lagom/configuration"
	"log"
	"net/http"
	"os"
)

const (
	POST = "POST"
	GET  = "GET"
)

func GetRequest(config *configuration.Http, payload *bytes.Buffer) *http.Request {
	request, err := http.NewRequest(config.Method, config.Url, payload)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if config.ContentType != "" {
		request.Header.Add("Content-Type", config.ContentType)
	}
	if config.Username != "" && config.Password != "" {
		request.SetBasicAuth(config.Username, config.Password)
	}

	return request
}

func Call(client *http.Client, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
	}
	return response, err
}
