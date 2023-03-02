package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func SetupMsgServerWithOneGameForPlayMove(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	return server, *k, context
}

func TestPlayMove(t *testing.T) {
	msgServer, _, context := SetupMsgServerWithOneGameForPlayMove(t)
	playMoveResponse, err := msgServer.PlayGame(context, &types.MsgPlayGame{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgPlayGameResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    "*",
	}, *playMoveResponse)
}

func TestPlayMoveSavedGame(t *testing.T) {
	msgServer, k, context := SetupMsgServerWithOneGameForPlayMove(t)

	// Make a move as black
	playGameResponse, err := msgServer.PlayGame(context, &types.MsgPlayGame{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgPlayGameResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    "*",
	}, playGameResponse)

	// Get the stored game
	storedGame, err := k.StoredGame(context, &types.QueryGetStoredGameRequest{
		Index: "1",
	})
	require.Nil(t, err)
	// Make sure the game is stored rightly
	require.EqualValues(t, &types.QueryGetStoredGameResponse{
		StoredGame: types.StoredGame{
			Index: "1",
			Turn:  "r",
			Black: bob,
			Red:   carol,
			Board: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
			MoveCount: 1,
			BeforeIndex: types.NoFifoIndex,
			AfterIndex: types.NoFifoIndex,
			Deadline: types.FormatDeadline(types.GetNextDeadline(sdk.UnwrapSDKContext(context))),
			Winner: "*",
		},
	}, storedGame)
}

func TestPlayGame1Emitted(t *testing.T) {
	msgServer, k, context := SetupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayGame(context, &types.MsgPlayGame{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	ctx := sdk.UnwrapSDKContext(context)
	require.NotNil(t, ctx)
	storedGame, _ := k.GetStoredGame(ctx, "1")
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "move-played",
		Attributes: []sdk.Attribute{
			{Key: "creator", Value: bob},
			{Key: "game-index", Value: "1"},
			{Key: "captured-x", Value: "-1"},
			{Key: "captured-y", Value: "-1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: storedGame.Board},
		},
	}, event)
}

func TestPlayGame2Emitted(t *testing.T) {
	msgServer, k, context := SetupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayGame(context, &types.MsgPlayGame{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayGame(context, &types.MsgPlayGame{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	ctx := sdk.UnwrapSDKContext(context)
	storedGame, _ := k.GetStoredGame(ctx, "1")
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.Equal(t, "move-played", event.Type)
	require.EqualValues(t, []sdk.Attribute{
		{Key: "creator", Value: carol},
		{Key: "game-index", Value: "1"},
		{Key: "captured-x", Value: "-1"},
		{Key: "captured-y", Value: "-1"},
		{Key: "winner", Value: "*"},
		{Key: "board", Value:  storedGame.Board},
	}, event.Attributes[6:])
}
