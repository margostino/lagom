package http

import (
	"bytes"
	"fmt"
	"github.com/margostino/lagom/common"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	POST = "POST"
	GET  = "GET"
)

func GetClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func GetRequest(config *common.Client, payload *bytes.Buffer) *http.Request {
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

func Call(client *http.Client, request *http.Request) *http.Response {
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(error.Error())
		//log.Fatal(error)
	}
	return response
}

func Print(response *http.Response) {
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
}
