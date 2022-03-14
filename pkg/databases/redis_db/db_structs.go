package redis_db

type sensorDayDB struct {
	Max   int32 `redis:"int1"`
	Min   int32 `redis:"int2"`
	Count int32 `redis:"int3"`
	Sum   int32 `redis:"int4"`
}

type sensorWeekDB struct {
	Week []sensorDayDB `redis:"arr1"`
}

//
//func (sd *sensorDayDB) MarshalBinary() ([]byte, error) {
//	return json.Marshal(sd)
//}
//
//func (sd *sensorDayDB) UnmarshalBinary(data []byte) error {
//	if err := json.Unmarshal(data, &sd); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (sw *sensorWeekDB) MarshalBinary() ([]byte, error) {
//	return json.Marshal(sw)
//}
//
//func (sw *sensorWeekDB) UnmarshalBinary(data []byte) error {
//	if err := json.Unmarshal(data, &sw); err != nil {
//		return err
//	}
//	return nil
//}
