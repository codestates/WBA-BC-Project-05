package model

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NFTWrapper struct {
	NFT NFT `json:"nft" bson:"nft" uri:"nft"`
}

type NFT struct {
	NFTId           *big.Int       `json:"nftId" bson:"nftId" uri:"nftId"`
	Creator          common.Address `json:"creator" bson:"creator" uri:"creator"`
	Title            string         `json:"title" bson:"title" uri:"title" binding:"required"`
}

type NFTForDB struct {
	NFTId           string `json:"nftId" bson:"nftId" uri:"nftId"`
	Creator          string `json:"creator" bson:"creator" uri:"creator"`
	Title            string `json:"title" bson:"title" uri:"title" binding:"required"`
}

func (p *nftModel) ConvertToDB(nft NFT) (NFTForDB, error) {
	nftForDb := NFTForDB{}
	nftForDb.NFTId = nft.NFTId.String()
	nftForDb.Creator = nft.Creator.String()
	nftForDb.Title = nft.Title
	return nftForDb, nil
}

func (p *nftModel) ConvertToSol(nft NFTForDB) (NFT, error) {
	nftForSol := NFT{}
	return nftForSol, nil
}

type nftModel struct {
	col *mongo.Collection
}

func NewNFTModel(col *mongo.Collection) *nftModel {
	m := new(nftModel)
	m.col = col
	return m
}

func (p *nftModel) Insert(nft NFTForDB) error {
	_, err := p.col.InsertOne(context.TODO(), nft)
	if err != nil {
		return err
	}
	return nil
}

func (p *nftModel) GetList(param string) ([]NFTForDB, error) {
	var results []NFTForDB
	var filter primitive.M
	cur, err := p.col.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		return results, err
	}
	return results, nil
}
