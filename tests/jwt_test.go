package tests

import (
	"auth_service/tools"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJwt(t *testing.T) {

	paidUUID := uuid.NewV4().String()
	jwtAccess, err := tools.GenerateJWT(paidUUID)
	require.NoError(t, err, "should be success token generation")

	isValid, paidUUIDV, err := tools.VerifyJWT(jwtAccess.AccessToken)
	require.NoError(t, err, "should be success token verification")
	require.Equal(t, isValid, true, "should be valid token")
	require.Equal(t, paidUUID, paidUUIDV, "should be equal paidUUID's")
}
