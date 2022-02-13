package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
	"log"
	"errors"
	"net/http"
)

var (
	APIEndpoint = "http://172.16.1.121:8000"
	APIKey      = "yeet"
)

const (
	INVALID_USER = "Invalid user"
)

func apiCheckUser(username string) error {
	request, err := http.NewRequest("GET", APIEndpoint+"/user/" + username, nil)
	request.Header.Set("X-Api-Key", "yeet")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	if strings.Contains(string(body), "valid\":false") {
		return errors.New(INVALID_USER)
	}
	return nil
}

func apiDeploy(template, catalog string, variants []string) error {
	deployObj := DeployRequest{
		Template:  template,
		Catalog:   catalog,
		Variants:  variants,
		MakeOwner: true,
	}
	jsonData, err := json.Marshal(deployObj)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}

	fmt.Println("json data is", string(jsonData))

	request, err := http.NewRequest("POST", APIEndpoint+"/deploy", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "yeet")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	if response.StatusCode != 200 {
		return errors.New("Deploy failed: " + string(body))
	}
	return nil
}
