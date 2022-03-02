package main

import (
	grpc_db "github.com/noamweiler/SnsorServer/pkg/grpc_db"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
)

//interface to represent DB functionalities
type sensorDB interface {
	//GetInfo(r *grpc_db.InfoReq) string
	//AddMeasure(measure *grpc_db.Measure)
	AddMeasure(string,int)
	GetInfo(string,int) string
	DayCleanup()
}

var (
	adminIsConnected = false
	sensorCount      = make(chan int64, 1)
	gs               *grpc.Server
	lis              net.Listener
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

	res :=
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
	dayCleanup(db)
	go addMeasure(db, in) //run in parallel
	return &grpc_db.MeasureRes{}, nil
}

//implementation of protocolServer interface
func (s *server) createServer() error {
	var err error
	lis, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// the DB is support only 1 week, so need to do know if day have been changed
	GlobalDay = time.Now().Weekday()

	gs = grpc.NewServer()
	grpc_db.RegisterSensorStreamServer(gs, &server{})
	grpc_db.RegisterClientInfoServer(gs, &server{})
	sensorCount <- 1

	//DB
	db = in_memory_db.SensorMap()

	log.Printf("server listening at %v", lis.Addr())
	return gs.Serve(lis)
}

func (s *server) cleanup() {
	adminIsConnected = false
	gs.GracefulStop()
	if err := lis.Close(); err != nil {
	}
	close(sensorCount)
}
