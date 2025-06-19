package errs

type Error[T any] struct {
	Code    int
	Data    T
	Message string
}
