package errors

// Public wraps the original error with a new error that has a
// `Public() string` method that will return a message that is
// acceptable to display to the public. This error can also be
// unwrapped using the traditional `errors` package approach.
func Public(err error, msg string) error {
	return publicError{err: err, msg: msg}
}

type publicError struct {
	err error
	msg string
}

func (r publicError) Public() string {
	return r.msg
}
func (r publicError) Error() string {
	return r.err.Error()
}
func (r publicError) Unwrap() error {
	return r.err
}
