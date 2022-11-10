package helpers

import (
	"errors"
	"fmt"
	"strconv"
)

//TODO(michel): rename this package to be more expressive. (100 common go mistakes #13 - creating common packages)

// MustParseBool takes in a boolean representation of a string and returns a boolean. Defaults to false if an invalid string is provided.
func MustParseBool(boolean string) bool {
	parsed, err := strconv.ParseBool(boolean)
	if err != nil && errors.Is(err, strconv.ErrSyntax) {
		return false
	} else if err != nil {
		panic(fmt.Errorf("parsing boolean string: %s", err))
	}
	return parsed
}
