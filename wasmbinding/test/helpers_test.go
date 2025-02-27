package wasmbinding

import (
	"testing"
	"time"

	"github.com/ingenuity-build/quicksilver/app"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateTestInput(t *testing.T) (*app.Quicksilver, sdk.Context) {
	quicksilverApp := app.Setup(t, false)
	ctx := quicksilverApp.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "quicksilver-1", Time: time.Now().UTC()})
	return quicksilverApp, ctx
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}
