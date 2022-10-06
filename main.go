package main

import (
	structure "CSC482/modules"
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func sendingLoggy(client *loggly.ClientType, msgType string, msg string) {
	err := client.EchoSend(msgType, msg)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var tag string
	tag = "My-Go-Demo"

	// Instantiate the client
	client := loggly.New(tag)

	// Get tokens from .env/ choose query for API
	accessKey := os.Getenv("ACCESS_KEY")
	location := "Oswego"

	// Create a new request using http
	url := "https://api.weatherstack.com/current"
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("access_key", accessKey)
	q.Add("query", location)
	req.URL.RawQuery = q.Encode()

	// Send req using http Client
	weatherClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	resp, sendErr := weatherClient.Do(req)
	if sendErr == nil {
		sendingLoggy(client, "info", "Successfully send the request to API")
	} else {
		sendingLoggy(client, "error", "Failed with error: "+sendErr.Error())
	}

	//Read the response body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		sendingLoggy(client, "info", "Successfully read the response body")
	} else {
		sendingLoggy(client, "error", "Failed with error: "+readErr.Error())
	}

	//Unmarshall the response into the data structure
	var data structure.Data
	unmarshallErr := json.Unmarshal(body, &data)
	if unmarshallErr != nil {
		return
	}

	// Print the data on the console
	formattedData, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(formattedData))

	// Send success message to loggly with response size
	var respSize = strconv.Itoa(len(body))
	sendingLoggy(client, "info", "Successful data collection of size: "+respSize)
}
