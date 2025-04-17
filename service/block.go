package service

import (
	"block-test/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hacpy/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

func (s *Service) CreateBlock(from, to, value string) *types.Block {
	var block *types.Block
	toBalance := "0"

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis block will be created")
			genesisMessage := "This is the genesis block"

			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx := createTransaction(genesisMessage, (common.Address{}).String(), pk, to, value, 1)
				block = createBlockInner([]*types.Transaction{tx}, "", 1)
			}
		}
	} else {
		var tx *types.Transaction

		if common.HexToAddress(from) == (common.Address{}) {

			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx = createTransaction("MintCoin", (common.Address{}).String(), pk, to, value, 1)
				toBalance = value

			}
		} else {
			if wallet, err := s.repository.GetWalletByPublickey(from); err != nil {
				panic(err)
			} else if toWallet, err := s.repository.GetWalletByPublickey(to); err != nil {
				if err == mongo.ErrNoDocuments {
					s.log.Debug("can`t Find to wallet", "to", to)
				} else {
					panic(err)
				}
				return nil
			} else {
				//from 주오세어 balance 적합하게 있는지 확인
				// value도 업데이트

				fromDecimalBalance, _ := decimal.NewFromString(wallet.Balance)
				valueDecimal, _ := decimal.NewFromString(value)
				toDecimalBalance, _ := decimal.NewFromString(toWallet.Balance)

				if fromDecimalBalance.Cmp(valueDecimal) == -1 {
					s.log.Debug("failed to transfer coin by From Balance", "from", from, "balance", wallet.Balance, "value", value)
					return nil
				} else {

					toDecimalBalance = toDecimalBalance.Add(valueDecimal)
					toBalance = toDecimalBalance.String()

					fromDecimalBalance = fromDecimalBalance.Sub(valueDecimal)
					toBalance = value
				}

				tx = createTransaction("TransferCoin", from, wallet.PrivateKey, to, value, 1)
			}
		}

		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)
	}
	pow := s.NewPow(block)
	block.Nonce, block.Hash = pow.RunMining()

	if err := s.repository.UpsertWalletsWhenTransfer(from, to, value, toBalance); err != nil {
		panic(err)
	} else if err := s.repository.SaveBlock(block); err != nil {
		panic(err)
	}
	return block
}

func createBlockInner(txs []*types.Transaction, prevHash string, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         "",
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, pk, to, amount string, block int64) *types.Transaction {
	data := struct {
		Message string `bson:"message"`
		From    string `bson:"from"`
		To      string `bson:"to"`
		Amount  string `bson:"amount"`
	}{
		Message: message,
		From:    from,
		To:      to,
		Amount:  amount,
	}

	dataToSign := fmt.Sprintf("%x\n", data)

	pk = strings.TrimPrefix(pk, "0x")

	if ecdsaPrivateKey, err := crypto.HexToECDSA(pk); err != nil {
		panic(err)
	} else if r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, []byte(dataToSign)); err != nil {
		panic(err)
	} else {
		signature := append(r.Bytes(), s.Bytes()...)
		return &types.Transaction{
			Block:   block,
			Time:    time.Now().Unix(),
			From:    from,
			To:      to,
			Amount:  amount,
			Message: message,
			Tx:      hex.EncodeToString(signature),
		}
	}
}

func HashTransaction(b *types.Block) []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		var encoded bytes.Buffer

		enc := gob.NewEncoder(&encoded)

		if err := enc.Encode(tx); err != nil {
			panic(err)
		} else {
			txHashes = append(txHashes, encoded.Bytes())
		}
	}

	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data

}
