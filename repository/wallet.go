package repository

import (
	"block-test/types"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r *Repository) CreateNewWallet(wallet *types.Wallet) error {

	ctx := context.Background()
	opt := options.Update().SetUpsert(true)
	wallet.Time = uint64(time.Now().Unix())

	filter := bson.M{"privateKey": wallet.PrivateKey}
	update := bson.M{"$set": bson.M{
		"privateKey": wallet.PrivateKey,
		"publicKey":  wallet.PublicKey,
		"time":       wallet.Time,
		"balance":    wallet.Balance,
	}}
	if _, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *Repository) GetWallet(pk string) (*types.Wallet, error) {
	ctx := context.Background()

	var wallet types.Wallet

	filter := bson.M{"privateKey": pk} // pk 값을 기준으로 가져온다

	if err := r.wallet.FindOne(ctx, filter, options.FindOne()).Decode(&wallet); err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}

func (r *Repository) GetWalletByPublickey(publicKey string) (*types.Wallet, error) {
	ctx := context.Background()

	var wallet types.Wallet

	filter := bson.M{"publicKey": publicKey} // pk 값을 기준으로 가져온다

	if err := r.wallet.FindOne(ctx, filter, options.FindOne()).Decode(&wallet); err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}

func (r *Repository) UpsertWalletsWhenTransfer(from, to, fromBalance, toBalance string) error {
	ctx := context.Background()

	opt := options.Update().SetUpsert(true)

	if from != (common.Address{}).String() {
		filter := bson.M{"publickey": from}
		update := bson.M{"$set": bson.M{
			"balance": fromBalance,
		}}
		_, err := r.wallet.UpdateOne(ctx, filter, update, opt)
		if err != nil {
			return err
		}
	}

	if to != (common.Address{}).String() {
		filter := bson.M{"publickey": to}
		update := bson.M{"$set": bson.M{
			"balance": toBalance,
		}}
		_, err := r.wallet.UpdateOne(ctx, filter, update, opt)
		if err != nil {
			return err
		}
	}
	return nil
}
