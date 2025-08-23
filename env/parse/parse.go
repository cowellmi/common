package parse

import "strconv"

type Parser[T any] func(s string) (T, error)

func Identity(s string) (string, error) {
	return s, nil
}

func Bool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func Int(s string) (int, error) {
	return strconv.Atoi(s)
}
