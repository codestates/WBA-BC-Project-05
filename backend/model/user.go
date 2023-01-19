package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserID     string `json:"user_id" bson:"user_id"`
	Pw         string `json:"pw" bson:"pw"`
	Wallet     string `json:"wallet" bson:"wallet"`
	IsManager  string `json:"is_manager" bson:"is_manager"`
	PrivateKey string `json:"privateKey" bson:"privateKey"`
}

type userModel struct {
	col *mongo.Collection
}

func NewUserModel(col *mongo.Collection) *userModel {
	m := new(userModel)
	m.col = col
	return m
}

func (p *userModel) SigninModel(id, pw string) error {
	opts := []*options.FindOneOptions{}
	var filter bson.M
	if id == "user_id" {
		filter = bson.M{"pw": pw}
	}

	var user User
	if err := p.col.FindOne(context.TODO(), filter, opts...).Decode(&user); err != nil {
		return err
	} else {
		return nil
	}
}

func (p *userModel) SignUpModel(user User) error {
	if _, err := p.col.InsertOne(context.TODO(), user); err != nil {
		fmt.Println("fail insert new user")
		return fmt.Errorf("fail, insert new user")
	}
	return nil
}
