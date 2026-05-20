# Uninstalls the ProxyX Windows service and removes its install directory.
# Run from an elevated PowerShell prompt.

$ErrorActionPreference = "Stop"

$ServiceName = "proxyx"
$InstallDir  = "C:\Program Files\ProxyX"

if (-not ([Security.Principal.WindowsPrincipal] `
    [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(`
    [Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw "uninstall.ps1 must be run from an elevated (Administrator) shell"
}

if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $ServiceName | Out-Null
    Start-Sleep -Seconds 2
    Write-Host "ProxyX service removed."
} else {
    Write-Host "ProxyX service was not installed."
}

if (Test-Path $InstallDir) {
    Remove-Item -Recurse -Force $InstallDir
    Write-Host "Removed $InstallDir."
}

Write-Host "Config under C:\ProgramData\ProxyX left intact (delete manually if no longer needed)."
