package pools

import (
	"github.com/thangle411/golang-web3-price-watcher/constants"
	"github.com/thangle411/golang-web3-price-watcher/utils"
)

type DexType string

const (
	UNIV2 DexType = "UNIV2"
	UNIV3 DexType = "UNIV3"
	CAMV3 DexType = "CAMV3"
)

type ABIStruct struct {
	Abi    string
	Method string
}

var ABIMap = map[DexType]ABIStruct{
	UNIV2: {
		Abi: constants.UniV2Abi,
	},
	UNIV3: {
		Abi:    constants.Univ3Abi,
		Method: "slot0",
	},
	CAMV3: {
		Abi:    constants.CamelotV3Abi,
		Method: "globalState",
	},
}

type Pool struct {
	PoolAddress  string
	TokenAddress string
	Denominator  constants.Token
	ChainID      int64
	DexType      DexType
	TokenName    string
	LastPrice    float64
	PriceDelta   func(last float64, current float64) bool
}

var Pools = []Pool{
	{
		PoolAddress:  "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		TokenAddress: "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
		Denominator:  constants.Denominator["WETH-ETH"],
		ChainID:      constants.ChainsList["ethereum"],
		DexType:      UNIV3,
		TokenName:    "DMT",
		LastPrice:    0,
		PriceDelta: func(last float64, current float64) bool {
			return utils.ShouldNotify(last, current, 5, 40, 0.5, 2.5)
		},
	},
	{
		PoolAddress:  "0x60451B6aC55E3C5F0f3aeE31519670EcC62DC28f",
		TokenAddress: "0x3d9907F9a368ad0a51Be60f7Da3b97cf940982D8",
		Denominator:  constants.Denominator["WETH-ARB"],
		ChainID:      constants.ChainsList["arbitrum"],
		DexType:      CAMV3,
		TokenName:    "GRAIL",
		LastPrice:    0,
		PriceDelta: func(last float64, current float64) bool {
			return utils.ShouldNotify(last, current, 700, 1600, 50, 2.5)
		},
	},
	{
		PoolAddress:  "0x0AD1e922e764dF5AB6D636F5D21Ecc2e41E827f0",
		TokenAddress: "0x772598E9e62155D7fDFe65FdF01EB5a53a8465BE",
		Denominator:  constants.Denominator["WETH-ARB"],
		ChainID:      constants.ChainsList["arbitrum"],
		DexType:      UNIV3,
		TokenName:    "EMP",
		LastPrice:    0,
		PriceDelta: func(last float64, current float64) bool {
			return utils.ShouldNotify(last, current, 15, 75, 10, 2.5)
		},
	},
	{
		PoolAddress:  "0x2556082B685593a652Dc6b6FCe523CaF6E5590fD",
		TokenAddress: "0x7865eC47bEF9823AD0010c4970Ed90A5E8107E53",
		Denominator:  constants.Denominator["WETH-ETH"],
		ChainID:      constants.ChainsList["ethereum"],
		DexType:      UNIV2,
		TokenName:    "NAAI",
		LastPrice:    0,
		PriceDelta: func(last float64, current float64) bool {
			return utils.ShouldNotify(last, current, 0.05, 0.3, 0.025, 2.5)
		},
	},
}
