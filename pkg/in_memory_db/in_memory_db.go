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
	db map[string]*sensorWeekDB
	sync.RWMutex
}

type sensorDayDB struct {
	max   int32
	min   int32
	count int32
	sum   int32
}

type sensorWeekDB struct {
	week []sensorDayDB
}

//day implementation
func (s *sensorDayDB) getDayAvg() float32 {
	count := s.count
	if count == 0 {
		return 0.0
	}
	return float32(s.sum) / float32(count)
}

func (s *sensorDayDB) getDayRes() (int32, int32, float32) {
	return s.max, s.min, s.getDayAvg()
}

func (s *sensorDayDB) AddMeasure(m int32) {
	/*
		s.count.Inc()
		s.sum.Add(m)
		min := func(a, b int32) int32 {
			if a < b {
				return a
			}
			return b
		}(s.min.Load(), m)
		s.min.Swap(min)
		max := func(a, b int32) int32 {
			if a < b {
				return b
			}
			return a
		}(s.max.Load(), m)
		s.max.Swap(max)

	*/
	s.count++
	s.sum += m
	s.min = func(a, b int32) int32 {
		if a < b {
			return a
		}
		return b
	}(s.min, m)
	s.max = func(a, b int32) int32 {
		if a > b {
			return a
		}
		return b
	}(s.max, m)
}

func (s *sensorDayDB) resetDay() {
	s.max = math.MinInt32
	s.min = math.MaxInt32
	s.count = 0
	s.sum = 0
}

// AddMeasure week implementation
func (sw *sensorWeekDB) AddMeasure(m int32) {
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

func (sw *sensorWeekDB) getInfoBySensorWeek(s string, d int32) string {

	var output strings.Builder
	switch d {
	case 0, 1, 2, 3, 4, 5, 6:
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&sw.week[d], d)); err != nil {
			log.Println(err)
		}
	case 8: //all week
		for i, d := range sw.week {
			if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&d, int32(i))); err != nil {
				log.Println(err)
			}
		}
	case 9: //today
		today := int32(time.Now().Weekday())
		if _, err := fmt.Fprintf(&output, ",%v", buildDayString(&sw.week[today], today)); err != nil {
			log.Println(err)
		}
	default:
		log.Println("getInfoBySensor - error: wrong day option:", d)
		return ""
	}

	return fmt.Sprintf("%s%s", s, output.String())
}

// AddMeasure - implementation of sensorDB interface
func (sm *sensormap) AddMeasure(serial string, measure int32) {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := sm.db[serial]; !ok {
		sm.addSensorToMap(serial)
	}
	sm.db[serial].AddMeasure(measure)
}

func (sm *sensormap) getInfoAllSensors(day int32) string {
	sm.RLock()
	defer sm.RUnlock()
	var output strings.Builder
	resChan := make(chan string, sm.len())
	var wg = &sync.WaitGroup{}

	for serial, sensorWeek := range sm.db {
		wg.Add(1)
		go func(c chan<- string, sensorWeek *sensorWeekDB, s string, d int32) {
			defer wg.Done()
			c <- sensorWeek.getInfoBySensorWeek(s, d)
		}(resChan, sensorWeek, serial, day)
	}

	//goroutine receiver for results from all the sensorWeeks
	go func(total int) {
		wg.Add(1)
		var counter = 0
		for sensorRes := range resChan {
			if _, err := fmt.Fprintf(&output, "%v", sensorRes); err != nil {
				log.Println(err)
			}
			counter++
			if counter == total {
				close(resChan)
				wg.Done()
			}
		}
	}(sm.len())

	wg.Wait()
	return output.String()
}

func buildDayString(day *sensorDayDB, d int32) string {
	a, b, c := day.getDayRes()
	//order: sensorSerial,day,max,min,avg
	return fmt.Sprintf("%v,%v,%v,%v,", time.Weekday(d), a, b, c)
}

func (sm *sensormap) getInfoBySensor(s string, d int32) string {
	if s == "" {
		return s
	}
	sm.RLock()
	defer sm.RUnlock()
	if _, ok := sm.db[s]; !ok {
		return ""
	}

	return sm.db[s].getInfoBySensorWeek(s, d)
}

func (sm *sensormap) GetInfo(serial string, daysBefore int32) string {
	if serial == "all" {
		return sm.getInfoAllSensors(daysBefore)
	}
	return sm.getInfoBySensor(serial, daysBefore)
}

func (sm *sensormap) addSensorToMap(s string) {
	sw := newSensorWeek()
	sm.db[s] = sw
}

func SensorMap() *sensormap {
	GlobalDay = time.Now().Weekday() //update global
	output := &sensormap{db: make(map[string]*sensorWeekDB, 1000)}
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
	sm.Lock()
	defer sm.Unlock()
	fname := "dayCleanup"
	var wg sync.WaitGroup
	now := time.Now().Weekday()

	//if same day - no need to cleanup day from all sensors
	if GlobalDay == now {
		return
	}

	log.Println(fname, "Starting day cleanup in DB")

	for _, sensorWeek := range sm.db {
		wg.Add(1)

		go func(s *sensorWeekDB) {
			defer wg.Done()
			s.cleanDay(now)
		}(sensorWeek)
	}

	wg.Wait()
	log.Println(fname, now)
	GlobalDay = now //update global
}

func (sm *sensormap) len() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.db)
}

//TODO
/*
1)

*/
