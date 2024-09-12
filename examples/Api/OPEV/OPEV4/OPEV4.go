package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// in this user story we will add only simulated CEM, and we will connect it to external EV

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("go run OPEV4.go <remoteSKI>")
		return
	}
	// reset the simulation
	// simulationReset()
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 1,
	}
	sendRequest("POST", "/sim", simulationData)

	resp := sendRequest("GET", "/cem", nil)
	cemSKI := resp.(map[string]any)["ski"]
	fmt.Printf("CEM ski is %s\n", cemSKI)

	ski := map[string]any{
		"remoteSki": os.Args[1],
	}
	// trusting with the remote SKI
	sendRequest("POST", "/cem/trust", ski)

	// getting external devices and its data
	var extDeviceAddress any
	for {
		resp := sendRequest("GET", "/cem/extDevices", nil)
		if resp != nil {
			for _, entity := range resp.([]map[string]any) {
				if !entity["isExternal"].(bool) {
					continue
				}
				switch entity["type"] {
				case "EVSE":
					extDeviceAddress = entity["address"]
					body := map[string]any{
						"entityAddressType": extDeviceAddress,
						"reqData":           "evseManufacturerData",
					}
					resp = sendRequest("POST", "/cem/extDevices/data", body)
					if resp != nil {
						fmt.Println(resp)
					}
				case "EV":
					extDeviceAddress = entity["address"]
					body := map[string]any{
						"entityAddressType": extDeviceAddress,
						"reqData":           "currentLimits",
					}
					resp = sendRequest("POST", "/cem/extDevices/data", body)
					if resp != nil {
						fmt.Println(resp)
					}
					extDeviceAddress = entity["address"]
					body = map[string]any{
						"entityAddressType": extDeviceAddress,
						"reqData":           "currentMeasurements",
					}
					resp = sendRequest("POST", "/cem/extDevices/data", body)
					if resp != nil {
						fmt.Println(resp)
					}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func sendRequest(method string, url string, payload any) any {
	client := &http.Client{}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, BaseIPAddress+url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result []map[string]any

	err = json.Unmarshal(body, &result)
	if err == nil {
		return result
	} else {
		var result2 map[string]any
		err = json.Unmarshal(body, &result2)
		if err != nil {
			return nil
		}
		return result2
	}

}
