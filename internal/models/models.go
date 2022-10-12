package models

type Block struct {
	Difficulty       string         `json:"difficulty"`
	ExtraData        string         `json:"extraData"`
	GasLimit         string         `json:"gasLimit"`
	GasUsed          string         `json:"gasUsed"`
	Hash             string         `json:"hash"`
	LogsBloom        string         `json:"logsBloom"`
	Miner            string         `json:"miner"`
	MixHash          string         `json:"mixHash"`
	Nonce            string         `json:"nonce"`
	Number           string         `json:"number"`
	ParentHash       string         `json:"parentHash"`
	ReceiptsRoot     string         `json:"receiptsRoot"`
	Sha3Uncles       string         `json:"sha3Uncles"`
	Size             string         `json:"size"`
	StateRoot        string         `json:"stateRoot"`
	Timestamp        string         `json:"timestamp"`
	TotalDifficulty  string         `json:"totalDifficulty"`
	Transactions     []*Transaction `json:"transactions"`
	TransactionsRoot string         `json:"transactionsRoot"`
	Uncles           []interface{}  `json:"uncles"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	PublicKey        string `json:"publicKey"`
	R                string `json:"r"`
	Raw              string `json:"raw"`
	S                string `json:"s"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	V                string `json:"v"`
	Value            string `json:"value"`
}
