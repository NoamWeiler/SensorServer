package redis_db

type sensorDayDB struct {
	Max   int32 `redis:"int1"`
	Min   int32 `redis:"int2"`
	Count int32 `redis:"int3"`
	Sum   int32 `redis:"int4"`
}

//type sensorWeekDB struct {
//	Week []sensorDayDB `redis:"arr1"`
//}
