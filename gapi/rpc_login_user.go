package gapi

import (
	"context"
	"database/sql"

	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/pb"
	"github.com/iamzay/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) loginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	// check password
	if err := util.CheckPassword(req.Password, []byte(user.HashedPassword)); err != nil {
		return nil, status.Error(codes.NotFound, "wrong password")
	}

	// make token
	token, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	// make refresh token and insert to sessions table
	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshTokenPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: "",
		ClientIp: "",
		IsBlocked: false,
		ExpiresAt: refreshTokenPayload.ExpiredAt,	
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           token,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}
	return rsp, nil
}
