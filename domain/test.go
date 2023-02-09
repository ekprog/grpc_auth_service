package domain

type Test struct {
	Message string `json:"message"`
}

type TestInteractor interface {
	MakeTest() (Test, error)
}

type TestRepository interface {
	FetchTest() (Test, error)
}
