package pools

import (
	"github.com/thangle411/golang-web3-price-watcher/constants"
)

type Pool struct {
	PoolAddress        string
	TokenAddress       string
	DenominatorAddress string
	ChainID            int64
	PoolAbi            string
	RPC                string
	StateMethod        string
	TokenName          string
	DenominatorName    string
	LastPrice          float64
	PriceDelta         func(last float64, current float64) bool
}

var Pools = []Pool{
	{
		PoolAddress:        "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		TokenAddress:       "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
		DenominatorAddress: constants.DenominatorAddress["WETH-ETH"],
		ChainID:            constants.ChainsList["ethereum"],
		PoolAbi:            constants.Univ3Abi,
		RPC:                "https://rpc.flashbots.net",
		StateMethod:        "slot0",
		TokenName:          "DMT",
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			thresholds := []float64{15, 20, 25, 30, 35}
			for _, threshold := range thresholds {
				if last < threshold && current > threshold {
					sendEmail = true
					break
				}
			}
			if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
	{
		PoolAddress:        "0x60451B6aC55E3C5F0f3aeE31519670EcC62DC28f",
		TokenAddress:       "0x3d9907F9a368ad0a51Be60f7Da3b97cf940982D8",
		DenominatorAddress: constants.DenominatorAddress["WETH-ARB"],
		ChainID:            constants.ChainsList["arbitrum"],
		PoolAbi:            constants.CamelotV3Abi,
		RPC:                "https://arb1.arbitrum.io/rpc",
		StateMethod:        "globalState",
		TokenName:          "GRAIL",
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			thresholds := []float64{1000, 1100, 1200, 1300, 1400, 1500}
			for _, threshold := range thresholds {
				if last < threshold && current > threshold {
					sendEmail = true
					break
				}
			}
			if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
	{
		PoolAddress:        "0x0AD1e922e764dF5AB6D636F5D21Ecc2e41E827f0",
		TokenAddress:       "0x772598E9e62155D7fDFe65FdF01EB5a53a8465BE",
		DenominatorAddress: constants.DenominatorAddress["WETH-ARB"],
		ChainID:            constants.ChainsList["arbitrum"],
		PoolAbi:            constants.Univ3Abi,
		RPC:                "https://arb1.arbitrum.io/rpc",
		StateMethod:        "slot0",
		TokenName:          "EMP",
		DenominatorName:    "WETH",
		LastPrice:          0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			thresholds := []float64{15, 25, 35}
			for _, threshold := range thresholds {
				if last < threshold && current > threshold {
					sendEmail = true
					break
				}
			}
			if percent > 2.5 {
				sendEmail = true
			}
			return sendEmail
		},
	},
}
