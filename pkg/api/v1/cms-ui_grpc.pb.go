// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// CmsUIClient is the client API for CmsUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CmsUIClient interface {
	// ListPageSpecs returns a list of Page(s) that can be started through the UI.
	ListPageSpecs(ctx context.Context, in *ListPageSpecsRequest, opts ...grpc.CallOption) (CmsUI_ListPageSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type cmsUIClient struct {
	cc grpc.ClientConnInterface
}

func NewCmsUIClient(cc grpc.ClientConnInterface) CmsUIClient {
	return &cmsUIClient{cc}
}

func (c *cmsUIClient) ListPageSpecs(ctx context.Context, in *ListPageSpecsRequest, opts ...grpc.CallOption) (CmsUI_ListPageSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &CmsUI_ServiceDesc.Streams[0], "/v1.CmsUI/ListPageSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &cmsUIListPageSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CmsUI_ListPageSpecsClient interface {
	Recv() (*ListPageSpecsResponse, error)
	grpc.ClientStream
}

type cmsUIListPageSpecsClient struct {
	grpc.ClientStream
}

func (x *cmsUIListPageSpecsClient) Recv() (*ListPageSpecsResponse, error) {
	m := new(ListPageSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cmsUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.CmsUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CmsUIServer is the server API for CmsUI service.
// All implementations must embed UnimplementedCmsUIServer
// for forward compatibility
type CmsUIServer interface {
	// ListPageSpecs returns a list of Page(s) that can be started through the UI.
	ListPageSpecs(*ListPageSpecsRequest, CmsUI_ListPageSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedCmsUIServer()
}

// UnimplementedCmsUIServer must be embedded to have forward compatible implementations.
type UnimplementedCmsUIServer struct {
}

func (UnimplementedCmsUIServer) ListPageSpecs(*ListPageSpecsRequest, CmsUI_ListPageSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListPageSpecs not implemented")
}
func (UnimplementedCmsUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedCmsUIServer) mustEmbedUnimplementedCmsUIServer() {}

// UnsafeCmsUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CmsUIServer will
// result in compilation errors.
type UnsafeCmsUIServer interface {
	mustEmbedUnimplementedCmsUIServer()
}

func RegisterCmsUIServer(s grpc.ServiceRegistrar, srv CmsUIServer) {
	s.RegisterService(&CmsUI_ServiceDesc, srv)
}

func _CmsUI_ListPageSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListPageSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CmsUIServer).ListPageSpecs(m, &cmsUIListPageSpecsServer{stream})
}

type CmsUI_ListPageSpecsServer interface {
	Send(*ListPageSpecsResponse) error
	grpc.ServerStream
}

type cmsUIListPageSpecsServer struct {
	grpc.ServerStream
}

func (x *cmsUIListPageSpecsServer) Send(m *ListPageSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _CmsUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CmsUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.CmsUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CmsUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CmsUI_ServiceDesc is the grpc.ServiceDesc for CmsUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CmsUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.CmsUI",
	HandlerType: (*CmsUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _CmsUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListPageSpecs",
			Handler:       _CmsUI_ListPageSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "cms-ui.proto",
}
