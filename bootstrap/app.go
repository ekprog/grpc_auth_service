package bootstrap

import (
	"auth_service/app"
	"auth_service/tools"
)

func Run(rootPath ...string) error {

	// ENV, etc
	err := app.InitApp(rootPath...)
	if err != nil {
		return err
	}

	// Logger
	logger, err := app.InitLogs(rootPath...)
	if err != nil {
		return err
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		return err
	}

	// Migrations
	err = app.RunMigrations(rootPath...)
	if err != nil {
		return err
	}

	// JWT
	err = tools.InitJWTTool()
	if err != nil {
		return err
	}

	// gRPC
	_, _, err = app.InitGRPCServer()
	if err != nil {
		return err
	}

	// DI
	if err := injectDependencies(db, logger); err != nil {
		return err
	}

	// Run gRPC and block
	app.RunGRPCServer()

	return nil
}
