package main

import (
	"fmt"
	"github.com/bojand/ghz/runner"
	"os"
)

func main() {
	report, err := runner.Run(
		"grpc_db.SensorStream.SensorMeasure",
		"localhost:50051",
		runner.WithProtoFile("./internal/grpc_db/grpc_db.proto", []string{}),
		runner.WithDataFromFile("./pkg/sensor/data.json"),
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

	fmt.Println(report)
}
