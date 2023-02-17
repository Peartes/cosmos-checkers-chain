package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPlayGame = "play_game"

var _ sdk.Msg = &MsgPlayGame{}

func NewMsgPlayGame(creator string, gameIndex string, fromX uint64, fromY uint64, toX uint64, toY uint64) *MsgPlayGame {
	return &MsgPlayGame{
		Creator:   creator,
		GameIndex: gameIndex,
		FromX:     fromX,
		FromY:     fromY,
		ToX:       toX,
		ToY:       toY,
	}
}

func (msg *MsgPlayGame) Route() string {
	return RouterKey
}

func (msg *MsgPlayGame) Type() string {
	return TypeMsgPlayGame
}

func (msg *MsgPlayGame) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPlayGame) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPlayGame) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
