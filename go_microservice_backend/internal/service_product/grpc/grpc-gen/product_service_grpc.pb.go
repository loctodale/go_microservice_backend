// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: internal/proto/product/product_service.proto

package grpc_gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProductService_CreateNewProduct_FullMethodName  = "/productservice.ProductService/CreateNewProduct"
	ProductService_CreateNewCategory_FullMethodName = "/productservice.ProductService/CreateNewCategory"
	ProductService_CreateNewSKU_FullMethodName      = "/productservice.ProductService/CreateNewSKU"
)

// ProductServiceClient is the client API for ProductService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServiceClient interface {
	CreateNewProduct(ctx context.Context, in *CreateProductInput, opts ...grpc.CallOption) (*SampleResponse, error)
	CreateNewCategory(ctx context.Context, in *CreateCategoryInput, opts ...grpc.CallOption) (*SampleResponse, error)
	CreateNewSKU(ctx context.Context, in *CreateSKUInput, opts ...grpc.CallOption) (*SampleResponse, error)
}

type productServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return &productServiceClient{cc}
}

func (c *productServiceClient) CreateNewProduct(ctx context.Context, in *CreateProductInput, opts ...grpc.CallOption) (*SampleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SampleResponse)
	err := c.cc.Invoke(ctx, ProductService_CreateNewProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) CreateNewCategory(ctx context.Context, in *CreateCategoryInput, opts ...grpc.CallOption) (*SampleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SampleResponse)
	err := c.cc.Invoke(ctx, ProductService_CreateNewCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) CreateNewSKU(ctx context.Context, in *CreateSKUInput, opts ...grpc.CallOption) (*SampleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SampleResponse)
	err := c.cc.Invoke(ctx, ProductService_CreateNewSKU_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServiceServer is the server API for ProductService service.
// All implementations must embed UnimplementedProductServiceServer
// for forward compatibility.
type ProductServiceServer interface {
	CreateNewProduct(context.Context, *CreateProductInput) (*SampleResponse, error)
	CreateNewCategory(context.Context, *CreateCategoryInput) (*SampleResponse, error)
	CreateNewSKU(context.Context, *CreateSKUInput) (*SampleResponse, error)
	mustEmbedUnimplementedProductServiceServer()
}

// UnimplementedProductServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProductServiceServer struct{}

func (UnimplementedProductServiceServer) CreateNewProduct(context.Context, *CreateProductInput) (*SampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewProduct not implemented")
}
func (UnimplementedProductServiceServer) CreateNewCategory(context.Context, *CreateCategoryInput) (*SampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewCategory not implemented")
}
func (UnimplementedProductServiceServer) CreateNewSKU(context.Context, *CreateSKUInput) (*SampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewSKU not implemented")
}
func (UnimplementedProductServiceServer) mustEmbedUnimplementedProductServiceServer() {}
func (UnimplementedProductServiceServer) testEmbeddedByValue()                        {}

// UnsafeProductServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServiceServer will
// result in compilation errors.
type UnsafeProductServiceServer interface {
	mustEmbedUnimplementedProductServiceServer()
}

func RegisterProductServiceServer(s grpc.ServiceRegistrar, srv ProductServiceServer) {
	// If the following call pancis, it indicates UnimplementedProductServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProductService_ServiceDesc, srv)
}

func _ProductService_CreateNewProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).CreateNewProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_CreateNewProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).CreateNewProduct(ctx, req.(*CreateProductInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_CreateNewCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCategoryInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).CreateNewCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_CreateNewCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).CreateNewCategory(ctx, req.(*CreateCategoryInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_CreateNewSKU_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSKUInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).CreateNewSKU(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_CreateNewSKU_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).CreateNewSKU(ctx, req.(*CreateSKUInput))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductService_ServiceDesc is the grpc.ServiceDesc for ProductService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "productservice.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNewProduct",
			Handler:    _ProductService_CreateNewProduct_Handler,
		},
		{
			MethodName: "CreateNewCategory",
			Handler:    _ProductService_CreateNewCategory_Handler,
		},
		{
			MethodName: "CreateNewSKU",
			Handler:    _ProductService_CreateNewSKU_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/product/product_service.proto",
}
