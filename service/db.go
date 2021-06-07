package service

import (
	"ocpp-client/boltdb"
)

var DB *boltdb.BoltManager

// ListChargeStation 调用boltdb获取数据库中所有的充电桩
func ListChargeStation() ([]*ChargeStation, error) {
	bucketName := new(ChargeStation).BucketName()
	_list, err := DB.List(bucketName)
	if err != nil {
		return nil, err
	}
	if _, ok := _list.([]*ChargeStation); !ok {
		return nil, nil
	}
	return _list.([]*ChargeStation), nil
}
