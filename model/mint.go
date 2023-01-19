package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type MintWrapper struct {
	Mint Mint `json:"mint" bson:"mint" uri:"mint"`
}

type Mint struct {
	From    common.Address `json:"from" bson:"from" uri:"from"`
	To      common.Address `json:"to" bson:"to" uri:"to"`
	TokenID *big.Int       `json:"tokenId" bson:"tokenId" uri:"tokenId"`
}

type MintForDB struct {
	From    string `json:"from" bson:"from" uri:"from"`
	To      string `json:"to" bson:"to" uri:"to"`
	TokenID string `json:"tokenId" bson:"tokenId" uri:"tokenId"`
}

type mintModel struct {
	col *mongo.Collection
}

func NewMintModel(col *mongo.Collection) *mintModel {
	m := new(mintModel)
	m.col = col
	return m
}

func (p *mintModel) ConvertToDB(mint Mint) (MintForDB, error) {
	mintForDB := MintForDB{}
	mintForDB.From = mint.From.String()
	mintForDB.To = mint.To.String()
	mintForDB.TokenID = mint.TokenID.String()
	return mintForDB, nil
}

func (p *mintModel) Insert(mint MintForDB) error {
	_, err := p.col.InsertOne(context.TODO(), mint)
	if err != nil {
		return err
	}
	return nil
}
