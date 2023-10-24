package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"conectivity-checker-wizard/models"
)

func getCiliumNetworkPolicy(namespace string) models.CiliumResponse {
	var res models.CiliumResponse
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
