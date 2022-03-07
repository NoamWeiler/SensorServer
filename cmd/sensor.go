package main

import (
	"fmt"
	"github.com/bojand/ghz/runner"
	"os"
)

func main() {
	_, err := runner.Run(
		"SensorServer.SensorStream.SensorMeasure",
		"localhost:50051",
		runner.WithProtoFile("./pkg/grpc_db/grpc_db.proto", []string{}),
		runner.WithDataFromFile("./internal/sensor/10k_test.json"),
		runner.WithInsecure(true),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//p := printer.ReportPrinter{
	//	Out:    os.Stdout,
	//	Report: report,
	//}

	//if err2 := p.Print("pretty"); err2 != nil {
	//	fmt.Println(err2)
	//}

	//fmt.Println(report)
}