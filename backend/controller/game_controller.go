package controller

import (
	"context"
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

const test_pk = "b4aa07e8d878757f1ac53ca9bc03ea1763b043f78adb8da3018b9f1eeb13673f"

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

	return
}

// 홈팀 베팅 함수
func (p *Controller) BetHome(c *gin.Context) {
	bet := model.Bet{}
	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	txOp, err := p.txPrepare(test_pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.betHome(txOp, bet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

func (p *Controller) betHome(txOp *bind.TransactOpts, bet model.Bet) (string, error) {
	tx, err := p.contract.BetHome(txOp, bet.GameId, bet.Amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// 원정팀 베팅 함수
func (p *Controller) BetAway(c *gin.Context) {
	bet := model.Bet{}
	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	txOp, err := p.txPrepare(test_pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "bet prepare error", Data: err.Error()})
		return
	}
	ret, err := p.betAway(txOp, bet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "betAway error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

func (p *Controller) betAway(txOp *bind.TransactOpts, bet model.Bet) (string, error) {
	tx, err := p.contract.BetAway(txOp, bet.GameId, bet.Amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// 베팅 리스트 반환 함수
func (p *Controller) GetBets(c *gin.Context) {

	return
}

// 홈팀 투표 함수
func (p *Controller) VoteHome(c *gin.Context) {
	vote := model.Vote{}
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	txOp, err := p.txPrepare(test_pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.voteHome(txOp, vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

func (p *Controller) voteHome(txOp *bind.TransactOpts, vote model.Vote) (string, error) {
	tx, err := p.contract.VoteHome(txOp, vote.GameId)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// 트랜잭션 옵션 준비 함수
func (p *Controller) txPrepare(senderPkHexStr string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(senderPkHexStr)
	if err != nil {
		return nil, err
	}
	chainID, err := p.client.NetworkID(context.Background())
	txOp, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	return txOp, nil
}
