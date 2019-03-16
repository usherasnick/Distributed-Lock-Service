package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

//1. 定义一个区块链迭代器的结构
type BlockChainIterator struct {
	db *bolt.DB

	//指向当前的区块
	current_point []byte
}
func NewBlockChainIterator(bc *BlockChain) BlockChainIterator {

	var it BlockChainIterator

	it.db = bc.db
	it.current_point = bc.tail

	return it
}

//一般叫Next(),迭代器的访问函数
func (it *BlockChainIterator)GetBlockAndMoveLeft() Block {
	var block Block
	//1. 获取block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		//如果是空的，报错。
		if bucket == nil {
			fmt.Println("bucket should not be nil!!")
			os.Exit(1)
		} else {

			//根据当前的current_pointer获取block
			//这是一个字节流，需要反序列化
			current_block_tmp := bucket.Get(it.current_point)
			//fmt.Println("current_block_tmp : ", current_block_tmp)
			current_block := Deserialize(current_block_tmp)

			//这就拿到了我们想要的区块数据，准备返回
			block = current_block

			//将游标（指针）左移
			//2. 向左移动
			it.current_point = current_block.PrevBlockHash
		}
		return nil
	})

	return block

}