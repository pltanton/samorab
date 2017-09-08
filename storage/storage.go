package storage

import (
	"log"

	"github.com/boltdb/bolt"

	"github.com/pltanton/samorab/configuration"
)

var db *bolt.DB

func InitStorage() {
	if db != nil {
		return
	}
	dbPath, err := configuration.GetCfg().String("database")
	if err != nil {
		log.Fatalln("Can't find `database` parameter in config file ")
	}

	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalln("Can't open database for", dbPath)
	}

	initSchema()

	log.Println("Database in", dbPath, "initialized successfully")
}

func initSchema() {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("configuration"))
		tx.CreateBucketIfNotExists([]byte("dictionary"))
		return nil
	})
}

func GetDb() *bolt.DB {
	if db == nil {
		log.Fatalln("Database requested, but not initialized!")
	}
	return db
}
