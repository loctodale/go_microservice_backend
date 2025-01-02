package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"go_microservice_backend_api/internal/model"
	sg "go_microservice_backend_api/internal/service_shop/grpc/grpc-gen"
	"go_microservice_backend_api/internal/service_shop/service"
	"go_microservice_backend_api/pkg/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type server struct {
	sg.UnimplementedShopServiceServer
}

func (s *server) ShopRegister(ctx context.Context, req *sg.ShopRegisterInput) (*sg.SampleResponse, error) {
	// Validate input
	if req.VerifyKey == "" || req.VerifyPurpose == "" || req.VerifyType == "" {
		return &sg.SampleResponse{
			StatusCode: response.ErrCodeParamInvalid,
			Message:    response.MSG[response.ErrCodeParamInvalid],
		}, nil
	}
	in := model.RegisterInput{VerifyKey: req.VerifyKey, VerifyPurpose: req.VerifyPurpose, VerifyType: 1}
	result, err := service.ShopRegisterService().Register(ctx, in)
	if err != nil {
		return &sg.SampleResponse{
			StatusCode: int32(result),
			Message:    response.MSG[result],
		}, nil
	}
	// Example successful response
	return &sg.SampleResponse{
		StatusCode: 200,
		Message:    "Shop registered successfully",
	}, nil
}

func (s *server) ShopVerifyOTP(ctx context.Context, req *sg.ShopVerifyOTPInput) (*sg.VerifyOTPResponse, error) {
	// Validate input
	if req.VerifyCode == "" || req.VerifyKey == "" {
		return &sg.VerifyOTPResponse{
			StatusCode: response.ErrCodeParamInvalid,
			Message:    response.MSG[response.ErrCodeParamInvalid],
		}, nil
	}
	in := model.VerifyInput{VerifyCode: req.VerifyCode, VerifyKey: req.VerifyKey}
	result, err := service.ShopRegisterService().VerifyOTP(ctx, in)
	if err != nil {
		return &sg.VerifyOTPResponse{
			StatusCode: 500,
			Message:    err.Error(),
		}, nil
	}
	// Example successful response
	return &sg.VerifyOTPResponse{
		StatusCode:   response.CodeSuccess,
		Message:      response.MSG[response.CodeSuccess],
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ShopId:       result.ShopId,
	}, nil
}

func (s *server) ShopChangePasswordVerify(ctx context.Context, req *sg.ShopChangePasswordVerifyInput) (*sg.StringResponse, error) {
	if req.Password == "" {
		return &sg.StringResponse{
			Message: response.MSG[response.ErrCodeParamInvalid],
		}, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &sg.StringResponse{
			Message: "Not found username",
		}, nil
	}
	userName := md.Get("X-Consumer-Username")[0]
	fmt.Println(userName)

	result, err := service.ShopRegisterService().ChangePasswordRegister(ctx, userName, req.Password)
	if err != nil {
		fmt.Println(err)
		return &sg.StringResponse{
			Message: err.Error(),
		}, err
	}
	return &sg.StringResponse{
		Message: result,
	}, nil
}

func (s *server) ShopLogin(ctx context.Context, req *sg.ShopLoginInput) (*sg.VerifyOTPResponse, error) {
	if req.Password == "" || req.Username == "" {
		return &sg.VerifyOTPResponse{
			StatusCode:   response.ErrCodeParamInvalid,
			Message:      response.MSG[response.ErrCodeParamInvalid],
			AccessToken:  "",
			ShopId:       "",
			RefreshToken: "",
		}, nil
	}
	in := model.ShopLoginInput{
		UserAccount:  req.Username,
		UserPassword: req.Password,
	}

	result, err := service.ShopRegisterService().LoginShop(ctx, in)
	if err != nil {
		fmt.Println(err)
		return &sg.VerifyOTPResponse{
			StatusCode:   response.ErrCodeAuthFailed,
			Message:      response.MSG[response.ErrCodeAuthFailed],
			AccessToken:  "",
			ShopId:       "",
			RefreshToken: "",
		}, err
	}
	return &sg.VerifyOTPResponse{
		StatusCode:   response.CodeSuccess,
		Message:      response.MSG[response.CodeSuccess],
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ShopId:       result.ShopId,
	}, nil
}

func ShopGrpcServer() error {
	// Load the TLS certificates
	cert, err := tls.LoadX509KeyPair("/certs/cert.crt", "/certs/key.pem")
	if err != nil {
		return fmt.Errorf("failed to load certificates: %v", err)
	}

	// Create TLS credentials
	creds := credentials.NewServerTLSFromCert(&cert)

	// Listen on TCP port
	lis, err := net.Listen("tcp", "crm-shop:6001")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create new gRPC server with TLS credentials
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Register your service
	sg.RegisterShopServiceServer(grpcServer, &server{})

	log.Println("starting secure grpc server on crm-shop:6001")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
