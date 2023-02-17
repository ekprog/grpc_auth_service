package interactors

import (
	"auth_service/app"
	"auth_service/domain"
	"auth_service/tools"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type AuthInteractor struct {
	log        app.Logger
	usersRepo  domain.UsersRepository
	tokensRepo domain.UserTokensRepository
}

func NewAuthUCase(log app.Logger, usersRepo domain.UsersRepository, tokensRepo domain.UserTokensRepository) domain.AuthInteractor {
	return &AuthInteractor{
		log:        log,
		usersRepo:  usersRepo,
		tokensRepo: tokensRepo,
	}
}

func (i *AuthInteractor) generateNewJWTPair(userId int64) (*domain.UserToken, error) {

	pairUUID := uuid.NewV4().String()
	jwtAccess, err := tools.GenerateJWT(pairUUID)
	if err != nil {
		return nil, err
	}
	token := &domain.UserToken{
		UserId:                userId,
		PairUUID:              pairUUID,
		AccessToken:           jwtAccess.AccessToken,
		RefreshToken:          jwtAccess.RefreshToken,
		AccessTokenExpiredAt:  jwtAccess.AccessTokenExpired,
		RefreshTokenExpiredAt: jwtAccess.RefreshTokenExpired,
	}
	err = i.tokensRepo.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (i *AuthInteractor) Register(username, password string) (domain.RegisterResponse, error) {

	// check exists
	existsUser, err := i.usersRepo.FindByUsername(username)
	if err != nil {
		return domain.RegisterResponse{}, errors.Wrap(err, "Cannot find user by username")
	}
	if existsUser != nil {
		i.log.Debug("Incorrect credentials: %s@%s", username, password)
		return domain.RegisterResponse{
			StatusCode: domain.AlreadyExists,
		}, nil
	}

	// make password hash
	pwdHash, err := tools.GenerateHashPassword(password)
	if err != nil {
		return domain.RegisterResponse{}, errors.Wrap(err, "Cannot generate pws hash")
	}

	// create user
	user := &domain.User{
		Username: username,
		PwdHash:  pwdHash,
	}
	err = i.usersRepo.Insert(user)
	if err != nil {
		return domain.RegisterResponse{}, errors.Wrap(err, "Cannot insert user in DB")
	}

	log.Infof("Register successful. Credentials: %s - %s", username, pwdHash)

	return domain.RegisterResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *AuthInteractor) Login(username, password string) (domain.LoginResponse, error) {

	// Check user exists
	user, err := i.usersRepo.FindByUsername(username)
	if err != nil {
		return domain.LoginResponse{}, errors.Wrap(err, "Cannot check user existing")
	}

	// Check password
	isPassValid := tools.CheckPasswordHash(password, user.PwdHash)
	if !isPassValid {
		i.log.Debug("Incorrect credentials: %s@%s", username, password)
		return domain.LoginResponse{
			StatusCode: domain.IncorrectCredentials,
		}, nil
	}

	// Generating new token
	jwtAccess, err := i.generateNewJWTPair(user.Id)
	if err != nil {
		return domain.LoginResponse{}, errors.Wrap(err, "Cannot generate jwt pair")
	}

	log.Infof("Login successful. Token = %s", jwtAccess.AccessToken)

	return domain.LoginResponse{
		StatusCode: domain.Success,
		UserToken:  jwtAccess,
	}, nil
}

func (i *AuthInteractor) Revoke(token string) (domain.RevokeResponse, error) {

	// Verify refresh token and get pairUUID
	ok, pairUUID, err := tools.VerifyJWT(token)
	if err != nil {
		return domain.RevokeResponse{}, errors.Wrap(err, "Error while verifying access token")
	}
	if !ok {
		i.log.Debug("Error while verifying access token: %s", err)
		return domain.RevokeResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}

	// Find valid token from db (if it is refresh token and not expired)
	jwtAccess, err := i.tokensRepo.FindValidPair(pairUUID)
	if err != nil {
		return domain.RevokeResponse{}, errors.Wrap(err, "Incorrect pair UUID")
	}
	// JWT was created but where is row in DB ?
	if jwtAccess == nil {
		i.log.Warn("JWT pair was not found in database", err)
		return domain.RevokeResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}
	if token != jwtAccess.AccessToken {
		return domain.RevokeResponse{}, errors.New("Need pass refresh token")
	}
	// Not necessary because tools.VerifyJWT already checked expiration
	if time.Now().After(jwtAccess.AccessTokenExpiredAt) {
		return domain.RevokeResponse{}, errors.New("Access token is expired")
	}

	// Updating updated_at field
	err = i.tokensRepo.RevokePair(pairUUID)
	if err != nil {
		return domain.RevokeResponse{}, errors.New("Cannot revoke token")
	}

	return domain.RevokeResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *AuthInteractor) RefreshToken(refreshToken string) (domain.RefreshTokenResponse, error) {

	// Verify refresh token and get pairUUID
	ok, pairUUID, err := tools.VerifyJWT(refreshToken)
	if err != nil {
		return domain.RefreshTokenResponse{}, errors.Wrap(err, "Error while verifying refresh token")
	}
	if !ok {
		i.log.Debug("Error while verifying refresh token: %s", err)
		return domain.RefreshTokenResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}

	// Find valid token from db (if it is refresh token and not expired)
	jwtAccess, err := i.tokensRepo.FindValidPair(pairUUID)
	if err != nil {
		return domain.RefreshTokenResponse{}, errors.Wrap(err, "Incorrect pair UUID")
	}

	// JWT was created but where is row in DB ?
	if jwtAccess == nil {
		i.log.Warn("JWT pair was not found in database", err)
		return domain.RefreshTokenResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}
	if refreshToken != jwtAccess.RefreshToken {
		return domain.RefreshTokenResponse{}, errors.New("Need pass refresh token")
	}
	// Not necessary because tools.VerifyJWT already checked expiration
	if time.Now().After(jwtAccess.RefreshTokenExpiredAt) {
		return domain.RefreshTokenResponse{}, errors.New("Refresh token is expired")
	}

	// Making paid invalid
	err = i.tokensRepo.RevokePair(jwtAccess.PairUUID)
	if err != nil {
		return domain.RefreshTokenResponse{}, errors.Wrap(err, "Cannot revoke access token")
	}

	// Generating token
	jwtAccess, err = i.generateNewJWTPair(jwtAccess.UserId)
	if err != nil {
		return domain.RefreshTokenResponse{}, errors.Wrap(err, "Cannot generate new jwt pair")
	}

	// Log
	i.log.Info("Refreshing token successful. PairUUID=%d", jwtAccess.PairUUID)

	return domain.RefreshTokenResponse{
		StatusCode: domain.Success,
		UserToken:  jwtAccess,
	}, nil
}

func (i *AuthInteractor) Extract(accessToken string) (domain.ExtractResponse, error) {

	// Verify refresh token and get pairUUID
	ok, pairUUID, err := tools.VerifyJWT(accessToken)
	if err != nil {
		return domain.ExtractResponse{}, errors.Wrap(err, "Error while verifying access token")
	}
	if !ok {
		i.log.Debug("Error while verifying access token: %s", err)
		return domain.ExtractResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}

	// Find valid token from db (if it is refresh token and not expired)
	jwtAccess, err := i.tokensRepo.FindValidPair(pairUUID)
	if err != nil {
		return domain.ExtractResponse{}, errors.Wrap(err, "Incorrect pair UUID")
	}
	// JWT was created but where is row in DB ?
	if jwtAccess == nil {
		i.log.Warn("JWT pair was not found in database", err)
		return domain.ExtractResponse{
			StatusCode: domain.IncorrectToken,
		}, nil
	}
	if accessToken != jwtAccess.AccessToken {
		return domain.ExtractResponse{}, errors.New("Need pass refresh token")
	}
	// Not necessary because tools.VerifyJWT already checked expiration
	if time.Now().After(jwtAccess.AccessTokenExpiredAt) {
		return domain.ExtractResponse{}, errors.New("Access token is expired")
	}

	// Updating updated_at field
	err = i.tokensRepo.UpdateTime(pairUUID)
	if err != nil {
		i.log.Warn("Cannot update time for jwt pair")
		// No need to return (not critical)
	}

	// Extract
	user, err := i.usersRepo.FindById(jwtAccess.UserId)
	if err != nil {
		return domain.ExtractResponse{}, errors.Wrap(err, "Error while finding user in database")
	}

	return domain.ExtractResponse{
		StatusCode: domain.Success,
		User:       user,
	}, nil
}
