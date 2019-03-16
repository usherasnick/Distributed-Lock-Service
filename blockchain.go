package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

//使用bolt改写
type BlockChain struct {
	//操作数据库的句柄
	db *bolt.DB

	//尾巴，存储最后一个区块的哈希
	tail []byte
}

const genesisInfo string = "创世区块"
const blockChainDb string = "my.db"
const blockBucket string = "myBucket"
const last string  =  "LastHashKey"

//- 提供一个创建BlockChain的方法

func NewBlockChain(address string) *BlockChain {  //<<---添加了创建人的地址，后面会用到

	var lastHash []byte

	/*
	1. 打开数据库(没有的话就创建)
	2. 找到抽屉（bucket），如果找到，就返回bucket，如果没有找到，我们要创建bucket，通过名字创建
		a. 找到了
			1. 通过"last"这个key找到我们最好一个区块的哈希。

		b. 没找到创建
			1. 创建bucket，通过名字
			2. 添加创世块数据
			3. 更新"last"这个key的value（创世块的哈希值）

	*/
	if dbExists() {   //<---这里
		fmt.Println("blockchain already exist!")
		os.Exit(1)
	}

	db, err := bolt.Open(blockChainDb, 0600, nil)

	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		var err error
		//如果是空的，表明这个bucket没有创建，我们就要去创建它，然后再写数据。
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				fmt.Println("createBucket failed!", err)
				os.Exit(1)
			}

			//抽屉准备好了，开始写区块数据，区块哪里来？？
			genesisBlock := NewBlock("Genesis Block!!!", []byte{})

			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize() /*block的字节流！*/) //TODO
			bucket.Put([]byte("last"), genesisBlock.Hash)

			//这个别忘了，我们需要返回它
			lastHash = genesisBlock.Hash
			return nil

			//抽屉已经存在，直接读取即可
		} else {
			//获取最后一个区块的哈希
			lastHash = bucket.Get([]byte("last"))
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

/*记得两件事，put数据，更新lastHash */
//向区块链添加区块方法
func (bc *BlockChain) AddBlock(data string) {

	/*
	//获取最后一个区块
	lastBlock := bc.blocks[len(bc.blocks) -1]
	//获取最后一个区块的哈希,作为最新（当前）区块的前哈希
	prevHash := lastBlock.Hash

	block := NewBlock(data, prevHash)
	bc.blocks = append(bc.blocks, &block)
	*/

	//获取最后区块的哈希值
	lastBlockHash := bc.tail

	//创建新区块
	newBlock := NewBlock(data, lastBlockHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		//如果是空的，表明这个bucket没有创建，我们就要去创建它，然后再写数据。
		if bucket == nil {
			fmt.Println("bucket should not be nil!!")
			os.Exit(1)
		} else {
			//添加区块
			bucket.Put(newBlock.Hash, newBlock.Serialize() /*block的字节流！*/) //TODO
			//更新最后区块的哈希值
			bucket.Put([]byte("last"), newBlock.Hash)

			//这个别忘了，我们需要返回它
			bc.tail = newBlock.Hash
			return nil
		}
		return nil
	})
}


func (bc *BlockChain) PrintChain() {

	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(blockBucket))

		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("last")) {  //-> "LastHashKey"
				return nil
			}

			block := Deserialize(v)
			//fmt.Printf("key=%x, value=%s\n", k, v)
			fmt.Printf("===========================\n\n")
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevBlockHash)
			fmt.Printf("梅克尔根: %x\n", block.MerKelRoot)
			fmt.Printf("时间戳: %d\n", block.TimeStamp)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficult)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Data)
			return nil
		})
		return nil
	})
}

//获取一个区块链实例
func GetBlockChainObj() *BlockChain {

	var lastHash []byte

	if !dbExists() {  //<---这里
		fmt.Println("blockchain not exist, pls create first!")
		os.Exit(1)
	}

	db, err := bolt.Open(blockChainDb, 0600, nil)

	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		var err error
		if bucket == nil {
			fmt.Println("bucket should not be nil!", err)
			os.Exit(1)

		}

		//抽屉已经存在，直接读取即可
		//获取最后一个区块的哈希
		lastHash = bucket.Get([]byte("last"))
		return nil
	})

	return &BlockChain{db, lastHash}
}
func dbExists() bool {
	if _, err := os.Stat(blockChainDb); os.IsNotExist(err) {
		return false
	}

	return true
}