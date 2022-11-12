package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/margostino/lagom/common"
	"github.com/margostino/lagom/configuration"
	"github.com/margostino/lagom/http-client"
	"github.com/margostino/lagom/io"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg *sync.WaitGroup
var opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "lagom_requests_total",
	Help: "The total number of requests",
})

func main() {
	go run()
	router := mux.NewRouter()
	router.Path("/metrics").Handler(promhttp.Handler())
	serverErr := http.ListenAndServe(":9000", router)
	log.Fatal(serverErr)
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func run() {
	var payload []byte
	var err error
	var config = configuration.GetConfig("./config.yml")
	wg = common.WaitGroup(1)

	if config.Http.Method == http_client.POST {
		requestFile := config.Http.RequestFile
		payload, err = io.OpenFile(requestFile)
		if err != nil {
			log.Println(fmt.Sprintf("Cannot open file %s", requestFile), err)
			os.Exit(1)
		}
	} else if config.Http.Method == http_client.GET {
		payload = nil
	}

	go execute(config, payload)
	wg.Wait()
}

func execute(config *configuration.Configuration, payload []byte) {
	end := config.Params.RunTime.Seconds()
	for start := time.Now(); time.Since(start).Seconds() <= end; {
		go call(config.Http, payload)
		recordMetrics()
		wait(100)
	}
	wg.Done()
}

func wait(maxStepTime int) {
	waitTime := time.Duration(rand.Intn(1) + maxStepTime)
	time.Sleep(waitTime * time.Millisecond)
}

func call(config *configuration.Http, data []byte) {
	payload := io.ReadAll(data)
	client := http_client.GetClient()
	request := http_client.GetRequest(config, payload)
	start := time.Now()
	//RegisterTime("Request", requestId)
	response := http_client.Call(client, request)
	if response != nil {
		end := time.Now()
		fmt.Printf("URL %s Elapsed time %s with status: %s\n", config.Url, end.Sub(start).String(), response.Status)
		if response.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Println(bodyString)
		}
	}
	//wg.Done()
}
