package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgGrantPermission_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgGrantPermission
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgGrantPermission{
				Creator: "invalid_address",
				Id:      1,
				Grantee: "cosmos1granteeaddressvalid",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid grantee address",
			msg: MsgGrantPermission{
				Creator: "cosmos1creatoraddressvalid",
				Id:      1,
				Grantee: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "creator and grantee are the same",
			msg: MsgGrantPermission{
				Creator: "cosmos1creatoraddressvalid",
				Id:      1,
				Grantee: "cosmos1creatoraddressvalid",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: MsgGrantPermission{
				Creator: "cosmos1creatoraddressvalid",
				Id:      1,
				Grantee: "cosmos1granteeaddressvalid",
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
