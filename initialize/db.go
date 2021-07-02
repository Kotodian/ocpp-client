package initialize

import (
	cmap "github.com/orcaman/concurrent-map"
	"ocpp-client/boltdb"
	"ocpp-client/service"
	"ocpp-client/websocket"
)

func initdb() {
	websocket.Cache = cmap.New()
	var err error
	service.DB, err = boltdb.New("ocpp.db", &service.ChargeStation{})
	if err != nil {
		panic(err)
	}

	err = service.DB.ForEach(new(service.ChargeStation).BucketName(), func(k string, v interface{}) error {
		_chargeStation := v.(*service.ChargeStation)
		chargeStation := service.NewChargeStation(_chargeStation.Sn)
		err := websocket.NewClient(chargeStation).ReConn()
		return err
	})
	if err != nil {
		service.DB.Close()
		panic(err)
	}
	service.DB.Close()
}
