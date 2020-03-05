SET varport=13945
echo varport
START SimElevatorServer --port %varport%

set varport=varport+1
echo varport
START SimElevatorServer --port %varport%

set varport=varport+1
echo varport
START SimElevatorServer --port %varport%