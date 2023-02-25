package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestReject1GameNoPlay(t *testing.T) {
	msgServer, k, ctx := SetupMsgServerWithOneGameForPlayMove(t)

	// Reject the game
	_, err := msgServer.RejectGame(ctx, &types.MsgRejectGame{
		Creator: bob,
		GameIndex: "1",
	})
	require.Nil(t, err)

	sdkContext := sdk.UnwrapSDKContext(ctx)
	// Load all games
	storedGames := k.GetAllStoredGame(sdkContext)
	require.Len(t, storedGames, 0)

	// Check emitted events
	events := sdk.StringifyEvents(sdkContext.EventManager().ABCIEvents())
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameRejectedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameRejectedEventCreator, Value: bob},
			{Key: types.GameRejectedEventGameIndex, Value: "1"},
		},
	}, events[0])
}

func TestGameRejected1GameBlackMove(t *testing.T) {
	msgServer, k, ctx := SetupMsgServerWithOneGameForPlayMove(t)
	
	// Let black make a move
	_, playErr := msgServer.PlayGame(ctx, &types.MsgPlayGame{
		Creator: bob,
		GameIndex: "1",
		FromX: 1,
		FromY: 2,
		ToX: 2,
		ToY: 3,
	})
	require.Nil(t, playErr)

	// Black reject a game
	_, rejectErr := msgServer.RejectGame(ctx, &types.MsgRejectGame{
		Creator: bob,
		GameIndex: "1",
	})
	require.Errorf(t, rejectErr, "black player has already played")

	// Red reject is fine
	_, redRejectErr := msgServer.RejectGame(ctx, &types.MsgRejectGame{
		Creator: carol,
		GameIndex: "1",
	})
	require.Nil(t, redRejectErr)

	// Make sure no more stored game
	sdkContext := sdk.UnwrapSDKContext(ctx)

	storedGames := k.GetAllStoredGame(sdkContext)
	require.Len(t, storedGames, 0)

	events := sdk.StringifyEvents(sdkContext.EventManager().ABCIEvents())
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameRejectedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameRejectedEventCreator, Value: carol},
			{Key: types.GameRejectedEventGameIndex, Value: "1"},
		},
	}, events[0])
}