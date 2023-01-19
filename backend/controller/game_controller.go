package controller

import (
	"context"
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// 게임 생성 함수
func (p *Controller) CreateGame(c *gin.Context) {
	game := model.Game{}
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.createGame(test_pk, game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

func (p *Controller) createGame(senderPkHexStr string, game model.Game) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPkHexStr)
	if err != nil {
		return "", err
	}
	chainID, err := p.client.NetworkID(context.Background())
	txOp, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	tx, err := p.contract.CreateGame(txOp, game.Title, game.Description, game.Home, game.Away,
		game.HomeOdd, game.AwayOdd, game.MaxRewardAmount, game.BetEndDate, game.VoteEndDate)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// 게임 리스트 반환 함수
func (p *Controller) GetGames(c *gin.Context) {
	filter := c.Query("filter")
	ret, err := p.md.GameModel.GetList(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}
