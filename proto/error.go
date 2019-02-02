package kvstore_methods

type UnexistsError struct {
	key string
}

func (err *UnexistsError) Error() string {
	return "Not exists " + err.key
}
