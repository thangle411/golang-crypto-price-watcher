package pools

import (
	"github.com/thangle411/golang-web3-price-watcher/constants"
)

type Pool struct {
	PoolAddress        string
	TokenAddress       string
	DenominatorAddress string
	IsUniswapV3        bool
	ChainID            int64
	PoolAbi            string
	RPC                string
	StateMethod        string
	TokenBalance       float64
	TokenName          string
	DenominatorBalance float64
	DenominatorName    string
	LastPrice          float64
	PriceDelta         func(last float64, current float64) bool
}

var Pools = []Pool{
	{
		PoolAddress:        "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		TokenAddress:       "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
		DenominatorAddress: constants.DenominatorAddress["WETH-ETH"],
		IsUniswapV3:        true,
		ChainID:            1,
		PoolAbi:            constants.Univ3Abi,
		RPC:                "https://rpc.flashbots.net",
		StateMethod:        "slot0",
		TokenBalance:       0,
		TokenName:          "DMT",
		DenominatorBalance: 0,
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			if last < 15 && current > 15 {
				sendEmail = true
			} else if last < 25 && current > 25 {
				sendEmail = true
			} else if last < 35 && current > 35 {
				sendEmail = true
			} else if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
	{
		PoolAddress:        "0x60451B6aC55E3C5F0f3aeE31519670EcC62DC28f",
		TokenAddress:       "0x3d9907F9a368ad0a51Be60f7Da3b97cf940982D8",
		DenominatorAddress: constants.DenominatorAddress["WETH-ARB"],
		IsUniswapV3:        true,
		ChainID:            42161,
		PoolAbi:            constants.CamelotV3Abi,
		RPC:                "https://arb1.arbitrum.io/rpc",
		StateMethod:        "globalState",
		TokenBalance:       0,
		TokenName:          "GRAIL",
		DenominatorBalance: 0,
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			if last < 1000 && current > 1000 {
				sendEmail = true
			} else if last < 1100 && current > 1100 {
				sendEmail = true
			} else if last < 1200 && current > 1200 {
				sendEmail = true
			} else if last < 1300 && current > 1300 {
				sendEmail = true
			} else if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
	{
		PoolAddress:        "0x0AD1e922e764dF5AB6D636F5D21Ecc2e41E827f0",
		TokenAddress:       "0x772598E9e62155D7fDFe65FdF01EB5a53a8465BE",
		DenominatorAddress: constants.DenominatorAddress["WETH-ARB"],
		IsUniswapV3:        true,
		ChainID:            42161,
		PoolAbi:            constants.Univ3Abi,
		RPC:                "https://arb1.arbitrum.io/rpc",
		StateMethod:        "slot0",
		TokenBalance:       0,
		TokenName:          "EMP",
		DenominatorBalance: 0,
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			if last < 15 && current > 15 {
				sendEmail = true
			} else if last < 25 && current > 25 {
				sendEmail = true
			} else if last < 35 && current > 35 {
				sendEmail = true
			} else if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
}
