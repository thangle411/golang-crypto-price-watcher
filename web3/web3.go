package web3

import (
	"fmt"
	"log"
	"math/big"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thangle411/golang-web3-price-watcher/constants"
)

type Token struct {
	Balance float64
	Name string
}

type PoolBalance struct {
	Token Token
	Eth Token
	Price float64
}

type PoolConfig struct {
	TokenAddress string;
	PoolAddress string
}

type Slot0Response struct {
	SqrtPriceX96              *big.Int 
	Tick                      int      
	ObservationIndex          uint16   
	ObservationCardinality    uint16   
	ObservationCardinalityNext uint16  
	FeeProtocol               uint8    
	Unlocked                  bool 
}

func WatchPoolBalance(ethPrice float64, config PoolConfig) PoolBalance {
	var rpcProviderURL = "https://rpc.flashbots.net"

	web3, err := web3.NewWeb3(rpcProviderURL)
	if err != nil {
		log.Fatal("Cannot initialize a web3 instance")
	}

	tokenBalance := getBalanceOfAddress(config.TokenAddress, config.PoolAddress, web3)
	ethBalance := getBalanceOfAddress(constants.WrappedEthAddress, config.PoolAddress, web3)
	price := getTokenPrice(ethPrice, config.PoolAddress, web3)

	return PoolBalance{
		Token: tokenBalance,
		Eth: ethBalance,
		Price: price,
	}
}

func initializeContract( abi string, address string, web3 *web3.Web3) (*eth.Contract, error) {
	contract, err := web3.Eth.NewContract(abi, address)
	return contract, err
}

/** 
* Equation to calculate token price, the price below is denominated in ETH
* sqrtRatioX96 can be queried from the pool contract
* (sqrtRatioX96 ** 2) / (2 ** 192)= price
* so price in $ would be price = ethPriceInUSD * (sqrtRatioX96 ** 2) / (2 ** 192)
*/
func getTokenPrice(ethPrice float64, poolAddress string, web3 *web3.Web3) float64 {
	contract, err := initializeContract(constants.PoolAbi, poolAddress, web3)
	if err != nil {
		log.Fatal("Cannot create contract", err)
	}

	slot0, err := contract.Call("slot0")
	if err != nil {
		log.Fatal("Cannot get slot0", err)
	}

	// Type assertion to treat slot0 as a slice of interfaces
	slot0Slice, ok := slot0.([]interface{})
	if !ok {
		log.Fatal("slot0 is not a []interface{}")
	}

	//Check if number is bigInt
	bigIntValue, ok := slot0Slice[0].(*big.Int)
	if !ok {
		log.Fatal("Not a bigint")
	}

	squared := new(big.Float).SetInt(new(big.Int).Mul(bigIntValue, bigIntValue))
	squared192 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil))
	result := new(big.Float).Quo(squared, squared192)
	ethPriceAsBigFloat := new(big.Float).SetFloat64(ethPrice)
	price := new(big.Float).Mul(ethPriceAsBigFloat, result)
	priceAsFloat64, _ := price.Float64()
	return priceAsFloat64
}

func getBalanceOfAddress[T Token] (contractAddress string, walletAddress string, web3 *web3.Web3) Token {
	zeroVal := Token{
		Balance: 0,
		Name: "",
	}
	contract, err := initializeContract(constants.TokenAbi, contractAddress, web3,)
	if err != nil {
		log.Fatal("Cannot create contract")
	}

	name, err := contract.Call("name")
	if err != nil {
		log.Fatal("Cannot get Name")
	}

	nameAsString, ok := name.(string)
	if !ok {
		return zeroVal
	}

	decimals, err := contract.Call("decimals")
	if err != nil {
		log.Fatal("Cannot get decimals")
	}

	fDecimals, ok := decimals.(uint8)
	if !ok {
		return zeroVal
	}

	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(fDecimals)), nil)

	balanceOfPool, err := contract.Call("balanceOf", common.HexToAddress(walletAddress))
	if err != nil {
		log.Fatal("Cannot get balance of pool ", err)
	}

	balanceAsBigInt, ok :=  balanceOfPool.(*big.Int)
	if !ok {
		return zeroVal
	}

	result := new(big.Int)
	result.Div(balanceAsBigInt, denominator)

	fmt.Printf("%v %v tokens in pool\n",result, name)
	floatValue := new(big.Float).SetInt(result)
	convertedResult, _ := floatValue.Float64()

	return Token{
		Balance: convertedResult,
		Name: nameAsString,
	}
}