package web3

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thangle411/golang-web3-price-watcher/coingecko"
	"github.com/thangle411/golang-web3-price-watcher/constants"
	"github.com/thangle411/golang-web3-price-watcher/pools"
)

type Token struct {
	Balance float64
	Name    string
}

type PoolInfo struct {
	Token       Token
	Denominator Token
	Price       float64
	LastPrice   float64
	SendEmail   bool
}

type PoolConfig struct {
	TokenAddress string
	PoolAddress  string
}

type Slot0Response struct {
	SqrtPriceX96               *big.Int
	Tick                       int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

func Start(receiverEmail string, senderEmail string, appPassword string) (html string, subject string) {
	ethPrice := coingecko.GetTickerPrice("ethereum")
	if ethPrice == nil {
		fmt.Println("Failed getting price from coingecko, setting it as 0")
		ethPrice = coingecko.CoingeckoResponse{
			"ethereum": {Usd: 0},
		}
	}
	fmt.Printf("\nEthereum price: $%v\n", ethPrice["ethereum"].Usd)
	fmt.Println("--------------------------------------")

	infoChannel := make(chan PoolInfo)
	var wg sync.WaitGroup
	for i := 0; i < len(pools.Pools); i++ {
		pool := &pools.Pools[i]
		var denomPrice float64
		if pool.Denominator.Symbol == "WETH" {
			denomPrice = ethPrice["ethereum"].Usd
		} else if pool.Denominator.Symbol == "USDC" {
			denomPrice = 1
		}
		wg.Add(1)
		go watchPool(denomPrice, pool, infoChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(infoChannel)
	}()

	htmlString := ""
	subjectString := ""
	for pool := range infoChannel {
		fmt.Println(pool)
		if pool.SendEmail {
			htmlString += fmt.Sprintf(`
			<br></br>
			<li>%s is $%f</li>
			<li>There is %f %s and %f %s in the pool</li>
			`, pool.Token.Name, pool.Price, pool.Denominator.Balance, pool.Denominator.Name, pool.Token.Balance, pool.Token.Name)
			subjectString += pool.Token.Name + " - "
		}
		fmt.Println("--------------------------------------")
	}

	return htmlString, subjectString
}

func watchPool(denomPrice float64, pool *pools.Pool, infoChannel chan PoolInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	web3, err := web3.NewWeb3(constants.RpcMap[pool.ChainID])
	if err != nil {
		fmt.Println("Cannot initialize a web3 instance - sending placeholder data")
		zeroVal := Token{
			Balance: 0.00,
			Name:    "placeholder",
		}
		infoChannel <- PoolInfo{
			Token:       zeroVal,
			Denominator: zeroVal,
			Price:       0.00,
		}
		return
	}
	web3.Eth.SetChainId(pool.ChainID)
	lastPrice := pool.LastPrice
	tokenBalance := getBalanceOfAddress(pool.TokenAddress, pool.PoolAddress, web3)
	denominatorBalance := getBalanceOfAddress(pool.Denominator.Address, pool.PoolAddress, web3)
	currentPrice := getTokenPrice(tokenBalance, denominatorBalance, denomPrice, pool, web3)
	pool.LastPrice = currentPrice //update price in memory

	infoChannel <- PoolInfo{
		Token:       tokenBalance,
		Denominator: denominatorBalance,
		Price:       currentPrice,
		LastPrice:   lastPrice,
		SendEmail:   pool.PriceDelta(lastPrice, currentPrice),
	}
}

func initializeContract(abi string, address string, web3 *web3.Web3) (*eth.Contract, error) {
	contract, err := web3.Eth.NewContract(abi, address)
	return contract, err
}

func getTokenPrice(tokenBalance Token, denominatorBalance Token, denomPrice float64, pool *pools.Pool, web3 *web3.Web3) float64 {
	if pool.DexType == pools.UNIV3 || pool.DexType == pools.CAMV3 {
		return getTokenPriceV3(pool, denomPrice, web3)
	} else if pool.DexType == pools.UNIV2 {
		return getTokenPriceV2(tokenBalance, denominatorBalance, denomPrice)
	}

	return 0
}

func getTokenPriceV2(tokenBalance Token, denominatorBalance Token, denomPrice float64) float64 {
	if denominatorBalance.Name == "WETH" {
		return (denomPrice * denominatorBalance.Balance) / tokenBalance.Balance
	} else if denominatorBalance.Name == "USDC" {
		return denominatorBalance.Balance / tokenBalance.Balance
	}

	return 0
}

/**
* Equation to calculate token price, the price below is denominated in ETH or USDC or other stablesy
* sqrtRatioX96 can be queried from the pool contract
* (sqrtRatioX96 ** 2) / (2 ** 192)= price
* so price in $ would be price = denominatorPriceInUSD * (sqrtRatioX96 ** 2) / (2 ** 192)
 */
func getTokenPriceV3(pool *pools.Pool, denomPrice float64, web3 *web3.Web3) float64 {
	contract, err := initializeContract(pools.ABIMap[pool.DexType].Abi, pool.PoolAddress, web3)
	if err != nil {
		fmt.Println("Cannot create contract", err)
		return 0.00
	}

	slot0, err := contract.Call(pools.ABIMap[pool.DexType].Method)
	if err != nil {
		fmt.Println("Cannot get "+pools.ABIMap[pool.DexType].Method, err)
		return 0.00
	}

	// Type assertion to treat slot0 as a slice of interfaces
	slot0Slice, ok := slot0.([]interface{})
	if !ok {
		fmt.Println("slot0 is not a []interface{}")
		return 0.00
	}

	//Check if number is bigInt
	bigIntValue, ok := slot0Slice[0].(*big.Int)
	if !ok {
		fmt.Println("Not a bigint")
		return 0.00
	}

	squared := new(big.Float).SetInt(new(big.Int).Mul(bigIntValue, bigIntValue))
	squared192 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil))
	result := new(big.Float).Quo(squared, squared192)
	denomPriceAsBigFloat := new(big.Float).SetFloat64(denomPrice)
	price := new(big.Float).Mul(denomPriceAsBigFloat, result)
	priceAsFloat64, _ := price.Float64()
	return priceAsFloat64
}

