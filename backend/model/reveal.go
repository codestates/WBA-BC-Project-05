package model

import (
	"context"
	"math/big"

	"go.mongodb.org/mongo-driver/mongo"
)

type RevealWrapper struct {
	Reveal Reveal `json:"reveal" bson:"reveal" uri:"reveal"`
}

type Reveal struct {
	TokenID      *big.Int `json:"tokenId" bson:"tokenId" uri:"tokenId"`
	RandomNumber *big.Int `json:"randomNumber" bson:"randomNumber" uri:"randomNumber"`
}

type RevealForDB struct {
	TokenID      string `json:"tokenId" bson:"tokenId" uri:"tokenId"`
	RandomNumber string `json:"randomNumber" bson:"randomNumber" uri:"randomNumber"`
}

type revealModel struct {
	col *mongo.Collection
}

func NewRevealModel(col *mongo.Collection) *revealModel {
	m := new(revealModel)
	m.col = col
	return m
}

func (p *revealModel) ConvertToDB(reveal Reveal) (RevealForDB, error) {
	revealForDB := RevealForDB{}
	revealForDB.TokenID = reveal.TokenID.String()
	revealForDB.RandomNumber = reveal.RandomNumber.String()
	return revealForDB, nil
}

func (p *revealModel) Insert(reveal RevealForDB) error {
	_, err := p.col.InsertOne(context.TODO(), reveal)
	if err != nil {
		return err
	}
	return nil
}
