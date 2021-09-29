package main

import (
	"encoding/json"
	"net/http"
	"log"
	"bytes"
	"io/ioutil"
	"fmt"
)

var (
	APIEndpoint = "http://localhost:8000"
	APIKey      = "yeet"
)

func apiDeploy(vapp string, variants []string) string {
	deployObj := DeployRequest{
		VApp:     vapp,
		Variants: variants,
	}
	jsonData, err := json.Marshal(deployObj)
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}

	request, err := http.NewRequest("POST", APIEndpoint+"/deploy", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "yeet")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return ""
}
