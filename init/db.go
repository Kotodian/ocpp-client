package init

import (
	"ocpp-client/boltdb"
	"ocpp-client/service"
)

func init() {
	var err error
	service.DB, err = boltdb.New("ocpp.db", &service.ChargeStation{})
	if err != nil {
		panic(err)
	}
}
