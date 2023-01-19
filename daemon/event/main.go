package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"wba-bc-project-05/backend/model"
	conf "wba-bc-project-05/config"
	"wba-bc-project-05/contracts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// config 초기화
	cf := conf.NewConfig("../../config/config.toml")

	// model 초기화
	md, err := model.NewModel(cf.DB.Host)
	if err != nil {
		panic(err)
	}

	// ethclint 초기화
	client, err := ethclient.Dial(cf.Blockchain.UrlWs)
	if err != nil {
		panic(err)
	}

	contractAddress := common.HexToAddress(cf.Blockchain.ContractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		panic(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ContractsMetaData.ABI)))
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// 이벤트 로그 DB에 저장
			e := model.EventLog{}
			e.Address = vLog.Address.Hex()
			e.BlockHash = vLog.BlockHash.Hex()
			e.BlockNumber = vLog.BlockNumber
			e.TxHash = vLog.TxHash.Hex()
			e.Signature = vLog.Topics[0].Hex()
			md.EventModel.InsertEventLog(e)

			// 이벤트에 따른 처리
			switch vLog.Topics[0].String() {
			case contractAbi.Events["EvCreateGame"].ID.String():
				saveGameEvent(md, contractAbi, "EvCreateGame", vLog.Data)
			case contractAbi.Events["EvBet"].ID.String():
				saveBetEvent(md, contractAbi, "EvBet", vLog.Data)
			case contractAbi.Events["EvVote"].ID.String():
				saveVoteEvent(md, contractAbi, "EvVote", vLog.Data)
			case contractAbi.Events["EvCalculate"].ID.String():
				saveResultEvent(md, contractAbi, "EvResult", vLog.Data)
			case contractAbi.Events["Transfer"].ID.String():
				saveTransferEvent(md, contractAbi, "Transfer", vLog.Topics, vLog.Data)
			default:
				fmt.Println("Unknown event")
			}

		case <-time.After(time.Second * 3):
			// fmt.Println(client.NetworkID(context.Background()))
		}
	}
}

func saveGameEvent(md *model.Model, abi abi.ABI, name string, data []byte) {
	game := model.GameWrapper{}
	err := abi.UnpackIntoInterface(&game, name, data)
	if err != nil {
		log.Fatal(err)
	}
	gameForDB, err := md.GameModel.ConvertToDB(game.Game)
	if err != nil {
		log.Fatal(err)
	}
	err = md.GameModel.Insert(gameForDB)
	if err != nil {
		log.Fatal(err)
	}
}

func saveBetEvent(md *model.Model, abi abi.ABI, name string, data []byte) {
	bet := model.BetWrapper{}
	err := abi.UnpackIntoInterface(&bet, name, data)
	if err != nil {
		log.Fatal(err)
	}
	betForDB, err := md.BetModel.ConvertToDB(bet.Bet)
	if err != nil {
		log.Fatal(err)
	}
	err = md.BetModel.Insert(betForDB)
	if err != nil {
		log.Fatal(err)
	}
}

func saveVoteEvent(md *model.Model, abi abi.ABI, name string, data []byte) {
	vote := model.VoteWrapper{}
	err := abi.UnpackIntoInterface(&vote, name, data)
	if err != nil {
		log.Fatal(err)
	}
	voteForDB, err := md.VoteModel.ConvertToDB(vote.Vote)
	if err != nil {
		log.Fatal(err)
	}
	err = md.VoteModel.Insert(voteForDB)
	if err != nil {
		log.Fatal(err)
	}
}

func saveResultEvent(md *model.Model, abi abi.ABI, name string, data []byte) {
	result := model.Result{}
	err := abi.UnpackIntoInterface(&result, name, data)
	if err != nil {
		log.Fatal(err)
	}
	resultForDB, err := md.ResultModel.ConvertForDB(result)
	if err != nil {
		log.Fatal(err)
	}
	err = md.ResultModel.Insert(resultForDB)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTransferEvent(md *model.Model, abi abi.ABI, name string, topics []common.Hash, data []byte) {
	te := model.TransferEvent{}
	err := abi.UnpackIntoInterface(&te, name, data)
	if err != nil {
		log.Fatal(err)
	}
	tedb := model.TransferEventForDB{}
	tedb.From = "0x" + topics[1].String()[26:]
	tedb.To = "0x" + topics[2].String()[26:]
	tedb.Value = te.Value.String()
	err = md.EventModel.InsertTransferEvent(tedb)
	if err != nil {
		log.Fatal(err)
	}
}
