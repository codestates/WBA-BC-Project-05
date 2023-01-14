package controller

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
)

// 코인 전송 함수
func (p *Controller) transferCoin(senderPkHexStr string, toAddressHexStr string, value int64) (string, error) {
	pk, err := crypto.HexToECDSA(senderPkHexStr)
	if err != nil {
		return "", err
	}
	// privatekey로부터 publickey를 거쳐 자신의 address 변환
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("fail convert, publickey")
	}
	// 보낼 address 설정
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 현재 계정의 nonce를 가져옴. 다음 트랜잭션에서 사용할 nonce
	nonce, err := p.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// gasLimit, gasPrice 설정. 추천되는 gasPrice를 가져옴
	// value := big.NewInt(70000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := p.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// 전송받을 상대방 address 설정
	toAddress := common.HexToAddress(toAddressHexStr)
	// 트랜잭션 생성
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(value), gasLimit, gasPrice, data)
	chainID, err := p.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// 트랜잭션 서명
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		return "", err
	}

	// RLP 인코딩 전 트랜잭션 묶음. 현재는 1개의 트랜잭션
	ts := types.Transactions{signedTx}
	// RLP 인코딩
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	rawTxHex := hex.EncodeToString(rawTxBytes)
	rTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		return "", err
	}

	// RLP 디코딩
	rlp.DecodeBytes(rTxBytes, &tx)
	// 트랜잭션 전송
	err = p.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// 컨트랙트 소유자 주소로부터 특정 주소로 코인을 전송하는 함수
func (p *Controller) TransferCoin(c *gin.Context) {
	transferReq := model.TransferReq{}
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.transferCoin(p.pk, transferReq.To, transferReq.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

// 전송자의 개인키로부터 특정 주소로 코인을 전송하는 함수
func (p *Controller) TransferCoinFrom(c *gin.Context) {
	transferFromReq := model.TransferFromReq{}
	if err := c.ShouldBindJSON(&transferFromReq); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.transferCoin(transferFromReq.Pk, transferFromReq.To, transferFromReq.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}
