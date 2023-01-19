package controller

import (
	"context"
	"wba-bc-project-05/backend/model"
	conf "wba-bc-project-05/config"
	"wba-bc-project-05/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var test_pk string

type ResultJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Controller struct {
	client       *ethclient.Client
	tokenAddress common.Address
	contract     *contracts.Contracts
	md           *model.Model
	pk           string
	ownerAddress string
}

func NewCTL(cf *conf.Config, md *model.Model) (*Controller, error) {
	var err error
	r := new(Controller)
	// 블록체인 네트워크와 연결할 클라이언트를 생성하기 위한 rpc url 연결
	r.client, err = ethclient.Dial(cf.Blockchain.UrlHttp)
	if err != nil {
		return r, err
	}
	// 컨트랙트 주소 및 객체 얻기
	r.tokenAddress = common.HexToAddress(cf.Blockchain.ContractAddr)
	r.contract, err = contracts.NewContracts(r.tokenAddress, r.client)
	if err != nil {
		return r, err
	}
	r.md = md
	// 개인키 및 컨트랙트 소유자 주소 저장
	r.pk = cf.Blockchain.OwnerPK
	r.ownerAddress = cf.Blockchain.OwnerAddr
	test_pk = cf.Blockchain.TestPK

	return r, nil
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
