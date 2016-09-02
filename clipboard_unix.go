// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux netbsd openbsd solaris

package clipboard

import (
	"errors"
	"os/exec"
)

const (
	xsel  = "xsel"
	xclip = "xclip"
)

var (
	internalClipboard string
)

func init() {
	pasteCmdArgs = xclipPasteArgs
	copyCmdArgs = xclipCopyArgs

	if _, err := exec.LookPath(xclip); err == nil {
		return
	}

	pasteCmdArgs = xselPasteArgs
	copyCmdArgs = xselCopyArgs

	if _, err := exec.LookPath(xsel); err == nil {
		return
	}

	Unsupported = true
}

func copyCommand(register string) []string {
	if _, err := exec.LookPath(xclip); err == nil {
		return []string{xclip, "-in", "-selection", register}
	}

	if _, err := exec.LookPath(xsel); err == nil {
		return []string{xsel, "--input", "--" + register}
	}
}
func pasteCommand(register string) []string {
	if _, err := exec.LookPath(xclip); err == nil {
		return []string{xclip, "-out", "-selection", register}
	}

	if _, err := exec.LookPath(xsel); err == nil {
		return []string{xsel, "--output", "--" + register}
	}
}

func getPasteCommand(register string) *exec.Cmd {
	pasteCmdArgs := pasteCommand(register)
	return exec.Command(pasteCmdArgs[0], pasteCmdArgs[1:]...)
}

func getCopyCommand(register string) *exec.Cmd {
	copyCmdArgs := copyCommand(register)
	return exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
}

func readAll(register string) (string, error) {
	if Unsupported {
		return internalClipboard, nil
	}
	pasteCmd := getPasteCommand(register)
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeAll(text string, register string) error {
	if Unsupported {
		internalClipboard = text
		return nil
	}
	copyCmd := getCopyCommand(register)
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
