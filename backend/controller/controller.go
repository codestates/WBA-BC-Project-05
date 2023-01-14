package controller

import (
	conf "wba-bc-project-05/backend/config"
	"wba-bc-project-05/backend/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ResultJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Controller struct {
	client       *ethclient.Client
	tokenAddress common.Address
	contract     *contracts.Contracts
	pk           string
	ownerAddress string
}

func NewCTL(cf *conf.Config) (*Controller, error) {
	var err error
	r := new(Controller)
	// 블록체인 네트워크와 연결할 클라이언트를 생성하기 위한 rpc url 연결
	r.client, err = ethclient.Dial(cf.Blockchain.RpcUrl)
	if err != nil {
		return r, err
	}
	// 컨트랙트 주소 및 객체 얻기
	r.tokenAddress = common.HexToAddress(cf.Blockchain.ContractAddr)
	r.contract, err = contracts.NewContracts(r.tokenAddress, r.client)
	if err != nil {
		return r, err
	}
	// 개인키 및 컨트랙트 소유자 주소 저장
	r.pk = cf.Blockchain.PrivateKey
	r.ownerAddress = cf.Blockchain.OwnerAddr

	return r, nil
}
