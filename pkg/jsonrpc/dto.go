package jsonrpc

import "encoding/json"

type rpcResp struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Block struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []Transaction `json:"transactions"`
	Uncles           []string      `json:"uncles"`
}

type Transaction struct {
	Hash             string  `json:"hash"`
	Nonce            string  `json:"nonce"`
	BlockHash        string  `json:"blockHash"`
	BlockNumber      string  `json:"blockNumber"`
	TransactionIndex string  `json:"transactionIndex"`
	From             string  `json:"from"`
	To               *string `json:"to"`
	Value            string  `json:"value"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Input            string  `json:"input"`
}
