package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayGame(goCtx context.Context, msg *types.MsgPlayGame) (*types.MsgPlayGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handling the message
	// Fetch the stored game at index msg.GameIndex
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	// Check if the player is legitimate
	var player rules.Player
	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	if !isBlack && !isRed {
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	} else if isBlack && isRed {
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isBlack {
		player = rules.BLACK_PLAYER
	} else {
		player = rules.RED_PLAYER
	}

	// Initialize the board
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}
	
	// Is it player turn ?
	if !game.TurnIs(player) {
		return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
	}

	// Make the move
	captured, err := game.Move(rules.Pos{
		X: int(msg.FromX),
		Y: int(msg.FromY),
	}, rules.Pos{
		X: int(msg.ToX),
		Y: int(msg.ToY),
	})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, err.Error())
	}

	storedGame.Board = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	return &types.MsgPlayGameResponse {
		CapturedX: int64(captured.X),
		CapturedY: int64(captured.Y),
		Winner: rules.PieceStrings[game.Winner()],
	}, nil
}
