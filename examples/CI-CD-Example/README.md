# Jenkins CI/CD pipeline To Test LPC device
   ![Jenkins_logo svg - Copy](https://github.com/user-attachments/assets/9124a82a-11b7-41fc-82e7-93cfb8a61394)

---

# 🎯 Overview

 This CI/CD example uses Jenkins to test the "Limitation of Power Consumption" on external EVSE.

### Jenkins
  - An open-source, Java-based, automation server that enables continuous integration/continuous delivery and deployment (CI/CD).
   
### Limitation of Power Consumption on External EVSE
- This example demonstrates how to automatically test the LPC failsafe scenario communication between a Controllable system (external EVSE) against an Energy Guard (EEBus-Hub CEM)..  

# ✨ Pre-Configurations

- Download EEBUS-Hub <https://www.coretech-innovations.com/products/eebushub/download>
- Ensure Go is installed on Windows machine <https://go.dev/doc/install>
- Ensure Python is Installed on Windows machine <https://www.python.org/downloads/> 
- **Jenkins** Tools and Plugins:
    - Make Sure Git is Installed in Tools [ Dashboard > Manage Jenkins > Tools ]
    - Python is added in Tools [ Dashboard > Manage Jenkins > Tools ]
    - [ Workspace Cleanup Plugin ] Plugin is Installed 
- Note: We are running on Windows Agent

# ⚙️ Configuration

- Navigate to Jenkins 

- Select New Item

![New Item](https://github.com/user-attachments/assets/bdc6ce62-1092-457b-9432-024f717a4cd2)

- Enter an item name
- Select Pipeline
- Press OK

![EEBUS Name](https://github.com/user-attachments/assets/543594c1-58f9-4fae-a0ff-05b7ca33f22b)

- From the Configure menu select Pipeline
- Pipeline > Definition > Pipeline script from SCM
- SCM > Select Git
- Type repository URL > ```https://github.com/Coretech-Innovations/EEBus-Hub.git```
- Branches to build > Branch Specifier (blank for 'any') > Type ```*/main```
- Script Path > Type ```examples/CI-CD-example/Jenkinsfile```
- Press Save

![LatestPipelineConfig](https://github.com/user-attachments/assets/5796b2f0-dcd8-4205-b2dd-aec040fc8b25)

# ➡️ Running Testcase

- Select Build NOW

![Build NOW](https://github.com/user-attachments/assets/31e0f17b-6e3e-440d-b474-1a2eddc06ae7)

- The Pipeline Will Fail for the First Time

![Failed1](https://github.com/user-attachments/assets/f28df571-f16d-4292-bdfd-6192aec87777)

- Refresh The Page Then Select [Build With Parmeter]

- Change The PORT if Jenkins is Running on <http://localhost:8080/>
- Select the Archived File Downloaded from <https://www.coretech-innovations.com/products/eebushub/download>
- If the agent running on a Windows machine has a different label, Please write the correct label 
- Click Build

![RunningLast](https://github.com/user-attachments/assets/e42d0a18-ab9c-4503-8779-8d1d2973ed07)

- Refresh the page after the build is complete
- The result of the Test case will appear as an Artifact After The Build Finishes
  - Note: The Artifact that appears is from the latest run  
  
![artifacts](https://github.com/user-attachments/assets/d97fc9ca-5bcf-4b6d-98dd-9118db9931e6)
