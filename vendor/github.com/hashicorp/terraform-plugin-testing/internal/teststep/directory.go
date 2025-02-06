// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
)

var _ Config = configurationDirectory{}

type configurationDirectory struct {
	directory string
}

// HasConfigurationFiles is used during validation to ensure that
// ExternalProviders are not declared at the TestCase or TestStep
// level when using TestStep.ConfigDirectory.
func (c configurationDirectory) HasConfigurationFiles() bool {
	return true
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
func (c configurationDirectory) HasProviderBlock(ctx context.Context) (bool, error) {
	configDirectory := c.directory

	if !filepath.IsAbs(configDirectory) {
		pwd, err := os.Getwd()

		if err != nil {
			return false, err
		}

		configDirectory = filepath.Join(pwd, configDirectory)
	}

	contains, err := filesContains(configDirectory, providerConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configurationDirectory) HasTerraformBlock(ctx context.Context) (bool, error) {
	configDirectory := c.directory

	if !filepath.IsAbs(configDirectory) {
		pwd, err := os.Getwd()

		if err != nil {
			return false, err
		}

		configDirectory = filepath.Join(pwd, configDirectory)
	}

	contains, err := filesContains(configDirectory, terraformConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

// Write copies all files from directory to destination.
func (c configurationDirectory) Write(ctx context.Context, dest string) error {
	configDirectory := c.directory

	if !filepath.IsAbs(configDirectory) {
		pwd, err := os.Getwd()

		if err != nil {
			return err
		}

		configDirectory = filepath.Join(pwd, configDirectory)
	}

	err := copyFiles(configDirectory, dest)

	if err != nil {
		return err
	}

	return nil
}
