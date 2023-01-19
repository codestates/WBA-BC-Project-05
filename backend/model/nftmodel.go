package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NftModel struct {
	client        *mongo.Client
	MintModel     *mintModel
	PurchaseModel *purchaseModel
	RevealModel   *revealModel
	colUser       *mongo.Collection
}

func NewNftModel(mgUrl string) (*NftModel, error) {
	r := &NftModel{}

	var err error
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("totoro")
		r.MintModel = NewMintModel(db.Collection("mint"))
		r.PurchaseModel = NewPurchaseModel(db.Collection("purchase"))
		r.RevealModel = NewRevealModel(db.Collection("reveal"))
		r.colUser = db.Collection("user")
	}

	return r, nil
}
