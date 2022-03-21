package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Withdrawal struct {
	To     common.Address
	Amount *big.Int
}

func main() {
	client, err := ethclient.Dial("ws://127.0.0.1:8545/")
	if err != nil {
		log.Fatal(err)
	}

	// contractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	contractAddress := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	abiString := `[
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_token",
          "type": "address"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "Deposit",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "Withdrawal",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "deposit",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "withdraw",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

	contractAbi, err := abi.JSON(strings.NewReader(string(abiString)))
	if err != nil {
		log.Fatal(err)
	}

	logWithdrawalSig := []byte("Withdrawal(address,uint256)")
	LogDepositSig := []byte("Deposit(address,uint256)")
	logWithdrawalSigHash := crypto.Keccak256Hash(logWithdrawalSig)
	logDepositSigHash := crypto.Keccak256Hash(LogDepositSig)

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
			fmt.Printf("Log Index: %d\n", vLog.Index)
			fmt.Printf("TxHash: %s\n", vLog.TxHash.String())
			fmt.Printf("Topic0: %s, %s\n", vLog.Topics[0].Hex(), logWithdrawalSigHash.Hex())
			switch vLog.Topics[0].Hex() {
			case logWithdrawalSigHash.Hex():

				var transferEvent Withdrawal

				err := contractAbi.UnpackIntoInterface(&transferEvent, "Withdrawal", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				// transferEvent.to = common.HexToAddress(vLog.Topics[1].Hex())
				// transferEvent.amount = vLog.Topics[2].Big()
				fmt.Printf("To: %s\n", transferEvent.To.Hex())
				fmt.Printf("Amount: %s\n", transferEvent.Amount)

			case logDepositSigHash.Hex():

			}
		}
	}
}
