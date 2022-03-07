package In_memo_db

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	GlobalDay time.Weekday
)

type sensormap struct {
	sync.Map //implements SensorDB interface
}

type sensorDayDB struct {
	max   int
	min   int
	count int
	sum   int
}

type sensorWeekDB struct {
	week []sensorDayDB
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

func (s *sensorDayDB) AddMeasure(m int) {
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

// AddMeasure week implementation
func (sw *sensorWeekDB) AddMeasure(m int) {
	dayIndex := int(time.Now().Weekday()) //Sunday=0
	sw.week[dayIndex].AddMeasure(m)
}

func (sw *sensorWeekDB) cleanDay(weekday time.Weekday) {
	d := int(weekday)
	sw.week[d].resetDay()
}

func newSensorWeek() *sensorWeekDB {
	sw := &sensorWeekDB{week: make([]sensorDayDB, 7)}
	sww := sw.week
	for i := range sww {
		sww[i].resetDay()
	}
	return sw
}

func (sw *sensorWeekDB) getInfoBySensorWeek(s string, d int) string {

	var output strings.Builder
	switch d {
	case 0, 1, 2, 3, 4, 5, 6:
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&sw.week[d], d)); err != nil {
			log.Println(err)
		}
	case 8: //all week
		for i, d := range sw.week {
			if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&d, i)); err != nil {
				log.Println(err)
			}
		}
	case 9: //today
		today := int(time.Now().Weekday())
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&sw.week[today], today)); err != nil {
			log.Println(err)
		}
	default:
		log.Println("getInfoBySensor - error: wrong day option:", d)
		return ""
	}

	return fmt.Sprintf("%s%s", s, output.String())
}

func (sm *sensormap) get(s string) *sensorWeekDB {
	interfaceValue, ok := sm.Load(s)
	if !ok {
		return nil
	}
	return interfaceValue.(*sensorWeekDB) //cast to sensorWeekDB, dou to LoadOrStore returns interface{}
}

// AddMeasure - implementation of sensorDB interface
func (sm *sensormap) AddMeasure(serial string, measure int) {
	sw := sm.get(serial)
	if sw == nil {
		sm.addSensorToMap(serial)
		sw = sm.get(serial)
	}
	sw.AddMeasure(measure)
}

func (sm *sensormap) getInfoAllSensors(day int) string {
	var output strings.Builder
	mapLen := sm.len()
	strChan := make(chan string, mapLen)
	sm.Range(func(k, v interface{}) bool {
		go func(c chan<- string, sensormapElem *sensorWeekDB, s string, d int) {
			c <- sensormapElem.getInfoBySensorWeek(s, d)
		}(strChan, v.(*sensorWeekDB), k.(string), day)
		return true
	})

	//get the results from all the sensorWeeks
	for i := 0; i < mapLen; i++ {
		sensorRes := <-strChan
		if _, err := fmt.Fprintf(&output, "%v", sensorRes); err != nil {
			log.Println(err)
		}
	}

	return output.String()
}

func buildDayString(day *sensorDayDB, d int) string {
	a, b, c := day.getDayRes()
	//order: sensorSerial,day,max,min,avg
	return fmt.Sprintf("%v,%v,%v,%v,", time.Weekday(d), a, b, c)
}

func (sm *sensormap) getInfoBySensor(s string, d int) string {
	if s == "" {
		return s
	}
	elem := sm.get(s)

	if elem == nil {
		return ""
	}

	return elem.getInfoBySensorWeek(s, d)
}

func (sm *sensormap) GetInfo(serial string, daysBefore int) string {
	if serial == "all" {
		return sm.getInfoAllSensors(daysBefore)
	}
	return sm.getInfoBySensor(serial, daysBefore)
}

func (sm *sensormap) addSensorToMap(s string) {
	sw := newSensorWeek()
	(*sm).Store(s, sw)
}

func SensorMap() *sensormap {
	GlobalDay = time.Now().Weekday() //update global
	output := &sensormap{}
	return output
}

/*
	Update that occur every AddMeasure and getInfo
	The design:
	Before client getInfo or sensor AddMeasure -
	Check if the day have been changed since last request
	If not - continue
	If so - need to clean the current day (run on parallel on all sensorWeekDB and tell then to reset the day)
*/
func (sm *sensormap) DayCleanup() {
	fname := "dayCleanup"
	var wg sync.WaitGroup
	now := time.Now().Weekday()

	//if same day - no need to cleanup day from all sensors
	if GlobalDay == now {
		return
	}

	log.Println(fname, "Starting day cleanup in DB")

	sm.Range(func(_, v interface{}) bool {
		sensorWeek := v.(*sensorWeekDB)
		wg.Add(1)

		go func(s *sensorWeekDB) {
			defer wg.Done()
			s.cleanDay(now)
		}(sensorWeek)

		return true
	})

	wg.Wait()
	log.Println(fname, now)
	GlobalDay = now //update global
}

func (sm *sensormap) len() int {
	length := 0
	sm.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
