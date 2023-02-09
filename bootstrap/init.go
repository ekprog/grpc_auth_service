package bootstrap

import (
	"Portfolio_Nodes/app"
	"Portfolio_Nodes/delivery/grpc_delivery"
	"Portfolio_Nodes/interactors"
	"Portfolio_Nodes/repos"
	"database/sql"
	"go.uber.org/dig"
)

func injectDependencies(db *sql.DB) error {

	// DI Auto example
	//diObj := dig.New()
	//err := diObj.Provide(func() *sql.DB {
	//	return db
	//})
	//if err != nil {
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
	usersRepo := repos.NewUsersRepo(db)
	userTokensRepo := repos.NewUserTokensRepo(db)
	authUCase := interactors.NewAuthUCase(usersRepo, userTokensRepo)

	// Delivery Init
	authDelivery := grpc_delivery.NewAuthDeliveryService(authUCase)

	err := app.InitDelivery(authDelivery)
	if err != nil {
		return err
	}
	return nil
}

func provide(diObj *dig.Container, list ...interface{}) error {
	for _, p := range list {
		if err := diObj.Provide(p); err != nil {
			return err
		}
	}
	return nil
}
