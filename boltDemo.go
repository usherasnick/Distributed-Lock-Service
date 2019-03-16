package main
import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

func testmain() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//2.写数据库

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))

		var err error
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("firstBucket"))
			if err != nil {
				fmt.Println("createBucket failed!", err)
				os.Exit(1)
			}
		}

		bucket.Put([]byte("aaaa"), []byte("HelloWorld!"))
		bucket.Put([]byte("bbbb"), []byte("HelloItcast!"))
		return nil
	})
	//3.读取数据库
	var value []byte

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))
		if bucket == nil {
			fmt.Println("Bucket is nil!")
			os.Exit(1)
		}

		value = bucket.Get([]byte("aaaa"))
		fmt.Println("aaaa => ", string(value))
		value = bucket.Get([]byte("bbbb"))
		fmt.Println("bbbb => ", string(value))

		return nil
	})
}