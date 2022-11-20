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
	ConfigChannel chan *configuration.Params
	httpClient    *http.Client
	request       *http.Request
}

func NewLoadGen(config *configuration.Configuration) *LoadGen {
	requestData := getPayload(config.Http.RequestFile)
	payload := io.ReadAll(requestData)

	loadgen := &LoadGen{
		waitGroup:     common.WaitGroup(1),
		httpConfig:    config.Http,
		perfConfig:    config.Params,
		ConfigChannel: make(chan *configuration.Params),
		updatedConfig: nil,
		httpClient: &http.Client{
			Timeout: time.Millisecond * 300,
		},
		request: buildRequest(config.Http, payload),
	}

	go loadgen.listenConfig()

	return loadgen
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

func (l *LoadGen) call(requestsCount int, calculatedRate float64, runtime float64, stepLoad int) {
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
		log.Printf("URL %s Elapsed time %dms with status %s (step load [rps]: %d, calculated rate [rps]: %.2f, runtime (sec): %.2f, total requests: %d)\n",
			l.httpConfig.Url,
			end.Sub(start).Milliseconds(),
			response.Status,
			stepLoad,
			calculatedRate,
			runtime,
			requestsCount)
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

func (l *LoadGen) listenConfig() {
	config := <-l.ConfigChannel
	if config.StepLoad != l.perfConfig.StepLoad {
		l.perfConfig.StepLoad = config.StepLoad
	}
}

func (l *LoadGen) Run() {
	var stepStart = time.Now()
	var requestsCount = 0
	var stepLoad = l.perfConfig.StepLoad

	var limiter = buildRateLimiter(stepLoad)
	var totalRuntime = l.perfConfig.RunTime

	for start := time.Now(); time.Since(start) <= totalRuntime; {
		if stepLoad != l.perfConfig.StepLoad {
			stepLoad = l.perfConfig.StepLoad
			limiter.SetLimit(rate.Limit(stepLoad))
			limiter.SetBurst(stepLoad)
		}

		reservation := limiter.ReserveN(time.Now(), 1)
		if !reservation.OK() {
			// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		} else {
			time.Sleep(reservation.Delay())
			// only report if call success
			monitoring.Report()
			requestsCount += 1
			runtime := time.Since(start).Seconds()
			stepRuntime := time.Since(stepStart)
			calculatedRate := float64(requestsCount) / runtime

			if stepRuntime >= l.perfConfig.StepTime && calculatedRate >= float64(stepLoad) {
				stepStart = time.Now()
				//stepLoad += 1
				stepLoad *= 2
				l.perfConfig.StepLoad = stepLoad
				limiter.SetLimit(rate.Limit(stepLoad))
				limiter.SetBurst(stepLoad)
			}

			go l.call(requestsCount, calculatedRate, runtime, stepLoad)
		}

	}

}
