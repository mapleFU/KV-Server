package storage

import "fmt"

type ArgumentError struct {
	Expected string
	Value interface{}
}

func (err *ArgumentError) Error() string {
	return fmt.Sprintf("Expected %s but got %v", err.Expected, err.Value)
}

