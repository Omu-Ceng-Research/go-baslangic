@echo off
rem Cift tiklayarak Docker ortamini baslatir. Ikinci kez calistirilirsa ikinci terminal acar.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0baslat.ps1" %*
if errorlevel 1 pause
