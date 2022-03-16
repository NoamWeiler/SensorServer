package main

import (
	"SensorServer/internal/sensor"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	numOfSensors = flag.Int("n", 1, "number of sensors in simulator")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}

	for i := 0; i < *numOfSensors; i++ {
		wg.Add(1)
		go func(str string) {
			s := sensor.Init(str)
			defer wg.Done()
			s.Run(ctx)
		}(fmt.Sprintf("sensor_%d", i))
	}

	//signal handler
	go func() {
		wg.Add(1)
		defer wg.Done()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		log.Println("\nsignal handler is up")
		sig := <-sigChan
		log.Println("\ninterrupted by:", sig)
		cancel()
	}()

	log.Println("All sensors are up")
	wg.Wait()
	log.Println("finished simulation")

}
