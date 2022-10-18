package main

import (
	"fmt"
	"github.com/margostino/lagom/common"
	"github.com/margostino/lagom/http"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg *sync.WaitGroup

func main() {
	var clients = common.GetConfig("./config.yml")
	delta := getDelta(clients)
	wg = common.WaitGroup(delta)
	for _, client := range clients {
		if !client.Enabled {
			log.Println(fmt.Sprintf("Client %s is not enabled", client.Url))
			wg.Add(-client.CallsNumber)
		} else if client.Method == http.POST {
			data, err := io.OpenFile(client.RequestFile)
			if err != nil {
				log.Println(fmt.Sprintf("Cannot open file %s", client.RequestFile), err)
				wg.Add(-client.CallsNumber)
			} else {
				go call(client, data)
			}
		} else {
			go call(client, nil)
		}
	}
	wg.Wait()
}

func getDelta(clients []*common.Client) int {
	var delta = 0
	for _, client := range clients {
		delta += client.CallsNumber
	}
	return delta
}

func progressiveCall(config *common.Client, data []byte) {
	calls := 0
	for i := 0; i < config.CallsNumber; i++ {
		calls += i + 1
		if calls > config.CallsNumber {
			break
		}
		for j := 0; j < calls; j++ {
			go execute(i, config, data)
			wait(50)
		}
		wait(config.MaxStepTime)
	}
}

func call(config *common.Client, data []byte) {
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
