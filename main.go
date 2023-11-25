package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thangle411/golang-web3-price-watcher/jsonHelper"

	"github.com/chenzhijie/go-web3"
	"github.com/joho/godotenv"
)

var abiString string = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
				{
						"name": "",
						"type": "string"
				}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_owner",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`

type Coin struct {
	Usd float64
}
type CoingeckoResponse map[string]Coin

type PoolConfig struct {
	tokenAddress string;
	poolAddress string
}

type PoolBalance struct {
	tokenBalance float64
	ethBalance float64
}

var ethAddress = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
var client *http.Client

func main() {

	i := 0
	for i < 5 {
		fmt.Printf("\n---------------------------\n")
		ethPrice := getTickerPrice("ethereum")
		if ethPrice == nil {
			log.Fatal("Failed getting price from coingecko")
		}

		fmt.Printf("\nEthereum price: %v\n", ethPrice["ethereum"].Usd)

		balances := watchPoolBalance(PoolConfig{
			tokenAddress: "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
			poolAddress: "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		})

		ethValuation := balances.ethBalance*ethPrice["ethereum"].Usd
		fmt.Printf("Eth value in pool: $%v\n", ethValuation)
		fmt.Printf("Price for each token: $%v\n", balances.tokenBalance/ethValuation)

		time.Sleep(30 * time.Second)
		i++
	}
}

func watchPoolBalance(config PoolConfig) PoolBalance {
	var rpcProviderURL = "https://rpc.flashbots.net"
	web3, err := web3.NewWeb3(rpcProviderURL)
	if err != nil {
		log.Fatal("Cannot initialize a web3 instance")
	}

	tokenBalance := getBalanceOfAddress(config.tokenAddress, config.poolAddress, web3)
	ethBalance := getBalanceOfAddress(ethAddress, config.poolAddress, web3)

	return PoolBalance{
		tokenBalance,
		ethBalance,
	}
}

func getBalanceOfAddress (contractAddress string, walletAddress string, web3 *web3.Web3) float64 {
	contract, err := web3.Eth.NewContract(abiString, contractAddress)
	if err != nil {
		log.Fatal("Cannot create contract")
	}

	name, err := contract.Call("name")
	if err != nil {
		log.Fatal("Cannot get Name")
	}

	decimals, err := contract.Call("decimals")
	if err != nil {
		log.Fatal("Cannot get decimals")
	}

	fDecimals, ok := decimals.(uint8)
	
	if !ok {
		return 0
	}
	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(fDecimals)), nil)

	balanceOfPool, err := contract.Call("balanceOf", common.HexToAddress(walletAddress))
	if err != nil {
		log.Fatal("Cannot get balance of pool ", err)
	}

	balanceAsBigInt, ok :=  balanceOfPool.(*big.Int)
	if !ok {
		return 0
	}

	result := new(big.Int)
	result.Div(balanceAsBigInt, denominator)

	fmt.Printf("%v %v tokens in pool\n",result, name)
	floatValue := new(big.Float).SetInt(result)
	convertedResult, _ := floatValue.Float64()

	return convertedResult
}

func getTickerPrice(ticker string) CoingeckoResponse {
	godotenv.Load()
	cgURL := os.Getenv("COINGECKO_URL")
	if cgURL == "" {
		log.Fatal("Coingecko URL not found in .env")
	}

	url := cgURL + "/simple/price?ids=" + ticker + "&vs_currencies=usd"

	var data CoingeckoResponse
	err := jsonHelper.GetJson(url, &data)
	if err != nil {
		log.Fatal("Error getting data from Coingecko")
	}

	return data
}