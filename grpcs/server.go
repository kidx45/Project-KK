package grpcs

import (
	db "github.com/kidx45/Project-KK/db/sqlc"
	"github.com/kidx45/Project-KK/pb"
	"github.com/kidx45/Project-KK/token"
	"github.com/kidx45/Project-KK/utils"
)

type Server struct {
	Store  db.Store
	Config utils.Config
	Token  token.Maker
	pb.UnimplementedProjectKKServer
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SYMMETRIC_SECRET_KEY)
	if err != nil {
		panic("cannot create token maker: %w")
	}

	server := &Server{Store: store, Config: config, Token: tokenMaker}
	
	return server, nil
}