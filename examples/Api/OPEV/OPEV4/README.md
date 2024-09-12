# OPEV4

## Description

In this user story, We will add an Energy Guard and connect it to external EV

The remoteSKI for the external EV should be passed to this example to connect the Energy Guard with the External EV

The Energy Guard shall supply the EV with its max current if the max current is below the allowed current (40A in this case)

## How to Run
```
go run ./examples/OPEV/OPEV4/OPEV4.go <remoteSKI>
```

