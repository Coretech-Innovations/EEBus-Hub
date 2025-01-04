# Show LPC use case on External EVSE

## Introduction

This example will help you understand and implement the "Limitation of Power Consumption" use case, detailing the setup and testing process for the EEBus Hub, Customer Energy Manager (CEM), and an external Electric Vehicle Supply Equipment (EVSE).

### Goal of this example

The focus of this example is to demonstrate how to connect the CEM to an external **EVSE** to observe and implement the "Limitation of Power Consumption" use case (scenario 1 and scenario 2). This use case ensures that energy consumption is managed effectively to prevent overloading or excessive power usage.

### Use Case Overview: Limitation of power consumption

This use case describes how an actor Energy Guard (CEM in our case) manages the maximum power consumption of an actor Controllable Systems (EVSE in our case) to achieve the following goals:  
- Grid stabilization  
- Prevention of overload in the low-voltage distribution network  
- Prevention of exceeding the maximum power consumption value of the grid connection point  

The following mechanisms are utilized within this Use Case:  
- Active Power Consumption Limit: The Active Power Consumption Limit allows to set a limit for the maximum active power consumption of a Controllable System. The Active Power Consumption Limit is typically used to improve grid stability by reducing the consumption of the Controllable System. The Active Power Consumption Limit may have a duration of validity.  
- Failsafe Consumption Active Power Limit: If the Controllable System does not receive any Heartbeats from the Energy Guard for more than 120 seconds (e.g. due to interrupted connectivity), the Failsafe Consumption Active Power Limit is used as fallback. It is intended to prevent overloads in case of connectivity problems or when a building or device is re-powered after loss of power.  

## Step-by-Step Guide to Execute the "Limitation of Power Consumption" Use Case

The following steps will show how to connect the external EVSE with the EEBus-Hub and apply Active power consumption Limit on the EVSE and see the response of the EVSE regarding the sent limits, also Sending Failsafe consumption active power limit will be shown in these steps.  

### Step 1: Run the EEBus-Hub  
- Open a terminal and navigate to the EEBus-Hub directory.
- Use the command:
```bash
go run ./cmd/eebus-hub/ -f
 ```
- This will open the UI page in your browser.  
![alt text](<Screenshot 2024-11-07 014536.png>)

### Step 2: Add external EVSE
- from the tab **Add Device** on the left side bar click on Add EVSE then choose the real EVSE option, this will open this window.  
![alt text](<Screenshot 2024-11-07 014912.png>)
- enter the device name and the device SKI, then add the device.

### Step 3: Start the Simulation
- after adding the EVSE, it will show up on the left side from the chart disconnected from the EEBus.
- start the simulation to make the CEM connect with it.

### Step 4: Trust the CEM on the external EVSE
- The SKI of the CEM could be found by clicking on its icon which will show its all information including the SKI.
- Pair the CEM SKI on the External EVSE to approve the connection.
- Once the connection is successfully made the device will show on the EEBus.  
![alt text](<Screenshot 2024-11-07 015616.png>)

### Step 5: Send Active Failsafe consumption active power limit
- After the connection is made, the test for the LPC use case could be started.
- Click the button (send LPC command) next to the CEM, and it will open this window.  
![alt text](<Screenshot 2024-11-07 020116.png>)
- From the Device drop down list select the EVSE.
- Uncheck Edit Active power limit, and check Edit Failsafe.
- Enter the Failsafe value and duration, then click Send.
- If the device accepted these values a confirmation of success will be shown after sending the command.  
![alt text](<Screenshot 2024-11-07 020920-1.png>)

### Step 6: Send active power consumption limit
- Click the button (send LPC command) next to the CEM, and it will open the window for sending LPC command.
- Check Edit active power limit, and uncheck Edit Failsafe.
- Enter the active power limit and its duration, then click Send.
- If the device accepted these values a confirmation message of success will be shown as the previous step.
- If the device rejected these values a failure message will be shown after the sending the command.  
![alt text](<Screenshot 2024-11-07 021358.png>)