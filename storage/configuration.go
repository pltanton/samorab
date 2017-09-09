package storage

import (
	"strconv"

	"github.com/boltdb/bolt"
)

func GetChance(id int) int {
	var result int64
	GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("configuration"))
		record := b.Get([]byte(strconv.Itoa(id)))

		if record == nil {
			result = 30
			return nil
		}
		result, _ = (strconv.ParseInt(string(record), 10, 32))
		return nil
	})
	return int(result)
}

func SetChance(id int, chance int) error {
	return GetDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("configuration"))
		return b.Put([]byte(strconv.Itoa(id)), []byte(strconv.Itoa(chance)))
	})
}
