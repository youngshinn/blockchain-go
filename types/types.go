package types

const (
	CreateWallet  = "CreateWallet"
	TransferCoin  = "TransferCoin"
	MintCoin      = "MintCoin"
	ChangeWallet  = "ChangeWallet"
	ConnectWallet = "ConnectWallet"
)

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Balance    string `json:"balance"`
	Time       uint64 `json:"time"`
}

type Block struct {
	Hash         string         `json:"hash"`
	PrevHash     string         `json:"prevHash"`
	Nonce        int64          `json:"nonce"`
	Height       int64          `json:"height"`
	Transactions []*Transaction `json:"transactions"`
	Time         int64          `json:"time"`
}

type Transaction struct {
	Block   int64  `json:"block"`
	Time    int64  `json:"time"`
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
	Message string `json:"message"`
	Tx      string `json:"tx"`
}
