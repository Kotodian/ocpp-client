package boltdb

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"reflect"
	"sync/atomic"
)

type Codec interface {
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(data []byte, value interface{}) error
}

type Interface interface {
	// BucketName 桶名称
	BucketName() string
	// Type 桶里数据类型
	Type() reflect.Type
	// SliceType 桶里的slice类型
	SliceType() reflect.Type
}

type BoltManager struct {
	path        string
	db          *bolt.DB
	connected   int32
	codec       Codec
	bucketsType map[string]Interface
}

type defaultCodec struct{}

func NewJsonCodec() Codec {
	return &defaultCodec{}
}

func (d *defaultCodec) Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (d *defaultCodec) Unmarshal(data []byte, value interface{}) error {
	return json.Unmarshal(data, value)
}

// New 创建库管理
func New(path string, data ...Interface) (*BoltManager, error) {
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	manager := &BoltManager{
		db:          db,
		path:        path,
		codec:       NewJsonCodec(),
		bucketsType: make(map[string]Interface),
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, v := range data {
			_, err = tx.CreateBucketIfNotExists([]byte(v.BucketName()))
			if err != nil {
				return err
			}
			manager.bucketsType[v.BucketName()] = v
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	manager.SetConnected(true)
	return manager, nil
}

// Close 关闭连接
func (b *BoltManager) Close() error {
	b.SetConnected(false)
	return b.db.Close()
}

// RemoveBucket 移除Bucket
func (b *BoltManager) RemoveBucket(bucketName string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
	if err != nil {
		return err
	}
	delete(b.bucketsType, bucketName)
	return nil
}

// Put 往Bucket里增加键值对
func (b *BoltManager) Put(key string, value Interface) (err error) {
	err = b.conn()
	if err != nil {
		return err
	}
	defer b.Close()
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(value.BucketName()))
		msg, err := b.codec.Marshal(value)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), msg)
	})
	return
}

func (b *BoltManager) List(bucketName string) (interface{}, error) {
	err := b.conn()
	if err != nil {
		return nil, err
	}
	defer b.Close()
	slice := reflect.MakeSlice(b.bucketsType[bucketName].SliceType(), 0, 100)
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.ForEach(func(k, v []byte) error {
			value := reflect.New(b.bucketsType[bucketName].Type()).Interface()
			err := b.codec.Unmarshal(v, value)
			if err != nil {
				return err
			}
			slice = reflect.Append(slice, reflect.ValueOf(value))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return slice.Interface(), nil
}

// ForEach 对每个键值对做处理
func (b *BoltManager) ForEach(bucketName string, handle func(k string, v interface{}) error) (err error) {
	err = b.conn()
	if err != nil {
		return err
	}
	defer b.Close()
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.ForEach(func(k, v []byte) error {
			value := reflect.New(b.bucketsType[bucketName].Type()).Interface()
			err := b.codec.Unmarshal(v, value)
			if err != nil {
				return err
			}
			return handle(string(k), value)
		})
		return err
	})
	return
}

func (b *BoltManager) Get(key string, value Interface) (err error) {
	err = b.conn()
	if err != nil {
		return err
	}
	defer b.Close()
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(value.BucketName()))
		return b.codec.Unmarshal(bucket.Get([]byte(key)), value)
	})
	return
}

// Remove 删除指定的键
func (b *BoltManager) Remove(bucketName string, key string) (err error) {
	err = b.conn()
	if err != nil {
		return err
	}
	defer b.Close()
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.DeleteBucket([]byte(key))
	})
	return
}

func (b *BoltManager) conn() error {
	if b.Connected() {
		return nil
	}
	db, err := bolt.Open(b.path, 0644, nil)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

func (b *BoltManager) Connected() bool {
	return atomic.LoadInt32(&b.connected) == 1
}

func (b *BoltManager) SetConnected(isConnected bool) {
	if isConnected {
		atomic.StoreInt32(&b.connected, 1)
	} else {
		atomic.StoreInt32(&b.connected, 0)
	}
}
