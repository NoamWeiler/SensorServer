package client

import (
	"bytes"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

var d1 = ""
var d2 = "sensor_1,Tuesday,30,10,20,"
var d3 = "sensor_1,Sunday,-,-,-,sensor_1,monday,-,-,-,sensor_1,Tuesday,30,10,20,sensor_1,Wednesday,-,-,-,sensor_1,Thursday,-,-,-,sensor_1,Friday,-,-,-,sensor_1,Saturday,-,-,-,1"

var res1 = "+---------+-----+-----+-----+-----+\n| #SERIAL | DAY | MIN | MAX | AVG |\n+---------+-----+-----+-----+-----+\n+---------+-----+-----+-----+-----+\n"
var res2 = "+----------+---------+-----+-----+-----+\n| #SERIAL  | DAY     | MIN | MAX | AVG |\n+----------+---------+-----+-----+-----+\n| sensor_1 | Tuesday | 30  | 10  | 20  |\n+----------+---------+-----+-----+-----+\n"
var res3 = "+----------+-----------+-----+-----+-----+\n| #SERIAL  | DAY       | MIN | MAX | AVG |\n+----------+-----------+-----+-----+-----+\n| sensor_1 | Sunday    | -   | -   | -   |\n| sensor_1 | monday    | -   | -   | -   |\n| sensor_1 | Tuesday   | 30  | 10  | 20  |\n| sensor_1 | Wednesday | -   | -   | -   |\n| sensor_1 | Thursday  | -   | -   | -   |\n| sensor_1 | Friday    | -   | -   | -   |\n| sensor_1 | Saturday  | -   | -   | -   |\n+----------+-----------+-----+-----+-----+\n"

func execToString(f func(s string), args string) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f(args)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	// back to normal state
	if err := w.Close(); err != nil {
		panic(err)
	}
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out
}

/*
	Test description:
	first:
		sanity test - first test that the function is working as expected.
	then:
		create DB and simulate measures sends to it
		get the res from the DB
		compare the tables
*/
func TestPrintTable(t *testing.T) {
	var tests = []struct {
		raw  string
		want string
	}{
		{d1, res1},
		{d2, res2},
		{d3, res3},
	}

	//first, test only the table output
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.raw)
		t.Run(testname, func(t *testing.T) {

			res := execToString(printResult, tt.raw)
			if res != tt.want {
				t.Errorf("\ngot\n %v\n\nwant\n %v", res, tt.want)
			}
		})
	}
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
	t.Render()
}

func printHelper(min, max, avg string) (string, string, string) {
	if min == strconv.Itoa(math.MinInt) {
		return "-", "-", "-"
	}
	return min, max, avg
}
