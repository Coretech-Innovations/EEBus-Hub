# OPEV

## Description

In this user story, We are adding Energy guard and EVSE, then connection them with each other.  
After connection is done, we are adding an EV and connecting it with the EVSE with current limits:
- min: 6
- max: 10

The EV should be supplied with the max current 10A as the available fuse limit in the Energy Guard is 40A and there is no other devices connected to the system

## How to Run
```
go run ./examples/Api/OPEV/OPEV.go
```