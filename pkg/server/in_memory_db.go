package main

import (
	pb "SensorServer/internal/mutual_db"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

type sensorMap map[string]sensorWeekDB

type sensorDayDB struct {
	max   int
	min   int
	count int
	sum   int
}

type sensorWeekDB struct {
	week []sensorDayDB
	mu   sync.Mutex
}

//day implementation
func (s *sensorDayDB) getDayAvg() float32 {
	if s.count == 0 {
		return 0.0
	}
	return float32(s.sum / s.count)
}

func (s *sensorDayDB) getDayRes() (int, int, float32) {
	return s.max, s.min, s.getDayAvg()
}

func (s *sensorDayDB) addMeasure(m int) {
	s.count++
	s.sum += m
	s.min = func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}(s.min, m)
	s.max = func(a, b int) int {
		if a < b {
			return b
		}
		return a
	}(s.max, m)
}

//week implementation
func (sw *sensorWeekDB) addMeasure(m int) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	dayIndex := int(time.Now().Weekday()) //Sunday=0
	sw.week[dayIndex].addMeasure(m)
}

func (sm sensorMap) addMeasure(m *pb.Measure) {
	serial := m.GetSerial()
	_, ok := sm[serial]
	if !ok {
		sm.addSensorToMap(serial)
	}
	elem := sm[serial]
	elem.addMeasure(int(m.GetM()))
}

//implementation of sensorDB interface
func (sm sensorMap) getInfoAllSensors(day int) string {
	var output strings.Builder
	for k, _ := range sm {
		str := sm.getInfoBySensor(k, day)
		if _, err := fmt.Fprintf(&output, str); err != nil {
			log.Println(err)
			return fmt.Sprintf("Error:%v", err)
		}
	}
	return output.String()
}

func buildDayString(s *strings.Builder, day *sensorDayDB, d int) {
	a, b, c := day.getDayRes()
	//order: sensorSerial,day,max,min,avg
	_, err := fmt.Fprintf(s, "%v,%v,%v,%v,%v,", s, time.Weekday(d), a, b, c)
	if err != nil {
		debug("getInfoBySensor:", fmt.Sprintf("%v", err))
	}
}

func (sm sensorMap) getInfoBySensor(s string, d int) string {
	//get element
	elem, ok := sm[s]
	if !ok {
		return ""
	}
	elem.mu.Lock()
	defer elem.mu.Unlock()

	var output strings.Builder
	switch d {
	case 0, 1, 2, 3, 4, 5, 6:
		buildDayString(&output, &elem.week[d], d)
	case 8:
		for i, d := range elem.week {
			buildDayString(&output, &d, i)
		}
	case 9:
		today := int(time.Now().Weekday())
		buildDayString(&output, &elem.week[today], today)
	default:
		log.Println("getInfoBySensor - error: wrong day option:", d)
		return ""
	}

	return fmt.Sprintf("%s%s", s, output.String())

}

func (sm sensorMap) getInfo(res *pb.InfoReq) string {
	d := int(res.GetDayBefore())
	s := res.GetSensorName()
	if s == "all" {
		return sm.getInfoAllSensors(d)
	}
	return sm.getInfoBySensor(s, d)
}

func (sm sensorMap) addSensorToMap(s string) {
	sw := sensorWeekDB{week: make([]sensorDayDB, 7)}
	sww := sw.week
	for i, _ := range sww {
		sww[i].min = math.MaxInt
		sww[i].max = math.MinInt
	}
	sm[s] = sw
}

func SensorMap() sensorMap {
	return sensorMap{}
}
