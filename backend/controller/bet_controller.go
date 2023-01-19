package controller

import (
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/gin-gonic/gin"
)

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
	tx, err := p.contract.BetHome(txOp, bet.GameId, bet.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
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
	tx, err := p.contract.BetAway(txOp, bet.GameId, bet.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "betAway error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
}

// 베팅 리스트 반환 함수
func (p *Controller) GetBets(c *gin.Context) {

	return
}
