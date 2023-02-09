package bootstrap

import (
	"Portfolio_Nodes/app"
	"Portfolio_Nodes/tools"
)

func Run(rootPath ...string) error {

	// ENV, LOGS, etc
	err := app.InitApp(rootPath...)
	if err != nil {
		return err
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		return err
	}
	err = app.RunMigrations(rootPath...)
	if err != nil {
		return err
	}

	// JWT
	err = tools.InitJWTTool()
	if err != nil {
		return err
	}

	_, _, err = app.InitGRPCServer()
	if err != nil {
		return err
	}

	// DI
	if err := injectDependencies(db); err != nil {
		return err
	}

	// Run gRPC and block
	app.RunGRPCServer()

	return nil
}
