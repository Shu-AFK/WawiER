@echo off
SET APPFOLDER=%~dp0

powershell.exe -Command "Add-MpPreference -ExclusionPath '%APPFOLDER%'"

powershell.exe -Command "Get-ChildItem -Path '%APPFOLDER%' -Filter *.exe | ForEach-Object { Add-MpPreference -ExclusionProcess $_.FullName }"

echo Exclusions applied for folder and all executables.
pause
