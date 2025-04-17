package repository

import (
	"block-test/config"
	"context"

	"github.com/inconshreveable/log15"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client //mongodb 연결
	db     *mongo.Database
	wallet *mongo.Collection
	block  *mongo.Collection
	config *config.Config
	log    log15.Logger
	tx     *mongo.Collection
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{
		config: config,
		log:    log15.New("module", "repository"),
	}

	//db 커넥트

	var err error
	ctx := context.Background()

	mConfig := config.Mongo

	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mConfig.Uri)); err != nil {
		r.log.Error("failed to connect to mongo", mConfig.Uri)
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		r.log.Error("Ping error to mongo", mConfig.Uri)
		return nil, err
	} else {
		db := r.client.Database(config.Mongo.DB, nil)

		r.wallet = db.Collection("wallet")
		r.tx = db.Collection("tx")
		r.block = db.Collection("block")

		//TODO 컬렉션 연결

		r.log.Info("success to connect")
		return r, nil
	}
}
