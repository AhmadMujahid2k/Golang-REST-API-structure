package util

import (
	"os"
	"regexp"
)

// Require that a environment variable exist and be non-empty.
// Returns the value associated with the variable.
func MustOsGetEnv(key string) string {
	res := os.Getenv(key)
	if res == "" {
		panic("env var " + key + " must not be empty")
	}
	return res
}

var (
	REGEX_EMAIL = regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9_.+-]*@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]{2,}$")
	REGEX_PHONE = regexp.MustCompile("^\\+[0-9]{2,15}$")
)

func IsEmail(email string) bool {
	return REGEX_EMAIL.MatchString(email)
}

func IsPhone(phone string) bool {
	return REGEX_PHONE.MatchString(phone)
}
