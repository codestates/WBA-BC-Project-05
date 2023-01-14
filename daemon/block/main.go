package main

import (
	"context"
	"fmt"
	"log"
	conf "wba-bc-project-05/daemon/block/config"
	"wba-bc-project-05/daemon/block/model"

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
		log.Fatal(err)
	}

	// subscribe
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())

			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(block.Hash().Hex())
			fmt.Println(block.Number().Uint64())
			fmt.Println(block.Time())
			fmt.Println(block.Nonce())
			fmt.Println(len(block.Transactions()))

			// TODO: 블록 구조체 생성
			b := model.Block{}
			b.BlockHash = block.Hash().Hex()
			b.BlockNumber = block.Number().Uint64()
			b.GasLimit = block.GasLimit()
			b.GasUsed = block.GasUsed()
			b.Nonce = block.Nonce()
			b.Time = block.Time()
			b.Transactions = make([]model.Transaction, 0)

			// TODO: 트랜잭션 추출
			for _, _tx := range block.Transactions() {
				// TODO: 트랜잭션 구조체 생성
				msg, err := _tx.AsMessage(types.LatestSignerForChainID(_tx.ChainId()), block.BaseFee())
				if err != nil {
					log.Fatal(err)
					break
				}
				tx := model.Transaction{}
				tx.Amount = _tx.Value().Uint64()
				tx.BlockHash = block.Hash().Hex()
				tx.BlockNumber = block.Number().Uint64()
				tx.From = msg.From().Hex()
				tx.GasLimit = _tx.GasPrice().Uint64()
				tx.GasPrice = _tx.Gas()
				tx.Nonce = _tx.Nonce()
				tx.To = ""
				tx.TxHash = _tx.Hash().Hex()
				if _tx.To() != nil {
					tx.To = _tx.To().Hex()
				}
				b.Transactions = append(b.Transactions, tx)
			}

			// DB 저장
			err = md.SaveBlock(&b)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
