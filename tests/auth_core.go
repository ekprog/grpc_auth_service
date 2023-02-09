package tests

import (
	"Portfolio_Nodes/app"
	"Portfolio_Nodes/domain"
	"Portfolio_Nodes/interactors"
	"Portfolio_Nodes/repos"
	"Portfolio_Nodes/tools"
	"github.com/goombaio/namegenerator"
	"github.com/sethvargo/go-password/password"
	"time"
)

var (
	testUsername string
	testPassword string
	testToken    string

	usersRepo      domain.UsersRepository
	userTokensRepo domain.UserTokensRepository
	authInteractor domain.AuthInteractor
)

func initCore() {

	// ENV, LOGS, etc
	err := app.InitApp("..")
	if err != nil {
		panic(err)
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		panic(err)
	}
	err = app.RunMigrations("..")
	if err != nil {
		panic(err)
	}

	// JWT
	err = tools.InitJWTTool()
	if err != nil {
		panic(err)
	}

	usersRepo = repos.NewUsersRepo(db)
	userTokensRepo = repos.NewUserTokensRepo(db)
	authInteractor = interactors.NewAuthUCase(usersRepo, userTokensRepo)
}

func generateCredentials() (string, string, error) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	username := nameGenerator.Generate()
	password, err := password.Generate(10, 4, 0, false, true)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
