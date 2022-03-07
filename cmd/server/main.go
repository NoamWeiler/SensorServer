package main

import (
	grpc_db "SensorServer/pkg/grpc_db"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

//struct to hold the GRPC's handlers
type server struct {
	grpc_db.UnimplementedClientInfoServer   //handle client request
	grpc_db.UnimplementedSensorStreamServer //hande sensors measures
}

// interface to represent a server
type protocolServer interface {
	runServer(ctx context.Context)
	cleanup()
}

var (
	verbose    = flag.Bool("v", false, "Verbose mode")
	grpcPort   = flag.Int("port", 50051, "The server port")
	grpcServer protocolServer
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	//GRPC server
	g.Go(func() error {
		grpcServer = &server{}
		grpcServer.runServer(gctx)
		return nil
	})

	// signal handler
	g.Go(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			close(sigChan)
			fmt.Printf("Received signal: %s\n", sig)
			cancel()
			break
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return gctx.Err()
		}
		return nil
	})

	//wait for all errgroup goroutines
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Print("context was canceled\n")
		} else {
			fmt.Printf("received error: %v", err)
		}
	} else {
		fmt.Println("finished shutdown")
	}
	fmt.Println("exit..")
}

//TODO
/*
	1)
	2)
*/
