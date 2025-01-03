package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// this Script aims to connect the EEBus-Hub CEM with External EVSE and send LPC command(Active power limit, failsafe consumption limit) to this device

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {
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
	time.Sleep(2 * time.Second)
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

	// step no 1: add real EVSE
	evseSKI := "41c98b1bbe5fc7657ce311981951f12d304ab419"
	addEvseResp, err := sendRequest[struct {
		Error  string `json:"error"`
		Status string `json:"status"`
		Id     int    `json:"id"`
	}]("POST", "/evse/addExt", map[string]any{
		"remoteSki":  evseSKI,
		"deviceName": "Demo EVSE",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if addEvseResp.Status != "OK" {
		log.Fatal(addEvseResp.Error)
	}
	// step no 2: start the simulation
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
	// step no 3: trust the EVSE from CEM side
	cemTrustResp, err := sendRequest[struct {
		Status string `json:"status"`
		Err    string `json:"err"`
	}]("POST", "/cem/trust", map[string]any{
		"remoteSKI": evseSKI,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if cemTrustResp.Status != "OK" {
		log.Fatal(cemTrustResp.Err)
	}
	OS := runtime.GOOS
	var cmd *exec.Cmd
	switch OS {
	case "windows":
		// init the submodule
		cmd = exec.Command("cmd", "/C", "git submodule init")
		cmd.Dir = "./devices/eebus-go/"
		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err.Error())
		}
		// update the submodule
		cmd = exec.Command("cmd", "/C", "git submodule update")
		cmd.Dir = "./devices/eebus-go/"
		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err.Error())
		}
		// running the External EVSE
		commandStr := fmt.Sprintf("go run examples/evse/main.go 4711 %s ./keys/evse.crt ./keys/evse.key", cemSKI)
		cmd = exec.Command("cmd", "/C", commandStr)
		cmd.Dir = "./devices/eebus-go/"
		// starting the simulated EVSE and trust the EEBus-Hub CEM
		cmd.Start()

		// shutdown the EVSE after the test case finishes
		defer func() {
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	case "linux":
		// init the submodule
		cmd := exec.Command("git submodule init")
		cmd.Dir = "./devices/eebus-go/"
		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err.Error())
		}
		// update the submodule
		cmd = exec.Command("git submodule update")
		cmd.Dir = "./devices/eebus-go/"
		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err.Error())
		}
		// running the External EVSE
		commandStr := fmt.Sprintf("go run examples/evse/main.go 4711 %s ./keys/evse.crt ./keys/evse.key", cemSKI)
		cmd = exec.Command(commandStr)
		cmd.Dir = "./devices/eebus-go/"
		// starting the simulated EVSE and trust the EEBus-Hub CEM
		cmd.Start()

		// shutdown the EVSE after the test case finishes
		defer func() {
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	fmt.Println("CEM Trying to connect to the EVSE")
	deviceAddress := ""
	// wait for connection
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
			if dev.DeviceSki == evseSKI {
				id = idx
				deviceAddress = dev.DeviceAddress
			}
		}
		if id >= 0 {
			fmt.Println("Successfully connected to the EVSE.")
			break
		}
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
	// send failsafe command and check the result
	route := fmt.Sprintf("/cem/FailsafeActivePowerLimit/%s", deviceAddress)
	FSValues, err := sendRequest[struct {
		Success bool `json:"success"`
		Data    struct {
			IsValueChangeable    bool    `json:"isValueChangeable"`
			Value                float64 `json:"value"`
			IsDurationChangeable bool    `json:"isDurationChangeable"`
			DurationSec          float64 `json:"durationSeconds"`
		} `json:"data"`
	}]("POST", route, map[string]any{
		"isValueChangeable":    true,
		"value":                6000,
		"isDurationChangeable": true,
		"durationSeconds":      10000,
	})
	if err != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Fatal(err.Error())
	}
	if FSValues.Success {
		fmt.Println("CEM successfully wrote new failsafe values on the EVSE")
		fmt.Printf("New data: failsafe value: %.2f W, failsafe duration: %.2f sec\n", FSValues.Data.Value, FSValues.Data.DurationSec)
	} else {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Fatal("CEM failed to update the failsafe values for the EVSE")
	}
	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(0)
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
