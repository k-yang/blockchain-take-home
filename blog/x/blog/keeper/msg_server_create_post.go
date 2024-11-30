package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"blog/x/blog/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	post := types.Post{
		Creator:       msg.Creator,
		Title:         msg.Title,
		Body:          msg.Body,
		CreatedAt:     sdk.UnwrapSDKContext(ctx).BlockTime().UTC().Format(time.RFC3339),
		LastUpdatedAt: sdk.UnwrapSDKContext(ctx).BlockTime().UTC().Format(time.RFC3339),
	}
	id := k.AppendPost(
		ctx,
		post,
	)
	return &types.MsgCreatePostResponse{
		Id: id,
	}, nil
}
