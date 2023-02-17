package app

import (
	"github.com/joho/godotenv"
)

func InitApp(rootDir ...string) error {

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	err := godotenv.Load(basePath + "/.env")
	if err != nil {
		return err
	}

	return nil
}
