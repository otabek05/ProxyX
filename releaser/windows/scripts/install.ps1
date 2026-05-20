# Installs ProxyX as a Windows service.
# Run from an elevated PowerShell prompt in the unpacked release directory.

$ErrorActionPreference = "Stop"

$ServiceName = "proxyx"
$InstallDir  = "C:\Program Files\ProxyX"
$ConfDir     = "C:\ProgramData\ProxyX\conf.d"
$WebDir      = "C:\ProgramData\ProxyX\web"

if (-not ([Security.Principal.WindowsPrincipal] `
    [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(`
    [Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw "install.ps1 must be run from an elevated (Administrator) shell"
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path $ConfDir    | Out-Null
New-Item -ItemType Directory -Force -Path $WebDir     | Out-Null

$source = Split-Path -Parent $PSScriptRoot
Copy-Item -Force "$source\proxyx.exe" "$InstallDir\proxyx.exe"
if (Test-Path "$source\web") {
    Copy-Item -Recurse -Force "$source\web\*" $WebDir
}

if (Get-Service -Name $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Service $ServiceName already exists; stopping before reinstall..."
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $ServiceName | Out-Null
    Start-Sleep -Seconds 2
}

sc.exe create $ServiceName `
    binPath= "`"$InstallDir\proxyx.exe`" start" `
    DisplayName= "ProxyX Reverse Proxy" `
    start= auto | Out-Null

sc.exe description $ServiceName "ProxyX advanced reverse proxy service" | Out-Null
sc.exe failure $ServiceName reset= 60 actions= restart/5000/restart/5000/restart/5000 | Out-Null

Start-Service -Name $ServiceName
Write-Host "ProxyX installed and started."
