package main

import "fmt"

func main(){
	blockChain:=NewBlockChain()
	blockChain.AddBlock("你好")
	blockChain.AddBlock("中国")
	blockChain.AddBlock("123456")
	//打印区块链
	for _, block := range blockChain.Blocks {
		fmt.Printf("++++++++++++++++++++++++++++++\n")
		fmt.Printf("前区块哈希值 : %x\n", block.PrevBlockHash)
		fmt.Printf("当前区块哈希值 : %x\n", block.Hash)
		fmt.Printf("区块数据 : %s\n", block.Data)
	}
}