package In_memo_db

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var emptydb = ""
var d0 = "sensor_1,Sunday,-9223372036854775808,9223372036854775807,0,"
var d1 = "sensor_1,monday,-9223372036854775808,9223372036854775807,0,"
var d2 = "sensor_1,Tuesday,-9223372036854775808,9223372036854775807,0,"
var d3 = "sensor_1,Wednesday,-9223372036854775808,9223372036854775807,0,"
var d4 = "sensor_1,Thursday,-9223372036854775808,9223372036854775807,0,"
var d5 = "sensor_1,Friday,-9223372036854775808,9223372036854775807,0,"
var d6 = "sensor_1,Saturday,-9223372036854775808,9223372036854775807,0,"

func TestGetInfoBySensor_emptyDB(t *testing.T) {
	testName := "TestGetInfoBySensor_emptyDB"
	mapDb := sensorMap()
	t.Run(testName, func(t *testing.T) {
		s := mapDb.getInfoBySensor("", 8)
		if s != emptydb {
			t.Errorf("got %v, want %v", s, emptydb)
		}
	})
}

func TestGetInfoBySensor_SensorsNoTraffic(t *testing.T) {
	testName := "TestGetInfoBySensor_SensorsNoTraffic"
	var tmgrpc_dbuff strings.Builder
	mapDB := sensorMap()
	for i := 1; i < 50; i++ {
		s := fmt.Sprintf("sensor_%d", i)
		mapDB.addSensorTomap(s)
		for i := 0; i < 7; i += 1 {
			if i != 0 {
				fmt.Fprintf(&tmgrpc_dbuff, "%v", fmt.Sprintf("%s,%s,-9223372036854775808,9223372036854775807,0,", "", time.Weekday(i)))
			} else {
				fmt.Fprintf(&tmgrpc_dbuff, "%v", fmt.Sprintf("%s,%s,-9223372036854775808,9223372036854775807,0,", s, time.Weekday(i)))
			}
		}
	}
	res := tmgrpc_dbuff.String()
	t.Run(testName, func(t *testing.T) {
		s := mapDB.GetInfo("all", 8)
		if len(s) != len(res) {
			t.Errorf("got %v, want %v", s, res)
		}
	})
}

// add measures to sernsor_1, then test all days ofsensor_1
func TestAddMeasure(t *testing.T) {
	testName, sname := "TestAddMeasure", "sensor_1"
	var tt string
	mapDB := sensorMap()
	mapDB.AddMeasure(sname, 10)
	mapDB.AddMeasure(sname, 30)

	var tests = []string{d0, d1, d2, d3, d4, d5, d6}
	now := time.Now().Weekday()
	today := int(now)

	for i := 0; i < 7; i++ {
		t.Run(testName, func(t *testing.T) {
			curDay := (today + i) % 7
			s := mapDB.getInfoBySensor(sname, curDay)
			//only today is having measures
			if curDay == today {
				tt = fmt.Sprintf("sensor_1,%s,30,10,20,", now)
			} else {
				tt = tests[curDay]
			}
			if s != tt {
				t.Errorf("go t:\t%v\nwant:\t%v\n\n", s, tt)
			}
		})
	}

}

func TestAddSensorTomap(t *testing.T) {
	mapDB := sensorMap()
	if len(*mapDB) != 0 {
		t.Errorf("got len=%v, want len=%v", len(*mapDB), 0)
	}
	for i := 1; i < 10; i++ {
		s := fmt.Sprintf("sensor_%d", i)
		mapDB.addSensorTomap(s)
		if len(*mapDB) != i {
			t.Errorf("got len=%v, want len=%v", len(*mapDB), i)
		}
	}
}
