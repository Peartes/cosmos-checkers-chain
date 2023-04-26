package keeper_test

import (
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) setupSuiteWithOneGameForPlayMove() {
	suite.setupSuiteWithBalances()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
		Wager:   45,
 	})
}

func (suite *IntegrationTestSuite) TestCreate1GameHasSaved() {
	suite.setupSuiteWithOneGameForPlayMove()
	keeper := suite.app.CheckersKeeper

	systemInfo, found := keeper.GetSystemInfo(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "1",
		FifoTailIndex: "1",
	}, systemInfo)
}

func (suite *IntegrationTestSuite) TestPlayMovePlayerPaid() {
	suite.setupSuiteWithOneGameForPlayMove()

	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)

	playGameResponse, err := suite.msgServer.PlayGame(sdk.WrapSDKContext(suite.ctx), &types.MsgPlayGame{
		Creator:   carol,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	suite.Require().Nil(err)
	suite.Require().EqualValues(types.MsgPlayGameResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner: "*",
	}, *playGameResponse)

	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol-45, carol)
	suite.RequireBankBalance(45, checkersModuleAddress)
}
