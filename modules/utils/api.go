package utils

import (
	"encoding/json"
	"github.com/jamespearly/loggly"
	"io/ioutil"
	"log"
	"main.go/structs"
	"net/http"
	"strconv"
)

func CreateRequest(url string, method string, accessKey string, location string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	q := req.URL.Query()
	q.Add("access_key", accessKey)
	q.Add("query", location)
	req.URL.RawQuery = q.Encode()
	return req
}

func sendingLoggy(client *loggly.ClientType, msgType string, msg string) {
	err := client.EchoSend(msgType, msg)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDataFromAPI(request *http.Request, weatherClient http.Client, logglyClient *loggly.ClientType) []byte {
	// Send req using http Client
	resp, sendErr := weatherClient.Do(request)
	HandleException(logglyClient, sendErr, "Successfully send the request to API")

	//Read the response body
	body, readErr := ioutil.ReadAll(resp.Body)
	HandleException(logglyClient, readErr, "Successfully collect the data and read the response body")

	return body
}

func UnmarshallData(body []byte, logglyClient *loggly.ClientType) structs.Data {
	//Unmarshall the response into the data structure
	var data structs.Data
	unmarshallErr := json.Unmarshal(body, &data)
	HandleException(logglyClient, unmarshallErr, "Successfully unmarshall the response body into the data structure")

	// Print the data on the console
	formattedData, _ := json.MarshalIndent(data, "", "    ")
	log.Println(string(formattedData))

	// Send success message to loggly with response size
	var respSize = strconv.Itoa(len(body))
	sendingLoggy(logglyClient, "info", "Successful data collection of size: "+respSize)

	return data
}
