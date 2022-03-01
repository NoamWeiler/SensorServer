package main

import (
	pb "SensorServer/internal/mutual_db"
	"errors"
	"fmt"
)

const userExit = "user Exit"

func dayOpt() int {
	var opt int
	for {
		fmt.Println("\nplease enter a number between 1-6 (day before today)\n10 - to exit")
		_, err := fmt.Scanf("%d", &opt)
		myPanic(err)
		switch {
		case 0 < opt && opt < 7:
			fmt.Println("OPT:", opt)
			return opt
		case opt == 10:
			return -1
		default:
			fmt.Println("Illegal option")
		}
	}
}

/*
	return values:
	1-6 :	number of days before today
	8	:	all week
	9	:	today
*/
func showDMenu() int {
	var opt int
	for {
		fmt.Printf("\nchoose day option\n1)\tshow by day - all past week\n2)\tshow by day - specific day\n3)\tshow by day - today\n5)\texit\n")
		_, err := fmt.Scanf("%d", &opt)
		myPanic(err)
		switch opt {
		case 1:
			return 8
		case 2:
			return dayOpt()
		case 3:
			return 9
		case 5:
			return -1
		default:
			fmt.Println("Illegal option")
		}
	}
}

func showMainMenu() int {
	var opt int
	for {
		fmt.Printf("\nchoose day option\n1)\tget info\n2)\tdisconnect\n3)\texit\n")
		_, err := fmt.Scanf("%d", &opt)
		myPanic(err)
		switch opt {
		case 1, 2, 3:
			return opt
		default:
			fmt.Println("Illegal option")
		}
	}
}

func sensorOpt() string {
	var output string
	fmt.Println("\nenter the sensor name (for example: sensor_1)")
	_, err := fmt.Scanf("%s", &output)
	myPanic(err)
	return output
}

func showSMenu() string {
	var opt int
	for {
		fmt.Printf("\nplease choose an option for sensors:\n1)\tshow by sensor - all sensors\n2)\tshow by sensor - specific sensor\n5)\texit\n")
		_, err := fmt.Scanf("%d", &opt)
		myPanic(err)
		switch opt {
		case 1:
			return "all"
		case 2:
			return sensorOpt()
		case 5:
			return userExit
		default:
			fmt.Println("Illegal option")
		}
	}
}

func showMenu() (int32, string) {
	d := showDMenu()
	if d == -1 { //if already want to quit - exit without further menu options
		return int32(d), ""
	}
	s := showSMenu()
	return int32(d), s
}

func createRequest() *pb.InfoReq {
	d, s := showMenu()
	return &pb.InfoReq{DayBefore: d, SensorName: s}
}

func clientMenu() (*pb.InfoReq, error) {
	cr := createRequest()
	if cr.SensorName == userExit || cr.DayBefore == -1 {
		return nil, errors.New(userExit)
	}
	return cr, nil
}
