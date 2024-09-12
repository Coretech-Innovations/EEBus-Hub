package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// in this use story we will add 2 EVs to the system, both max current under the available limits

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	// reset simulation session by deleting all components
	simulationReset()

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 15,
	}
	sendRequest("POST", "/sim", simulationData)

	resp := sendRequest("GET", "/cem", nil)
	cemSKI := resp.(map[string]any)["ski"]
	fmt.Printf("CEM ski is %s\n", cemSKI)

	// adding EVSE
	evseInfo := map[string]any{
		"deviceName": "Coretech EVSE WallBox",
		"deviceCode": "0001",
		"vendor": map[string]any{
			"name": "Coretech-Innovations",
			"code": "60745",
		},
		"softwareRev": "0",
		"hardwareRev": "0",
		"brandName":   "Coretech-Innovations",
		"Manufacturer": map[string]any{
			"label":       "Coretech-Innovations",
			"description": "Charging Station",
		},
		"deviceModel":  "EVSE",
		"serialNumber": "00000002",
		"port":         4712,
	}
	// sending request to add EVSE
	resp = sendRequest("POST", "/evse/add", evseInfo)
	evseID := int(resp.(map[string]any)["id"].(float64))

	fmt.Printf("A new EVSE Device is added with ID %d\n", evseID)

	endPoint := fmt.Sprintf("/evse/%d/cem", evseID)
	// Running the EVSE to connect with the CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": cemSKI,
	})
	evseSKI := resp.(map[string]any)["ski"]

	// trusting the EVSE from the CEM Side
	sendRequest("POST", "/cem/trust", map[string]any{
		"remoteSKI": evseSKI,
	})

	nonControllableDevice := map[string]any{
		"current": map[string]any{
			"a": 24.0,
			"b": 24.0,
			"c": 24.0,
		},
		"powerStateOn": true,
	}
	sendRequest("POST", "/uncontrollabledevice", nonControllableDevice)

	// add EV
	var EV1 map[string]any = map[string]any{
		"asymmetricCharging": false,
		"currentLimits": map[string]int{
			"min": 6,
			"max": 10,
		},
		"dischargingEnable": false,
		"chargingEnable":    true,
		"chargingCapacity":  80,
		"charged":           20,
		"batteryHealth":     100,
	}
	resp = sendRequest("POST", "/ev/add", EV1)
	ev1ID := int(resp.(map[string]any)["id"].(float64))
	fmt.Printf("A new EV Device is added with ID %d\n", ev1ID)
	// connect EV to EVSE
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", ev1ID, evseID)
	sendRequest("POST", endPoint, nil)

	time.Sleep(5 * time.Second)
	// creating second EV
	var EV2 map[string]any = map[string]any{
		"asymmetricCharging": false,
		"currentLimits": map[string]int{
			"min": 10,
			"max": 20,
		},
		"dischargingEnable": false,
		"chargingEnable":    true,
		"chargingCapacity":  80,
		"charged":           20,
		"batteryHealth":     100,
	}
	resp = sendRequest("POST", "/ev/add", EV2)
	ev2ID := int(resp.(map[string]any)["id"].(float64))
	fmt.Printf("A new EV Device is added with ID %d\n", ev2ID)
	// connect EV to EVSE
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", ev2ID, evseID)
	sendRequest("POST", endPoint, nil)
	for {
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", ev1ID)
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.([]map[string]any))
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", ev2ID)
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.([]map[string]any))
		time.Sleep(2 * time.Second)
	}
}

func simulationReset() {
	// delete all the EVSEs
	resp := sendRequest("GET", "/evse/list", nil)
	evses := resp.([]map[string]any)
	for _, evse := range evses {
		endpoint := fmt.Sprintf("/evse/%d", int(evse["evseId"].(float64)))
		sendRequest("DELETE", endpoint, nil)
	}
	// delete all the EVs
	resp = sendRequest("GET", "/ev/list", nil)
	evs := resp.([]map[string]any)
	for _, ev := range evs {
		endpoint := fmt.Sprintf("/ev/%d", int(ev["evId"].(float64)))
		sendRequest("DELETE", endpoint, nil)
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
			log.Fatal("Can not convert to map")
			return nil
		}
		return result2
	}

}
