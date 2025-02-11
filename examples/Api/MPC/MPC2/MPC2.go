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
	Energy    *float64                 `json:"energy,omitempty"`
	Power     []MeasurementDataElement `json:"power,omitempty"`
	Current   []MeasurementDataElement `json:"current,omitempty"`
	Voltage   []MeasurementDataElement `json:"voltage,omitempty"`
	Frequency *float64                 `json:"frequency,omitempty"`
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
	var energy string
	var frequency string
	if m.Energy == nil {
		energy = "[]"
	} else {
		energy = fmt.Sprintf("%.3f", *m.Energy)
	}
	if m.Frequency == nil {
		frequency = "[]"
	} else {
		frequency = fmt.Sprintf("%.3f", *m.Frequency)
	}
	return fmt.Sprintf(`=================
Power: %v
Energy: %v
Current: %v
Voltage: %v
Frequency: %v
=================`,
		m.Power, energy, m.Current, m.Voltage, frequency)
}

func (m MeasurementDataElement) String() string {
	if m.Value == nil {
		return "[]"
	}
	return fmt.Sprintf("[Phase: %v, Value: %.3f]", m.Phase, *m.Value)
}

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
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
