package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client   *mongo.Client
	colBlock *mongo.Collection
}

type Block struct {
	BlockHash    string        `bson:"blockHash"`
	BlockNumber  uint64        `bson:"blockNumber"`
	GasLimit     uint64        `bson:"gasLimit"`
	GasUsed      uint64        `bson:"gasUsed"`
	Time         uint64        `bson:"timestamp"`
	Nonce        uint64        `bson:"nonce"`
	Transactions []Transaction `bson:"transactions"`
}

type Transaction struct {
	TxHash      string `bson:"hash"`
	From        string `bson:"from"`
	To          string `bson:"to"` // 컨트랙트의 경우 nil 반환
	Nonce       uint64 `bson:"nonce"`
	GasPrice    uint64 `bson:"gasPrice"`
	GasLimit    uint64 `bson:"gasLimit"`
	Amount      uint64 `bson:"amount"`
	BlockHash   string `bson:"blockHash"`
	BlockNumber uint64 `bson:"blockNumber"`
}

type EventLog struct {
	Address     string `bson:"address"`
	BlockHash   string `bson:"blackHash"`
	BlockNumber uint64 `bson:"blackNumber"`
	TxHash      string `bson:"txHash"`
}

func NewModel(mgUrl string) (*Model, error) {
	r := &Model{}

	var err error
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("daemon")
		r.colBlock = db.Collection("block")
	}

	return r, nil
}

func (p *Model) SaveBlock(block *Block) error {
	result, err := p.colBlock.InsertOne(context.TODO(), *block)
	if err != nil {
		return err
	}
	fmt.Println(result.InsertedID)

	return nil
}

func (p *Model) SaveEvent(event *EventLog) error {
	result, err := p.colBlock.InsertOne(context.TODO(), *event)
	if err != nil {
		return err
	}
	fmt.Println(result.InsertedID)

	return nil
}
