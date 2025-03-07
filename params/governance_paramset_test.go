package params

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGovParamSet_ParseValue(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	zeroAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	mintingAmount := "9600000000000000000"
	mintingAmountBig, _ := new(big.Int).SetString(mintingAmount, 10)

	testcases := []struct {
		ty     *govParamType
		value  interface{}
		parsed interface{} // If ok, expected value. Ignored if not ok.
		ok     bool        // Expected 'ok'
	}{
		{govParamTypeGovMode, "none", "none", true},
		{govParamTypeGovMode, "single", "single", true},
		{govParamTypeGovMode, "ballot", "ballot", true},
		{govParamTypeGovMode, "asdf", nil, false},
		{govParamTypeGovMode, "", nil, false},
		{govParamTypeGovMode, 1, nil, false},

		{govParamTypeAddress, zeroAddr.Hex(), zeroAddr, true},
		{govParamTypeAddress, zeroAddr, zeroAddr, true},
		{govParamTypeAddress, 1, nil, false},

		{govParamTypeUint64, int(7), uint64(7), true},
		{govParamTypeUint64, uint(7), uint64(7), true},
		{govParamTypeUint64, uint64(7), uint64(7), true},
		{govParamTypeUint64, float64(1e9), uint64(1e9), true},
		{govParamTypeUint64, "123", nil, false},
		{govParamTypeUint64, -1, nil, false},
		{govParamTypeUint64, -12.0, nil, false},

		{govParamTypeBigInt, mintingAmount, mintingAmount, true},
		{govParamTypeBigInt, mintingAmountBig, mintingAmount, true},
		{govParamTypeBigInt, "123", "123", true},
		{govParamTypeBigInt, "-123", nil, false},
		{govParamTypeBigInt, "abc", nil, false},
		{govParamTypeBigInt, "", nil, false},

		{govParamTypeRatio, "100/0/0", "100/0/0", true},
		{govParamTypeRatio, "30/30/40", "30/30/40", true},
		{govParamTypeRatio, "10/20/30/40", nil, false},
		{govParamTypeRatio, "0/0/0", nil, false},
		{govParamTypeRatio, "1/2/3", nil, false},
		{govParamTypeRatio, "", nil, false},

		{govParamTypeBool, true, true, true},
		{govParamTypeBool, 0, nil, false},
		{govParamTypeBool, "", nil, false},
	}

	for _, tc := range testcases {
		parsed, ok := tc.ty.ParseValue(tc.value)
		assert.Equal(t, tc.ok, ok)
		if ok {
			assert.Equal(t, tc.parsed, parsed)
		}
	}
}

func TestGovParamSet_ParseBytes(t *testing.T) {
	zeroAddrHex := "0x0000000000000000000000000000000000000000"
	zeroAddr := common.HexToAddress(zeroAddrHex)
	mintingAmountStr := "9600000000000000000"

	testcases := []struct {
		ty     *govParamType
		bytes  []byte
		parsed interface{} // If ok, expected value. Ignored if not ok.
		ok     bool        // Expected 'ok'
	}{
		{govParamTypeGovMode, []byte("single"), "single", true},
		{govParamTypeGovMode, []byte(""), nil, false},

		{govParamTypeAddress, zeroAddr.Bytes(), zeroAddr, true},
		{govParamTypeAddress, []byte(zeroAddr.Hex()), nil, false},
		{govParamTypeAddress, []byte(""), nil, false},

		{govParamTypeUint64, []byte{0x12, 0x34}, uint64(0x1234), true},
		{govParamTypeUint64, []byte{}, uint64(0), true},
		{govParamTypeUint64, []byte{1, 2, 3, 4, 5, 6, 7, 8}, uint64(0x0102030405060708), true},
		{govParamTypeUint64, []byte{1, 1, 2, 3, 4, 5, 6, 7, 8}, nil, false},

		{govParamTypeBigInt, []byte(mintingAmountStr), mintingAmountStr, true},
		{govParamTypeBigInt, []byte("abcd"), nil, false},
		{govParamTypeBigInt, []byte(""), nil, false},

		{govParamTypeRatio, []byte("100/0/0"), "100/0/0", true},
		{govParamTypeRatio, []byte("10/20/30/40"), nil, false},
		{govParamTypeRatio, []byte(""), nil, false},

		{govParamTypeBool, []byte{0x01}, true, true},
		{govParamTypeBool, []byte{0x00}, false, true},
		{govParamTypeBool, []byte{0x99}, nil, false},
		{govParamTypeBool, []byte{}, nil, false},
	}

	for _, tc := range testcases {
		parsed, ok := tc.ty.ParseBytes(tc.bytes)
		assert.Equal(t, tc.ok, ok)
		if ok {
			assert.Equal(t, tc.parsed, parsed)
		}
	}
}

func TestGovParamSet_GlobalMaps(t *testing.T) {
	// Check that govParam* maps hold the same set of parameters.

	assert.Equal(t, len(govParamTypes), len(govParamNames))
	for _, key := range govParamNames {
		assert.NotNil(t, govParamTypes[key])
	}
}

