# LPC3

## Description

In this user story, We are connecting the CEM to an external EVSE which listed in the **devices** directory.
After connection, the CEM will send LPC command to set the Failsafe consumption active power limit and observe how the EVSE will respond.
Then, the CEM wil send LPC command to set the Active power consumption limit and observe how the EVSE will respond.


## How to run 
``` bash
go run .\examples\Api\LPC\LPC3\LPC3.go -f
```