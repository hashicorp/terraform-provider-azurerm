// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var _ Config = configurationString{}

type configurationString struct {
	raw string
}

// HasConfigurationFiles is used during validation to allow declaration
// of ExternalProviders at the TestCase or TestStep level when using
// TestStep.Config.
func (c configurationString) HasConfigurationFiles() bool {
	return false
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
func (c configurationString) HasProviderBlock(ctx context.Context) (bool, error) {
	return providerConfigBlockRegex.MatchString(c.raw), nil
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configurationString) HasTerraformBlock(ctx context.Context) (bool, error) {
	return terraformConfigBlockRegex.MatchString(c.raw), nil
}

// Write creates a file and writes c.raw into it.
func (c configurationString) Write(ctx context.Context, dest string) error {
	outFilename := filepath.Join(dest, rawConfigFileName)
	rmFilename := filepath.Join(dest, rawConfigFileNameJSON)

	bCfg := []byte(c.raw)

	if json.Valid(bCfg) {
		outFilename, rmFilename = rmFilename, outFilename
	}

	if err := os.Remove(rmFilename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to remove %q: %w", rmFilename, err)
	}

	err := os.WriteFile(outFilename, bCfg, 0700)

	if err != nil {
		return err
	}

	return nil
}
