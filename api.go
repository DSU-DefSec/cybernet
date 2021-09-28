package main

import (
	"encoding/json"
	"fmt"
)

var (
	APIEndpoint = ""
	APIKey      = "yeet"
)

func apiDeploy(vapp string, variants []string) string {
	deployObj := DeployRequest{
		VApp:     vapp,
		Variants: variants,
	}
	var jsonData = json.Encode(deployObj)

	request, error := http.NewRequest("POST", APIEndpoint+"/deploy", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "yeet")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Println("ERROR:", error)
		return ""
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}
