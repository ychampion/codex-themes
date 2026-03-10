@echo off
setlocal
powershell -ExecutionPolicy Bypass -File "%~dp0scripts\install.ps1" %*
