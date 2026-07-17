// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"os"
	"path/filepath"
)

var _ Config = configurationFile{}

type configurationFile struct {
	file           string
	appendedConfig string
}

// HasConfigurationFiles is used during validation to ensure that
// ExternalProviders are not declared at the TestCase or TestStep
// level when using TestStep.ConfigFile.
func (c configurationFile) HasConfigurationFiles() bool {
	return true
}

// HasProviderBlock returns true if the Config has declared a provider
// configuration block, e.g. provider "examplecloud" {...}
func (c configurationFile) HasProviderBlock(ctx context.Context) (bool, error) {
	configFile := c.file

	if !filepath.IsAbs(configFile) {
		pwd, err := os.Getwd()

		if err != nil {
			return false, err
		}

		configFile = filepath.Join(pwd, configFile)
	}

	contains, err := fileContains(configFile, providerConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

// HasTerraformBlock returns true if the Config has declared a terraform
// configuration block, e.g. terraform {...}
func (c configurationFile) HasTerraformBlock(ctx context.Context) (bool, error) {
	configFile := c.file

	if !filepath.IsAbs(configFile) {
		pwd, err := os.Getwd()

		if err != nil {
			return false, err
		}

		configFile = filepath.Join(pwd, configFile)
	}

	contains, err := fileContains(configFile, terraformConfigBlockRegex)

	if err != nil {
		return false, err
	}

	return contains, nil
}

func (c configurationFile) WriteQuery(ctx context.Context, dest string) error {
	panic("WriteQuery not supported for configurationFile")
}

// Write copies file from c.file to destination.
func (c configurationFile) Write(ctx context.Context, dest string) error {
	configFile := c.file

	if !filepath.IsAbs(configFile) {
		pwd, err := os.Getwd()

		if err != nil {
			return err
		}

		configFile = filepath.Join(pwd, configFile)
	}

	destPath, err := copyFile(configFile, dest)
	if err != nil {
		return err
	}

	if len(c.appendedConfig) > 0 {
		err := appendToFile(destPath, c.appendedConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c configurationFile) Append(config string) Config {
	return configurationFile{
		file:           c.file,
		appendedConfig: config,
	}
}
