package auth

type User interface {
	IsAuthenticated() bool
	UserName() string
}

type UnauthenticatedUser interface {
	IsAuthenticated()
	UserName() string
}
