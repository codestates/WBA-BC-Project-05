package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client     *mongo.Client
	GameModel  *gameModel
	BetModel   *betModel
	VoteModel  *voteModel
	UserModel  *userModel
	EventModel *eventModel
}

func NewModel(mgUrl string) (*Model, error) {
	r := &Model{}

	var err error
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("totoro")
		r.GameModel = NewGameModel(db.Collection("game"))
		r.BetModel = NewBetModel(db.Collection("bet"))
		r.VoteModel = NewVoteModel(db.Collection("vote"))
		r.UserModel = NewUserModel(db.Collection("user"))
		r.EventModel = NewEventModel(db.Collection("event"))
	}

	return r, nil
}
