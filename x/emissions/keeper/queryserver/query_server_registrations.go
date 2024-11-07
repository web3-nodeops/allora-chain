package queryserver

import (
	"context"
	"time"

	"github.com/allora-network/allora-chain/x/emissions/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-chain/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (qs queryServer) GetWorkerNodeInfo(ctx context.Context, req *types.GetWorkerNodeInfoRequest) (_ *types.GetWorkerNodeInfoResponse, err error) {
	defer metrics.RecordMetrics("GetWorkerNodeInfo", time.Now(), &err)

	node, err := qs.k.GetWorkerInfo(sdk.UnwrapSDKContext(ctx), req.Address)
	if err != nil {
		return nil, err
	}

	return &types.GetWorkerNodeInfoResponse{NodeInfo: &node}, nil
}

func (qs queryServer) GetReputerNodeInfo(ctx context.Context, req *types.GetReputerNodeInfoRequest) (_ *types.GetReputerNodeInfoResponse, err error) {
	defer metrics.RecordMetrics("GetReputerNodeInfo", time.Now(), &err)

	node, err := qs.k.GetReputerInfo(sdk.UnwrapSDKContext(ctx), req.Address)
	if err != nil {
		return nil, err
	}

	return &types.GetReputerNodeInfoResponse{NodeInfo: &node}, nil
}

func (qs queryServer) IsWorkerRegisteredInTopicId(ctx context.Context, req *types.IsWorkerRegisteredInTopicIdRequest) (_ *types.IsWorkerRegisteredInTopicIdResponse, err error) {
	defer metrics.RecordMetrics("IsWorkerRegisteredInTopicId", time.Now(), &err)

	if err := qs.k.ValidateStringIsBech32(req.Address); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", err)
	}

	topicExists, err := qs.k.TopicExists(ctx, req.TopicId)
	if !topicExists {
		return nil, status.Errorf(codes.NotFound, "topic %v not found", req.TopicId)
	} else if err != nil {
		return nil, err
	}

	isRegistered, err := qs.k.IsWorkerRegisteredInTopic(sdk.UnwrapSDKContext(ctx), req.TopicId, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.IsWorkerRegisteredInTopicIdResponse{IsRegistered: isRegistered}, nil
}

func (qs queryServer) IsReputerRegisteredInTopicId(ctx context.Context, req *types.IsReputerRegisteredInTopicIdRequest) (_ *types.IsReputerRegisteredInTopicIdResponse, err error) {
	defer metrics.RecordMetrics("IsReputerRegisteredInTopicId", time.Now(), &err)

	if err := qs.k.ValidateStringIsBech32(req.Address); err != nil {
		return nil, err
	}
	isRegistered, err := qs.k.IsReputerRegisteredInTopic(sdk.UnwrapSDKContext(ctx), req.TopicId, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.IsReputerRegisteredInTopicIdResponse{IsRegistered: isRegistered}, nil
}
