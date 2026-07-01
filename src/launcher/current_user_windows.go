//go:build windows

package main

import (
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func CurrentUser() (string, error) {
	// call GetUserNameW to get the current user in Windows "DOMAIN\\user" or "user"
	modadvapi32 := syscall.NewLazyDLL("advapi32.dll")
	procGetUserNameW := modadvapi32.NewProc("GetUserNameW")

	// allocate buffer for UTF-16 chars
	var buf [256]uint16
	size := uint32(len(buf))

	ret, _, err := procGetUserNameW.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if ret == 0 {
		return "", err
	}

	// convert UTF-16 to Go string
	name := string(utf16.Decode(buf[:size-1])) // size includes null terminator

	// Normalize: uppercase and strip "NUTH\\" prefix if present (case-insensitive)
	nameUp := strings.ToUpper(name)
	const prefix = "NUTH\\"
	nameUp = strings.TrimPrefix(nameUp, prefix)
	return nameUp, nil
}
