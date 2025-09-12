@echo off
SETLOCAL

:: --- Set folder to exclude (folder where the exe is) ---
SET "APPFOLDER=%~dp0"
SET "APPFOLDER=%APPFOLDER:~0,-1%"  :: Remove trailing backslash

:: --- PowerShell check ---
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command ^
"$folder = '%APPFOLDER%';" ^
"$exists = $false;" ^
"try { $pref = Get-MpPreference; } catch { exit 1 }" ^
"if ($pref -and $pref.ExclusionPath) {" ^
"foreach ($ep in $pref.ExclusionPath) {" ^
"if (-not [string]::IsNullOrWhiteSpace($ep)) {" ^
"try { $normalized = [System.IO.Path]::GetFullPath($ep).TrimEnd('\').ToLower() } catch { continue }" ^
"if ($normalized -eq $folder.ToLower()) { $exists = $true; break }" ^
"}" ^
"}" ^
"}" ^
"if (-not $exists) {" ^
"Add-MpPreference -ExclusionPath $folder;" ^
"Get-ChildItem -Path $folder -Filter *.exe | ForEach-Object { Add-MpPreference -ExclusionProcess $_.FullName };" ^
"echo Exclusions applied for folder and all executables." ^
"} else { echo Folder already excluded. }"

ENDLOCAL