package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCreate3GamesHaveSavedFifo(t *testing.T) {
	msgServer, k, ctx := setupMsgServerCreateGame(t)
	sdkContext := sdk.UnwrapSDKContext(ctx)

	msgServer.CreateGame(ctx,  &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})

	// Second game
	msgServer.CreateGame(ctx, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
	
	systemInfo2, found := k.GetSystemInfo(sdkContext)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "1",
		FifoTailIndex: "2",
	}, systemInfo2)

	game1, found := k.GetStoredGame(sdkContext, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(0),
		BeforeIndex: "-1",
		AfterIndex:  "2",
		Deadline: types.FormatDeadline(types.GetNextDeadline(sdkContext)),
		Winner: "*",
	}, game1)
	game2, found := k.GetStoredGame(sdkContext, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		MoveCount:   uint64(0),
		BeforeIndex: "1",
		AfterIndex:  "-1",
		Deadline: types.FormatDeadline(types.GetNextDeadline(sdkContext)),
		Winner: "*",
	}, game2)

	// Third game
	msgServer.CreateGame(ctx, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
	})
	systemInfo3, found := k.GetSystemInfo(sdkContext)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "1",
		FifoTailIndex: "3",
	}, systemInfo3)
	game1, found = k.GetStoredGame(sdkContext, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(0),
		BeforeIndex: "-1",
		AfterIndex:  "2",
		Deadline: types.FormatDeadline(types.GetNextDeadline(sdkContext)),
		Winner: "*",
	}, game1)
	game2, found = k.GetStoredGame(sdkContext, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		MoveCount:   uint64(0),
		BeforeIndex: "1",
		AfterIndex:  "3",
		Deadline: types.FormatDeadline(types.GetNextDeadline(sdkContext)),
		Winner: "*",
	}, game2)
	game3, found := k.GetStoredGame(sdkContext, "3")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "3",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       alice,
		Red:         bob,
		MoveCount:   uint64(0),
		BeforeIndex: "2",
		AfterIndex:  "-1",
		Deadline: types.FormatDeadline(types.GetNextDeadline(sdkContext)),
		Winner: "*",
	}, game3)
}