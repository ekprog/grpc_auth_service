package bootstrap

import (
	"auth_service/app"
	"auth_service/delivery"
	"auth_service/interactors"
	"auth_service/repos"
	"database/sql"
)

func injectDependencies(db *sql.DB, logger app.Logger) error {

	// DI Auto example
	//diObj := dig.New()

	// logger
	//if err := diObj.Provide(func() app.Logger {
	//	return logger
	//}); err != nil {
	//	return err
	//}
	//
	//// db
	//if err := diObj.Provide(func() *sql.DB {
	//	return db
	//}); err != nil {
	//	return err
	//}

	//
	//err = provide(diObj,
	//	repos.NewUsersRepo,
	//	repos.NewUserTokensRepo,
	//	interactors.NewAuthUCase,
	//	interactors.NewAuthUCase,
	//	grpc_delivery.NewAuthDeliveryService)
	//if err != nil {
	//	return err
	//}

	// DI Manual
	usersRepo := repos.NewUsersRepo(logger, db)
	userTokensRepo := repos.NewUserTokensRepo(logger, db)
	authUCase := interactors.NewAuthUCase(logger, usersRepo, userTokensRepo)

	// Delivery Init
	authDelivery := delivery.NewAuthDeliveryService(logger, authUCase)

	err := app.InitDelivery(authDelivery)
	if err != nil {
		return err
	}
	return nil
}

//func provide(diObj *dig.Container, list ...interface{}) error {
//	for _, p := range list {
//		if err := diObj.Provide(p); err != nil {
//			return err
//		}
//	}
//	return nil
//}
