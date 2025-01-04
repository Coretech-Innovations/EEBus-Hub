# LPC3

## Introduction

This example will help you understand and implement the "Limitation of Power Consumption" use case, detailing the setup and testing process for the EEBus Hub Customer Energy Manager (CEM), and an external Electric Vehicle Supply Equipment (EVSE).

### Goal of this example

The focus of this example is to demonstrate how to connect the CEM to an external **EVSE** to observe and implement the "Limitation of Power Consumption" use case (scenario 2). This use case ensures that energy consumption is managed effectively to prevent overloading or excessive power usage.

### Use Case Overview: Limitation of power consumption

This use case describes how an actor Energy Guard (CEM in our case) manages the maximum power consumption of an actor Controllable Systems (EVSE in our case) to achieve the following goals:  
- Grid stabilization  
- Prevention of overload in the low-voltage distribution network  
- Prevention of exceeding the maximum power consumption value of the grid connection point  

The following mechanisms are utilized within this Use Case:  
- Active Power Consumption Limit: The Active Power Consumption Limit allows to set a limit for the maximum active power consumption of a Controllable System. The Active Power Consumption Limit is typically used to improve grid stability by reducing the consumption of the Controllable System. The Active Power Consumption Limit may have a duration of validity.  
- Failsafe Consumption Active Power Limit: If the Controllable System does not receive any Heartbeats from the Energy Guard for more than 120 seconds (e.g. due to interrupted connectivity), the Failsafe Consumption Active Power Limit is used as fallback. It is intended to prevent overloads in case of connectivity problems or when a building or device is re-powered after loss of power.

## Example implementation details

In this example, We are connecting the CEM to an external EVSE which listed in the **devices** directory.
After connection, the CEM will send LPC command to set the Failsafe consumption active power limit and observe how the EVSE will respond.


## How to run 
``` bash
go run .\examples\Api\LPC\LPC3\LPC3.go -f
```