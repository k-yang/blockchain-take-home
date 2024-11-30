package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blog/x/blog/types"
)

func (k msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	val, found := k.GetPost(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	// Check if the message sender is authorized
	isAuthorized := msg.Creator == val.Creator
	for _, updater := range val.GrantedUpdaters {
		if msg.Creator == updater {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "you do not have permission to delete this post")
	}
	k.RemovePost(ctx, msg.Id)
	return &types.MsgDeletePostResponse{}, nil
}
