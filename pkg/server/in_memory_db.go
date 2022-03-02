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

type sensorMap map[string]*sensorWeekDB

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
	return float32(s.sum) / float32(s.count)
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

func (s *sensorDayDB) resetDay() {
	s.max = math.MinInt
	s.min = math.MaxInt
	s.count = 0
	s.sum = 0
}

//week implementation
func (sw *sensorWeekDB) addMeasure(m int) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	dayIndex := int(time.Now().Weekday()) //Sunday=0
	sw.week[dayIndex].addMeasure(m)
}

func (sw *sensorWeekDB) cleanDay(weekday time.Weekday) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	d := int(weekday)
	sw.week[d].resetDay()
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

func buildDayString(day *sensorDayDB, d int) string {
	a, b, c := day.getDayRes()
	//order: sensorSerial,day,max,min,avg
	return fmt.Sprintf("%v,%v,%v,%v,", time.Weekday(d), a, b, c)
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
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&elem.week[d], d)); err != nil {
			log.Println(err)
		}
	case 8: //all week
		for i, d := range elem.week {
			if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&d, i)); err != nil {
				log.Println(err)
			}
		}
	case 9: //today
		today := int(time.Now().Weekday())
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&elem.week[today], today)); err != nil {
			log.Println(err)
		}
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
		sww[i].resetDay()
	}
	sm[s] = &sw
}

func SensorMap() sensorMap {
	return sensorMap{}
}

/*
	Update that occur every addMeasure and getInfo
	The design:
	Before client getInfo or sensor addMeasure -
	Check if the day have been changed since last request
	If not - continue
	If so - need to clean the current day (run on parallel on all sensorWeekDB and tell then to reset the day)
*/
func (sm sensorMap) dayCleanup() {
	fname := "dayCleanup"
	var wg sync.WaitGroup
	now := time.Now().Weekday()

	//if same day - no need to cleanup day from all sensors
	if GlobalDay == now {
		return
	}

	debug(fname, "Starting day cleanup in DB")
	wg.Add(len(sm))
	for _, v := range sm {
		go func(s *sensorWeekDB) {
			defer wg.Done()
			s.cleanDay(now)
		}(v)
	}

	wg.Wait()
	debug(fname, fmt.Sprintf("finished update,GlobalDay=%v\n", now))
	GlobalDay = now

}
