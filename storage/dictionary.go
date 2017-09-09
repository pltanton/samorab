package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type DictionaryRecord struct {
	Original    string
	Synonim     []string
	Alternative []string
}

func FindAlternatives(key string) *DictionaryRecord {
	var result *DictionaryRecord
	GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("dictionary"))
		record := b.Get([]byte(key))

		if record == nil {
			return nil
		}

		json.Unmarshal(record, &result)

		return nil
	})
	return result
}
