package interactors

import (
	"Portfolio_Nodes/app"
	"Portfolio_Nodes/domain"
	"Portfolio_Nodes/tools"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthInteractor struct {
	usersRepo  domain.UsersRepository
	tokensRepo domain.UserTokensRepository
}

func NewAuthUCase(usersRepo domain.UsersRepository, tokensRepo domain.UserTokensRepository) domain.AuthInteractor {
	return &AuthInteractor{usersRepo: usersRepo, tokensRepo: tokensRepo}
}

func (i *AuthInteractor) invokeNewToken(userId int64) (*domain.UserToken, error) {
	tokenString, expired, err := tools.GenerateJWT(userId)
	if err != nil {
		return nil, err
	}
	token := &domain.UserToken{
		UserId:    userId,
		Token:     tokenString,
		ExpiredAt: expired,
	}
	err = i.tokensRepo.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (i *AuthInteractor) Register(username, password string) error {

	if username == "" || password == "" {
		return app.UCaseError("validation_error")
	}

	// check exists
	existsUser, err := i.usersRepo.FindByUsername(username)
	if err != nil {
		return err
	}
	if existsUser != nil {
		return status.Error(codes.AlreadyExists, "username has already taken")
	}

	// make password hash
	pwdHash, err := tools.GenerateHashPassword(password)
	if err != nil {
		return err
	}

	// create user
	user := &domain.User{
		Username: username,
		PwdHash:  pwdHash,
	}
	err = i.usersRepo.Insert(user)
	if err != nil {
		return err
	}

	log.Infof("Register successful. Credentials: %s - %s", username, pwdHash)

	return nil
}

func (i *AuthInteractor) Login(username, password string) (*domain.UserToken, error) {

	user, err := i.usersRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	isPassValid := tools.CheckPasswordHash(password, user.PwdHash)
	if !isPassValid {
		return nil, app.UCaseError("incorrect_credentials")
	}

	// generate token
	token, err := i.invokeNewToken(user.Id)
	if err != nil {
		return nil, err
	}

	log.Infof("Login successful. Token = %s", token.Token)

	return token, nil
}

func (i *AuthInteractor) ValidateAndExtract(tokenString string) (*domain.User, error) {
	// Find valid token from db
	_, err := i.tokensRepo.FindValid(tokenString)
	if err != nil {
		return nil, err
	}
	err = i.tokensRepo.UpdateTime(tokenString)
	if err != nil {
		log.Errorf("cannot update updated_at at token (%s)", tokenString)
		// no need return (not critical error)
	}

	// Verify token
	isValid, userId, err := tools.VerifyJWT(tokenString)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, app.UCaseError("incorrect_credentials")
	}

	user, err := i.usersRepo.FindById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.Errorf("cannot find user after getting userId (%d) from token.", userId)
		return nil, app.UCaseError("incorrect_credentials")
	}

	log.Infof("Validate successful. UserId=%d, Username=%s", user.Id, user.Username)

	return user, nil
}
