package keeper_test

import (
	"testing"
	"time"

	"github.com/alice/checkers/testutil"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPlayMoveUpToWinner(t *testing.T) {
	msgServer, k, ctx, ctrl, escrow := SetupMsgServerWithOneGameForPlayMove(t)
	sdkContext := sdk.UnwrapSDKContext(ctx)
	
	defer ctrl.Finish()

	payBob := escrow.ExpectPay(ctx, bob, 45).AnyTimes()
    payCarol := escrow.ExpectPay(ctx, carol, 45).After(payBob).AnyTimes()
    escrow.ExpectRefund(ctx, bob, 90).Times(1).After(payCarol)


	testutil.PlayAllMoves(t, msgServer, ctx, "1",  testutil.Game1Moves)

	systemInfo, found := k.GetSystemInfo(sdkContext)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId: 2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)

	game, found := k.GetStoredGame(sdkContext, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(len(testutil.Game1Moves)),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Deadline:    types.FormatDeadline(time.Time(sdkContext.BlockTime().Add(types.MaxTurnDuration))),
		Winner:      "b",
		Wager:        45,
	}, game)
	

	events := sdk.StringifyEvents(sdkContext.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.Equal(t, event.Type, "move-played")
	require.EqualValues(t, []sdk.Attribute{
        {Key: "creator", Value: bob},
        {Key: "game-index", Value: "1"},
        {Key: "captured-x", Value: "2"},
        {Key: "captured-y", Value: "5"},
        {Key: "winner", Value: "b"},
        {Key: "board", Value: "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"},
    }, event.Attributes[(len(testutil.Game1Moves)-1)*6:])
}
func TestPlayMoveUpToWinnerCalledBank(t *testing.T) {
    msgServer, _, context, ctrl, escrow := SetupMsgServerWithOneGameForPlayMove(t)
    defer ctrl.Finish()
    payBob := escrow.ExpectPay(context, bob, 45).AnyTimes()
    payCarol := escrow.ExpectPay(context, carol, 45).After(payBob).AnyTimes()
    escrow.ExpectRefund(context, bob, 90).Times(1).After(payCarol)
	
    testutil.PlayAllMoves(t, msgServer, context, "1", testutil.Game1Moves)

}
