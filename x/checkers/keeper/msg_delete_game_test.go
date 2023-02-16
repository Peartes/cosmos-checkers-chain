package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestCanDelete1Game(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServerCreateGame(t)

	// create a game
	gameResponse, err := msgServer.CreateGame(ctx, &types.MsgCreateGame{
		Creator: alice,
		Black: alice,
		Red: bob,
	})
	require.Nil(t, err)
	// make sure game is created and stored
	storedGames, err := keeper.StoredGameAll(ctx, &types.QueryAllStoredGameRequest{})
	require.Nil(t, err)
	require.EqualValues(t, len(storedGames.StoredGame), 1)

	// delete the newly created game
	deletedResponse, err := msgServer.DeleteGame(ctx, &types.MsgDeleteGame{
		Creator: alice,
		Index: gameResponse.GameIndex,
	})
	require.Nil(t, err)
	
	// make sure we have no more games in store
	afterStoredGames, err := keeper.StoredGameAll(ctx, &types.QueryAllStoredGameRequest{})
	require.Nil(t, err)
	require.EqualValues(t, deletedResponse.GameIndex, gameResponse.GameIndex)
	require.EqualValues(t, len(afterStoredGames.StoredGame), 0)
	// require.Errorf(t, err, "creator has initialized an unauthorized tx")

}