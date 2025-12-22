// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/cmd"
)

func main() {
	c := cmd.Make()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
