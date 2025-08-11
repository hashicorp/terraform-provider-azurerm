// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// CallTerraform calls the locally installed terraform binary with the specified options
// returns the raw []byte output for the caller to process.
func CallTerraform(opts ...string) ([]byte, error) {
	cmd := exec.Command("terraform", opts...)
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("missing terraform binary, please ensure Terraform is installed and discoverable in your PATH")
		}
		return nil, err
	}

	return out, nil
}

// GoFmt calls `gofmt -w` over the specified file (including path)
func GoFmt(file string) error {
	cmd := exec.Command("gofmt", "-w", fmt.Sprintf("./%s", file))
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

// Terrafmt calls (if installed) katbyte/terrafmt to format Terraform
// configurations in the specified file
func Terrafmt(path string) error {
	cmd := exec.Command("terrafmt", "fmt", "-f", path)
	if _, err := cmd.Output(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing terrafmt, please ensure katbyte/terrafmt is installed (`go install github.com/katbyte/terrafmt@latest`) and discoverable in your PATH")
		}
		return err
	}

	return nil
}
