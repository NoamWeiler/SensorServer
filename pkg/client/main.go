package main

import (
	pb "SensorServer/internal/mutual_db"
	"context"
	"flag"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/grpc"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultName           = "yochbad"
	loginConnectedMessage = "Connected successfully"
	loginCredentialsError = "Wrong credentials"
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	name    = flag.String("name", defaultName, "Name to greet")
	verbose = flag.Bool("v", true, "Verbose mode")
)

func myPanic(e error) {
	if e != nil {
		panic(fmt.Sprintf("%v", e))
	}
}

func newConnReq() *pb.ConnReq {
	fields := make([]string, 2)
	for i := 0; i < 2; i++ {
		switch i {
		case 0:
			fmt.Println("enter user name")
		case 1:
			fmt.Println("enter password")
		}
		_, err := fmt.Scanf("%s", &fields[i])
		myPanic(err)
	}
	return &pb.ConnReq{UserName: fields[0], Password: fields[1]}
}

func verifyLogin(r *pb.ConnRes, err error) bool {
	res := ""
	if err != nil {
		if e := unpackError(err); e == loginCredentialsError {
			log.Println("Error:", e)
		} else {
			log.Fatalf("Error:%s", e)
		}
	} else {
		res = r.GetRes()
		log.Printf("Res: %s", res)
	}
	return res == loginConnectedMessage
}

func unpackError(e error) string {
	s := fmt.Sprintf("%v", e)
	return s[strings.LastIndex(s, "=")+2:]
}

func printHelper(min, max, avg string) (string, string, string) {
	if min == strconv.Itoa(math.MinInt) {
		return "-", "-", "-"
	}
	return min, max, avg
}

func printResult(s string) {
	arr := strings.Split(s, ",")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"#SERIAL", "DAY", "MIN", "MAX", "AVG"})
	for i := 0; i < len(arr)-1; i += 5 {
		if arr[i] != "" && i > 0 {
			t.AppendSeparator()
		}
		a, b, c := printHelper(arr[i+2], arr[i+3], arr[i+4])
		t.AppendRows([]table.Row{
			{arr[i], arr[i+1], a, b, c},
		})
	}
	//t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()
}

func main() {
	flag.Parse()
	isConnected := false
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewClientInfoClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer func() { //cleanup
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
		defer cancel()
	}()

forLoop:
	for {
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		if !isConnected {
			cr := newConnReq()
			r, err := c.ConnectClient(ctx, cr)
			isConnected = verifyLogin(r, err)
		} else {
			switch showMainMenu() {
			case 1: //got info
				ir, err := clientMenu()
				if err != nil && fmt.Sprintf("%v", err) == userExit {
					continue //got error if userExit from menu
				}
				res, err := c.GetInfo(ctx, ir)
				if err != nil {
					fmt.Println(err)
				} else { //got response from server
					fmt.Println("RAW_RES:", res)
					printResult(res.GetResponce())
				}
			case 2: //disconnect
				r, err := c.DisconnectClient(ctx, &pb.DisConnReq{})
				if err != nil {
					fmt.Println("Error:", unpackError(err))
				} else {
					log.Println(r.GetRes())
					isConnected = false
				}
			case 3: //exit
				_, _ = c.DisconnectClient(ctx, &pb.DisConnReq{}) //disconnect by default on exit - my design
				break forLoop
			default:
			}

		}
	}

	fmt.Println("exit..")
}

//TODO
/*


 */
