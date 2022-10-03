package main

import (
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

type Request struct {
	Type     string `json:"type"`
	Query    string `json:"query"`
	Language string `json:"language"`
	Unit     string `json:"unit"`
}

type Location struct {
	Name           string `json:"name"`
	Country        string `json:"country"`
	Region         string `json:"region"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
	TimezoneID     string `json:"timezone_id"`
	Localtime      string `json:"localtime"`
	LocaltimeEpoch int    `json:"localtime_epoch"`
	UtcOffset      string `json:"utc_offset"`
}

type Current struct {
	ObservationTime     string   `json:"observation_time"`
	Temperature         int      `json:"temperature"`
	WeatherCode         int      `json:"weather_code"`
	WeatherIcons        []string `json:"weather_icons"`
	WeatherDescriptions []string `json:"weather_descriptions"`
	WindSpeed           int      `json:"wind_speed"`
	WindDegree          int      `json:"wind_degree"`
	WindDir             string   `json:"wind_dir"`
	Pressure            int      `json:"pressure"`
	Precip              int      `json:"precip"`
	Humidity            int      `json:"humidity"`
	Cloudcover          int      `json:"cloudcover"`
	Feelslike           int      `json:"feelslike"`
	UvIndex             int      `json:"uv_index"`
	Visibility          int      `json:"visibility"`
	IsDay               string   `json:"is_day"`
}

type Data struct {
	Request  Request  `json:"request"`
	Location Location `json:"location"`
	Current  Current  `json:"current"`
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
	url := "http://api.weatherstack.com/current"
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
		err := client.EchoSend("info", "Successfully send the request to API")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := client.EchoSend("error", "Failed with error: "+sendErr.Error())
		if err != nil {
			log.Fatalln(err)
		}
	}

	//Read the response body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		err := client.EchoSend("info", "Successfully read the response body")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := client.EchoSend("error", "Failed with error: "+readErr.Error())
		if err != nil {
			log.Fatalln(err)
		}
	}

	//Unmarshall the response into the data structure
	var data Data
	unmarshallErr := json.Unmarshal(body, &data)
	if unmarshallErr != nil {
		return
	}

	// Print the data on the console
	formattedData, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(formattedData))

	// Send success message to loggly with response size
	var respSize = strconv.Itoa(len(body))
	logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if logErr != nil {
		fmt.Println("err: ", logErr)
	}
}
