package initialize

import (
	"ocpp-client/boltdb"
	"ocpp-client/service"
	"ocpp-client/websocket"
)

func init() {
	var err error
	service.DB, err = boltdb.New("ocpp.db", &service.ChargeStation{})
	if err != nil {
		panic(err)
	}
	err = service.DB.ForEach(new(service.ChargeStation).BucketName(), func(k string, v interface{}) error {
		chargeStation := v.(*service.ChargeStation)
		err := websocket.NewClient(chargeStation).ReConn()
		return err
	})
	if err != nil {
		panic(err)
	}
}
