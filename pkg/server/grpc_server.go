package main

import (
	pb "SensorServer/internal/mutual_db"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	adminIsConnected = false
	sensorCount      = make(chan int64, 1)
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
	debug("returnError", "err=nil")
	return nil
}

// ClientInfo implementation

func (s *server) ConnectClient(ctx context.Context, in *pb.ConnReq) (*pb.ConnRes, error) {
	f := "ConnectClient"
	debug(f, fmt.Sprintf("%v", in))

	if adminIsConnected { //can't connect twice
		debug(f, "adminIsConnected is true")
		return &pb.ConnRes{Res: ""}, returnError("yochbad is already connected")
	}
	if in.UserName != adminName || in.Password != adminPass {
		debug(f, fmt.Sprintf("Wrong credentials:\tin.UserName:%v, in.Password:%v", in.UserName, in.Password))
		return &pb.ConnRes{Res: ""}, returnError("Wrong credentials")
	}
	debug(f, "Connect Success!")
	adminIsConnected = true
	return &pb.ConnRes{Res: "Connected successfully"}, nil
}

func (s *server) DisconnectClient(ctx context.Context, in *pb.DisConnReq) (*pb.ConnRes, error) {
	f := "DisconnectClient"
	debug(f, fmt.Sprintf("%v", "enter"))

	if !adminIsConnected { //can't disconnect is not connected first
		debug(f, "adminIsConnected is false, DisconnectClient error")
		return &pb.ConnRes{Res: ""}, returnError("yochbad is not connected")
	}

	debug(f, "Disconnected successfully")
	adminIsConnected = false
	return &pb.ConnRes{Res: "Disconnected successfully"}, nil
}

func (s *server) GetInfo(ctx context.Context, in *pb.InfoReq) (*pb.InfoRes, error) {
	f := "GetInfo"
	debug(f, fmt.Sprintf("enter, args:%v", in))
	return &pb.InfoRes{Responce: "GOT YOUR REQUEST"}, nil
}

// SensorStream implementation

func (s *server) ConnectSensor(ctx context.Context, in *pb.ConnSensorReq) (*pb.ConnSensorRes, error) {
	f := "ConnectSensor"
	var num int64
	debug(f, fmt.Sprintf("enter, args:%v", in))
	//get the next serial number and increase by 1 the value to the next
	num = <-sensorCount
	sensorCount <- num + 1
	return &pb.ConnSensorRes{Serial: fmt.Sprintf("sensor_%d", num)}, nil
}

func (s *server) SensorMeasure(ctx context.Context, in *pb.Measure) (*pb.MeasureRes, error) {
	f := "SensorMeasure"
	debug(f, fmt.Sprintf("got measure=%d from %s", in.GetM(), in.GetSerial()))
	return &pb.MeasureRes{}, nil
}

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterSensorStreamServer(s, &server{})
	pb.RegisterClientInfoServer(s, &server{})
	sensorCount <- 1
	return s
}
