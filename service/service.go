package service

import (
	"ocpp-client/boltdb"
	"reflect"
)

var (
	DB                  *boltdb.BoltManager
	ChargeStationBucket = "chargeStation"
)

func init() {
	var err error
	buckets := make(map[string]reflect.Type)
	buckets[ChargeStationBucket] = reflect.TypeOf(&ChargeStation{})
	DB, err = boltdb.New("ocpp", buckets)
	if err != nil {
		panic(err)
	}
}
