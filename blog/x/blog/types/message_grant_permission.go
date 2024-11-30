package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgGrantPermission{}

// NewMsgGrantPermission creates a new MsgGrantPermission instance.
func NewMsgGrantPermission(creator string, id uint64, grantee string) *MsgGrantPermission {
	return &MsgGrantPermission{
		Creator: creator,
		Id:      id,
		Grantee: grantee,
	}
}

// ValidateBasic performs basic validation of MsgGrantPermission fields.
func (msg *MsgGrantPermission) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate grantee address
	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	// Check that creator and grantee are different
	if msg.Creator == msg.Grantee {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "creator and grantee cannot be the same")
	}

	return nil
}

// GetSigners returns the required signers for MsgGrantPermission.
func (msg *MsgGrantPermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
