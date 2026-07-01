//go:build !windows

package main

import (
	"os"
)

func CurrentUser() (string, error) {
	return os.Getenv("USER"), nil
}
