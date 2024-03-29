package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRejectSecondGameHasSavedFifo(t *testing.T) {
	msgServer, keeper, context, _, _ := SetupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "2",
		FifoTailIndex: "2",
	}, systemInfo)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		MoveCount:   uint64(0),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Deadline: types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner: "*",
	}, game2)
}

func TestRejectMiddleGameHasSavedFifo(t *testing.T) {
    msgServer, keeper, context, _, _ := SetupMsgServerWithOneGameForPlayMove(t)
    ctx := sdk.UnwrapSDKContext(context)
    msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: bob,
        Black:   carol,
        Red:     alice,
        Wager:  0,
    })
    msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: carol,
        Black:   alice,
        Red:     bob,
        Wager:  0,
    })
    msgServer.RejectGame(context, &types.MsgRejectGame{
        Creator:   carol,
        GameIndex: "2",
    })
    systemInfo, found := keeper.GetSystemInfo(ctx)
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId:        4,
        FifoHeadIndex: "1",
        FifoTailIndex: "3",
    }, systemInfo)
    game1, found := keeper.GetStoredGame(ctx, "1")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "1",
        Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "b",
        Black:       bob,
        Red:         carol,
        MoveCount:   uint64(0),
        BeforeIndex: "-1",
        AfterIndex:  "3",
		Deadline: types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner: "*",
        Wager:  45,
    }, game1)
    game3, found := keeper.GetStoredGame(ctx, "3")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "3",
        Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "b",
        Black:       alice,
        Red:         bob,
        MoveCount:   uint64(0),
        BeforeIndex: "1",
        AfterIndex:  "-1",
		Deadline: types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner: "*",
        Wager:  0,
    }, game3)
}
