package model

import (
	"context"
	"math/big"

	"go.mongodb.org/mongo-driver/mongo"
)

const RESULT_WIN_HOME = 0
const RESULT_WIN_AWAY = 1
const RESULT_WIN_VOID = 2

type Result struct {
	GameId *big.Int `json:"gameId" bson:"gameId" uri:"gameId"`
	Win    uint8    `json:"result" bson:"result" uri:"result"`
}

type ResultgForDB struct {
	GameId string `json:"gameId" bson:"gameId" uri:"gameId"`
	Win    uint8  `json:"result" bson:"result" uri:"result"`
}

type resultModel struct {
	col *mongo.Collection
}

func NewResultModel(col *mongo.Collection) *resultModel {
	m := new(resultModel)
	m.col = col
	return m
}

func (p *resultModel) ConvertForDB(result Result) (ResultgForDB, error) {
	resultForDB := ResultgForDB{}
	resultForDB.GameId = result.GameId.String()
	resultForDB.Win = result.Win
	return resultForDB, nil
}

func (p *resultModel) Insert(result ResultgForDB) error {
	_, err := p.col.InsertOne(context.TODO(), result)
	if err != nil {
		return err
	}
	return nil
}
