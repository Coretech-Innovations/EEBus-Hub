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

type MeasurementData struct {
	EnergyConsumed *float64                 `json:"energyConsumed,omitempty"`
	EnergyProduced *float64                 `json:"energyProduced,omitempty"`
	Power          []MeasurementDataElement `json:"power,omitempty"`
	Current        []MeasurementDataElement `json:"current,omitempty"`
	Voltage        []MeasurementDataElement `json:"voltage,omitempty"`
	Frequency      *float64                 `json:"frequency,omitempty"`
}

type MeasurementDataElement struct {
	Phase string   `json:"phase"`
	Value *float64 `json:"value,omitempty"`
}

type MPCResponseBody struct {
	DeviceAddress string          `json:"deviceAddress"`
	Measurements  MeasurementData `json:"measurements"`
}

func (m MPCResponseBody) String() string {
	return fmt.Sprintf(`
-----------------
Device Address: %v
%v
-----------------
`, m.DeviceAddress, m.Measurements)
}

func (m MeasurementData) String() string {
	var energyConsumed string
	var energyProduced string
	var frequency string
	if m.EnergyConsumed == nil {
		energyConsumed = "[]"
	} else {
		energyConsumed = fmt.Sprintf("%.3f", *m.EnergyConsumed)
	}
	if m.EnergyProduced == nil {
		energyProduced = "[]"
	} else {
		energyProduced = fmt.Sprintf("%.3f", *m.EnergyProduced)
	}
	if m.Frequency == nil {
		frequency = "[]"
	} else {
		frequency = fmt.Sprintf("%.3f", *m.Frequency)
	}
	return fmt.Sprintf(`=================
Power: %v
Energy Consumed: %v
Energy Produced: %v
Current: %v
Voltage: %v
Frequency: %v
=================`,
		m.Power, energyConsumed, energyProduced, m.Current, m.Voltage, frequency)
}

func (m MeasurementDataElement) String() string {
	if m.Value == nil {
		return "[]"
	}
	return fmt.Sprintf("[Phase: %v, Value: %.3f]", m.Phase, *m.Value)
}

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
	// Step 1: Reset simulation environment
	fmt.Println("Resetting the simulation Environment")
	SimResetResp, err := sendRequest[struct {
		Status string `json:"status"`
		Err    string `json:"err"`
	}]("POST", "/sim", map[string]any{
		"action":      "reset",
		"speedFactor": 1,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if SimResetResp.Status != "OK" {
		log.Fatal(SimResetResp.Err)
	}
	time.Sleep(3 * time.Second)

	// getting CEM Ski
	cemSkiResp, err := sendRequest[struct {
		Ski string `json:"ski"`
	}]("GET", "/cem", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(cemSkiResp.Ski) == 0 {
		log.Fatal("no CEM in the system")
	}
	cemSKI := cemSkiResp.Ski
	fmt.Printf("the CEM Ski is %s\n", cemSKI)

	// Step 2: Adding new EVSE
	evseInfo := map[string]any{
		"deviceName":  "Coretech EVSE WLBX",
		"deviceCode":  "037d42e1",
		"deviceModel": "Wallbox",
		"brandName":   "Coretech",
		"vendor": map[string]any{
			"name": "Coretech",
			"code": "60745",
		},
		"softwareRev": "0",
		"hardwareRev": "0",
		"Manufacturer": map[string]any{
			"label":       "Coretech",
			"description": "Charging Station",
		},
		"serialNumber":        "de07c278",
		"failsafeValue":       4320,
		"failsafeDuration":    2,
		"failSafeDurationMax": 24,
		"nominalPower": map[string]any{
			"min": 4320,
			"max": 23000,
		},
		"nominalCurrent": map[string]any{
			"min": 6,
			"max": 32,
		},
		"approveWriteLimit": true,
	}
	evseAddResponse, err := sendRequest[struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
		Error  any    `json:"err"`
	}]("POST", "/evse/add", evseInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	if evseAddResponse.Status != "OK" {
		log.Fatal(evseAddResponse.Error)
	}
	fmt.Printf("Adding new EVSE with id %d\n", evseAddResponse.ID)

	// Step 3: Add EV and connect it to the EVSE
	ev := map[string]any{
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
		"dischargeCurrent":  5,
	}
	evAddResponse, err := sendRequest[struct {
		Status string `json:"status"`
		ID     uint   `json:"id"`
		Error  string `json:"err"`
	}]("POST", "/ev/add", ev)
	if err != nil {
		log.Fatal(err.Error())
	}
	if evAddResponse.Status != "OK" {
		log.Fatal(evAddResponse.Error)
	}
	fmt.Printf("Adding new EV with id %d\n", evAddResponse.ID)

	endPoint := fmt.Sprintf("/ev/%d/evse/%d", evAddResponse.ID, evseAddResponse.ID)
	evConnectResponse, err := sendRequest[struct {
		Error  string `json:"err"`
		Status string `json:"status"`
	}]("POST", endPoint, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	if evConnectResponse.Status != "OK" {
		log.Fatal(evConnectResponse.Error)
	}
	fmt.Println("Connecting the EV to the EVSE")

	// Step 4: Pair the two devices with each other
	endPoint = fmt.Sprintf("/evse/%d/cem", evseAddResponse.ID)
	evsePairCEM, err := sendRequest[struct {
		Status string `json:"status"`
		Error  string `json:"err"`
		SKI    string `json:"ski"`
	}]("POST", endPoint, map[string]string{
		"remoteSki": cemSKI,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	if evsePairCEM.Status != "OK" {
		log.Fatal(evsePairCEM.Error)
	}
	fmt.Println("Trusting the CEM from the EVSE side")

	cemPairEVSE, err := sendRequest[struct {
		Status string `json:"status"`
		Error  string `json:"err"`
	}]("POST", "/cem/trust", map[string]string{
		"remoteSki": evsePairCEM.SKI,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	if cemPairEVSE.Status != "OK" {
		log.Fatal(cemPairEVSE.Error)
	}
	fmt.Println("Trusting the EVSE from the CEM side")

	// Step 5: start the simulation
	simStartResp, err := sendRequest[struct {
		Status string `json:"status"`
		Err    string `json:"err"`
	}]("POST", "/sim", map[string]any{
		"action":      "start",
		"speedFactor": 1,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if simStartResp.Status != "OK" {
		log.Fatal(simStartResp.Err)
	}
	fmt.Println("Starting the Simulation")

	// Wait till connection made
	for {
		devices, err := sendRequest[[]struct {
			DeviceAddress string `json:"deviceAddress"`
			DeviceSki     string `json:"ski"`
		}]("GET", "/cem/LpcDevicesData", nil)
		if err != nil {
			log.Fatal(err.Error())
		}
		id := -1
		for idx, dev := range devices {
			if dev.DeviceSki == evsePairCEM.SKI {
				id = idx
			}
		}
		if id >= 0 {
			fmt.Println("CEM successfully connected to the EVSE")
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Infinitely query MPC measurements
	for {
		response, err := sendRequest[[]MPCResponseBody]("GET", "/cem/MpcMeasurementData", nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf(`
=================
MPC values @ %v
%v
=================
`, time.Now(), response)

		time.Sleep(1 * time.Second)
	}

}

func sendRequest[T any](method string, url string, payload any) (T, error) {
	client := &http.Client{}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, BaseIPAddress+url, bytes.NewBuffer(jsonData))
	if err != nil {
		var x T
		return x, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		var x T
		return x, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		var x T
		return x, err
	}
	var result2 T
	err = json.Unmarshal(body, &result2)
	if err != nil {
		var x T
		return x, err
	}
	return result2, nil
}
