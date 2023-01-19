package controller

import (
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/gin-gonic/gin"
)

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
	tx, err := p.contract.VoteHome(txOp, vote.GameId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
}

// 원정팀 투표 함수
func (p *Controller) VoteAway(c *gin.Context) {
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
	tx, err := p.contract.VoteAway(txOp, vote.GameId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
}

// 무효 투표 함수
func (p *Controller) VoteVoid(c *gin.Context) {
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
	tx, err := p.contract.VoteVoid(txOp, vote.GameId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
}
