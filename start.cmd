SET varport=15657
echo varport
START .\elevator\SimElevatorServer --port %varport%
START go run main.go %varport%

set varport=varport+1
echo varport
START .\elevator\SimElevatorServer --port %varport%
START go run main.go %varport%