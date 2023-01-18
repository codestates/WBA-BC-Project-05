package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameWrapper struct {
	Game Game `json:"game" bson:"game" uri:"game"`
}

type Game struct {
	GameId           *big.Int       `json:"gameId" bson:"gameId" uri:"gameId"`
	Creator          common.Address `json:"creator" bson:"creator" uri:"creator"`
	Title            string         `json:"title" bson:"title" uri:"title" binding:"required"`
	Description      string         `json:"description" bson:"description" uri:"description" binding:"required"`
	Home             string         `json:"home" bson:"home" uri:"home" binding:"required"`
	Away             string         `json:"away" bson:"away" uri:"away" binding:"required"`
	HomeOdd          *big.Int       `json:"homeOdd" bson:"homeOdd" uri:"homeOdd" binding:"required"`
	AwayOdd          *big.Int       `json:"awayOdd" bson:"awayOdd" uri:"awayOdd" binding:"required"`
	VoteHomeCount    *big.Int       `json:"voteHomeCount" bson:"voteHomeCount" uri:"voteHomeCount"`
	VoteAwayCount    *big.Int       `json:"voteAwayCount" bson:"voteAwayCount" uri:"voteAwayCount"`
	VoteVoidCount    *big.Int       `json:"voteVoidCount" bson:"voteVoidCount" uri:"voteVoidCount"`
	MaxRewardAmount  *big.Int       `json:"maxRewardAmount" bson:"maxRewardAmount" uri:"maxRewardAmount" binding:"required"`
	MaxRewardHomeAcc *big.Int       `json:"maxRewardHomeAcc" bson:"maxRewardHomeAcc" uri:"maxRewardHomeAcc"`
	MaxRewardAwayAcc *big.Int       `json:"maxRewardAwayAcc" bson:"maxRewardAwayAcc" uri:"maxRewardAwayAcc"`
	HomeAccReward    *big.Int       `json:"homeAccReward" bson:"homeAccReward" uri:"homeAccReward"`
	AwayAccReward    *big.Int       `json:"awayAccReward" bson:"awayAccReward" uri:"awayAccReward"`
	CreateDate       uint32         `json:"createDate" bson:"createDate" uri:"createDate"`
	BetEndDate       uint32         `json:"betEndDate" bson:"betEndDate" uri:"betEndDate" binding:"required"`
	VoteEndDate      uint32         `json:"voteEndDate" bson:"voteEndDate" uri:"voteEndDate" binding:"required"`
}

type GameForDB struct {
	GameId           string `json:"gameId" bson:"gameId" uri:"gameId"`
	Creator          string `json:"creator" bson:"creator" uri:"creator"`
	Title            string `json:"title" bson:"title" uri:"title" binding:"required"`
	Description      string `json:"description" bson:"description" uri:"description" binding:"required"`
	Home             string `json:"home" bson:"home" uri:"home" binding:"required"`
	Away             string `json:"away" bson:"away" uri:"away" binding:"required"`
	HomeOdd          string `json:"homeOdd" bson:"homeOdd" uri:"homeOdd" binding:"required"`
	AwayOdd          string `json:"awayOdd" bson:"awayOdd" uri:"awayOdd" binding:"required"`
	VoteHomeCount    string `json:"voteHomeCount" bson:"voteHomeCount" uri:"voteHomeCount"`
	VoteAwayCount    string `json:"voteAwayCount" bson:"voteAwayCount" uri:"voteAwayCount"`
	VoteVoidCount    string `json:"voteVoidCount" bson:"voteVoidCount" uri:"voteVoidCount"`
	MaxRewardAmount  string `json:"maxRewardAmount" bson:"maxRewardAmount" uri:"maxRewardAmount" binding:"required"`
	MaxRewardHomeAcc string `json:"maxRewardHomeAcc" bson:"maxRewardHomeAcc" uri:"maxRewardHomeAcc"`
	MaxRewardAwayAcc string `json:"maxRewardAwayAcc" bson:"maxRewardAwayAcc" uri:"maxRewardAwayAcc"`
	HomeAccReward    string `json:"homeAccReward" bson:"homeAccReward" uri:"homeAccReward"`
	AwayAccReward    string `json:"awayAccReward" bson:"awayAccReward" uri:"awayAccReward"`
	CreateDate       uint32 `json:"createDate" bson:"createDate" uri:"createDate"`
	BetEndDate       uint32 `json:"betEndDate" bson:"betEndDate" uri:"betEndDate" binding:"required"`
	VoteEndDate      uint32 `json:"voteEndDate" bson:"voteEndDate" uri:"voteEndDate" binding:"required"`
}

func (p *gameModel) ConvertToDB(game Game) (GameForDB, error) {
	gameForDb := GameForDB{}
	gameForDb.GameId = game.GameId.String()
	gameForDb.Creator = game.Creator.String()
	gameForDb.Title = game.Title
	gameForDb.Description = game.Description
	gameForDb.Home = game.Home
	gameForDb.Away = game.Away
	gameForDb.HomeOdd = game.HomeOdd.String()
	gameForDb.AwayOdd = game.AwayOdd.String()
	gameForDb.VoteHomeCount = game.VoteHomeCount.String()
	gameForDb.VoteAwayCount = game.VoteAwayCount.String()
	gameForDb.VoteVoidCount = game.VoteVoidCount.String()
	gameForDb.MaxRewardAmount = game.MaxRewardAmount.String()
	gameForDb.MaxRewardHomeAcc = game.MaxRewardHomeAcc.String()
	gameForDb.MaxRewardAwayAcc = game.MaxRewardAwayAcc.String()
	gameForDb.HomeAccReward = game.HomeAccReward.String()
	gameForDb.AwayAccReward = game.AwayAccReward.String()
	gameForDb.CreateDate = game.CreateDate
	gameForDb.BetEndDate = game.BetEndDate
	gameForDb.VoteEndDate = game.VoteEndDate
	return gameForDb, nil
}

func (p *gameModel) ConvertToSol(game GameForDB) (Game, error) {
	gameForSol := Game{}
	return gameForSol, nil
}

type gameModel struct {
	col *mongo.Collection
}

func NewGameModel(col *mongo.Collection) *gameModel {
	m := new(gameModel)
	m.col = col
	return m
}

func (p *gameModel) Insert(game GameForDB) error {
	_, err := p.col.InsertOne(context.TODO(), game)
	if err != nil {
		return err
	}
	return nil
}
