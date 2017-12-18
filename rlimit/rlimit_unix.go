// +build !windows

package rlimit

import (
	"fmt"
	"syscall"
)

func Setrlimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}

	rLimit.Max = 65535
	rLimit.Cur = 65535

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
}