func getBalanceOfAddress[T Token](contractAddress string, walletAddress string, web3 *web3.Web3) Token {
	zeroVal := Token{
		Balance: 0,
		Name:    "",
	}
	contract, err := initializeContract(constants.TokenAbi, contractAddress, web3)
	if err != nil {
		fmt.Println("Cannot create contract")
		return zeroVal
	}

	name, err := contract.Call("symbol")
	if err != nil {
		fmt.Println("Cannot get Name")
		return zeroVal
	}

	nameAsString, ok := name.(string)
	if !ok {
		return zeroVal
	}

	decimals, err := contract.Call("decimals")
	if err != nil {
		fmt.Println("Cannot get decimals")
		return zeroVal
	}

	fDecimals, ok := decimals.(uint8)
	if !ok {
		return zeroVal
	}

	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(fDecimals)), nil)

	balanceOfPool, err := contract.Call("balanceOf", common.HexToAddress(walletAddress))
	if err != nil {
		fmt.Println("Cannot get balance of pool ")
		return zeroVal
	}

	balanceAsBigInt, ok := balanceOfPool.(*big.Int)
	if !ok {
		fmt.Println("Cannot convert to bigInt")
		return zeroVal
	}

	result := new(big.Int)
	result.Div(balanceAsBigInt, denominator)

	// fmt.Printf("%v %v tokens in pool\n", result, nameAsString)
	floatValue := new(big.Float).SetInt(result)
	convertedResult, _ := floatValue.Float64()

	return Token{
		Balance: convertedResult,
		Name:    nameAsString,
	}
}
