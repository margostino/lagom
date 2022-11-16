package loader

import (
	"bytes"
	"fmt"
	"github.com/margostino/lagom/common"
	"github.com/margostino/lagom/configuration"
	"github.com/margostino/lagom/io"
	"github.com/margostino/lagom/monitoring"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type LoadGen struct {
	waitGroup     *sync.WaitGroup
	httpConfig    *configuration.Http
	perfConfig    *configuration.Params
	updatedConfig *configuration.Params
	httpClient    *http.Client
	request       *http.Request
}

func NewLoadGen(config *configuration.Configuration) *LoadGen {
	requestData := getPayload(config.Http.RequestFile)
	payload := io.ReadAll(requestData)
	return &LoadGen{
		waitGroup:     common.WaitGroup(1),
		httpConfig:    config.Http,
		perfConfig:    config.Params,
		updatedConfig: nil,
		httpClient: &http.Client{
			Timeout: time.Millisecond * 300,
		},
		request: buildRequest(config.Http, payload),
	}
}

func buildRequest(config *configuration.Http, payload *bytes.Buffer) *http.Request {
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

func buildRateLimiter(rps int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(rps), rps)
}

func getPayload(requestFile string) []byte {
	var payload []byte
	var err error

	if requestFile != "" {
		payload, err = io.OpenFile(requestFile)
		if err != nil {
			log.Println(fmt.Sprintf("Cannot open file %s", requestFile), err)
			os.Exit(1)
		}
	} else {
		payload = nil
	}

	return payload
}

func (l *LoadGen) call(requestsCount int) {
	start := time.Now()
	response, err := l.httpClient.Do(l.request)
	if err != nil {
		log.Println(err.Error())
	}
	if err != nil {
		log.Printf("failure call (request #%d): %s\n", requestsCount, err.Error())
	}

	if response != nil {
		end := time.Now()
		log.Printf("URL %s Elapsed time %s with status: %s (total requests: %d)\n", l.httpConfig.Url, end.Sub(start).String(), response.Status, requestsCount)
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

func (l *LoadGen) Run() {
	var spawnRate int
	var requestsCount = 0

	var limiter = buildRateLimiter(l.perfConfig.SpawnRate)
	var runTime = l.perfConfig.RunTime.Seconds()

	for start := time.Now(); time.Since(start).Seconds() <= runTime; {
		if l.updatedConfig.SpawnRate != 0 {
			spawnRate = l.updatedConfig.SpawnRate
			limiter.SetLimit(rate.Limit(spawnRate))
			limiter.SetBurst(spawnRate)
		} else {
			spawnRate = l.perfConfig.SpawnRate
		}

		reservation := limiter.ReserveN(time.Now(), 1)
		if !reservation.OK() {
			// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		} else {
			time.Sleep(reservation.Delay())
			monitoring.Report()
			requestsCount += 1
			runtime := time.Since(start).Seconds()
			partialRate := int(float64(requestsCount) / runtime)
			if partialRate <= spawnRate {
				println(partialRate)
			}

			go l.call(requestsCount)
		}

	}

}
