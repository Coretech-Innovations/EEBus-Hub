
# LPC

## Description

User Story: We are integrating an Energy Guard with an EVSE and establishing their connection. Post-connection, we will monitor the active limit initial values in the EVSE.

Steps:

- Send Active Limit: Specify the value and duration.
- Update EVSE State: Ensure the EVSE reflects its new limit.
- Revert State: Once the duration elapses, the EVSE returns to its unlimited/controlled state.

## How to run 

go run ./examples/Api/LPC/LPC.go
