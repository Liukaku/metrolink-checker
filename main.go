package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	message string
	apiUrl string
}

type MetroStop struct {
	StationLocation string `json:"StationLocation"`
	AtcoCode        string `json:"AtcoCode"`
	Direction       string `json:"Direction"`
	Dest0           string `json:"Dest0"`
}

type Result struct {
	Value        map[string]string 
}


func main() {
	envs := loadEnv()
	fmt.Println(envs.message)

	var realRequest string
	stopId := "130818"
	
	fmt.Println("Send a real request? Y/n")
	fmt.Scanln(&realRequest)

	if realRequest == "n" {
		stopId = "1"
	}
	
	respBody := apiRequest(envs.apiUrl, stopId)

	handleResponse(respBody)

}

func checkForError(resp map[string]interface{}, count int16) string {
	if count >= 10 {
		return "Error fetching HTTP error cause"
	}

	if resp["errorMessage"] != nil {
		return resp["errorMessage"].(string)
	} 
	
	newResp := resp["cause"].(map[string]interface{})
	return checkForError(newResp, count + 1)

}

func apiRequest(url string, stopId string)[]byte{
	pb, _ := json.Marshal(map[string]string {
		// "stop": "130211",
		"stop": stopId,
	})

	postBody := bytes.NewBuffer(pb)

	resp, err := http.Post(url, "application/json", postBody)

	if err != nil {
		fmt.Print("http error")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print("body reading error")
	}

	return respBody
}

func handleResponse(respBody []byte){
	var str string
	var dat MetroStop
	var apiErr map[string]interface {}

	json.Unmarshal(respBody, &str)

	jsonErr := json.Unmarshal([]byte(str), &dat)
	
	if jsonErr != nil {
		json.Unmarshal(respBody, &apiErr)
		apiErrRes := checkForError(apiErr, 0)
		fmt.Printf("Error: %s ",apiErrRes)
	} else {
		fmt.Printf("Tram Stop: %s", dat.StationLocation)
	}
}

func loadEnv() Envs{
	err := godotenv.Load()

	if err != nil {
		fmt.Print(".env error")
	}

	var message string = os.Getenv("TEST")
	var url string = os.Getenv("API_URL")

	retVal := Envs{message: message, apiUrl: url}

	return retVal
}