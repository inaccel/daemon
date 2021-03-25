package env

import (
	"os"
)

func Disabled(key string) bool {
	return os.Getenv(key) == "disabled"
}

func Enabled(key string) bool {
	return os.Getenv(key) == "enabled"
}
