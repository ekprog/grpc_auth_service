package tests

import (
	"Portfolio_Nodes/tools"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJwt(t *testing.T) {

	userId := int64(5)

	tokenString, _, err := tools.GenerateJWT(userId)
	require.NoError(t, err, "should be success token generation")

	isValid, userIdV, err := tools.VerifyJWT(tokenString)
	require.NoError(t, err, "should be success token verification")
	require.Equal(t, isValid, true, "should be valid token")
	require.Equal(t, userId, userIdV, "should be equal userId's")
}
