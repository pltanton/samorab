package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kljensen/snowball"

	"github.com/pltanton/samorab/utils"
)

type DictionarySynchronizer struct {
	clock *time.Ticker
}

func NewDictionarySynchronizer() DictionarySynchronizer {
	synchronizer := DictionarySynchronizer{
		time.NewTicker(12 * time.Hour),
	}
	return synchronizer
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func (d DictionarySynchronizer) Synchronize() {
	response, err := downloadDictionary()
	if err != nil {
		log.Println("Can't download dictionary ", err)
		return
	}
	r := csv.NewReader(response)

	GetDb().Update(func(tx *bolt.Tx) error {
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Can't parse dictionary")
				return err
			}

			dbKey, err := snowball.Stem(record[1], "russian", true)

			dbRecord := DictionaryRecord{
				record[1],
				deleteEmpty(utils.SPLIT_REGEX.Split(record[6], -1)),
				deleteEmpty(utils.SPLIT_REGEX.Split(record[7], -1)),
			}

			b := tx.Bucket([]byte("dictionary"))
			marshaledRecord, _ := json.Marshal(dbRecord)

			b.Put([]byte(dbKey), marshaledRecord)
		}
		return nil
	})

	log.Println("Successfully update dictionary")
}

func downloadDictionary() (io.ReadCloser, error) {
	id := "924233218"
	key := "1IkzUV-1Lq8bI264gGGz5JzI0EEEVio8OBtw-XORtIdk"
	response, err := http.Get(
		fmt.Sprintf("https://docs.google.com/spreadsheets/d/%v/gviz/tq?tqx=out:csv&sheet=%v", key, id),
	)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (d DictionarySynchronizer) Start() {
	for {
		d.Synchronize()
		_ = <-d.clock.C
	}
}
