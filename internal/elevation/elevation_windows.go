//go:build windows
// +build windows

package elevation

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func isWindowsElevated() bool {
	var sid *windows.SID

	// Get the built-in administrators group SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	// Check if current token is a member of the administrators group
	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

func requestWindowsElevation() error {
	verb := "runas"
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := ""
	for i, arg := range os.Args[1:] {
		if i > 0 {
			args += " "
		}
		args += arg
	}

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 // SW_NORMAL

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return fmt.Errorf("failed to elevate privileges: %w", err)
	}

	// Exit current process after starting elevated one
	os.Exit(0)
	return nil
}

func requestDarwinElevation() error {
	return fmt.Errorf("not on macOS")
}

func requestLinuxElevation() error {
	return fmt.Errorf("not on Linux")
}
