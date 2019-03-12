package main
type BlockChain struct {
	//数组实现的区块链结构
	Blocks []*Block

}

const  genesisInfo  string  = "创世区块"
//创建一个区块链
func NewBlockChain() *BlockChain {

		//创建BlockChain并且添加一个最初的区块：genesisBlock (创世块)
		genesisBlock := NewBlock(genesisInfo, []byte{})

		bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

		return &bc

}



//向区块链添加区块方法
func (bc *BlockChain) AddBlock(data string) {

	//添加新区块的时候，我们向区块链的最后一个区块后面追加，我们需要最后一个区块的哈希
	//我们可以通过数组的最后一个元素，获取最后一个区块的哈希值blocks[len -1]

	//0. 获取最后一个区块
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	//最后一个区块哈希值，是当前新区块的前哈希
	prevHash := lastBlock.Hash

	//1. 创建一个block
	newBlock := NewBlock(data, prevHash)

	//2. 将block追加到Blocks数组里面
	bc.Blocks = append(bc.Blocks, newBlock)
}
