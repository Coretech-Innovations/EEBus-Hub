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

// in this use story we are adding CEM Device, Heatpump and connecting them with each other
// After the connection is done, we create an OHPCF announcement with power limits and timing constraints:
// - power(goodApproximation): 15000 W
// - powerMaximum: 16000 W
// - activeDurationMinimum: 5 seconds
// - pauseDurationMinimum: 5 seconds
// - startTime: 5 seconds

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	fmt.Println("Starting the simulation session example for OHPCF use case")
	// reset simulation session by deleting all components
	simulationReset()

	time.Sleep(5 * time.Second)

	// Get CEM Ski
	resp := sendRequest("GET", "/cem", nil)
	cemSKI := resp.(map[string]any)["ski"]

	// Add EVSE
	HeatPumpInfo := map[string]any{
		"deviceName":  "Coretech Heat Pump",
		"deviceCode":  "a7171fc2",
		"deviceModel": "HeatPump-1",
		"brandName":   "Coretech",
		"vendor": map[string]any{
			"name": "Coretech",
			"code": "CT",
		},
		"softwareRev": "0",
		"hardwareRev": "0",
		"manufacturer": map[string]any{
			"label":       "Coretech Innovations",
			"description": "",
		},
		"serialNumber":        "2a7ae968",
		"approveWriteLimit":   true,
		"failsafeValue":       5000,
		"failsafeDuration":    2,
		"failSafeDurationMax": 24,
		"nominalPower": map[string]any{
			"min": 0,
			"max": 23000,
		},
		"nominalCurrent": map[string]any{
			"min": 0,
			"max": 32,
		},
		"manufacturerDescription": "Heat Pump",
	}

	// sending request to add Heat Pump
	resp = sendRequest("POST", "/heatpump/add", HeatPumpInfo)
	heatPumpID := resp.(map[string]any)["id"]
	fmt.Printf("A new Heat Pump Device is added with ID %d\n", int(heatPumpID.(float64)))

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 60,
	}
	sendRequest("POST", "/sim", simulationData)

	endPoint := fmt.Sprintf("/heatpump/%d/cem", int(heatPumpID.(float64)))
	// Running the Heat Pump to connect with the CEM
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": cemSKI,
	})
	heatPumpSKI := resp.(map[string]any)["ski"]

	// trusting the Heat Pump from the CEM Side
	sendRequest("POST", "/cem/trust", map[string]any{
		"remoteSKI": heatPumpSKI,
	})
	fmt.Println("Ski = ", heatPumpSKI, " is trusted by the CEM")

	var announcmentInfo map[string]any = map[string]any{
		"isPausable":               true,
		"isStoppable":              true,
		"power":                    15000,
		"powerMax":                 16000,
		"powerType":                "both",
		"activeDurationMinSeconds": 5,
		"pauseDurationMinSeconds":  5,
	}

	time.Sleep(5 * time.Second)

	// send request to create an OHPCF announcement
	endPoint = fmt.Sprintf("/heatpump/%d/announceOptional", int(heatPumpID.(float64)))
	resp = sendRequest("POST", endPoint, announcmentInfo)
	fmt.Println("Announcement is created with following parameters :")
	for key, value := range announcmentInfo {
		fmt.Printf("%s = %v\n", key, value)
	}

	resp = sendRequest("GET", "/heatpump/list", nil)
	state := resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))

	// send request to make CEM send a start time to the Heat Pump
	resp = sendRequest("GET", "/heatpump/list", nil)
	deviceAddress := resp.([]map[string]any)[0]["deviceAddress"]

	time.Sleep(3 * time.Second)

	startime := 5
	endPoint = fmt.Sprintf("/cem/SendOhpcfStartTime/%s", deviceAddress.(string))
	resp = sendRequest("POST", endPoint, map[string]any{
		"startTime": float64(startime),
	})

	fmt.Println("CEM sent start time with value 5 seconds to the Heat Pump")
	resp = sendRequest("GET", "/heatpump/list", nil)
	state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))

	for i := 0; i < startime; i++ {
		resp = sendRequest("GET", "/heatpump/list", nil)
		state := resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
		fmt.Printf("Heat Pump OHPCF current state: %s, %d seconds passed\n", state.(string), i+1)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	resp = sendRequest("GET", "/heatpump/list", nil)
	state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
	fmt.Printf("Start time is reached, Heat Pump OHPCF current state: %s\n", state.(string))

	// send request to pause the Heat Pump OHPCF session
	fmt.Println("CEM can not pause the Heat Pump OHPCF session because minimum active duration is not reached")
	for i := 0; i < 5; i++ {
		fmt.Println("Waiting minimum active duration to be reached...", i+1, "seconds passed")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	endPoint = fmt.Sprintf("/cem/ChangeOhpcfState/%s", deviceAddress.(string))
	resp = sendRequest("PATCH", endPoint, map[string]any{
		"state": "paused",
	})
	fmt.Println("CEM sent pause command to the Heat Pump")


	resp = sendRequest("GET", "/heatpump/list", nil)
	state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))

	// send request to resume the Heat Pump OHPCF session
	fmt.Println("CEM can not resume the Heat Pump OHPCF session because minimum pause duration is not reached")
	for i := 0; i < 5; i++ {
		fmt.Println("Waiting minimum pause duration to be reached...", i+1, "seconds passed")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	endPoint = fmt.Sprintf("/cem/ChangeOhpcfState/%s", deviceAddress.(string))
	resp = sendRequest("PATCH", endPoint, map[string]any{
		"state": "running",
	})

	fmt.Println("CEM sent resume command to the Heat Pump")

	resp = sendRequest("GET", "/heatpump/list", nil)
	state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))

	// wait for the heat pump to complete its optional power consumption process and transition to completed state
	fmt.Println("Waiting the Heat Pump to complete its optional power consumption process...")

	for {
		resp = sendRequest("GET", "/heatpump/list", nil)
		state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
		if state.(string) == "completed" {
			break
		}
	}
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))

	// wait for the heat pump to clear its power sequence and transition to inactive state
	fmt.Println("Waiting the Heat Pump to clear its power sequence and transition to inactive state...")

	for {
		resp = sendRequest("GET", "/heatpump/list", nil)
		state = resp.([]map[string]any)[0]["DeviceInfo"].(map[string]any)["state"]
		if state.(string) == "inactive" {
			break
		}
	}
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state.(string))
}

func simulationReset() {
	// clear simulation
	simulationData := map[string]any{
		"action":      "reset",
		"speedFactor": 0,
	}
	sendRequest("POST", "/sim", simulationData)
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
	if len(body) == 0 {
		fmt.Printf("Empty response body (status=%d) for %s %s\n", resp.StatusCode, method, url)
		return nil
	}
	var result []map[string]any

	err = json.Unmarshal(body, &result)
	if err == nil {
		return result
	} else {
		var result2 map[string]any
		err = json.Unmarshal(body, &result2)
		if err != nil {
			fmt.Printf("Non-JSON response (status=%d) for %s %s: %s\n", resp.StatusCode, method, url, string(body))
			return nil
		}
		return result2
	}

}
