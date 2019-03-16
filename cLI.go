package main

import (
	"fmt"
	"os"
)

type CLI struct {
	//bc *BlockChain
}

const Usage = `
	createBlockChain --address ADDRESS "create a blockchian" //<--这里
	addBlock --data DATA "add a block"
	printChain "print block Chain"
`

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "createBlockChain":
		if len(os.Args) > 3 && os.Args[2] == "--address" {
			address := os.Args[3]
			if address == "" {
				fmt.Println("address should not be empty!")
				os.Exit(1)
			}
			cli.CreateBlockChain(address)
		} else {
			fmt.Println(Usage)
		}
	case "addBlock":
		if len(os.Args) > 3 && os.Args[2] == "--data" {
			data := os.Args[3]
			if data == "" {
				fmt.Println("data should not be empty!")
				os.Exit(1)
			}
			cli.addBlock(data)
		} else {
			fmt.Println(Usage)
		}
	case "printChain":
		cli.printChain()
	default:
		fmt.Println(Usage)
	}
}

func (cli *CLI)addBlock(data string)  {
	//cli.bc.AddBlock(data)
	bc := GetBlockChainObj()
	bc.AddBlock(data)
	bc.db.Close()
	fmt.Println("add block successfully: ", data)
}

func (cli *CLI)printChain()  {

	//定义迭代器
	//it := NewBlockChainIterator(cli.bc)
	bc := GetBlockChainObj()
	it := NewBlockChainIterator(bc)
	for {

		block := it.GetBlockAndMoveLeft()

		fmt.Println(" ============== =============")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKelRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficuty : %d\n", block.Difficult)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)
		pow := NewProofOfWork(&block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())


		if len(block.PrevBlockHash)  == 0 {
			fmt.Println("print over!")
			break
		}
	}
}

func (cli *CLI)CreateBlockChain(address string) {
	bc := NewBlockChain(address)
	err := bc.db.Close()

	if err != nil {
		//删除数据库
		if dbExists() {
			os.Remove(blockChainDb)
		}
		fmt.Println("create blockchain failed!")
	}

	fmt.Println("create blockchain successfully!")
}