package gapi

import (
	"fmt"

	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/pb"
	"github.com/iamzay/simplebank/token"
	"github.com/iamzay/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store db.Store
	tokenMaker token.Maker
	config util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	server := &Server{ store: store, config: config, tokenMaker: tokenMaker}
	return server, nil
}
