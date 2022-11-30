package main

import (
	"github.com/jamespearly/loggly"
	"log"
	"main.go/utils"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load variables
	config, err := utils.LoadConfig("./", "config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate the loggly client and the http client
	logglyClient := loggly.New("Weather-App")
	weatherClient := http.Client{Timeout: time.Second * 2}
	awsClient := utils.CreateDynamoDBClient()

	// Set up the AWS table
	utils.SetUpTableAWS(config.TableName, awsClient, logglyClient)

	count := 0
	for true {
		for _, city := range config.Query {
			// Create a new request using http
			request := utils.CreateRequest(config.Url, config.Method, os.Getenv("ACCESS_KEY"), city)

			// Send get requests to the provided API
			body := utils.GetDataFromAPI(request, weatherClient, logglyClient)

			// Unmarshall the body response, print the data on the console, and send success message to loggly with response size
			data := utils.UnmarshallData(body, logglyClient)

			// Put the data into the DynamoDB on AWS cloud
			utils.PutItemInput(config.TableName, data, awsClient, logglyClient)

			// Count the request sending times
			count++
			log.Printf("This is the time %v the GET request is sent\n", count)
		}
		time.Sleep(config.Time * time.Hour)
	}
}
