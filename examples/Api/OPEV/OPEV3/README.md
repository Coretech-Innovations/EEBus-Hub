# OPEV3

## Description

In this user story, We are adding Energy guard and EVSE, then connection them with each other.  
After connection is done, we will add 2 EVs.

first EV is with current limits:
- min: 6
- max: 18  

and the second EV is with current limits:
- min: 10
- max: 30

The Energy Guard shall apply Round Robin algorithm to distribute the current over the 2 EVs

## How to Run
```
go run ./examples/OPEV/OPEV3/OPEV3.go
```