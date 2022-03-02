package main

import (
	pb "SensorServer/internal/mutual_db"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//struct to hold the GRPC's handlers
type server struct {
	pb.UnimplementedClientInfoServer   //handle client request
	pb.UnimplementedSensorStreamServer //hande sensors measures
}

//interface to represent DB functionalities
type sensorDB interface {
	addMeasure(measure *pb.Measure)
	getInfo(r *pb.InfoReq) string
	dayCleanup()
}

// interface to represent a server
type protocolServer interface {
	createServer() error
	cleanup()
}

var (
	verbose  = flag.Bool("v", true, "Verbose mode")
	grpcPort = flag.Int("port", 50051, "The server port")
)

func debug(f string, s string) {
	if *verbose {
		log.Printf("[%s]: %v", f, s)
	}
}

func cleanupProtocolServer(ps protocolServer) {
	ps.cleanup()
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer close(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	grpcs := new(server)

	g.Go(grpcs.createServer)
	//defer grpcs.cleanup() -> after the select (using the interface)

	select {
	case in := <-interrupt:
		fmt.Println(in)
		break
	case <-ctx.Done():
		break
	}

	//put the server up for 10 minutes, then close it gracefully
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer shutdownCancel()

	cleanupProtocolServer(grpcs)

	err := g.Wait()
	if err != nil {
		log.Fatalf("server returning an error.\nerror:%v", err)
	}

	fmt.Println("exit..")
	fmt.Println("exit..")
}

//TODO
/*
	1)
	2)
*/
