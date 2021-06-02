package service

import (
	"ocpp-client/boltdb"
	"reflect"
)

var (
	DB                  *boltdb.BoltManager
	ChargeStationBucket = "chargeStation"
	TransactionBucket   = "transaction"
)

func init() {
	var err error
	buckets := make(map[string]reflect.Type)
	buckets[ChargeStationBucket] = reflect.TypeOf(&ChargeStation{})
	buckets[TransactionBucket] = reflect.TypeOf(&Transaction{})
	DB, err = boltdb.New("ocpp", buckets)
	if err != nil {
		panic(err)
	}
}
