package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPlayMove2Games2MovesHasSavedFifo(t *testing.T) {
    msgServer, keeper, context := SetupMsgServerWithOneGameForPlayMove(t)
    ctx := sdk.UnwrapSDKContext(context)
    msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: bob,
        Black:   carol,
        Red:     alice,
    })
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
        GameIndex: "2",
        FromX:     1,
        FromY:     2,
        ToX:       2,
        ToY:       3,
    })
    systemInfo1, found := keeper.GetSystemInfo(ctx)
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId:        3,
        FifoHeadIndex: "1",
        FifoTailIndex: "2",
    }, systemInfo1)
    game1, found := keeper.GetStoredGame(ctx, "1")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "1",
        Board:       "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "r",
        Black:       bob,
        Red:         carol,
        MoveCount:   uint64(1),
        BeforeIndex: "-1",
        AfterIndex:  "2",
		Winner: "*",
		Deadline: types.FormatDeadline(types.GetNextDeadline(ctx)),
    }, game1)
    game2, found := keeper.GetStoredGame(ctx, "2")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "2",
        Board:       "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "r",
        Black:       carol,
        Red:         alice,
        MoveCount:   uint64(1),
        BeforeIndex: "1",
        AfterIndex:  "-1",
		Deadline: types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner: "*",
    }, game2)
}
