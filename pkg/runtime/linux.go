package runtime

import (
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

var LinuxVersionCode uint32

func init() {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return
	}

	versionCore := string(uname.Release[:])
	if index := strings.IndexFunc(versionCore, func(c rune) bool {
		return c == '-' || c == '+'
	}); index != -1 {
		versionCore = string(uname.Release[:index])
	}

	fields := strings.FieldsFunc(versionCore, func(c rune) bool {
		return c == '.'
	})
	if len(fields) != 3 {
		return
	}

	major, err := strconv.ParseUint(fields[0], 10, 16)
	if err != nil {
		return
	}

	minor, err := strconv.ParseUint(fields[1], 10, 8)
	if err != nil {
		return
	}

	patch, err := strconv.ParseUint(fields[2], 10, 8)
	if err != nil {
		return
	}

	LinuxVersionCode = KernelVersion(uint16(major), uint8(minor), uint8(patch))
}

func KernelVersion(a uint16, b, c uint8) uint32 {
	return (uint32(a) << 16) + (uint32(b) << 8) + uint32(c)
}
