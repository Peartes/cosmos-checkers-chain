package types_test

import (
	"strings"
	"testing"
	"time"

	"github.com/alice/checkers/testutil"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black:       alice,
		Red:         bob,
		Index:       "1",
		Board:       rules.New().String(),
		Turn:        "b",
		MoveCount:   0,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.DeadlineLayout,
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	// BLACK PLAYER SHOULD BE ALICE ADDRESS
	aliceAddress, err1 := sdk.AccAddressFromBech32(alice)
	require.Nil(t, err1)

	blackAddress, err2 := GetStoredGame1().GetBlackAddress()
	require.Nil(t, err2)

	require.EqualValues(t, aliceAddress, blackAddress)
}

func TestGetAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4" // Bad last digit
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(t,
		err,
		"black address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCanGetAddressRed(t *testing.T) {
	// BLACK PLAYER SHOULD BE ALICE ADDRESS
	bobAdress, err1 := sdk.AccAddressFromBech32(bob)
	require.Nil(t, err1)

	redAddress, err2 := GetStoredGame1().GetRedAddress()
	require.Nil(t, err2)

	require.EqualValues(t, bobAdress, redAddress)
}

func TestGetAddressWrongRed(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4" // Bad last digit
	red, err := storedGame.GetRedAddress()
	require.Nil(t, red)
	require.EqualError(t,
		err,
		"red address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCanParseGame(t *testing.T) {
	game, err := GetStoredGame1().ParseGame()
	require.Nil(t, err)
	require.EqualValues(t, rules.New().Pieces, game.Pieces)

	storedGame := GetStoredGame1()
	// CHANGE THE GAME INDEX, SHOULD STILL PARSE RIGHTLY
	storedGame.Index = "2"
	_, err2 := storedGame.ParseGame()
	require.Nil(t, err2)
	// CHANGE THE PIECE COLORS, SHOULD STILL PARSE RIGHT
	storedGame.Board = strings.Replace(storedGame.Board, "r", "b", -1)
	_, err3 := storedGame.ParseGame()
	require.Nil(t, err3)
	require.NotEqualValues(t, storedGame.Board, rules.New().String())

	// CHANGE THE TURN OF PLAYER, PARSING WILL FAIL
	storedGame.Turn = "r"

	_, err4 := storedGame.ParseGame()
	require.Nil(t, err4)
}

func TestParseDeadlineCorrect(t *testing.T) {
	deadline, err := GetStoredGame1().GetDeadlineAsTime()
	require.Nil(t, err)
	require.Equal(t, time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.UTC)), deadline)
}

func TestParseDeadlineMissingMonth(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Deadline = "2006-02 15:04:05.999999999 +0000 UTC"
	_, err := storedGame.GetDeadlineAsTime()
	require.EqualError(t,
		err,
		"deadline cannot be parsed: 2006-02 15:04:05.999999999 +0000 UTC: parsing time \"2006-02 15:04:05.999999999 +0000 UTC\" as \"2006-01-02 15:04:05.999999999 +0000 UTC\": cannot parse \" 15:04:05.999999999 +0000 UTC\" as \"-\"")
	require.EqualError(t, storedGame.Validate(), err.Error())
}