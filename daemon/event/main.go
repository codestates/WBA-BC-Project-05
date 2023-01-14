package event

import (
	"context"
	"fmt"
	"log"
	"time"
	conf "wba-bc-project-05/daemon/event/config"
	"wba-bc-project-05/daemon/event/model"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// config 초기화
	cf := conf.GetConfig("./config/config.toml")

	// model 초기화
	md, err := model.NewModel(cf.DB.Host)
	if err != nil {
		log.Fatal(err)
	}

	// ethclint 초기화
	client, err := ethclient.Dial(cf.Network.URL)
	if err != nil {
		panic(err)
	}

	contractAddress := common.HexToAddress("0xFa88b4faC0a7Ab78dec4679DeAB8A05CB254b039")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// 이벤트 로그 데이터 추출
			fmt.Println("address:", vLog.Address.Hex())
			fmt.Println("blockHash:", vLog.BlockHash)
			fmt.Println("blockNumber:", vLog.BlockNumber)
			fmt.Println("tx hash:", vLog.TxHash)
			fmt.Println("data:", vLog.Data)
			e := model.EventLog{}
			e.Address = vLog.Address.Hex()
			e.BlockHash = vLog.BlockHash.Hex()
			e.BlockNumber = vLog.BlockNumber
			e.TxHash = vLog.TxHash.Hex()
			// DB 저장
			err = md.SaveEvent(&e)
			if err != nil {
				log.Fatal(err)
			}
		case <-time.After(time.Second * 3):
			fmt.Println(client.NetworkID(context.Background()))
		}
	}
}
