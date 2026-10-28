package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func ctxWithToken(token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("Token %s", token))
	return metautils.NiceMD(md).ToIncoming(context.Background())
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(42)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	ctx := ctxWithToken(token)
	userID, err := GetUserID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(42), userID)
}

func TestGenerateTokenWithTime(t *testing.T) {
	// Token generated in the future / within expiry window (72 hours)
	validTime := time.Now().Add(-1 * time.Hour)
	token, err := GenerateTokenWithTime(100, validTime)
	require.NoError(t, err)

	ctx := ctxWithToken(token)
	userID, err := GetUserID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(100), userID)

	// Token generated 80 hours ago (expired)
	expiredTime := time.Now().Add(-80 * time.Hour)
	expiredToken, err := GenerateTokenWithTime(100, expiredTime)
	require.NoError(t, err)

	expiredCtx := ctxWithToken(expiredToken)
	_, err = GetUserID(expiredCtx)
	require.Error(t, err)
	assert.Equal(t, "token expired", err.Error())
}

func TestGetUserID_MissingMetadata(t *testing.T) {
	_, err := GetUserID(context.Background())
	require.Error(t, err)
}

func TestGetUserID_MalformedToken(t *testing.T) {
	ctx := ctxWithToken("not-a-valid-jwt-token")
	_, err := GetUserID(ctx)
	require.Error(t, err)
	assert.Equal(t, "invalid token: it's not even a token", err.Error())
}
