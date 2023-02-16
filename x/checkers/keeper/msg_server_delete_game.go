package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteGame(goCtx context.Context, msg *types.MsgDeleteGame) (*types.MsgDeleteGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handling the message
	// Get the game from storage
	game, found := k.Keeper.GetStoredGame(ctx, msg.Index)
	if !found {
		return nil, types.ErrGameNotFound
	}
	// make sure creator is red or black player
	if  game.Black != msg.Creator && game.Red != msg.Creator {
		return nil, types.ErrUnAuthorizedOperation
	}
	// delete the game 
	k.Keeper.RemoveStoredGame(ctx, msg.Index)

	return &types.MsgDeleteGameResponse{
		GameIndex: msg.Index,
	}, nil
}
