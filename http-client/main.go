package main

import (
	"fmt"
	"github.com/margostino/lagom/common"
	"github.com/margostino/lagom/configuration"
	"github.com/margostino/lagom/http"
	"github.com/margostino/lagom/io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg *sync.WaitGroup

func main() {
	var config = configuration.GetConfig("./config.yml")
	if !config.Enabled {
		log.Println(fmt.Sprintf("Client %s is not enabled", config.Http.Url))
	} else if config.Http.Method == http.POST {
		requestFile := config.Http.RequestFile
		data, err := io.OpenFile(requestFile)
		if err != nil {
			log.Println(fmt.Sprintf("Cannot open file %s", requestFile), err)
		} else {
			go call(config, data)
		}
	} else {
		go call(config, nil)
	}
	wg.Wait()
}

func call(config *configuration.Configuration, data []byte) {
	for i := 0; i < config.CallsNumber; i++ {
		go execute(i, config, data)
		wait(config.MaxStepTime + 500)
	}
}

func wait(maxStepTime int) {
	waitTime := time.Duration(rand.Intn(1) + maxStepTime)
	time.Sleep(waitTime * time.Millisecond)
}

func execute(requestId int, config *common.Client, data []byte) {
	payload := io.ReadAll(data)
	client := http.GetClient()
	request := http.GetRequest(config, payload)
	start := time.Now()
	//RegisterTime("Request", requestId)
	response := http.Call(client, request)
	if response != nil {
		end := time.Now()
		fmt.Printf("URL %s Request #%d Elapsed time %s with status: %s\n", config.Url, requestId, end.Sub(start).String(), response.Status)
		http.Print(response)
	}
	wg.Done()
}
