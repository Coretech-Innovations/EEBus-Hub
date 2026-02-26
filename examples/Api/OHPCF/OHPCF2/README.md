# OHPCF2

## Description

In this use case, we add a CEM device and a heat pump, then connect them with each other.
After the connection is done, we create an OHPCF announcement with power limits and timing constraints:
- power(goodApproximation): 15000 W
- powerMaximum: 16000 W
- activeDurationMinimum: 5 seconds
- pauseDurationMinimum: 5 seconds
- startTime: 5 seconds

Then the CEM sends a start time and exercises pause and resume transitions after the minimum
active/pause durations are reached. The simulation waits for the heat pump to complete its optional
power consumption process and transition to the completed state, then to the inactive state.

## How to Run
```
go run ./examples/OHPCF/OHPCF2/OHPCF2.go
```
