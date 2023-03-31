package keeper

import (
	"fmt"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	if storedGame.MoveCount == 0 {
		// black is to play
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrBlackCannotPay.Error())
		}
	} else if storedGame.MoveCount == 1 {
		// red is to play
		red, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrRedAlreadyPlayed.Error())
		}
	}
	return nil
}
func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {
	winner, found, err := storedGame.GetWinnerAddress()
	if err != nil {
		panic(err.Error())
	}
	if !found {
		panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
	}

	winnings := storedGame.GetWagerCoin()
	if storedGame.MoveCount == 0 {
		panic(types.ErrNothingToPay.Error())
	} else if storedGame.MoveCount > 1 {
		winnings = winnings.Add(winnings)
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winner, sdk.NewCoins(winnings))
	if err != nil {
		panic(fmt.Sprintf(types.ErrCannotPayWinnings.Error(), err.Error()))
	}
}
func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {
	if storedGame.MoveCount == 1 {
		// Refund the black player
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), err.Error()))
		}
	} else if storedGame.MoveCount == 0 {

	} else {
		// TODO: Implement a DRAW mechanism
		panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
	}
}
