package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blog/x/blog/types"
)

func (k msgServer) UpdatePost(goCtx context.Context, msg *types.MsgUpdatePost) (*types.MsgUpdatePostResponse, error) {
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
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "you do not have permission to update this post")
	}
	// Update only the mutable fields
	val.Title = msg.Title
	val.Body = msg.Body
	val.LastUpdatedAt = sdk.UnwrapSDKContext(ctx).BlockTime().UTC().Format(time.RFC3339) // Update last_updated_at

	k.SetPost(ctx, val)
	return &types.MsgUpdatePostResponse{}, nil
}
