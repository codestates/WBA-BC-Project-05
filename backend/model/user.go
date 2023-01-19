package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserID     string `json:"user_id" bson:"user_id"`
	Pw         string `json:"pw" bson:"pw"`
	Wallet     string `json:"wallet" bson:"wallet"`
	IsManager  string `json:"is_manager" bson:"is_manager"`
	PrivateKey string `json:"privateKey" bson:"privateKey"`
}

func (p *Model) SigninModel(id, pw string) error {
	opts := []*options.FindOneOptions{}
	var filter bson.M
	if id == "user_id" {
		filter = bson.M{"pw": pw}
	}

	var user User
	if err := p.colUser.FindOne(context.TODO(), filter, opts...).Decode(&user); err != nil {
		return err
	} else {
		return nil
	}
}

func (p *Model) SignUpModel(user User) error {
	if _, err := p.colUser.InsertOne(context.TODO(), user); err != nil {
		fmt.Println("fail insert new user")
		return fmt.Errorf("fail, insert new user")
	}
	return nil
}
