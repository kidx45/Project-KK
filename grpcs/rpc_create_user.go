package grpcs

import (
	"context"

	db "github.com/kidx45/Project-KK/db/sqlc"
	"github.com/kidx45/Project-KK/pb"
	"github.com/kidx45/Project-KK/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserToPB(user db.User) *pb.User {
	return &pb.User{
		UserName:        user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   timestamppb.New(user.CreatedAt),
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashed, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	user, err := s.Store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.GetUserName(),
		HashedPassword: hashed,
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		PhoneNumber:    req.GetPhoneNumber(),
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Error(codes.AlreadyExists, "Unique violation error")
			}
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserResponse{User: convertUserToPB(user)}, nil
}
