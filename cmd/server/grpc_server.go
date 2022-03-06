package main

import (
	grpc_db "SensorServer/pkg/grpc_db"
	In_memo_db "SensorServer/pkg/in_memory_db"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

//interface to represent DB functionalities
type sensorDB interface {
	AddMeasure(string, int)
	GetInfo(string, int) string
	DayCleanup()
}

var (
	adminIsConnected = false
	sensorCount      = make(chan int64, 1)
	gs               *grpc.Server
	lis              net.Listener
	db               sensorDB
)

const (
	adminName = "yochbad"
	adminPass = "123"
)

func returnError(s string) error {
	err := status.Error(codes.Unimplemented, s)
	if err != nil {
		return err
	}
	debug("returnError", s)
	return nil
}

func debug(f string, s string) {
	if *verbose {
		log.Printf("[%s]: %v", f, s)
	}
}

// ClientInfo implementation

func (s *server) ConnectClient(ctx context.Context, in *grpc_db.ConnReq) (*grpc_db.ConnRes, error) {
	f := "ConnectClient"
	debug(f, fmt.Sprintf("%v", in))

	if adminIsConnected { //can't connect twice
		debug(f, "adminIsConnected is true")
		return &grpc_db.ConnRes{Res: ""}, returnError("yochbad is already connected")
	}
	if in.UserName != adminName || in.Password != adminPass {
		debug(f, fmt.Sprintf("Wrong credentials:\tin.UserName:%v, in.Password:%v", in.UserName, in.Password))
		return &grpc_db.ConnRes{Res: ""}, returnError("Wrong credentials")
	}
	debug(f, "Connect Success!")
	adminIsConnected = true
	return &grpc_db.ConnRes{Res: "Connected successfully"}, nil
}

func (s *server) DisconnectClient(ctx context.Context, in *grpc_db.DisConnReq) (*grpc_db.ConnRes, error) {
	f := "DisconnectClient"
	debug(f, fmt.Sprintf("%v", "enter"))

	if !adminIsConnected { //can't disconnect is not connected first
		debug(f, "adminIsConnected is false, DisconnectClient error")
		return &grpc_db.ConnRes{Res: ""}, returnError("yochbad is not connected")
	}

	debug(f, "Disconnected successfully")
	adminIsConnected = false
	return &grpc_db.ConnRes{Res: "Disconnected successfully"}, nil
}

func (s *server) GetInfo(ctx context.Context, in *grpc_db.InfoReq) (*grpc_db.InfoRes, error) {
	f := "GetInfo"
	debug(f, fmt.Sprintf("args:%v", in))
	//unpack request for sensorDB interface
	res := db.GetInfo(in.GetSensorName(), int(in.GetDayBefore()))

	return &grpc_db.InfoRes{Responce: res}, nil
}

// SensorStream implementation

func (s *server) ConnectSensor(ctx context.Context, in *grpc_db.ConnSensorReq) (*grpc_db.ConnSensorRes, error) {
	f := "ConnectSensor"
	var num int64
	debug(f, fmt.Sprintf("args:%v", in))
	//get the next serial number and increase by 1 the value to the next
	num = <-sensorCount
	sensorCount <- num + 1
	return &grpc_db.ConnSensorRes{Serial: fmt.Sprintf("sensor_%d", num)}, nil
}

func (s *server) SensorMeasure(ctx context.Context, in *grpc_db.Measure) (*grpc_db.MeasureRes, error) {
	f := "SensorMeasure"
	debug(f, fmt.Sprintf("got measure=%d from %s", in.GetM(), in.GetSerial()))
	db.DayCleanup()
	//unpack request for sensorDB interface
	go db.AddMeasure(in.GetSerial(), int(in.GetM())) //run in parallel
	return &grpc_db.MeasureRes{}, nil
}

//implementation of protocolServer interface
func (s *server) runServer() {
	var err error
	lis, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs = grpc.NewServer()
	grpc_db.RegisterSensorStreamServer(gs, &server{})
	grpc_db.RegisterClientInfoServer(gs, &server{})
	sensorCount <- 1

	//DB - used sensorDB interface
	db = In_memo_db.SensorMap()

	log.Printf("server listening at %v", lis.Addr())
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func (s *server) cleanup() {
	adminIsConnected = false
	gs.GracefulStop()
	close(sensorCount)
}
