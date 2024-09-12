# OPEV2

## Description

In this user story, We are adding Energy guard and EVSE, then connection them with each other.  
After connection is done, we are adding an EV and connecting it with the EVSE with current limits:
- min: 6
- max: 10

We will add a non-controllable device to the system that will absorb all the fuse limits (40A), this will cause the Energy Guard to set the current of the EV to 0

## How to Run

import the configuration file in the simulation and start start simulation

![alt text](<Screenshot 2024-09-12 112604.png>)