package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"go_microservice_backend_api/internal/model"
	pb "go_microservice_backend_api/internal/service_product/grpc/grpc-gen"
	"go_microservice_backend_api/internal/service_product/service"
	"go_microservice_backend_api/pkg/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedProductServiceServer
}

func (s *server) CreateNewProduct(ctx context.Context, req *pb.CreateProductInput) (*pb.SampleResponse, error) {
	in := model.CreateProductInput{
		CategoryID: int(req.CategoryID),
		BrandID:    int(req.BrandID),
		SPUName:    req.SPUName,
		SPUDesc:    req.SPUDesc,
		SPUImg:     req.SPUImg,
		SPUVideo:   req.SPUVideo,
		SPUSort:    int(req.SPUSort),
		SPUPrice:   int(req.SPUPrice),
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.SampleResponse{
			StatusCode: 1,
			Message:    "Test",
		}, nil
	}
	shopId := md.Get("x-consumer-custom-id")[0]

	result, err := service.ProductService().AddNewProduct(ctx, in, shopId)
	if err != nil {
		return &pb.SampleResponse{
			StatusCode: int32(result),
			Message:    response.MSG[result],
		}, err
	}
	return &pb.SampleResponse{
		StatusCode: int32(result),
		Message:    response.MSG[result],
	}, nil
}
func (s *server) CreateNewCategory(ctx context.Context, req *pb.CreateCategoryInput) (*pb.SampleResponse, error) {
	in := model.CreateCategoryInput{
		ParentId:            int(req.ParentID),
		CategoryName:        req.CategoryName,
		HasActiveChildren:   req.HasActiveChildren,
		CategorySPUCount:    int(req.CategorySPUCount),
		CategoryStatus:      int(req.CategoryStatus),
		CategoryDescription: req.CategoryDescription,
		CategoryIcon:        req.CategoryIcon,
		CategorySort:        int(req.CategorySort),
	}

	result, err := service.CategoryService().AddNewCategory(ctx, in)
	if err != nil {
		return &pb.SampleResponse{
			StatusCode: int32(result),
			Message:    response.MSG[result],
		}, nil
	}
	return &pb.SampleResponse{
		StatusCode: int32(result),
		Message:    response.MSG[result],
	}, nil
}
func (s *server) CreateNewSKU(ctx context.Context, req *pb.CreateSKUInput) (*pb.SampleResponse, error) {
	in := model.CreateProductSKUInput{
		SPUID:             int(req.SpuID),
		SKUPrice:          int(req.SPUPrice),
		SKUStock:          int(req.SKUStock),
		SKUAttributeValue: req.SKUAttributeValue,
	}

	result, err := service.ProductService().AddNewSKUProduct(ctx, in)
	if err != nil {
		return &pb.SampleResponse{
			StatusCode: int32(result),
			Message:    response.MSG[result],
		}, err
	}

	return &pb.SampleResponse{
		StatusCode: int32(result),
		Message:    response.MSG[result],
	}, nil
}
func ProductGrpcServer() error {
	// Load the TLS certificates
	cert, err := tls.LoadX509KeyPair("/certs/cert.crt", "/certs/key.pem")
	if err != nil {
		return fmt.Errorf("failed to load certificates: %v", err)
	}

	// Create TLS credentials
	creds := credentials.NewServerTLSFromCert(&cert)

	// Listen on TCP port
	lis, err := net.Listen("tcp", "crm-product:6002")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create new gRPC server with TLS credentials
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Register your service
	pb.RegisterProductServiceServer(grpcServer, &server{})

	log.Println("starting secure grpc server on crm-product:6002")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
