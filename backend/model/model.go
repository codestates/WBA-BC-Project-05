package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client   *mongo.Client
	colUser *mongo.Collection
}

type TransferReq struct {
	To    string `json:"to" bson:"to" uri:"to" binding:"required"`
	Value int64  `json:"value" bson:"value" uri:"value" binding:"required"`
}

type TransferFromReq struct {
	To    string `json:"to" bson:"to" uri:"to" binding:"required"`
	Value int64  `json:"value" bson:"value" uri:"value" binding:"required"`
	Pk    string `json:"pk" bson:"pk" uri:"pk" binding:"required"`
}

type User struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Pw        string    `json:"pw" bson:"pw"`
	Wallet    string    `json:"wallet" bson:"wallet"`
	IsManager string    `json:"is_manager" bson:"is_manager"`
}

func NewModel(mgUrl string) (*Model, error) {
	r := &Model{}

	var err error
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("backend")
		r.colUser = db.Collection("user")
	}

	return r, nil
}

func (p *Model) SigninModel(id, pw string) (User, error) {
	opts := []*options.FindOneOptions{}
	var filter bson.M
	if id == "user_id" {
		filter = bson.M{"user_id": id}
	}

	var user User
	if err := p.colUser.FindOne(context.TODO(), filter, opts...).Decode(&user); err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func (p *Model) SignUpModel(user User) error {
	if _, err := p.colUser.InsertOne(context.TODO(), user); err != nil {
		fmt.Println("fail insert new user")
		return fmt.Errorf("fail, insert new user")
	}
	return nil
}