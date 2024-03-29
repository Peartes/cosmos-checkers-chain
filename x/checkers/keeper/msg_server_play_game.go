package keeper

import (
	"context"
	"strconv"

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

	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
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

	err = k.CollectWager(ctx, &storedGame)
	if err != nil {
		return nil, err
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

	storedGame.Winner = rules.PieceStrings[game.Winner()]
	lastBoard := game.String()

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if (!found) {
		panic("SystemInfo not found")
	}

	if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
		k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)
		storedGame.Board = lastBoard
	} else {
		k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)
		storedGame.Board = ""
		k.Keeper.MustPayWinnings(ctx, &storedGame)
	}
	storedGame.MoveCount++
	storedGame.Turn = rules.PieceStrings[game.Turn]
	storedGame.Deadline = types.FormatDeadline(types.GetNextDeadline(ctx))

	k.Keeper.SetStoredGame(ctx, storedGame)
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.MovePlayedEventType,
			sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
			sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
			sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
			sdk.NewAttribute(types.MovePlayedEventBoard, lastBoard),
		),
	)

	return &types.MsgPlayGameResponse{
		CapturedX: int64(captured.X),
		CapturedY: int64(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil
}
