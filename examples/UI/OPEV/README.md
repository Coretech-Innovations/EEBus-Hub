# OPEV User Stories

This folder contains different examples to see the Overload Protection by EV Charging Curtailment use case (OPEV) in action

## Actors 

- EV 
- Energy Guard

## Supported Scenarios

All three scenarios represented in the EEBus documentaion are supported:
- Energy Guard curtails charging current of EV
- EV checks Energy Guard availability
- Energy Guard sends error state

## Examples List

| Example | Description |
| :---    | :---        |
| OPEV1   | Connecting simulated CEM with simulated EVSE with EV attached and Applying Current Limits on this EV |
| OPEV2   | Connecting simulated CEM with simulated EVSE with EV attached and un-controllable device in the system that absorbs all the available current to curtail the EV's current to 0 |
| OPEV3   | Connecting simulated CEM with simulated EVSE with 2 EVs attached to it, and see how the RoundRobin algorithm will distribute the current between these 2 EVs |
| OPEV4   | Connecting simulated CEM with external EVSE and if there is EV connected to this EVSE the CEM shall apply current limits on this EV |
| OPEV5   |   Connecting simulated CEM with simulated EVSE with 2 EVs attached to it, also Un-controllable device is added to the system, the CEM shall distribute the current limits between the 2 EVs and the un-controllable device absorbs its current all the time |
| OPEV6   | Connecting simulated EVSE with external CEM, an EV is attached to the EVSE. the EV shall receive the limits from the external CEM |
| OPEV7   | Connecting 2 simulated EVSEs with external CEM, each EVSE is connected to EV, the CEM shall apply the limits on the EVs and CEM algorithm may be noticed |
