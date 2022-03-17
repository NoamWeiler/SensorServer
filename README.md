# OnBoardProject

## Run server
Few  running options, all of them should be executed from root run: </br>
* Run single server:</br>make run_server</br></br>
* Run single server (debug mode):</br>make run_server_debug</br></br>
* Run multiple servers:</br>make run_multiple_servers</br></br>
* Run multiple servers (debug mode):</br>make run_multiple_servers_debug</br></br>

### Important note </br>
If Running multiple servers (both regular and  debug modes) need to shut them down:</br>
make shutdown_multiple_servers
</br></br></br>
## Run client
From root run: </br>
make run_client </br>
username/password: yochbad/123

## Run stream of sensors
* Run stream for a single server: </br>
make run_stream </br></br>
* Run sensors simulator (with client load-balancer) : </br>
make run_sensors_simulator </br></br>
* Testing responses per second:</br>
ghz --insecure --proto ./pkg/grpc_db/grpc_db.proto --call SensorServer.SensorStream.SensorMeasure -d '[{"m":13,"serial":"ser345345"}]' 0.0.0.0:50051 -x 1s
