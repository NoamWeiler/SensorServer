package main

import (
	pb "SensorServer/internal/mutual_db"
	"fmt"
	"strings"
	"testing"
	"time"
)

var emptydb = ""

//r1 = 1 sensor, empty
//r2 = 1 sensor, 1 min=10, max=30, avg=20
//r3 = 1 sensor with stream
var r1 = "sensor_1,MONDAY,0,0,0,sensor_1,TUESDAY,0,0,0,sensor_1,WEDNESDAY,0,0,0,sensor_1,THURSDAY,0,0,0,sensor_1,FRIDAY,0,0,0,sensor_1,SATURDAY,0,0,0,sensor_1,SUNDAY,0,0,0,sensor_1,MONDAY,0,0,0,sensor_1,TUESDAY,0,0,0,sensor_1,WEDNESDAY,0,0,0,"
var r3 = "sensor_1,MONDAY,10,20,15,sensor_1,TUESDAY,10,20,15,sensor_1,WEDNESDAY,10,20,15,sensor_1,THURSDAY,10,20,15,sensor_1,FRIDAY,10,20,15,sensor_1,SATURDAY,10,20,15,sensor_1,SUNDAY,10,20,15,sensor_1,MONDAY,10,20,15,sensor_1,TUESDAY,10,20,15,sensor_1,WEDNESDAY,10,20,15,"
var r2 = "sensor_1,Sunday,-9223372036854775808,9223372036854775807,0,sensor_1,Monday,-9223372036854775808,9223372036854775807,0,sensor_1,Tuesday,30,10,20,sensor_1,Wednesday,-9223372036854775808,9223372036854775807,0,sensor_1,Thursday,-9223372036854775808,9223372036854775807,0,sensor_1,Friday,-9223372036854775808,9223372036854775807,0,sensor_1,Saturday,-9223372036854775808,9223372036854775807,0,1"
var d0 = "sensor_1,Sunday,-9223372036854775808,9223372036854775807,0,"
var d1 = "sensor_1,Monday,-9223372036854775808,9223372036854775807,0,"
var d2 = "sensor_1,Tuesday,30,10,20,"
var d3 = "sensor_1,Wednesday,-9223372036854775808,9223372036854775807,0,"
var d4 = "sensor_1,Thursday,-9223372036854775808,9223372036854775807,0,"
var d5 = "sensor_1,Friday,-9223372036854775808,9223372036854775807,0,"
var d6 = "sensor_1,Saturday,-9223372036854775808,9223372036854775807,0,"

func TestGetInfoBySensor_emptyDB(t *testing.T) {
	testName := "TestGetInfoBySensor_emptyDB"
	mapDb := SensorMap()
	t.Run(testName, func(t *testing.T) {
		s := mapDb.getInfoBySensor("", 8)
		if s != emptydb {
			t.Errorf("got %v, want %v", s, emptydb)
		}
	})
}

func TestGetInfoBySensor_SensorsNoTraffic(t *testing.T) {
	testName := "TestGetInfoBySensor_SensorsNoTraffic"
	var tmpBuff strings.Builder
	mapDB := SensorMap()
	for i := 1; i < 50; i++ {
		s := fmt.Sprintf("sensor_%d", i)
		mapDB.addSensorToMap(s)
		for i := 0; i < 7; i += 1 {
			fmt.Fprintf(&tmpBuff, fmt.Sprintf("%s,%s,0,0,0,", s, time.Weekday(i)))
		}
	}
	res := tmpBuff.String()
	t.Run(testName, func(t *testing.T) {
		s := mapDB.getInfoAllSensors(8)
		if len(s) != len(res) {
			t.Errorf("got %v, want %v", s, res)
		}
	})
}

// add measures to sernsor_1, then test all days ofsensor_1
func TestAddMeasure(t *testing.T) {
	testName, sname := "TestAddMeasure", "sensor_1"
	mapDB := SensorMap()
	mapDB.addMeasure(&pb.Measure{Serial: sname, M: 10})
	mapDB.addMeasure(&pb.Measure{Serial: sname, M: 30})

	var tests = []string{d0, d1, d2, d3, d4, d5, d6}
	today := int(time.Now().Weekday())

	for i := 1; i < 7; i++ {
		t.Run(testName, func(t *testing.T) {
			curDay := (today + i) % 7
			s := mapDB.getInfoBySensor(sname, curDay)
			tt := tests[curDay]
			if s != tt {
				t.Errorf("got:\t%v\nwant:\t%v\n\n", s, tt)
			}
		})
	}

}

func TestAddSensorToMap(t *testing.T) {
	mapDB := SensorMap()
	if len(mapDB) != 0 {
		t.Errorf("got len=%v, want len=%v", len(mapDB), 0)
	}
	for i := 1; i < 10; i++ {
		s := fmt.Sprintf("sensor_%d", i)
		mapDB.addSensorToMap(s)
		if len(mapDB) != i {
			t.Errorf("got len=%v, want len=%v", len(mapDB), i)
		}
	}

}
