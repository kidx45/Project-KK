package grpcs

import (
	"context"
	"database/sql"

	db "github.com/kidx45/Project-KK/db/sqlc"
	"github.com/kidx45/Project-KK/pb"
	"github.com/kidx45/Project-KK/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.Store.GetUser(ctx, req.GetUserName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "failed to login user")
	}
	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Invalid credentials")
	}

	accessToken, accessTokenPayload, err := s.Token.CreateToken(
		user.Username,
		s.Config.ACCESS_TOKEN_DURATION,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshTokenPayload, err := s.Token.CreateToken(
		user.Username,
		s.Config.REFRESH_TOKEN_DURATION,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	session, err := s.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	return &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUserToPB(user),
	}, nil
}
