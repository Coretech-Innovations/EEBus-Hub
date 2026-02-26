# EEBus-Hub

[EEBUS](https://eebus.org) framework to test integrating your device over an EEBUS network. It uses the versatile open-source stack [eebus-go](https://github.com/enbility/eebus-go) for EEBUS interactions and implementing different actors and use cases.

The EEBus-Hub provides APIs to control different actors participating that may interact in an EEBUS environment e.g. EV, EVSE, HEMS, Energy Guard, SMGW,...
The simulation allows plugging in real devices besides the numerous simulated devices to ease the testing of an EEBUS device.


📩 Licensing: To obtain a license for EEBUS Hub, contact
eebus.hub@coretech-innovations.com

Also, for further info or support, you can contact us at: <eebus.hub@coretech-innovations.com>

## Supported Use Cases
| UseCase                                            | Scenarios | Server | Client | Usecase Category |
| :------------------------------------------------- |:---------:|:------:|:------:|:----------------:|
| Monitoring of Grid Connection Point (MGCP)       | Scenario 1  | - | - | Grid |
|                                                  | Scenario 2  | ✅ | - |  |
|                                                  | Scenario 3  | - | - |  |
|                                                  | Scenario 4  | ✅ | - |  |
|                                                  | Scenario 5  | ✅ | - |  |
|                                                  | Scenario 6  | ✅ | - |  |
|                                                  | Scenario 7  | ✅ | - |  |
|                                                  |  |  |   |  |
| Limitation of Power Consumption (LPC)            | Scenario 1  | ✅ | ✅ | Grid |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| Limitation of Power Production (LPP)             | Scenario 1  | ✅ | ✅ | Grid |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| Monitoring of Power Consumption (MPC)            | Scenario 1  | ✅ | ✅ | Grid |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
|                                                  | Scenario 5  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| EV Commissioning and Configuration (EVCC)        | Scenario 1  | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | - | - |  |
|                                                  | Scenario 3  | - | - |  |
|                                                  | Scenario 4  | - | - |  |
|                                                  | Scenario 5  | - | - |  |
|                                                  | Scenario 6  | ✅ | ✅ |  |
|                                                  | Scenario 7  | ✅ | ✅ |  |
|                                                  | Scenario 8  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| EVSE Commissioning and Configuration (EVSECC)    | Scenario 1  | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  |  |  |  |  |
| EV Charging Electricity Measurement (EVEM)       | Scenario 1  | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | - | - |  |
|                                                  | Scenario 3  | - | - |  |
|                                                  |  |  |   |  |
| Overload Protection by EV current curtailment (OPEV)| Scenario 1 | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| EV State of Charge (EVSOC)                       | Scenario 1  | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
| EV Charging Electricity Measurement (EVCEM)      | Scenario 1  | ✅ | ✅ | E-mobility |
|                                                  | Scenario 2  | - | - |  |
|                                                  |  |  |   |  |
| Monitoring of Inverter (MOI)                     | Scenario 1  | ✅ | ✅ | Inverter |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
|                                                  | Scenario 5  | ✅ | ✅ |  |
|                                                  | Scenario 6  | ✅ | ✅ |  |
|                                                  | Scenario 7  | ✅ | ✅ |  |
|                                                  |  |  |   |  |
| Monitoring of Battery (MOB)                      | Scenario 1  | ✅ | ✅ | Inverter |
|                                                  | Scenario 2  | ✅ | ✅ |  |
|                                                  | Scenario 3  | ✅ | ✅ |  |
|                                                  | Scenario 4  | ✅ | ✅ |  |
|                                                  | Scenario 5  | ✅ | ✅ |  |
|                                                  | Scenario 6  | ✅ | ✅ |  |
|                                                  | Scenario 7  | ✅ | ✅ |  |
|                                                  | Scenario 8  | ✅ | ✅ |  |
|                                                  | Scenario 9  | ✅ | ✅ |  |
| Optimization of Self-Consumption by Heat Pump Compressor Flexibility (OHPCF) | Scenario 1  | ✅ | ✅ | HVAC |
|                                                  | Scenario 2  | ✅ | ✅ |  |


"✅" - Supported

"-"  - Not Supported yet

## Clone the project

```bash
git clone https://github.com/Coretech-Innovations/EEBus-Hub.git

```

## How to Use

This is a simple API calls for adding EVSE and EV and connecting them with the CEM in the system

```bash
# adding new EVSE
curl -X POST http://localhost:8080/api/v1/evse/add  -H 'Content-Type: application/json'  -d '{"deviceName":"Coretech EVSE WLBX", "deviceCode":"0002","deviceModel":"Charging Station","brandName":"Coretech Innovations","vendor":{"name":"Coretech Innovations","code":"60745"},"serialNumber":"SN7640"}'

# Trusting the created EVSE from the CEM side
curl -X POST http://localhost:8080/api/v1/cem/trust -H 'Content-Type: application/json' -d '{"remoteSki": <EVSE Ski>}'

# Get the SKI of the CEM
curl -X GET http://localhost:8080/api/v1/cem 

# Pairing the EVSE with the CEM
curl -X POST http://localhost:8080/api/v1/evse/<EVSE ID>/cem -H 'Content-Type: application/json' -d '{"remoteski": <CEM Ski>}'

# adding new EV
curl -X POST http://localhost:8080/api/v1/ev/add -H 'Content-Type: application/json' -d '{"device": {"name": "Taycan", "code": "0003", "serialNumber": "SN1235"},"currentLimits": {"min": 5, "max": 10}, "asymmetricCharging": false}'

# adding the created EV to the EVSE we created before
curl -X POST http://localhost:8080/api/v1/ev/<EV ID>/evse/<EVSE ID>

# starting the simulation
curl -X POST http://localhost:8080/api/v1/sim -H 'Content-Type: application/json' -d '{"action": "start","tickRate": 1000,"simTimePerTick": 10}'
```

EVSE Ski: The Ski returned from creating the EVSE  
EVSE ID: The id returned from creating the EVSE  
CEM Ski: The Ski returned from calling the GET request on the CEM  
EV ID: The id returned from creating the EV

## How to load the docker image

### Step 1: Import the image from downloaded archive

```bash
docker load < docker-eebus-hub-vx.y.z.tar.xz
```

### Step 2: Create a container

* using bridged networking mode:

```bash
docker run -p $PORT:8080 eebus-hub:vx.y.z
```

where $PORT is the desired published port.

* using host networking mode:

```bash
docker run --network=host eebus-hub:vx.y.z
```
