package controller

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"wba-bc-project-05/backend/model"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
)

// 토큰 심볼 반환 함수
func (p *Controller) GetTokenSymbol(c *gin.Context) {
	symbol, err := p.contract.Symbol(&bind.CallOpts{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: symbol})
	return
}

// 특정 주소의 보유 토큰 수 반환 함수
func (p *Controller) GetTokenBalance(c *gin.Context) {
	addressHex := c.Param("address")
	if addressHex == "" {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: "Invalid address"})
		return
	}
	address := common.HexToAddress(addressHex)
	bal, err := p.contract.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: bal})
	return
}

// 토큰 전송 함수
func (p *Controller) transferToken(senderPkHexStr string, toAddressHexStr string, value int64) (string, error) {
	privateKey, err := crypto.HexToECDSA(senderPkHexStr)
	if err != nil {
		return "", err
	}
	// privatekey로부터 publickey를 거쳐 자신의 address 변환
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("fail convert, publickey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 현재 계정의 nonce를 가져옴. 다음 트랜잭션에서 사용할 nonce
	nonce, err := p.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// gasLimit, gasPrice 설정. 추천되는 gasPrice를 가져옴
	gasPrice, err := p.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// 보낼 주소
	toAddress := common.HexToAddress(toAddressHexStr)

	// 컨트랙트 전송시 사용할 함수명
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(big.NewInt(value).Bytes(), 32)
	zvalue := big.NewInt(0)

	//컨트랙트 전송 정보 입력
	var pdata []byte
	pdata = append(pdata, methodID...)
	pdata = append(pdata, paddedAddress...)
	pdata = append(pdata, paddedAmount...)

	gasLimit := uint64(200000)

	// 트랜잭션 생성
	tx := types.NewTransaction(nonce, p.tokenAddress, zvalue, gasLimit, gasPrice, pdata)
	chainID, err := p.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// 트랜잭션 서명
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}
	// 트랜잭션 전송
	err = p.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// 컨트랙트 소유자 주소로부터 특정 주소로 토큰을 전송하는 함수
func (p *Controller) TransferToken(c *gin.Context) {
	transferReq := model.TransferReq{}
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.transferToken(p.pk, transferReq.To, transferReq.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}

// 전송자의 개인키로부터 특정 주소로 토큰을 전송하는 함수
func (p *Controller) TransferTokenFrom(c *gin.Context) {
	transferFromReq := model.TransferFromReq{}
	if err := c.ShouldBindJSON(&transferFromReq); err != nil {
		c.JSON(http.StatusBadRequest, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	ret, err := p.transferToken(transferFromReq.Pk, transferFromReq.To, transferFromReq.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultJSON{Message: "error", Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultJSON{Message: "success", Data: ret})
	return
}
