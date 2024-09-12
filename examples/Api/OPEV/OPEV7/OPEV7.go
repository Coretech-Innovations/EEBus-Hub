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

// In this user story, We will add two EVs and connect them with external CEM, each EV will be connected to an EVSE

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("go run OPEV7.go <remoteSKI>")
		return
	}
	// reset simulation session by deleting all components
	simulationReset()

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 15,
	}
	sendRequest("POST", "/sim", simulationData)

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
	resp := sendRequest("POST", "/evse/add", evseInfo)
	evseID := int(resp.(map[string]any)["id"].(float64))

	endPoint := fmt.Sprintf("/evse/%d/cem", evseID)
	// Running the EVSE to connect with the external CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": os.Args[1],
	})
	evseSKI := resp.(map[string]any)["ski"]
	fmt.Printf("A new EVSE Device is added with ID %d\nlocal SKI: %d\n", evseID, evseSKI)

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

	// adding 2nd EVSE

	evseInfo = map[string]any{
		"deviceName": "Coretech EVSE WallBox 2",
		"deviceCode": "0002",
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
		"serialNumber": "00000003",
		"port":         4715,
	}
	resp = sendRequest("POST", "/evse/add", evseInfo)
	evse2ID := int(resp.(map[string]any)["id"].(float64))

	// connecting 2nd EVSE with CEM
	endPoint = fmt.Sprintf("/evse/%d/cem", evse2ID)
	// Running the EVSE to connect with the external CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": os.Args[1],
	})
	evse2SKI := resp.(map[string]any)["ski"]
	fmt.Printf("A new EVSE Device is added with ID %d\nlocal SKI: %d\n", evse2ID, evse2SKI)

	// creating second EV
	var EV2 map[string]any = map[string]any{
		"asymmetricCharging": false,
		"currentLimits": map[string]int{
			"min": 8,
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
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", ev2ID, evse2ID)
	sendRequest("POST", endPoint, nil)
	time.Sleep(5 * time.Second)
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
	// delete the CEM
	sendRequest("DELETE", "/cem", nil)
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
