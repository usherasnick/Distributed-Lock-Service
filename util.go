package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

)
//这个文件中实现一些工具函数

//数字转换为byte

// go内置的包， gob

//gob编码和解码
type Vector struct {
	x, y, z int
}
// This example transmits a value that implements the custom encoding and decoding methods.
func Example_encodeDecode() {
	var network bytes.Buffer // Stand-in for the network.
	// Create an encoder and send a value.
	enc := gob.NewEncoder(&network)    //--  ??需要的不是一个io.Writer吗
	err := enc.Encode(Vector{3, 4, 5})
	if err != nil {
		log.Fatal("encode:", err)
	}
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&network)
	var v Vector
	err = dec.Decode(&v)
	if err != nil {
		log.Fatal("decode:", err)
	}
	fmt.Println(v)
	// Output:
	// {3 4 5}
}

func MyEncode(data uint64) []byte {
	//用来装东西的篮子 buffer
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	//将编码后的东西放篮子里
	 encoder.Encode(data)
	return buffer.Bytes()

}
func MyDecode(b bytes.Buffer){
		var buffer bytes.Buffer
	// Create a decoder and receive a value.

	dec := gob.NewDecoder(&buffer)

	var data uint64
	//将编码后的东西放到 data里
	dec.Decode(&data)
    //var i int32
	//s:=strconv.Itoa(int(i))

}

func uintToByte(num uint64) []byte {
	//编码过程
	var buffer bytes.Buffer

	//1. 创建编码器
	encoder := gob.NewEncoder(&buffer)   //设置一个装东西的 buffer篮子

	//2. 编码器执行编码
	encoder.Encode(&num)    //将num编码后 的东西  ，自动放入buffer那个篮子

	return buffer.Bytes()
}

//解码过程
func ByteToUint(src []byte) uint64 {
	//1. 创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(src))

	var value uint64

	//2. 解码器执行解码
	decoder.Decode(&value)

	return value
}
