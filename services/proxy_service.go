package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CiliumResponse struct {
	IsHostname  bool `json:"isHostname"`
	IsIPAddress bool `json:"isIPAddress"`
}

func getCiliumNetworkPolicy(namespace string) CiliumResponse {
	var res CiliumResponse
	url := fmt.Sprintf("http://localhost:9090/api-sever/%s/my-cnp", namespace)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	response, _ := client.Do(req)
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	fmt.Println("Response:", res)
	return res
}
