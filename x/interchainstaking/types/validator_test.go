package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var (
	acc1 = utils.GenerateAccAddressForTest().String()
	v1   = utils.GenerateValAddressForTest().String()
	v2   = utils.GenerateValAddressForTest().String()
	v3   = utils.GenerateValAddressForTest().String()
	v4   = utils.GenerateValAddressForTest().String()
)

func TestNormalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 0)
}

func TestOrdinalizeIntentWithZeroLength(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di = di.Ordinalize(sdk.NewDec(10000))
	require.Equal(t, len(di.Intents), 0)
}

func TestNormalizeIntentWithOneIntent(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{}}
	di.Intents = append(di.Intents, &types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDec(1000)})
	di = di.Normalize()
	require.Equal(t, len(di.Intents), 1)
	require.Equal(t, di.Intents[0].Weight, sdk.OneDec())
}

func TestNormalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: acc1, Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.NewDec(1000)},
		{ValoperAddress: v2, Weight: sdk.NewDec(1000)},
		{ValoperAddress: v3, Weight: sdk.NewDec(1000)},
	}}

	di = di.Normalize()
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.Intents[0].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents[1].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
	require.Equal(t, di.Intents[2].Weight, sdk.OneDec().Quo(sdk.NewDec(3)))
}

func TestNormalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: utils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.NewDec(5).Quo(sdk.NewDec(50))},
		{ValoperAddress: v2, Weight: sdk.NewDec(10).Quo(sdk.NewDec(50))},
		{ValoperAddress: v3, Weight: sdk.NewDec(35).Quo(sdk.NewDec(50))},
	}}

	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, di.MustIntentForValoper(v1).Weight, sdk.NewDecWithPrec(1, 1))
	require.Equal(t, di.MustIntentForValoper(v2).Weight, sdk.NewDecWithPrec(2, 1))
	require.Equal(t, di.MustIntentForValoper(v3).Weight, sdk.NewDecWithPrec(7, 1))
}

func TestOrdinalizeIntentWithEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: utils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v2, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v3, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
	}}
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, len(di.Intents), 3)
	require.Equal(t, sdk.NewInt(1000), di.MustIntentForValoper(v1).Weight.RoundInt())
}

func TestOrdinalizeIntentWithNonEqualIntents(t *testing.T) {
	di := types.DelegatorIntent{Delegator: utils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.NewDec(5).Quo(sdk.NewDec(50))},
		{ValoperAddress: v2, Weight: sdk.NewDec(10).Quo(sdk.NewDec(50))},
		{ValoperAddress: v3, Weight: sdk.NewDec(35).Quo(sdk.NewDec(50))},
	}}
	di = di.Ordinalize(sdk.NewDec(3000))
	require.Equal(t, di.MustIntentForValoper(v1).Weight.RoundInt(), sdk.NewInt(300))
	require.Equal(t, di.MustIntentForValoper(v2).Weight.RoundInt(), sdk.NewInt(600))
	require.Equal(t, di.MustIntentForValoper(v3).Weight.RoundInt(), sdk.NewInt(2100))
}

func TestAddOrdinal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: utils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v2, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v3, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
	}}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: v1, Weight: sdk.NewDec(1000)},
		{ValoperAddress: v2, Weight: sdk.NewDec(2000)},
	}

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 3, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(v1).Weight, sdk.NewDec(3).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(v2).Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(v3).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
}

func TestAddOrdinalWithNewVal(t *testing.T) {
	di := types.DelegatorIntent{Delegator: utils.GenerateAccAddressForTest().String(), Intents: []*types.ValidatorIntent{
		{ValoperAddress: v1, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v2, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
		{ValoperAddress: v3, Weight: sdk.OneDec().Quo(sdk.NewDec(3))},
	}}

	newIntents := types.ValidatorIntents{
		{ValoperAddress: v4, Weight: sdk.NewDec(1000)},
		{ValoperAddress: v3, Weight: sdk.NewDec(2000)},
	}

	di = di.AddOrdinal(sdk.NewDec(6000), newIntents)

	require.Equal(t, 4, len(di.Intents))

	require.Equal(t, di.MustIntentForValoper(v1).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(v2).Weight, sdk.NewDec(2).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(v3).Weight, sdk.NewDec(4).QuoTruncate(sdk.NewDec(9)))
	require.Equal(t, di.MustIntentForValoper(v4).Weight, sdk.NewDec(1).QuoTruncate(sdk.NewDec(9)))
}
