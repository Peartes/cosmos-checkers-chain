package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrInvalidBlack     = sdkerrors.Register(ModuleName, 1100, "black address is invalid: %s")
	ErrInvalidRed       = sdkerrors.Register(ModuleName, 1101, "red address is invalid: %s")
	ErrGameNotParseable = sdkerrors.Register(ModuleName, 1102, "game cannot be parsed")
	ErrGameNotFound = sdkerrors.Register(ModuleName, 1103, "game at index not found: %s")
	ErrUnAuthorizedOperation = sdkerrors.Register(ModuleName, 1104, "creator has initialized an unauthorized tx")
)
