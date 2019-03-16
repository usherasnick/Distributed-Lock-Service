package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct {
	//大写
	Name string
	Age int
}

func test2main()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	lily := Person{ "Lily", 28}

	err := encoder.Encode(&lily)    //.Encode（）会把lily加密 然后 放到buffer里面  。  如果手动填写的话，要用buffer.Write(data)
	if err != nil {
		fmt.Println("encode failed!", err)
	}

	fmt.Println("after serialize :", buffer)  //相当于把它的成员属性都打印出来  显示{[buf],off,lastRead]}
	fmt.Println("after serialize :", buffer.Bytes())  //[37 255 129 3 1 1 6 80 1]
	var LILY Person

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&LILY)
	if err != nil {
		fmt.Println("decode failed!", err)
	}

	fmt.Println(LILY)
}