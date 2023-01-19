package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type BetWrapper struct {
	Bet Bet `json:"bet" bson:"bet" uri:"bet"`
}

type Bet struct {
	BetId  *big.Int       `json:"betId" bson:"betId" uri:"betId"`
	GameId *big.Int       `json:"gameId" bson:"gameId" uri:"gameId" binding:"required"`
	Amount *big.Int       `json:"amount" bson:"amount" uri:"amount" binding:"required"`
	Bettor common.Address `json:"bettor" bson:"bettor" uri:"bettor"`
	Target uint8          `json:"target" bson:"target" uri:"target"`
	Hit    bool           `json:"hit" bson:"hit" uri:"hit"`
}

type BetForDB struct {
	BetId  string `json:"betId" bson:"betId" uri:"betId"`
	GameId string `json:"gameId" bson:"gameId" uri:"gameId" binding:"required"`
	Amount string `json:"amount" bson:"amount" uri:"amount" binding:"required"`
	Bettor string `json:"bettor" bson:"bettor" uri:"bettor"`
	Target uint8  `json:"target" bson:"target" uri:"target"`
	Hit    bool   `json:"hit" bson:"hit" uri:"hit"`
}

type betModel struct {
	col *mongo.Collection
}

func NewBetModel(col *mongo.Collection) *betModel {
	m := new(betModel)
	m.col = col
	return m
}

func (p *betModel) ConvertToDB(bet Bet) (BetForDB, error) {
	betForDB := BetForDB{}
	betForDB.BetId = bet.BetId.String()
	betForDB.GameId = bet.GameId.String()
	betForDB.Amount = bet.Amount.String()
	betForDB.Bettor = bet.Bettor.String()
	betForDB.Target = bet.Target
	return betForDB, nil
}

func (p *betModel) Insert(bet BetForDB) error {
	_, err := p.col.InsertOne(context.TODO(), bet)
	if err != nil {
		return err
	}
	return nil
}
