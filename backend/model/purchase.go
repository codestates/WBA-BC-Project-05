package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseWrapper struct {
	Purchase Purchase `json:"purchase" bson:"purchase" uri:"purchase"`
}

type Purchase struct {
	From    common.Address `json:"from" bson:"from" uri:"from"`
	TokenID *big.Int       `json:"tokenId" bson:"tokenId" uri:"tokenId"`
}

type PurchaseForDB struct {
	From    string `json:"from" bson:"from" uri:"from"`
	TokenID string `json:"tokenId" bson:"tokenId" uri:"tokenId"`
}

type purchaseModel struct {
	col *mongo.Collection
}

func NewPurchaseModel(col *mongo.Collection) *purchaseModel {
	m := new(purchaseModel)
	m.col = col
	return m
}

func (p *purchaseModel) ConvertToDB(purchase Purchase) (PurchaseForDB, error) {
	purchaseForDB := PurchaseForDB{}
	purchaseForDB.From = purchase.From.String()
	purchaseForDB.TokenID = purchase.TokenID.String()
	return purchaseForDB, nil
}

func (p *purchaseModel) Insert(purchase PurchaseForDB) error {
	_, err := p.col.InsertOne(context.TODO(), purchase)
	if err != nil {
		return err
	}
	return nil
}
