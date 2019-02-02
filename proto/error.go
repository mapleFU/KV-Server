package kvstore_methods

type UnexistsError struct {
	Key string
}

func (err *UnexistsError) Error() string {
	return "Not exists " + err.Key
}
