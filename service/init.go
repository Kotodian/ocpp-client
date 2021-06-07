package service

import (
	"ocpp-client/boltdb"
)

var (
	DB                  *boltdb.BoltManager
	ChargeStationBucket = "chargeStation"
)

func init() {
	var err error
	buckets := make(map[string]boltdb.BoltType)
	buckets[ChargeStationBucket] = boltdb.NewBoltType(ChargeStation{}, []*ChargeStation{})
	DB, err = boltdb.New("ocpp", buckets)
	if err != nil {
		panic(err)
	}
}
