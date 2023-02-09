package tests

import (
	"Portfolio_Nodes/tools"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestAuth_Core(t *testing.T) {

	initCore()

	var err error
	testUsername, testPassword, err = generateCredentials()
	require.NoError(t, err, "should be success credentials generation process")

	err = authInteractor.Register(testUsername, testPassword)
	require.NoError(t, err, "should be success register ucase")

	token, err := authInteractor.Login(testUsername, testPassword)
	require.NoError(t, err, "should be success login ucase")
	testToken = token.Token

	// Valid validating
	user, err := authInteractor.ValidateAndExtract(testToken)
	require.NoError(t, err, "should be success validate ucase")
	require.Equal(t, user.Username, testUsername, "should be equal userNames")

	// Invalid Validation
	usedId := int64(rand.Uint64())
	tokenString, _, err := tools.GenerateJWT(usedId)
	require.NoError(t, err, "should be success jwt generation process")
	_, err = authInteractor.ValidateAndExtract(tokenString)
	require.Error(t, err, "should be error validate ucase")
}
