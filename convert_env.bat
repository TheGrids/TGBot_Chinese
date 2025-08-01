@echo off
powershell -Command "[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Content .env -Raw)))"
echo.
pause