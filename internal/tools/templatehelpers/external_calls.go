// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package templatehelpers

import (
	"context"
	"fmt"
	"os/exec"
)

// GoImports calls `goimports -w` over the specified file (including path)
func GoImports(file string) error {
	cmd := exec.CommandContext(context.Background(), "goimports", "-w", file)

	if combined, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, combined)
	}

	return nil
}
