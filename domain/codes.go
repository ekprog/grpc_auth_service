package domain

type Code string

const (
	Success              string = "success"
	AlreadyExists        string = "already_exists"
	IncorrectCredentials string = "incorrect_credentials"
	IncorrectToken       string = "incorrect_token"
)
