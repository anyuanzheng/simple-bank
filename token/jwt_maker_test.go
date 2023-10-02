package token

import (
	"testing"
	"time"

	"github.com/iamzay/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	username := util.RandomName()
	secretKey := util.RandomString(32)
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	token, err := jwtMaker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(time.Minute), time.Second)
}

func TestInvalidToken(t *testing.T) {
	username := util.RandomName()
	secretKey := util.RandomString(32)	
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)	

	token, err := jwtMaker.CreateToken(username, -time.Minute)	
	require.NoError(t, err)	
	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}
