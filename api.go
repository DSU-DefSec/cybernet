package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	APIEndpoint = "http://172.16.1.121:8000"
	APIKey      = "yeet"
)

func apiDeploy(template, catalog string, variants []string) string {
	deployObj := DeployRequest{
		Template:  template,
		Catalog:   catalog,
		Variants:  variants,
		MakeOwner: true,
	}
	jsonData, err := json.Marshal(deployObj)
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}

	fmt.Println("json data is", string(jsonData))

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
