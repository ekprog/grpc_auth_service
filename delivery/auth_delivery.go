package delivery

import (
	"auth_service/app"
	"auth_service/domain"
	pb "auth_service/pkg/pb/api"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthDeliveryService struct {
	pb.UnsafeAuthServiceServer
	log       app.Logger
	authUCase domain.AuthInteractor
}

func NewAuthDeliveryService(log app.Logger, authUCase domain.AuthInteractor) *AuthDeliveryService {
	return &AuthDeliveryService{
		log:       log,
		authUCase: authUCase,
	}
}

func (d *AuthDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterAuthServiceServer, pb.AuthServiceServer(d))
	return nil
}

func (d *AuthDeliveryService) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	uCaseRes, err := d.authUCase.Register(r.Username, r.Password)
	if err != nil {
		return nil, err
	}

	response := &pb.RegisterResponse{
		Status: &pb.StatusResponse{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *AuthDeliveryService) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	uCaseRes, err := d.authUCase.Login(r.Username, r.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Error while Login")
	}

	response := &pb.LoginResponse{
		Status: &pb.StatusResponse{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		JwtAccess: nil,
	}

	if uCaseRes.StatusCode == domain.Success {
		response.JwtAccess = &pb.JWTAccess{
			AccessToken:      uCaseRes.UserToken.AccessToken,
			RefreshToken:     uCaseRes.UserToken.RefreshToken,
			AccessExpiredAt:  timestamppb.New(uCaseRes.UserToken.AccessTokenExpiredAt),
			RefreshExpiredAt: timestamppb.New(uCaseRes.UserToken.RefreshTokenExpiredAt),
		}
	}

	return response, nil
}

func (d *AuthDeliveryService) Revoke(ctx context.Context, r *pb.RevokeRequest) (*pb.RevokeResponse, error) {

	uCaseRes, err := d.authUCase.Revoke(r.AccessToken)
	if err != nil {
		return nil, err
	}

	return &pb.RevokeResponse{
		Status: &pb.StatusResponse{
			Code: uCaseRes.StatusCode,
		},
	}, nil
}

func (d *AuthDeliveryService) Verify(ctx context.Context, r *pb.VerifyRequest) (*pb.VerifyResponse, error) {

	uCaseRes, err := d.authUCase.Extract(r.AccessToken)
	if err != nil {
		return nil, err
	}

	response := &pb.VerifyResponse{
		Status: &pb.StatusResponse{
			Code: uCaseRes.StatusCode,
		},
		User: nil,
	}

	if uCaseRes.StatusCode == domain.Success {
		response.User = &pb.User{
			Id:        uCaseRes.User.Id,
			Username:  uCaseRes.User.Username,
			CreatedAt: timestamppb.New(uCaseRes.User.CreatedAt),
			UpdatedAt: timestamppb.New(uCaseRes.User.UpdatedAt),
		}
	}

	return response, nil
}

func (d *AuthDeliveryService) Refresh(ctx context.Context, r *pb.RefreshRequest) (*pb.RefreshResponse, error) {

	uCaseRes, err := d.authUCase.RefreshToken(r.RefreshToken)
	if err != nil {
		return nil, err
	}

	response := &pb.RefreshResponse{
		Status: &pb.StatusResponse{
			Code: uCaseRes.StatusCode,
		},
		JwtAccess: nil,
	}

	if uCaseRes.StatusCode == domain.Success {
		response.JwtAccess = &pb.JWTAccess{
			AccessToken:      uCaseRes.UserToken.AccessToken,
			RefreshToken:     uCaseRes.UserToken.RefreshToken,
			AccessExpiredAt:  timestamppb.New(uCaseRes.UserToken.AccessTokenExpiredAt),
			RefreshExpiredAt: timestamppb.New(uCaseRes.UserToken.RefreshTokenExpiredAt),
		}
	}

	return response, nil
}
