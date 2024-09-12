# OPEV5

## Description

In this user story, We are adding Energy guard and EVSE, then connection them with each other.  
After connection is done, we will add 2 EVs.

first EV is with current limits:
- min: 6
- max: 10 

and the second EV is with current limits:
- min: 10
- max: 20

Also we will add a non-controllable device that absorb 24A from the available fuse limit to make the available 16A

The Energy Guard shall apply Round Robin algorithm to distribute the current over the 2 EVs

## How to Run

import the configuration file in the simulation and start start simulation

![alt text](<Screenshot 2024-09-12 122235.png>)
![alt text](<Screenshot 2024-09-12 122300.png>)