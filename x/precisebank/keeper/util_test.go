package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
)

// MintToAccount mints coins to an account with the x/precisebank methods. This
// must be used when minting extended coins, ie. aatom coins. This depends on
// the methods to be properly tested to be implemented correctly.
func (suite *KeeperIntegrationTestSuite) MintToAccount(addr sdk.AccAddress, amt sdk.Coins) {
	accBalancesBefore := suite.GetAllBalances(addr)

	err := suite.network.App.PreciseBankKeeper.MintCoins(suite.network.GetContext(), minttypes.ModuleName, amt)
	suite.Require().NoError(err)

	err = suite.network.App.PreciseBankKeeper.SendCoinsFromModuleToAccount(suite.network.GetContext(), minttypes.ModuleName, addr, amt)
	suite.Require().NoError(err)

	// Double check balances are correctly minted and sent to account
	accBalancesAfter := suite.GetAllBalances(addr)

	netIncrease := accBalancesAfter.Sub(accBalancesBefore...)
	suite.Require().Equal(ConvertCoinsToExtendedCoinDenom(amt), netIncrease)

	suite.T().Logf("minted %s to %s", amt, addr)
}

// GetAllBalances returns all the account balances for the given account address.
// This returns the extended coin balance if the account has a non-zero balance,
// WITHOUT the integer coin balance.
func (suite *KeeperIntegrationTestSuite) GetAllBalances(addr sdk.AccAddress) sdk.Coins {
	// Get all balances for an account
	bankBalances := suite.network.App.BankKeeper.GetAllBalances(suite.network.GetContext(), addr)

	// Remove integer coins from the balance
	for _, coin := range bankBalances {
		if coin.Denom == precisebanktypes.IntegerCoinDenom {
			bankBalances = bankBalances.Sub(coin)
		}
	}

	// Replace the integer coin with the extended coin, from x/precisebank
	extendedBal := suite.network.App.PreciseBankKeeper.GetBalance(suite.network.GetContext(), addr, precisebanktypes.ExtendedCoinDenom)

	return bankBalances.Add(extendedBal)
}

// ConvertCoinsToExtendedCoinDenom converts sdk.Coins that includes Integer denoms
// to sdk.Coins that includes Extended denoms of the same amount. This is useful
// for testing to make sure only extended amounts are compared instead of double
// counting balances.
func ConvertCoinsToExtendedCoinDenom(coins sdk.Coins) sdk.Coins {
	integerCoinAmt := coins.AmountOf(precisebanktypes.IntegerCoinDenom)
	if integerCoinAmt.IsZero() {
		return coins
	}

	// Remove the integer coin from the coins
	integerCoin := sdk.NewCoin(precisebanktypes.IntegerCoinDenom, integerCoinAmt)

	// Add the equivalent extended coin to the coins
	extendedCoin := sdk.NewCoin(precisebanktypes.ExtendedCoinDenom, integerCoinAmt.Mul(precisebanktypes.ConversionFactor()))

	return coins.Sub(integerCoin).Add(extendedCoin)
}
