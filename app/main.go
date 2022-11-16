package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/margostino/lagom/common"
	"github.com/margostino/lagom/configuration"
	"github.com/margostino/lagom/io"
	"github.com/margostino/lagom/loadgen"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg *sync.WaitGroup

var totalRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "lagom_requests_total",
	Help: "The total number of requests",
})

var config = configuration.GetConfig("./config.yml")
var updatedParams configuration.Params

func main() {
	go run()
	router := mux.NewRouter()
	router.Path("/metrics").Handler(promhttp.Handler()).Methods("GET")
	router.HandleFunc("/configuration", updateConfiguration).Methods("POST")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}

func updateConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewDecoder(request.Body).Decode(&updatedParams)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}
}

func run() {
	var payload []byte
	var err error
	wg = common.WaitGroup(1)

	if config.Http.Method == loadgen.POST {
		requestFile := config.Http.RequestFile
		payload, err = io.OpenFile(requestFile)
		if err != nil {
			log.Println(fmt.Sprintf("Cannot open file %s", requestFile), err)
			os.Exit(1)
		}
	} else if config.Http.Method == loadgen.GET {
		payload = nil
	}

	go execute(config, payload)
	wg.Wait()
}

func execute(config *configuration.Configuration, payload []byte) {
	end := config.Params.RunTime.Seconds()
	//var spawnRate, waitTime int
	var spawnRate int
	var requestsCount = 0
	var limiter = rate.NewLimiter(rate.Limit(config.Params.SpawnRate), config.Params.SpawnRate)
	var client = &http.Client{
		Timeout: time.Millisecond * 300,
	}

	for start := time.Now(); time.Since(start).Seconds() <= end; {

		if updatedParams.SpawnRate != 0 {
			spawnRate = updatedParams.SpawnRate
			limiter.SetLimit(rate.Limit(spawnRate))
			limiter.SetBurst(spawnRate)
		}

		//if updatedParams.BufferTime != 0 {
		//	waitTime = updatedParams.BufferTime
		//} else {
		//	waitTime = config.Params.BufferTime
		//}

		//err := limiter.WaitN(context.TODO(), spawnRate)
		//if err != nil {
		//	log.Printf("rate limitted when requests %d: %s\n", requestsCount, err.Error())
		//} else {
		//	totalRequests.Inc()
		//	requestsCount += 1
		//	go call(config.Http, payload, requestsCount, client)
		//}

		reservation := limiter.ReserveN(time.Now(), 1)
		if !reservation.OK() {
			// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		} else {
			println(reservation.Delay())
			time.Sleep(reservation.Delay())
			totalRequests.Inc()
			requestsCount += 1
			go call(config.Http, payload, requestsCount, client)
		}

		//if limiter.AllowN(time.Now(), spawnRate) == true {
		//	totalRequests.Inc()
		//	requestsCount += 1
		//	go call(config.Http, payload, requestsCount, client)
		//	wait(100)
		//} else {
		//	log.Printf("rate limitted when requests %d\n", requestsCount)
		//}
	}
	wg.Done()
}

func wait(bufferTime int) {
	//waitTime := time.Duration(rand.Intn(1) + maxStepTime)
	waitTime := time.Duration(bufferTime)
	time.Sleep(waitTime * time.Millisecond)
}

func call(config *configuration.Http, data []byte, requestsCount int, client *http.Client) {
	start := time.Now()
	payload := io.ReadAll(data)
	request := loadgen.GetRequest(config, payload)
	response, err := loadgen.Call(client, request)

	if err != nil {
		log.Printf("failure call (request #%d): %s\n", requestsCount, err.Error())
	}

	if response != nil {
		end := time.Now()
		log.Printf("URL %s Elapsed time %s with status: %s (total requests: %d)\n", config.Url, end.Sub(start).String(), response.Status, requestsCount)
		//if response.StatusCode == http.StatusOK {
		//	bodyBytes, err := ioutil.ReadAll(response.Body)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	bodyString := string(bodyBytes)
		//	log.Println(bodyString)
		//}
	}
}
