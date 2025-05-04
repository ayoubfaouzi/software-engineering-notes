# Save the current directory
$originalDir = Get-Location

Import-Module "C:\Program Files (x86)\Microsoft Visual Studio\2019\Professional\Common7\Tools\Microsoft.VisualStudio.DevShell.dll"
Enter-VsDevShell -VsInstallPath "C:\Program Files (x86)\Microsoft Visual Studio\2019\Professional" -DevCmdArguments '-arch=x64'

# Restore the original directory at the end
Set-Location $originalDir

