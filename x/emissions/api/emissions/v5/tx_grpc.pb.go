// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: emissions/v5/tx.proto

package emissionsv5

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MsgService_UpdateParams_FullMethodName              = "/emissions.v5.MsgService/UpdateParams"
	MsgService_CreateNewTopic_FullMethodName            = "/emissions.v5.MsgService/CreateNewTopic"
	MsgService_Register_FullMethodName                  = "/emissions.v5.MsgService/Register"
	MsgService_RemoveRegistration_FullMethodName        = "/emissions.v5.MsgService/RemoveRegistration"
	MsgService_AddStake_FullMethodName                  = "/emissions.v5.MsgService/AddStake"
	MsgService_RemoveStake_FullMethodName               = "/emissions.v5.MsgService/RemoveStake"
	MsgService_CancelRemoveStake_FullMethodName         = "/emissions.v5.MsgService/CancelRemoveStake"
	MsgService_DelegateStake_FullMethodName             = "/emissions.v5.MsgService/DelegateStake"
	MsgService_RewardDelegateStake_FullMethodName       = "/emissions.v5.MsgService/RewardDelegateStake"
	MsgService_RemoveDelegateStake_FullMethodName       = "/emissions.v5.MsgService/RemoveDelegateStake"
	MsgService_CancelRemoveDelegateStake_FullMethodName = "/emissions.v5.MsgService/CancelRemoveDelegateStake"
	MsgService_FundTopic_FullMethodName                 = "/emissions.v5.MsgService/FundTopic"
	MsgService_AddToWhitelistAdmin_FullMethodName       = "/emissions.v5.MsgService/AddToWhitelistAdmin"
	MsgService_RemoveFromWhitelistAdmin_FullMethodName  = "/emissions.v5.MsgService/RemoveFromWhitelistAdmin"
	MsgService_InsertWorkerPayload_FullMethodName       = "/emissions.v5.MsgService/InsertWorkerPayload"
	MsgService_InsertReputerPayload_FullMethodName      = "/emissions.v5.MsgService/InsertReputerPayload"
)

// MsgServiceClient is the client API for MsgService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgServiceClient interface {
	UpdateParams(ctx context.Context, in *UpdateParamsRequest, opts ...grpc.CallOption) (*UpdateParamsResponse, error)
	CreateNewTopic(ctx context.Context, in *CreateNewTopicRequest, opts ...grpc.CallOption) (*CreateNewTopicResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	RemoveRegistration(ctx context.Context, in *RemoveRegistrationRequest, opts ...grpc.CallOption) (*RemoveRegistrationResponse, error)
	AddStake(ctx context.Context, in *AddStakeRequest, opts ...grpc.CallOption) (*AddStakeResponse, error)
	RemoveStake(ctx context.Context, in *RemoveStakeRequest, opts ...grpc.CallOption) (*RemoveStakeResponse, error)
	CancelRemoveStake(ctx context.Context, in *CancelRemoveStakeRequest, opts ...grpc.CallOption) (*CancelRemoveStakeResponse, error)
	DelegateStake(ctx context.Context, in *DelegateStakeRequest, opts ...grpc.CallOption) (*DelegateStakeResponse, error)
	RewardDelegateStake(ctx context.Context, in *RewardDelegateStakeRequest, opts ...grpc.CallOption) (*RewardDelegateStakeResponse, error)
	RemoveDelegateStake(ctx context.Context, in *RemoveDelegateStakeRequest, opts ...grpc.CallOption) (*RemoveDelegateStakeResponse, error)
	CancelRemoveDelegateStake(ctx context.Context, in *CancelRemoveDelegateStakeRequest, opts ...grpc.CallOption) (*CancelRemoveDelegateStakeResponse, error)
	FundTopic(ctx context.Context, in *FundTopicRequest, opts ...grpc.CallOption) (*FundTopicResponse, error)
	AddToWhitelistAdmin(ctx context.Context, in *AddToWhitelistAdminRequest, opts ...grpc.CallOption) (*AddToWhitelistAdminResponse, error)
	RemoveFromWhitelistAdmin(ctx context.Context, in *RemoveFromWhitelistAdminRequest, opts ...grpc.CallOption) (*RemoveFromWhitelistAdminResponse, error)
	InsertWorkerPayload(ctx context.Context, in *InsertWorkerPayloadRequest, opts ...grpc.CallOption) (*InsertWorkerPayloadResponse, error)
	InsertReputerPayload(ctx context.Context, in *InsertReputerPayloadRequest, opts ...grpc.CallOption) (*InsertReputerPayloadResponse, error)
}

type msgServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgServiceClient(cc grpc.ClientConnInterface) MsgServiceClient {
	return &msgServiceClient{cc}
}

func (c *msgServiceClient) UpdateParams(ctx context.Context, in *UpdateParamsRequest, opts ...grpc.CallOption) (*UpdateParamsResponse, error) {
	out := new(UpdateParamsResponse)
	err := c.cc.Invoke(ctx, MsgService_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) CreateNewTopic(ctx context.Context, in *CreateNewTopicRequest, opts ...grpc.CallOption) (*CreateNewTopicResponse, error) {
	out := new(CreateNewTopicResponse)
	err := c.cc.Invoke(ctx, MsgService_CreateNewTopic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, MsgService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) RemoveRegistration(ctx context.Context, in *RemoveRegistrationRequest, opts ...grpc.CallOption) (*RemoveRegistrationResponse, error) {
	out := new(RemoveRegistrationResponse)
	err := c.cc.Invoke(ctx, MsgService_RemoveRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) AddStake(ctx context.Context, in *AddStakeRequest, opts ...grpc.CallOption) (*AddStakeResponse, error) {
	out := new(AddStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_AddStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) RemoveStake(ctx context.Context, in *RemoveStakeRequest, opts ...grpc.CallOption) (*RemoveStakeResponse, error) {
	out := new(RemoveStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_RemoveStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) CancelRemoveStake(ctx context.Context, in *CancelRemoveStakeRequest, opts ...grpc.CallOption) (*CancelRemoveStakeResponse, error) {
	out := new(CancelRemoveStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_CancelRemoveStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) DelegateStake(ctx context.Context, in *DelegateStakeRequest, opts ...grpc.CallOption) (*DelegateStakeResponse, error) {
	out := new(DelegateStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_DelegateStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) RewardDelegateStake(ctx context.Context, in *RewardDelegateStakeRequest, opts ...grpc.CallOption) (*RewardDelegateStakeResponse, error) {
	out := new(RewardDelegateStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_RewardDelegateStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) RemoveDelegateStake(ctx context.Context, in *RemoveDelegateStakeRequest, opts ...grpc.CallOption) (*RemoveDelegateStakeResponse, error) {
	out := new(RemoveDelegateStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_RemoveDelegateStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) CancelRemoveDelegateStake(ctx context.Context, in *CancelRemoveDelegateStakeRequest, opts ...grpc.CallOption) (*CancelRemoveDelegateStakeResponse, error) {
	out := new(CancelRemoveDelegateStakeResponse)
	err := c.cc.Invoke(ctx, MsgService_CancelRemoveDelegateStake_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) FundTopic(ctx context.Context, in *FundTopicRequest, opts ...grpc.CallOption) (*FundTopicResponse, error) {
	out := new(FundTopicResponse)
	err := c.cc.Invoke(ctx, MsgService_FundTopic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) AddToWhitelistAdmin(ctx context.Context, in *AddToWhitelistAdminRequest, opts ...grpc.CallOption) (*AddToWhitelistAdminResponse, error) {
	out := new(AddToWhitelistAdminResponse)
	err := c.cc.Invoke(ctx, MsgService_AddToWhitelistAdmin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) RemoveFromWhitelistAdmin(ctx context.Context, in *RemoveFromWhitelistAdminRequest, opts ...grpc.CallOption) (*RemoveFromWhitelistAdminResponse, error) {
	out := new(RemoveFromWhitelistAdminResponse)
	err := c.cc.Invoke(ctx, MsgService_RemoveFromWhitelistAdmin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) InsertWorkerPayload(ctx context.Context, in *InsertWorkerPayloadRequest, opts ...grpc.CallOption) (*InsertWorkerPayloadResponse, error) {
	out := new(InsertWorkerPayloadResponse)
	err := c.cc.Invoke(ctx, MsgService_InsertWorkerPayload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) InsertReputerPayload(ctx context.Context, in *InsertReputerPayloadRequest, opts ...grpc.CallOption) (*InsertReputerPayloadResponse, error) {
	out := new(InsertReputerPayloadResponse)
	err := c.cc.Invoke(ctx, MsgService_InsertReputerPayload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServiceServer is the server API for MsgService service.
// All implementations must embed UnimplementedMsgServiceServer
// for forward compatibility
type MsgServiceServer interface {
	UpdateParams(context.Context, *UpdateParamsRequest) (*UpdateParamsResponse, error)
	CreateNewTopic(context.Context, *CreateNewTopicRequest) (*CreateNewTopicResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	RemoveRegistration(context.Context, *RemoveRegistrationRequest) (*RemoveRegistrationResponse, error)
	AddStake(context.Context, *AddStakeRequest) (*AddStakeResponse, error)
	RemoveStake(context.Context, *RemoveStakeRequest) (*RemoveStakeResponse, error)
	CancelRemoveStake(context.Context, *CancelRemoveStakeRequest) (*CancelRemoveStakeResponse, error)
	DelegateStake(context.Context, *DelegateStakeRequest) (*DelegateStakeResponse, error)
	RewardDelegateStake(context.Context, *RewardDelegateStakeRequest) (*RewardDelegateStakeResponse, error)
	RemoveDelegateStake(context.Context, *RemoveDelegateStakeRequest) (*RemoveDelegateStakeResponse, error)
	CancelRemoveDelegateStake(context.Context, *CancelRemoveDelegateStakeRequest) (*CancelRemoveDelegateStakeResponse, error)
	FundTopic(context.Context, *FundTopicRequest) (*FundTopicResponse, error)
	AddToWhitelistAdmin(context.Context, *AddToWhitelistAdminRequest) (*AddToWhitelistAdminResponse, error)
	RemoveFromWhitelistAdmin(context.Context, *RemoveFromWhitelistAdminRequest) (*RemoveFromWhitelistAdminResponse, error)
	InsertWorkerPayload(context.Context, *InsertWorkerPayloadRequest) (*InsertWorkerPayloadResponse, error)
	InsertReputerPayload(context.Context, *InsertReputerPayloadRequest) (*InsertReputerPayloadResponse, error)
	mustEmbedUnimplementedMsgServiceServer()
}

// UnimplementedMsgServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServiceServer struct {
}

func (UnimplementedMsgServiceServer) UpdateParams(context.Context, *UpdateParamsRequest) (*UpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServiceServer) CreateNewTopic(context.Context, *CreateNewTopicRequest) (*CreateNewTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewTopic not implemented")
}
func (UnimplementedMsgServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedMsgServiceServer) RemoveRegistration(context.Context, *RemoveRegistrationRequest) (*RemoveRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveRegistration not implemented")
}
func (UnimplementedMsgServiceServer) AddStake(context.Context, *AddStakeRequest) (*AddStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStake not implemented")
}
func (UnimplementedMsgServiceServer) RemoveStake(context.Context, *RemoveStakeRequest) (*RemoveStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveStake not implemented")
}
func (UnimplementedMsgServiceServer) CancelRemoveStake(context.Context, *CancelRemoveStakeRequest) (*CancelRemoveStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelRemoveStake not implemented")
}
func (UnimplementedMsgServiceServer) DelegateStake(context.Context, *DelegateStakeRequest) (*DelegateStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelegateStake not implemented")
}
func (UnimplementedMsgServiceServer) RewardDelegateStake(context.Context, *RewardDelegateStakeRequest) (*RewardDelegateStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RewardDelegateStake not implemented")
}
func (UnimplementedMsgServiceServer) RemoveDelegateStake(context.Context, *RemoveDelegateStakeRequest) (*RemoveDelegateStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDelegateStake not implemented")
}
func (UnimplementedMsgServiceServer) CancelRemoveDelegateStake(context.Context, *CancelRemoveDelegateStakeRequest) (*CancelRemoveDelegateStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelRemoveDelegateStake not implemented")
}
func (UnimplementedMsgServiceServer) FundTopic(context.Context, *FundTopicRequest) (*FundTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FundTopic not implemented")
}
func (UnimplementedMsgServiceServer) AddToWhitelistAdmin(context.Context, *AddToWhitelistAdminRequest) (*AddToWhitelistAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToWhitelistAdmin not implemented")
}
func (UnimplementedMsgServiceServer) RemoveFromWhitelistAdmin(context.Context, *RemoveFromWhitelistAdminRequest) (*RemoveFromWhitelistAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFromWhitelistAdmin not implemented")
}
func (UnimplementedMsgServiceServer) InsertWorkerPayload(context.Context, *InsertWorkerPayloadRequest) (*InsertWorkerPayloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertWorkerPayload not implemented")
}
func (UnimplementedMsgServiceServer) InsertReputerPayload(context.Context, *InsertReputerPayloadRequest) (*InsertReputerPayloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertReputerPayload not implemented")
}
func (UnimplementedMsgServiceServer) mustEmbedUnimplementedMsgServiceServer() {}

// UnsafeMsgServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServiceServer will
// result in compilation errors.
type UnsafeMsgServiceServer interface {
	mustEmbedUnimplementedMsgServiceServer()
}

func RegisterMsgServiceServer(s grpc.ServiceRegistrar, srv MsgServiceServer) {
	s.RegisterService(&MsgService_ServiceDesc, srv)
}

func _MsgService_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).UpdateParams(ctx, req.(*UpdateParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_CreateNewTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNewTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).CreateNewTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_CreateNewTopic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).CreateNewTopic(ctx, req.(*CreateNewTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_RemoveRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).RemoveRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_RemoveRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).RemoveRegistration(ctx, req.(*RemoveRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_AddStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).AddStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_AddStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).AddStake(ctx, req.(*AddStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_RemoveStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).RemoveStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_RemoveStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).RemoveStake(ctx, req.(*RemoveStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_CancelRemoveStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelRemoveStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).CancelRemoveStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_CancelRemoveStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).CancelRemoveStake(ctx, req.(*CancelRemoveStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_DelegateStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelegateStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).DelegateStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_DelegateStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).DelegateStake(ctx, req.(*DelegateStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_RewardDelegateStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RewardDelegateStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).RewardDelegateStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_RewardDelegateStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).RewardDelegateStake(ctx, req.(*RewardDelegateStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_RemoveDelegateStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveDelegateStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).RemoveDelegateStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_RemoveDelegateStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).RemoveDelegateStake(ctx, req.(*RemoveDelegateStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_CancelRemoveDelegateStake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelRemoveDelegateStakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).CancelRemoveDelegateStake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_CancelRemoveDelegateStake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).CancelRemoveDelegateStake(ctx, req.(*CancelRemoveDelegateStakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_FundTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FundTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).FundTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_FundTopic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).FundTopic(ctx, req.(*FundTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_AddToWhitelistAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToWhitelistAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).AddToWhitelistAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_AddToWhitelistAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).AddToWhitelistAdmin(ctx, req.(*AddToWhitelistAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_RemoveFromWhitelistAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFromWhitelistAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).RemoveFromWhitelistAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_RemoveFromWhitelistAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).RemoveFromWhitelistAdmin(ctx, req.(*RemoveFromWhitelistAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_InsertWorkerPayload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertWorkerPayloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).InsertWorkerPayload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_InsertWorkerPayload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).InsertWorkerPayload(ctx, req.(*InsertWorkerPayloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_InsertReputerPayload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertReputerPayloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).InsertReputerPayload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgService_InsertReputerPayload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).InsertReputerPayload(ctx, req.(*InsertReputerPayloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MsgService_ServiceDesc is the grpc.ServiceDesc for MsgService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MsgService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "emissions.v5.MsgService",
	HandlerType: (*MsgServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _MsgService_UpdateParams_Handler,
		},
		{
			MethodName: "CreateNewTopic",
			Handler:    _MsgService_CreateNewTopic_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _MsgService_Register_Handler,
		},
		{
			MethodName: "RemoveRegistration",
			Handler:    _MsgService_RemoveRegistration_Handler,
		},
		{
			MethodName: "AddStake",
			Handler:    _MsgService_AddStake_Handler,
		},
		{
			MethodName: "RemoveStake",
			Handler:    _MsgService_RemoveStake_Handler,
		},
		{
			MethodName: "CancelRemoveStake",
			Handler:    _MsgService_CancelRemoveStake_Handler,
		},
		{
			MethodName: "DelegateStake",
			Handler:    _MsgService_DelegateStake_Handler,
		},
		{
			MethodName: "RewardDelegateStake",
			Handler:    _MsgService_RewardDelegateStake_Handler,
		},
		{
			MethodName: "RemoveDelegateStake",
			Handler:    _MsgService_RemoveDelegateStake_Handler,
		},
		{
			MethodName: "CancelRemoveDelegateStake",
			Handler:    _MsgService_CancelRemoveDelegateStake_Handler,
		},
		{
			MethodName: "FundTopic",
			Handler:    _MsgService_FundTopic_Handler,
		},
		{
			MethodName: "AddToWhitelistAdmin",
			Handler:    _MsgService_AddToWhitelistAdmin_Handler,
		},
		{
			MethodName: "RemoveFromWhitelistAdmin",
			Handler:    _MsgService_RemoveFromWhitelistAdmin_Handler,
		},
		{
			MethodName: "InsertWorkerPayload",
			Handler:    _MsgService_InsertWorkerPayload_Handler,
		},
		{
			MethodName: "InsertReputerPayload",
			Handler:    _MsgService_InsertReputerPayload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emissions/v5/tx.proto",
}
