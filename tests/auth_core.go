package tests

import (
	"github.com/goombaio/namegenerator"
	"github.com/sethvargo/go-password/password"
	"time"
)

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
