@echo off
setlocal

set /p VERSION="Enter version (e.g 1.0.0): "

echo Building Windows x86_64...
set GOOS=windows
set GOARCH=amd64

go build -trimpath -ldflags="-s -w" -o builds\catapult-%VERSION%-windows-x86_64.exe main.go

echo Building Linux x86_64...
set GOOS=linux
set GOARCH=amd64

go build -trimpath -ldflags="-s -w" -o builds\catapult-%VERSION%-linux-x86_64 main.go

echo.
echo --------------------------------------------------
echo Done! Files are in the \builds folder.
echo --------------------------------------------------
pause