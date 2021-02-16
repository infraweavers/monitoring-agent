@echo off
echo "Installing `Monitoring Agent` Service"
echo "Running sc create"
sc create "Monitoring Agent" start= auto error= critical binpath= %~dp0\monitoring-agent.exe