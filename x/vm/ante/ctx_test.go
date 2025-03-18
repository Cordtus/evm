package ante_test

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/evm/testutil/integration/os/network"
	evmante "github.com/cosmos/evm/x/vm/ante"
)

func (suite *EvmAnteTestSuite) TestBuildEvmExecutionCtx() {
	network := network.New()

	ctx := evmante.BuildEvmExecutionCtx(network.GetContext())

	suite.Equal(storetypes.GasConfig{}, ctx.KVGasConfig())
	suite.Equal(storetypes.GasConfig{}, ctx.TransientKVGasConfig())
}
