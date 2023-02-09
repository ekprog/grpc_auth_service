package grpc_delivery

import (
	"Portfolio_Nodes/app"
	"Portfolio_Nodes/domain"
	pb "Portfolio_Nodes/pkg/pb/api"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthDeliveryService struct {
	pb.UnsafeAuthServiceServer
	authUCase domain.AuthInteractor
}

func NewAuthDeliveryService(authUCase domain.AuthInteractor) *AuthDeliveryService {
	return &AuthDeliveryService{
		authUCase: authUCase,
	}
}

func (d *AuthDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterAuthServiceServer, pb.AuthServiceServer(d))
	return nil
}

func (d *AuthDeliveryService) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.GenericResponse, error) {
	err := d.authUCase.Register(r.Username, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.GenericResponse{
		Status:  "ok",
		Message: "Successful registration. Now you can login!",
	}, nil
}

func (d *AuthDeliveryService) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := d.authUCase.Login(r.Username, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		AccessToken: token.Token,
		ExpiredAt:   timestamppb.New(token.ExpiredAt),
	}, nil
}

func (d *AuthDeliveryService) Verify(ctx context.Context, r *pb.VerifyRequest) (*pb.VerifyResponse, error) {

	user, err := d.authUCase.ValidateAndExtract(r.AccessToken)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyResponse{
		Status: true,
		User: &pb.User{
			Id:        user.Id,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
