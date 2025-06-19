package auth

type User interface {
	IsAuthenticated() bool
	UserName() string
}