func TestGovParamSet_Get(t *testing.T) {
	num := uint64(123456)
	p, _ := NewGovParamSetStrMap(map[string]interface{}{
		"istanbul.epoch": num,
	})

	// Exists
	v, ok := p.Get(Epoch)
	assert.True(t, ok)
	assert.Equal(t, num, v)
	assert.Equal(t, num, p.MustGet(Epoch))

	// Not exists
	v, ok = p.Get(CommitteeSize)
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestGovParamSet_Nominal(t *testing.T) {
	c := CypressChainConfig.Copy()
	c.Governance.KIP71 = &KIP71Config{
		LowerBoundBaseFee:         12340000,
		UpperBoundBaseFee:         56780000,
		GasTarget:                 3000,
		MaxBlockGasUsedForBaseFee: 6000,
		BaseFeeDenominator:        100,
	}
	p, err := NewGovParamSetChainConfig(c)
	assert.Nil(t, err)

	assert.Equal(t, c.Istanbul.Epoch, p.Epoch())
	assert.Equal(t, c.Istanbul.ProposerPolicy, p.Policy())
	assert.Equal(t, c.Istanbul.SubGroupSize, p.CommitteeSize())
	assert.Equal(t, c.UnitPrice, p.UnitPrice())
	assert.Equal(t, c.DeriveShaImpl, p.DeriveShaImpl())
	assert.Equal(t, c.Governance.GovernanceMode, p.GovernanceModeStr())
	assert.Equal(t, c.Governance.GoverningNode, p.GoverningNode())
	assert.Equal(t, c.Governance.Reward.MintingAmount.String(), p.MintingAmountStr())
	assert.Equal(t, c.Governance.Reward.MintingAmount, p.MintingAmountBig())
	assert.Equal(t, c.Governance.Reward.Ratio, p.Ratio())
	assert.Equal(t, c.Governance.Reward.UseGiniCoeff, p.UseGiniCoeff())
	assert.Equal(t, c.Governance.Reward.DeferredTxFee, p.DeferredTxFee())
	assert.Equal(t, c.Governance.Reward.MinimumStake.String(), p.MinimumStakeStr())
	assert.Equal(t, c.Governance.Reward.MinimumStake, p.MinimumStakeBig())
	assert.Equal(t, c.Governance.Reward.StakingUpdateInterval, p.StakeUpdateInterval())
	assert.Equal(t, c.Governance.Reward.ProposerUpdateInterval, p.ProposerRefreshInterval())
	assert.Equal(t, c.Governance.KIP71.LowerBoundBaseFee, p.LowerBoundBaseFee())
	assert.Equal(t, c.Governance.KIP71.UpperBoundBaseFee, p.UpperBoundBaseFee())
	assert.Equal(t, c.Governance.KIP71.GasTarget, p.GasTarget())
	assert.Equal(t, c.Governance.KIP71.MaxBlockGasUsedForBaseFee, p.MaxBlockGasUsedForBaseFee())
	assert.Equal(t, c.Governance.KIP71.BaseFeeDenominator, p.BaseFeeDenominator())
}

func TestGovParamSet_New(t *testing.T) {
	p, err := NewGovParamSetStrMap(map[string]interface{}{
		"istanbul.epoch": 604800,
	})
	assert.Nil(t, err)
	v, ok := p.Get(Epoch)
	assert.Equal(t, uint64(604800), v)
	assert.True(t, ok)

	p, err = NewGovParamSetIntMap(map[int]interface{}{
		Epoch: 604800,
	})
	assert.Nil(t, err)
	v, ok = p.Get(Epoch)
	assert.Equal(t, uint64(604800), v)
	assert.True(t, ok)

	p, err = NewGovParamSetBytesMap(map[string][]byte{
		"istanbul.epoch": {0x12, 0x34},
	})
	assert.Nil(t, err)
	v, ok = p.Get(Epoch)
	assert.Equal(t, uint64(0x1234), v)
	assert.True(t, ok)

	c := CypressChainConfig
	p, err = NewGovParamSetChainConfig(c)
	assert.Nil(t, err)
	v, ok = p.Get(Epoch)
	assert.Equal(t, c.Istanbul.Epoch, v)
	assert.True(t, ok)

	// Error cases
	_, err = NewGovParamSetStrMap(map[string]interface{}{
		"istanbul.epoch": "asdf",
	})
	assert.NotNil(t, err)

	_, err = NewGovParamSetIntMap(map[int]interface{}{
		Epoch: "asdf",
	})
	assert.NotNil(t, err)

	_, err = NewGovParamSetBytesMap(map[string][]byte{
		"istanbul.epoch": {1, 1, 2, 3, 4, 5, 6, 7, 8},
	})
	assert.NotNil(t, err)
}

func TestGovParamSet_Merged(t *testing.T) {
	base, err := NewGovParamSetStrMap(map[string]interface{}{
		"istanbul.epoch":         123456,
		"istanbul.committeesize": 77,
	})
	assert.Nil(t, err)

	update, err := NewGovParamSetStrMap(map[string]interface{}{
		"istanbul.committeesize": 99,
		"istanbul.policy":        2,
	})
	assert.Nil(t, err)

	p := NewGovParamSetMerged(base, update)

	// Was only in base
	v, ok := p.Get(Epoch)
	assert.Equal(t, uint64(123456), v)
	assert.True(t, ok)

	// Was only in update
	v, ok = p.Get(Policy)
	assert.Equal(t, uint64(2), v)
	assert.True(t, ok)

	// Was in both - prefers the value in update
	v, ok = p.Get(CommitteeSize)
	assert.Equal(t, uint64(99), v)
	assert.True(t, ok)
}

func TestGovParamSet_RegressDb(t *testing.T) {
	// MiscDB stores governance data as JSON strings. The value types can be
	// slightly wrong during unmarshal because we unmarshal into interface{}.
	// Namely, JSON integers can be converted to float64.

	c := CypressChainConfig
	p, err := NewGovParamSetChainConfig(c)
	assert.Nil(t, err)

	// Simulate database write then read
	j, _ := json.Marshal(p.StrMap())
	var data map[string]interface{}
	json.Unmarshal(j, &data)

	pp, err := NewGovParamSetStrMap(data)
	assert.Nil(t, err)
	assert.Equal(t, p.items, pp.items)
}

func TestGovParamSet_GetMap(t *testing.T) {
	c := CypressChainConfig
	p, err := NewGovParamSetChainConfig(c)
	assert.Nil(t, err)

	sm := p.StrMap()
	psm, err := NewGovParamSetStrMap(sm)
	assert.Nil(t, err)
	assert.Equal(t, p.items, psm.items)

	im := p.IntMap()
	pim, err := NewGovParamSetIntMap(im)
	assert.Nil(t, err)
	assert.Equal(t, p.items, pim.items)
}

func TestGovParamSet_ChainConfig(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	testcases := []struct {
		pset     map[int]interface{}
		expected *ChainConfig
	}{
		// partial chainConfig
		{
			pset: map[int]interface{}{
				UnitPrice:               25e9,
				Epoch:                   30,
				GoverningNode:           common.HexToAddress("0x0000000000000000000000000000000000000400"),
				MintingAmount:           new(big.Int).SetUint64(9.6e18),
				Ratio:                   "34/54/12",
				UseGiniCoeff:            true,
				ProposerRefreshInterval: 3600,
				LowerBoundBaseFee:       10000,
			},
			expected: &ChainConfig{
				UnitPrice: 25e9,
				Istanbul: &IstanbulConfig{
					Epoch: 30,
				},
				Governance: &GovernanceConfig{
					GoverningNode: common.HexToAddress("0x0000000000000000000000000000000000000400"),
					Reward: &RewardConfig{
						MintingAmount:          new(big.Int).SetUint64(9.6e18),
						Ratio:                  "34/54/12",
						UseGiniCoeff:           true,
						ProposerUpdateInterval: 3600,
					},
					KIP71: &KIP71Config{
						LowerBoundBaseFee: 10000,
					},
				},
			},
		},
		// complete chainConfig
		{
			pset: map[int]interface{}{
				UnitPrice:                 25e9,
				Epoch:                     30,
				Policy:                    1,
				CommitteeSize:             27,
				GoverningNode:             common.HexToAddress("0x0000000000000000000000000000000000000400"),
				GovernanceMode:            "single",
				MintingAmount:             new(big.Int).SetUint64(9.6e18),
				Ratio:                     "34/54/12",
				UseGiniCoeff:              true,
				DeferredTxFee:             true,
				StakeUpdateInterval:       86400,
				ProposerRefreshInterval:   3600,
				MinimumStake:              big.NewInt(5e6),
				LowerBoundBaseFee:         25000000000,
				UpperBoundBaseFee:         750000000000,
				GasTarget:                 30000000,
				MaxBlockGasUsedForBaseFee: 60000000,
				BaseFeeDenominator:        20,
				Kip82Ratio:                "20/80",
			},
			expected: &ChainConfig{
				UnitPrice: 25e9,
				Istanbul: &IstanbulConfig{
					Epoch:          30,
					ProposerPolicy: 1,
					SubGroupSize:   27,
				},
				Governance: &GovernanceConfig{
					GoverningNode:  common.HexToAddress("0x0000000000000000000000000000000000000400"),
					GovernanceMode: "single",
					Reward: &RewardConfig{
						MintingAmount:          new(big.Int).SetUint64(9.6e18),
						Ratio:                  "34/54/12",
						Kip82Ratio:             "20/80",
						UseGiniCoeff:           true,
						DeferredTxFee:          true,
						StakingUpdateInterval:  86400,
						ProposerUpdateInterval: 3600,
						MinimumStake:           big.NewInt(5e6),
					},
					KIP71: &KIP71Config{
						LowerBoundBaseFee:         25000000000,
						UpperBoundBaseFee:         750000000000,
						GasTarget:                 30000000,
						MaxBlockGasUsedForBaseFee: 60000000,
						BaseFeeDenominator:        20,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		pset, err := NewGovParamSetIntMap(tc.pset)
		require.Nil(t, err)
		config := pset.ToChainConfig()
		assert.Equal(t, tc.expected, config)
	}
}
