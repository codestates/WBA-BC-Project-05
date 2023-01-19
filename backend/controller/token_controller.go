package controller

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func privateKeyToPublicKey(pk string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return common.Address{}, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("fail convert, publickey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, nil
}

func (p *Controller) GetBalance(c *gin.Context) {
	publicAddress, err := privateKeyToPublicKey(test_pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	balance, err := p.contract.BalanceOf(&bind.CallOpts{}, publicAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: balance})
	return
}

// 가입 축하금 요청 함수
func (p *Controller) Welcome(c *gin.Context) {
	txOp, err := p.txPrepare(test_pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	tx, err := p.contract.WelcomeToken(txOp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: tx.Hash().Hex()})
	return
}
