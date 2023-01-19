package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventLog struct {
	Address     string `bson:"address"`
	BlockHash   string `bson:"blackHash"`
	BlockNumber uint64 `bson:"blackNumber"`
	TxHash      string `bson:"txHash"`
	Signature   string `bson:"signature"`
}

type TransferEvent struct {
	From  common.Address `json:"from" bson:"from" uri:"from" binding:"required"`
	To    common.Address `json:"to" bson:"to" uri:"to"`
	Value *big.Int       `json:"value" bson:"value" uri:"value"`
}

type TransferEventForDB struct {
	From  string `json:"from" bson:"from" uri:"from" binding:"required"`
	To    string `json:"to" bson:"to" uri:"to"`
	Value string `json:"value" bson:"value" uri:"value"`
}

type eventModel struct {
	eventCol    *mongo.Collection
	transferCol *mongo.Collection
}

func NewEventModel(eventCol, transferCol *mongo.Collection) *eventModel {
	m := new(eventModel)
	m.eventCol = eventCol
	m.transferCol = transferCol
	return m
}

func (p *eventModel) InsertEventLog(event EventLog) error {
	_, err := p.eventCol.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}
	return nil
}

func (p *eventModel) InsertTransferEvent(te TransferEventForDB) error {
	_, err := p.transferCol.InsertOne(context.TODO(), te)
	if err != nil {
		return err
	}
	return nil
}
