package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/kljensen/snowball"
)

type DictionaryRecord struct {
	Original    string
	Synonim     []string
	Alternative []string
}

func FindAlternatives(key string) *DictionaryRecord {
	stemmedKey, _ := snowball.Stem(key, "russian", true)
	var result *DictionaryRecord
	GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("dictionary"))
		record := b.Get([]byte(stemmedKey))

		if record == nil {
			return nil
		}

		json.Unmarshal(record, &result)

		return nil
	})
	return result
}
