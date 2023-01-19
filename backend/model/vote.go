package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type VoteWrapper struct {
	Vote Vote `json:"vote" bson:"vote" uri:"vote"`
}

type Vote struct {
	VoteId *big.Int       `json:"voteId" bson:"voteId" uri:"voteId"`
	GameId *big.Int       `json:"gameId" bson:"gameId" uri:"gameId" binding:"required"`
	Voter  common.Address `json:"voter" bson:"voter" uri:"voter"`
	Target uint8          `json:"target" bson:"target" uri:"target"`
}

type VoteForDB struct {
	VoteId string `json:"voteId" bson:"voteId" uri:"voteId"`
	GameId string `json:"gameId" bson:"gameId" uri:"gameId" binding:"required"`
	Voter  string `json:"voter" bson:"voter" uri:"voter"`
	Target uint8  `json:"target" bson:"target" uri:"target"`
}

type voteModel struct {
	col *mongo.Collection
}

func NewVoteModel(col *mongo.Collection) *voteModel {
	m := new(voteModel)
	m.col = col
	return m
}

func (p *voteModel) ConvertToDB(vote Vote) (VoteForDB, error) {
	voteForDB := VoteForDB{}
	voteForDB.VoteId = vote.VoteId.String()
	voteForDB.GameId = vote.GameId.String()
	voteForDB.Voter = vote.Voter.String()
	voteForDB.Target = vote.Target
	return voteForDB, nil
}

func (p *voteModel) Insert(vote VoteForDB) error {
	_, err := p.col.InsertOne(context.TODO(), vote)
	if err != nil {
		return err
	}
	return nil
}
