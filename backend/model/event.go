package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type EventLog struct {
	Address     string `bson:"address"`
	BlockHash   string `bson:"blackHash"`
	BlockNumber uint64 `bson:"blackNumber"`
	TxHash      string `bson:"txHash"`
}

type eventModel struct {
	col *mongo.Collection
}

func NewEventModel(col *mongo.Collection) *eventModel {
	m := new(eventModel)
	m.col = col
	return m
}

func (p *eventModel) Save(event *EventLog) error {
	result, err := p.col.InsertOne(context.TODO(), *event)
	if err != nil {
		return err
	}
	fmt.Println(result.InsertedID)

	return nil
}
