package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteGame = "delete_game"

var _ sdk.Msg = &MsgDeleteGame{}

func NewMsgDeleteGame(creator string, index string) *MsgDeleteGame {
	return &MsgDeleteGame{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteGame) Route() string {
	return RouterKey
}

func (msg *MsgDeleteGame) Type() string {
	return TypeMsgDeleteGame
}

func (msg *MsgDeleteGame) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteGame) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteGame) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
