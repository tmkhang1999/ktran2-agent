package main

import (
	"CSC482/utils"
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func sendingLoggy(client *loggly.ClientType, msgType string, msg string) {
	err := client.EchoSend(msgType, msg)
	if err != nil {
		log.Fatalln(err)
	}
}

func createRequest(url string, method string, accessKey string, location string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	q := req.URL.Query()
	q.Add("access_key", accessKey)
	q.Add("query", location)
	req.URL.RawQuery = q.Encode()
	return req
}

func main() {
	// Load variables
	config, err := utils.LoadConfig("./", "config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate the loggly client and the http client
	logglyClient := loggly.New("Weather-App")
	weatherClient := http.Client{Timeout: time.Second * 2}

	// Create a new request using http
	request := createRequest(config.Url, config.Method, config.AccessKey, config.Location)

	count := 0
	for true {
		// Send req using http Client
		resp, sendErr := weatherClient.Do(request)
		if sendErr == nil {
			sendingLoggy(logglyClient, "info", "Successfully send the request to API")
		} else {
			sendingLoggy(logglyClient, "error", "Failed with error: "+sendErr.Error())
		}

		//Read the response body
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr == nil {
			sendingLoggy(logglyClient, "info", "Successfully read the response body")
		} else {
			sendingLoggy(logglyClient, "error", "Failed with error: "+readErr.Error())
		}

		//Unmarshall the response into the data structure
		var data utils.Data
		unmarshallErr := json.Unmarshal(body, &data)
		if unmarshallErr != nil {
			return
		}

		// Print the data on the console
		formattedData, _ := json.MarshalIndent(data, "", "    ")
		fmt.Println(string(formattedData))

		// Send success message to loggly with response size
		var respSize = strconv.Itoa(len(body))
		sendingLoggy(logglyClient, "info", "Successful data collection of size: "+respSize)

		// Count the request sending times
		count++
		fmt.Printf("This is the time %v the GET request is sent\n", count)
		time.Sleep(config.Time * time.Minute)
	}
}
