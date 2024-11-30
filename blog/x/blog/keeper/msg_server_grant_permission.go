package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blog/x/blog/types"
)

func (k msgServer) GrantPermission(goCtx context.Context, msg *types.MsgGrantPermission) (*types.MsgGrantPermissionResponse, error) {
	// Get the current context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve the post
	post, found := k.GetPost(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("post with ID %d not found", msg.Id))
	}

	// Ensure the message sender is the post creator
	if msg.Creator != post.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the creator can grant permissions")
	}

	// Check if the grantee is already in the list
	for _, updater := range post.GrantedUpdaters {
		if updater == msg.Grantee {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "address already has permissions")
		}
	}

	// Add the new grantee to the list
	post.GrantedUpdaters = append(post.GrantedUpdaters, msg.Grantee)
	k.SetPost(ctx, post)

	return &types.MsgGrantPermissionResponse{}, nil
}
