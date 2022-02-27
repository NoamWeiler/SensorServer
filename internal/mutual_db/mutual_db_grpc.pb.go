// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: internal/mutual_db/mutual_db.proto

package SensorServer

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ClientInfoClient is the client API for ClientInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientInfoClient interface {
	ConnectClient(ctx context.Context, in *ConnReq, opts ...grpc.CallOption) (*ConnRes, error)
	GetInfo(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoRes, error)
	DisconnectClient(ctx context.Context, in *DisConnReq, opts ...grpc.CallOption) (*ConnRes, error)
}

type clientInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewClientInfoClient(cc grpc.ClientConnInterface) ClientInfoClient {
	return &clientInfoClient{cc}
}

func (c *clientInfoClient) ConnectClient(ctx context.Context, in *ConnReq, opts ...grpc.CallOption) (*ConnRes, error) {
	out := new(ConnRes)
	err := c.cc.Invoke(ctx, "/SensorServer.ClientInfo/ConnectClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientInfoClient) GetInfo(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoRes, error) {
	out := new(InfoRes)
	err := c.cc.Invoke(ctx, "/SensorServer.ClientInfo/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientInfoClient) DisconnectClient(ctx context.Context, in *DisConnReq, opts ...grpc.CallOption) (*ConnRes, error) {
	out := new(ConnRes)
	err := c.cc.Invoke(ctx, "/SensorServer.ClientInfo/DisconnectClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientInfoServer is the server API for ClientInfo service.
// All implementations must embed UnimplementedClientInfoServer
// for forward compatibility
type ClientInfoServer interface {
	ConnectClient(context.Context, *ConnReq) (*ConnRes, error)
	GetInfo(context.Context, *InfoReq) (*InfoRes, error)
	DisconnectClient(context.Context, *DisConnReq) (*ConnRes, error)
	mustEmbedUnimplementedClientInfoServer()
}

// UnimplementedClientInfoServer must be embedded to have forward compatible implementations.
type UnimplementedClientInfoServer struct {
}

func (UnimplementedClientInfoServer) ConnectClient(context.Context, *ConnReq) (*ConnRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectClient not implemented")
}
func (UnimplementedClientInfoServer) GetInfo(context.Context, *InfoReq) (*InfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedClientInfoServer) DisconnectClient(context.Context, *DisConnReq) (*ConnRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisconnectClient not implemented")
}
func (UnimplementedClientInfoServer) mustEmbedUnimplementedClientInfoServer() {}

// UnsafeClientInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientInfoServer will
// result in compilation errors.
type UnsafeClientInfoServer interface {
	mustEmbedUnimplementedClientInfoServer()
}

func RegisterClientInfoServer(s grpc.ServiceRegistrar, srv ClientInfoServer) {
	s.RegisterService(&ClientInfo_ServiceDesc, srv)
}

func _ClientInfo_ConnectClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientInfoServer).ConnectClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorServer.ClientInfo/ConnectClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientInfoServer).ConnectClient(ctx, req.(*ConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientInfo_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientInfoServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorServer.ClientInfo/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientInfoServer).GetInfo(ctx, req.(*InfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientInfo_DisconnectClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientInfoServer).DisconnectClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorServer.ClientInfo/DisconnectClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientInfoServer).DisconnectClient(ctx, req.(*DisConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientInfo_ServiceDesc is the grpc.ServiceDesc for ClientInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SensorServer.ClientInfo",
	HandlerType: (*ClientInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectClient",
			Handler:    _ClientInfo_ConnectClient_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _ClientInfo_GetInfo_Handler,
		},
		{
			MethodName: "DisconnectClient",
			Handler:    _ClientInfo_DisconnectClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/mutual_db/mutual_db.proto",
}

// SensorStreamClient is the client API for SensorStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SensorStreamClient interface {
	ConnectSensor(ctx context.Context, in *ConnSensorReq, opts ...grpc.CallOption) (*ConnSensorRes, error)
	SensorMeasure(ctx context.Context, in *Measure, opts ...grpc.CallOption) (*MeasureRes, error)
}

type sensorStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewSensorStreamClient(cc grpc.ClientConnInterface) SensorStreamClient {
	return &sensorStreamClient{cc}
}

func (c *sensorStreamClient) ConnectSensor(ctx context.Context, in *ConnSensorReq, opts ...grpc.CallOption) (*ConnSensorRes, error) {
	out := new(ConnSensorRes)
	err := c.cc.Invoke(ctx, "/SensorServer.SensorStream/ConnectSensor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sensorStreamClient) SensorMeasure(ctx context.Context, in *Measure, opts ...grpc.CallOption) (*MeasureRes, error) {
	out := new(MeasureRes)
	err := c.cc.Invoke(ctx, "/SensorServer.SensorStream/SensorMeasure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SensorStreamServer is the server API for SensorStream service.
// All implementations must embed UnimplementedSensorStreamServer
// for forward compatibility
type SensorStreamServer interface {
	ConnectSensor(context.Context, *ConnSensorReq) (*ConnSensorRes, error)
	SensorMeasure(context.Context, *Measure) (*MeasureRes, error)
	mustEmbedUnimplementedSensorStreamServer()
}

// UnimplementedSensorStreamServer must be embedded to have forward compatible implementations.
type UnimplementedSensorStreamServer struct {
}

func (UnimplementedSensorStreamServer) ConnectSensor(context.Context, *ConnSensorReq) (*ConnSensorRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectSensor not implemented")
}
func (UnimplementedSensorStreamServer) SensorMeasure(context.Context, *Measure) (*MeasureRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SensorMeasure not implemented")
}
func (UnimplementedSensorStreamServer) mustEmbedUnimplementedSensorStreamServer() {}

// UnsafeSensorStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SensorStreamServer will
// result in compilation errors.
type UnsafeSensorStreamServer interface {
	mustEmbedUnimplementedSensorStreamServer()
}

func RegisterSensorStreamServer(s grpc.ServiceRegistrar, srv SensorStreamServer) {
	s.RegisterService(&SensorStream_ServiceDesc, srv)
}

func _SensorStream_ConnectSensor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnSensorReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorStreamServer).ConnectSensor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorServer.SensorStream/ConnectSensor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorStreamServer).ConnectSensor(ctx, req.(*ConnSensorReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _SensorStream_SensorMeasure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Measure)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorStreamServer).SensorMeasure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorServer.SensorStream/SensorMeasure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorStreamServer).SensorMeasure(ctx, req.(*Measure))
	}
	return interceptor(ctx, in, info, handler)
}

// SensorStream_ServiceDesc is the grpc.ServiceDesc for SensorStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SensorStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SensorServer.SensorStream",
	HandlerType: (*SensorStreamServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectSensor",
			Handler:    _SensorStream_ConnectSensor_Handler,
		},
		{
			MethodName: "SensorMeasure",
			Handler:    _SensorStream_SensorMeasure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/mutual_db/mutual_db.proto",
}
