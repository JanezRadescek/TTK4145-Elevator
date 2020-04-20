SET varport=15657
echo varport
START .\elevator\SimElevatorServer.exe --port %varport%
START go run main.go %varport%

timeout 1

SET varport=15658
echo varport
START .\elevator\SimElevatorServer.exe --port %varport%
START go run main.go %varport%

:: pause