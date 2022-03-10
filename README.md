# OnBoardProject

## run server
From root run: </br>
make run_server

## run client
From root run: </br>
make run_client </br>
username/password: yochbad/123

## run stream of sensors
From root run: </br>
make run_stream </br>
Testing responces per second:
ghz --insecure --proto ./pkg/grpc_db/grpc_db.proto --call SensorServer.SensorStream.SensorMeasure -d '[{"m":13,"serial":"ser345345"}]' 0.0.0.0:50051 -x 1s
