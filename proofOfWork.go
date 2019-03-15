package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)




type ProofOfWork struct {
	//- 区块数据
	block *Block

	//- 难度值
	target *big.Int
	//setString(), string转数字,
	//setBytes(), bytes转数字
	//Cmp(), 同类型的比较
}

//提供一个创建ProofOfWork的方法
//数据由外部传入
//难度值由系统提供
func NewProofOfWork(b *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: b,
	}

	//难度值是多少？？
	//难度值不是固定的，是由一个数字推导而出，
	//我们这里先写成固定的，后面再推导
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"

	//创建一个临时辅助变量, 存储我们的难度值
	var bigIntTmp big.Int
	bigIntTmp.SetString(targetStr, 16) //把字符串转成大整形 16进制

	pow.target = &bigIntTmp

	return &pow
}

//提供一个计算随机数的方法, 返回随机值nonce和当前区块的哈希值
func (pow *ProofOfWork) Run() (uint64, []byte) {

	fmt.Printf("计算哈希值....\n")

	//随机数
	var nonce uint64

	//当前哈希值
	var currHash [32]byte

	//定义一个临时的bigIntTmp，用于接收当前的哈希值, 从而与目标值进行比较
	var currHashBigInt big.Int

	for {
		fmt.Printf("%x\r", currHash[:])

		//1. 获取区块数据 与 nonce进行拼接
		//把区块数据转换为字节流
		//info := pow.block + nonce
		info := pow.prepareData(nonce)

		//2. 对数据做sha256运算
		currHash /*[32]byte*/ = sha256.Sum256(info)

		currHashBigInt.SetBytes(currHash[:])   //将hash 字节切片 转成一个big Int 类型

		/*
		//3. 与难度值进行比较, 两个big.int如何比较大小？？
		//if 当前哈希 < 难度值 {
		//if 当前哈希 < pow.target {

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		//func (x *Int) Cmp(y *Int) (r int) {
		*/

		//-1 表述当前的哈希值小于目标难度值
		if currHashBigInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功! hash : %x, nonce : %d\n", currHash[:], nonce)
			break
		} else {
			//随机数加1
			nonce ++
		}
	}

	return nonce, currHash[:]
}

//准备数据，将block转换为字节流，把nonce也传递进来
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	b := pow.block
	var info []byte

	//1. 拼接数据
	info = append(info, uintToByte(b.Version)...)
	info = append(info, b.PrevBlockHash...)
	info = append(info, b.MerKelRoot...)
	info = append(info, uintToByte(b.TimeStamp)...)

	/*encodeInfo := uintToByte(b.TimeStamp)
	fmt.Printf("timestamp origin value : %d\n", b.TimeStamp)  //1552617092
	fmt.Printf("timestamp encode value : %x\n", encodeInfo)   //070600fc5c8b0e84
	fmt.Printf("timestamp decode value: %d\n", ByteToUint(encodeInfo)) //1552617092*/

	info = append(info, uintToByte(b.Difficult)...)

	//一定要记得修改nonce，否则哈希值不会受到nonce的影响
	info = append(info, uintToByte(nonce)...)
	info = append(info, b.Data...)

	return info
}
//校验即对求出来的哈希和随机数进行验证，只需要对求出来的值进行反向的计算比较即可
//- 提供一个校验函数
//IsValid()

func (pow *ProofOfWork)IsValid()  bool{
	hash := sha256.Sum256(pow.prepareData(pow.block.Nonce))
	fmt.Printf("is valid hash : %x, %d\n", hash[:], pow.block.Nonce)

	tTmp := big.Int{}
	tTmp.SetBytes(hash[:])
	if tTmp.Cmp(pow.target)  == -1 {
		return true
	}

	return false

	//return tTmp.Cmp(&pow.target)  == -1
}