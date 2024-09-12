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

// in this user story we will add only simulated EV, and connect it to external CEM

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("go run OPEV6.go <remoteSKI>")
		return
	}

	// reset the simulation
	simulationReset()

	// add EVSE
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
	evseID := resp.(map[string]any)["id"]

	// Connecting the EVSE with external CEM
	endPoint := fmt.Sprintf("/evse/%d/cem", int(evseID.(float64)))
	// Running the EVSE to connect with the CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": os.Args[1],
	})
	evseSKI := resp.(map[string]any)["SKI"]
	fmt.Printf("A new EVSE Device is added with ID %d\nLocal Ski: %s\n", int(evseID.(float64)), evseSKI.(string))

	// add EV

	var EVEntity map[string]any = map[string]any{
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
	resp = sendRequest("POST", "/ev/add", EVEntity)
	evID := resp.(map[string]any)["id"]
	fmt.Printf("A new EV Device is added with ID %d\n", int(evID.(float64)))
	// connect EV to EVSE
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", int(evID.(float64)), int(evseID.(float64)))
	sendRequest("POST", endPoint, nil)

	for {
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", int(evID.(float64)))
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
