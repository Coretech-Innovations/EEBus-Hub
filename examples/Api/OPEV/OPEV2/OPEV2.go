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

// in this use story we are adding CEM Device, EVSE and connecting them with each other
// after connection is done we are adding an EV with current limits (min:6, max:10) on the three phases
// the fuse limit is set to 40A
// the EV should be supplied with current 10A on all its phases as the available current is greater than 10A
// after that we will add nonControllable device that will absorb all the current which will lead to stop charging the EV

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	// reset the simulation
	simulationReset()

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 100,
	}
	sendRequest("POST", "/sim", simulationData)

	// Get CEM SKI
	resp := sendRequest("GET", "/cem", nil)
	cemSKI := resp.(map[string]any)["ski"]

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
	evseID := resp.(map[string]any)["id"]

	fmt.Printf("A new EVSE Device is added with ID %d\n", int(evseID.(float64)))

	endPoint := fmt.Sprintf("/evse/%d/cem", int(evseID.(float64)))
	// Running the EVSE to connect with the CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": cemSKI,
	})

	evseSKI := resp.(map[string]any)["ski"]
	// trusting the EVSE from the CEM Side
	sendRequest("POST", "/cem/trust", map[string]any{
		"remoteSKI": evseSKI,
	})

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

	time.Sleep(3 * time.Second)
	for i := 0; i < 3; i++ {
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", int(evID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.([]map[string]any))
		time.Sleep(2 * time.Second)
	}
	nonControllableDevice := map[string]any{
		"deviceName": "Home Appliances",
		"current": map[string]any{
			"a": 40.0,
			"b": 40.0,
			"c": 40.0,
		},
		"powerFactor":  1.0,
		"powerStateOn": true,
	}
	sendRequest("POST", "/uncontrollabledevice", nonControllableDevice)
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
