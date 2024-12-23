# Jenkins CI/CD pipeline to test LPC device
   ![Jenkins_logo svg - Copy](https://github.com/user-attachments/assets/da3154b3-0f7c-42f2-b423-d7f95e3110c6)

---

# 🎯 Overview

 This CI/CD example uses Jenkins to test the "Limitation of Power Consumption" on external EVSE.

### Jenkins
  - An open-source, Java-based, automation server that enables continuous integration/continuous delivery and deployment (CI/CD).
   
### Limitation of Power Consumption
- This example demonstrates how to automatically test the LPC failsafe scenario communication between a Controllable system (external EVSE) against an Energy Guard (EEBus-Hub CEM).

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

  ![New Item](https://github.com/user-attachments/assets/ece08a1f-7ff5-4c39-b9a5-26054681e548)


- Enter an item name
- Select Pipeline
- Press OK

   ![EEBUS Name](https://github.com/user-attachments/assets/da4ae5e9-2e1d-4b47-b79b-39d3c47b9f74)

- From the Configure menu select Pipeline
- Pipeline > Definition > Pipeline script from SCM
- SCM > Select Git
- Type repository URL > ```https://github.com/Coretech-Innovations/EEBus-Hub.git```
- Branches to build > Branch Specifier (blank for 'any') > Type ```*/main```
- Script Path > Type ```examples/CI-CD-Example/Jenkinsfile```
- Press Save
  
   ![LatestPipelineConfig](https://github.com/user-attachments/assets/f0293011-b014-4a0f-81f4-5ed60bc098de)


# ➡️ Running Testcase

- Select Build NOW

   ![Build NOW](https://github.com/user-attachments/assets/b5897293-c236-401e-97a3-b2312472c095)

- The Pipeline Will Fail for the First Time

   ![Failed1](https://github.com/user-attachments/assets/583c17f5-bb9f-4bca-9767-41afa5d2eec9)

- Refresh The Page Then Select [Build With Parmeter]

- Change The PORT if Jenkins is Running on <http://localhost:8080/>
- Select the Archived File Downloaded from <https://www.coretech-innovations.com/products/eebushub/download>
- If the agent running on a Windows machine has a different label, Please write the correct label 
- Click Build
  
![RunningLast](https://github.com/user-attachments/assets/c95ffad3-7238-409d-b988-3dd615dc7a55)


- Refresh the page after the build is complete
- The result of the Test case will appear as an Artifact After The Build Finishes
  - Note: The Artifact that appears is from the latest run  
  
![artifacts](https://github.com/user-attachments/assets/66f31bfc-15a9-40c2-a7c1-c896bf9b02c6)
