package boltdb

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"reflect"
)

type Codec interface {
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(data []byte, value interface{}) error
}

type BoltManager struct {
	db              *bolt.DB
	codec           Codec
	bucketValueType map[string][]reflect.Type
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

func NewBucket(structType interface{}, structSliceType interface{}) []reflect.Type {
	bucketValueType := make([]reflect.Type, 0)
	bucketValueType = append(bucketValueType, reflect.TypeOf(structType), reflect.TypeOf(structSliceType))
	return bucketValueType
}

// New 创建库管理
func New(path string, bucket map[string][]reflect.Type) (*BoltManager, error) {
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	if len(bucket) == 0 {
		goto create
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for k, _ := range bucket {
			_, err = tx.CreateBucketIfNotExists([]byte(k))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
create:
	return &BoltManager{db: db, codec: NewJsonCodec(), bucketValueType: bucket}, nil
}

// Close 关闭连接
func (b *BoltManager) Close() error {
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
	delete(b.bucketValueType, bucketName)
	return nil
}

// Put 往Bucket里增加键值对
func (b *BoltManager) Put(bucketName string, key string, value interface{}) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		msg, err := b.codec.Marshal(value)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), msg)
	})
	return
}

func (b *BoltManager) List(bucketName string) (interface{}, error) {
	slice := reflect.MakeSlice(b.bucketValueType[bucketName][1], 0, 100)
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.ForEach(func(k, v []byte) error {
			value := reflect.New(b.bucketValueType[bucketName][0]).Interface()
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
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.ForEach(func(k, v []byte) error {
			value := reflect.New(b.bucketValueType[bucketName][0]).Interface()
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

func (b *BoltManager) Get(bucketName string, key string, value interface{}) (err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return b.codec.Unmarshal(bucket.Get([]byte(key)), value)
	})
	return
}

// Remove 删除指定的键
func (b *BoltManager) Remove(bucketName string, key string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.DeleteBucket([]byte(key))
	})
	return
}
