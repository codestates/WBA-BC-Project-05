package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	conf "wba-bc-project-05/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	contract "wba-bc-project-05/contracts"
)

func main() {
	cf := conf.NewConfig(os.Args[1])

	client, err := ethclient.Dial(cf.Blockchain.UrlHttp)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("53994feaabfb3dcc3c34c087d45d60799ccdfdf5c02b05285022d0aa763984bc")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(6721975) // in units
	auth.GasPrice = gasPrice

	// input := "1.0"
	address, tx, instance, err := contract.DeployContracts(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract address:", address.Hex())
	fmt.Println("transaction hash:", tx.Hash().Hex())

	_ = instance
}
