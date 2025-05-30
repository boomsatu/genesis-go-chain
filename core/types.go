
package core

import (
	"math/big"
	"time"

	"blockchain-node/crypto"
)

// Block represents a block in the blockchain
type Block struct {
	Header       *BlockHeader   `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	Hash         crypto.Hash    `json:"hash"`
}

// BlockHeader represents the header of a block
type BlockHeader struct {
	PreviousHash     crypto.Hash    `json:"previousHash"`
	StateRoot        crypto.Hash    `json:"stateRoot"`
	TransactionsRoot crypto.Hash    `json:"transactionsRoot"`
	ReceiptsRoot     crypto.Hash    `json:"receiptsRoot"`
	LogsBloom        [256]byte      `json:"logsBloom"`
	Number           *big.Int       `json:"number"`
	GasLimit         uint64         `json:"gasLimit"`
	GasUsed          uint64         `json:"gasUsed"`
	Timestamp        uint64         `json:"timestamp"`
	Nonce            uint64         `json:"nonce"`
	Difficulty       *big.Int       `json:"difficulty"`
	Coinbase         crypto.Address `json:"coinbase"`
	ExtraData        []byte         `json:"extraData"`
}

// Transaction represents a transaction
type Transaction struct {
	Nonce    uint64          `json:"nonce"`
	GasPrice *big.Int        `json:"gasPrice"`
	GasLimit uint64          `json:"gasLimit"`
	To       *crypto.Address `json:"to"` // nil means contract creation
	Value    *big.Int        `json:"value"`
	Data     []byte          `json:"data"`
	V        *big.Int        `json:"v"`
	R        *big.Int        `json:"r"`
	S        *big.Int        `json:"s"`
	Hash     crypto.Hash     `json:"hash"`
	From     crypto.Address  `json:"from"`
}

// TransactionReceipt represents the receipt of a transaction
type TransactionReceipt struct {
	TransactionHash   crypto.Hash     `json:"transactionHash"`
	TransactionIndex  uint64          `json:"transactionIndex"`
	BlockHash         crypto.Hash     `json:"blockHash"`
	BlockNumber       *big.Int        `json:"blockNumber"`
	From              crypto.Address  `json:"from"`
	To                *crypto.Address `json:"to"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   *crypto.Address `json:"contractAddress"`
	Logs              []*Log          `json:"logs"`
	Status            uint64          `json:"status"` // 0 = failure, 1 = success
}

// Log represents an event log
type Log struct {
	Address     crypto.Address `json:"address"`
	Topics      []crypto.Hash  `json:"topics"`
	Data        []byte         `json:"data"`
	BlockNumber uint64         `json:"blockNumber"`
	TxHash      crypto.Hash    `json:"transactionHash"`
	TxIndex     uint           `json:"transactionIndex"`
	BlockHash   crypto.Hash    `json:"blockHash"`
	Index       uint           `json:"logIndex"`
	Removed     bool           `json:"removed"`
}

// Account represents an account in the world state
type Account struct {
	Nonce       uint64      `json:"nonce"`
	Balance     *big.Int    `json:"balance"`
	CodeHash    crypto.Hash `json:"codeHash"`
	StorageRoot crypto.Hash `json:"storageRoot"`
}

// Genesis represents the genesis block configuration
type Genesis struct {
	Config      *ChainConfig                    `json:"config"`
	Nonce       uint64                          `json:"nonce"`
	Timestamp   uint64                          `json:"timestamp"`
	ExtraData   []byte                          `json:"extraData"`
	GasLimit    uint64                          `json:"gasLimit"`
	Difficulty  *big.Int                        `json:"difficulty"`
	Coinbase    crypto.Address                  `json:"coinbase"`
	Alloc       map[crypto.Address]Account      `json:"alloc"`
}

// ChainConfig represents the chain configuration
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"`
}

// NewBlock creates a new block
func NewBlock(header *BlockHeader, txs []*Transaction) *Block {
	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	block.Hash = block.CalculateHash()
	return block
}

// CalculateHash calculates the hash of the block
func (b *Block) CalculateHash() crypto.Hash {
	// Serialize header and calculate hash
	data := b.Header.Serialize()
	return crypto.Keccak256Hash(data)
}

// Serialize serializes the block header
func (h *BlockHeader) Serialize() []byte {
	// Simple serialization - in production, use more robust encoding
	data := append(h.PreviousHash.Bytes(), h.StateRoot.Bytes()...)
	data = append(data, h.TransactionsRoot.Bytes()...)
	data = append(data, h.Number.Bytes()...)
	data = append(data, big.NewInt(int64(h.Timestamp)).Bytes()...)
	data = append(data, big.NewInt(int64(h.Nonce)).Bytes()...)
	return data
}

// NewTransaction creates a new transaction
func NewTransaction(nonce uint64, to *crypto.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return &Transaction{
		Nonce:    nonce,
		To:       to,
		Value:    amount,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
}

// CalculateHash calculates the hash of the transaction
func (tx *Transaction) CalculateHash() crypto.Hash {
	// Simple serialization for hash calculation
	data := append(big.NewInt(int64(tx.Nonce)).Bytes(), tx.GasPrice.Bytes()...)
	data = append(data, big.NewInt(int64(tx.GasLimit)).Bytes()...)
	if tx.To != nil {
		data = append(data, tx.To.Bytes()...)
	}
	data = append(data, tx.Value.Bytes()...)
	data = append(data, tx.Data...)
	return crypto.Keccak256Hash(data)
}

// IsContractCreation returns whether the transaction is a contract creation
func (tx *Transaction) IsContractCreation() bool {
	return tx.To == nil
}

// NewGenesisBlock creates a new genesis block
func NewGenesisBlock(genesis *Genesis) *Block {
	header := &BlockHeader{
		PreviousHash: crypto.Hash{},
		Number:       big.NewInt(0),
		GasLimit:     genesis.GasLimit,
		GasUsed:      0,
		Timestamp:    genesis.Timestamp,
		Nonce:        genesis.Nonce,
		Difficulty:   genesis.Difficulty,
		Coinbase:     genesis.Coinbase,
		ExtraData:    genesis.ExtraData,
	}

	return NewBlock(header, []*Transaction{})
}

// DefaultGenesis returns the default genesis configuration
func DefaultGenesis() *Genesis {
	return &Genesis{
		Config: &ChainConfig{
			ChainID: big.NewInt(1337),
		},
		Nonce:      0,
		Timestamp:  uint64(time.Now().Unix()),
		ExtraData:  []byte("Genesis Block"),
		GasLimit:   8000000,
		Difficulty: big.NewInt(4),
		Coinbase:   crypto.Address{},
		Alloc:      make(map[crypto.Address]Account),
	}
}
