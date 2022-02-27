package main

import (
	pb "SensorServer/internal/mutual_db"
	"flag"
	"fmt"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedClientInfoServer   //handle client request
	pb.UnimplementedSensorStreamServer //hande sensors measures
}

var (
	verbose = flag.Bool("v", false, "Verbose mode")
	port    = flag.Int("port", 50051, "The server port")
)

func debug(f string, s string) {
	if *verbose {
		log.Printf("[%s]: %v", f, s)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := NewGRPCServer()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	defer func() { //cleanup
		adminIsConnected = false
		s.Stop()
		if err := lis.Close(); err != nil {
			log.Println(err)
		}
	}()
}

//TODO
/*
	1)	add signal interrupt for proper closing ctl+c
	2)	start working on the measures handling (incllude DB)
*/
