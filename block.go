package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type Block struct {
	//- 版本号：区块链的版本
	Version uint64

	//- 父区块哈希值（重要）
	PrevBlockHash []byte

	//-. 当前区块哈希(为了方便编写), 比特币的结构中不包含这个字段
	Hash []byte

	//- 梅克尔根（忽略）一个哈希值
	MerKelRoot []byte

	//- 时间戳（区块产生的时间）, 从1970年至今的秒数
	TimeStamp uint64

	//- 难度值（挖矿的难度，是一个数字，可以推导出一个哈希值）
	Difficult uint64

	//- Nonce （挖矿随机值）
	Nonce uint64

	//-. 数据
	Data []byte
}
func NewBlock(data string ,prevBlockHash []byte) *Block{
	newBlock:=Block{
		Version:       00,
		PrevBlockHash: prevBlockHash, //前区块哈希值
		Hash:          nil,           //当前区块哈希值，可以自己计算出来
		MerKelRoot:    nil,           //v4版本修改
		TimeStamp:     uint64(time.Now().Unix()),
		Difficult:     0, //先填0，v2版本修改
		Nonce:         0, //v2版本修改

		Data: []byte(data), //外部传入
	}
	//计算自己哈希值
	//newBlock.setHash()
	//在这里调用pow相关函数
	pow:=NewProofOfWork(&newBlock)
	nonce,hash :=pow.Run()
	//将挖矿产生的数据赋值给区块
	newBlock.Nonce=nonce
	newBlock.Hash=hash



	return &newBlock
}

func (b *Block)setHash(){
	var info []byte

	//1. 拼接数据
	info = append(info, uintToByte(b.Version)...)
	info = append(info, b.PrevBlockHash...)
	info = append(info, b.MerKelRoot...)
	info = append(info, uintToByte(b.TimeStamp)...)

	//encodeInfo := uintToByte(b.TimeStamp)
	//fmt.Printf("timestamp origin value : %d\n", b.TimeStamp)
	//fmt.Printf("timestamp encode value : %x\n", encodeInfo)
	//fmt.Printf("timestamp decode value: %d\n", ByteToUint(encodeInfo))

	info = append(info, uintToByte(b.Difficult)...)
	info = append(info, uintToByte(b.Nonce)...)
	info = append(info, b.Data...)

	//2.hash运算
	hash := sha256.Sum256(info)

	//3. 将得到的哈希赋值给Hash字段
	b.Hash = hash[:]   // ？？？？[32]byte -->>[]byte

}
//序列化:将结构转化成字节流，在网络上按字节传输。
func (block *Block)Serialize() []byte {
	//将block数据转换成字节流

	var buffer bytes.Buffer
	//创建一个编码器
	encoder := gob.NewEncoder(&buffer)

	//编码，将block编码成buffer
	err := encoder.Encode(block)

	if err != nil {
		fmt.Println("encode failed!", err)
		os.Exit(1)
	}

	return buffer.Bytes()
}
// 反序列化:将接受到的字节流转换成目标结构。
func Deserialize(data []byte) Block {
	var block Block
	var buffer bytes.Buffer

	//将data写入buffer
	_, err := buffer.Write(data)
	if err != nil {
		fmt.Println("buffer.Read failed!", err)
		os.Exit(1)
	}

	//创建decoder
	decoder := gob.NewDecoder(&buffer)

	//将buffer数据转换成block
	err = decoder.Decode(&block)
	if err != nil {
		fmt.Println("decode failed!", err)
		os.Exit(1)
	}

	return block
}