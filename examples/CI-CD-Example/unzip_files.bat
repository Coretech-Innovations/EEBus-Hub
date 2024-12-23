@echo off
setlocal enabledelayedexpansion

rem Set the destination folder
set DEST_FOLDER=.

rem Loop through all zip files matching the pattern
for %%f in (eebus-hub-windows-amd64*.zip) do (
    echo Unzipping %%f
    powershell -command "Expand-Archive -Path '%%f' -DestinationPath '%DEST_FOLDER%' -Force"
)

endlocal
