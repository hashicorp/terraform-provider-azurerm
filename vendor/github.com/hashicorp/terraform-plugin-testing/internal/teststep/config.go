// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

const (
	rawConfigFileName     = "terraform_plugin_test.tf"
	rawConfigFileNameJSON = rawConfigFileName + ".json"
)

var (
	providerConfigBlockRegex  = regexp.MustCompile(`provider "?[a-zA-Z0-9_-]+"? {`)
	terraformConfigBlockRegex = regexp.MustCompile(`terraform {`)
)

// Config defines an interface implemented by all types
// that represent Terraform configuration:
//
//   - [config.configurationDirectory]
//   - [config.configurationFile]
//   - [config.configurationString]
type Config interface {
	HasConfigurationFiles() bool
	HasProviderBlock(context.Context) (bool, error)
	HasTerraformBlock(context.Context) (bool, error)
	Write(context.Context, string) error
}

// PrepareConfigurationRequest is used to simplify the generation of
// a ConfigurationRequest which is required when calling the
// Configuration func.
type PrepareConfigurationRequest struct {
	Directory             config.TestStepConfigFunc
	File                  config.TestStepConfigFunc
	Raw                   string
	TestStepConfigRequest config.TestStepConfigRequest
}

// Exec returns a Configuration request which is required when
// calling the Configuration func.
func (p PrepareConfigurationRequest) Exec() ConfigurationRequest {
	directory := Pointer(p.Directory.Exec(p.TestStepConfigRequest))
	file := Pointer(p.File.Exec(p.TestStepConfigRequest))
	raw := Pointer(p.Raw)

	return ConfigurationRequest{
		Directory: directory,
		File:      file,
		Raw:       raw,
	}
}

// ConfigurationRequest is used by the Configuration func to determine
// the underlying type to instantiate.
type ConfigurationRequest struct {
	Directory *string
	File      *string
	Raw       *string
}

// Validate ensures that only one of Directory, File or Raw are non-empty.
func (c ConfigurationRequest) Validate() error {
	var configSet []string

	if c.Directory != nil && *c.Directory != "" {
		configSet = append(configSet, "directory")
	}

	if c.File != nil && *c.File != "" {
		configSet = append(configSet, "file")
	}

	if c.Raw != nil && *c.Raw != "" {
		configSet = append(configSet, "raw")
	}

	if len(configSet) > 1 {
		configSetStr := strings.Join(configSet, `, `)

		i := strings.LastIndex(configSetStr, ", ")

		if i != -1 {
			configSetStr = configSetStr[:i] + " and " + configSetStr[i+len(", "):]
		}

		return fmt.Errorf(`%s are populated, only one of "directory", "file", or "raw"  is allowed`, configSetStr)
	}

	return nil
}

// Configuration uses the supplied ConfigurationRequest to determine
// which of the types that implement Config to instantiate. If none
// of the fields in ConfigurationRequest are populated nil is returned.
func Configuration(req ConfigurationRequest) Config {
	if req.Directory != nil && *req.Directory != "" {
		return configurationDirectory{
			directory: *req.Directory,
		}
	}

	if req.File != nil && *req.File != "" {
		return configurationFile{
			file: *req.File,
		}
	}

	if req.Raw != nil && *req.Raw != "" {
		return configurationString{
			raw: *req.Raw,
		}
	}

	return nil
}

// copyFiles accepts a path to a directory and a destination. Only
// files in the path directory are copied, any nested directories
// are ignored.
func copyFiles(path string, dstPath string) error {
	infos, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, info := range infos {
		srcPath := filepath.Join(path, info.Name())

		if info.IsDir() {
			continue
		} else {
			err = copyFile(srcPath, dstPath)

			if err != nil {
				return err
			}
		}

	}
	return nil
}

// copyFile accepts a path to a file and a destination,
// copying the file from path to destination.
func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)

	if err != nil {
		return err
	}

	defer srcF.Close()

	di, err := os.Stat(dstPath)

	if err != nil {
		return err
	}

	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.Create(dstPath)

	if err != nil {
		return err
	}

	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}

// filesContains accepts a string representing a directory and a
// regular expression. For each file that is found within the
// directory fileContains func is called. Any nested directories
// within the directory specified by dir are ignored.
func filesContains(dir string, find *regexp.Regexp) (bool, error) {
	dirEntries, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		path := filepath.Join(dir, dirEntry.Name())

		contains, err := fileContains(path, find)

		if err != nil {
			return false, err
		}

		if contains {
			return true, nil
		}
	}

	return false, nil
}

// fileContains accepts a path and a regular expression. The
// file is read and the supplied regular expression is used
// to determine whether the file contains the specified string.
func fileContains(path string, find *regexp.Regexp) (bool, error) {
	f, err := os.ReadFile(path)

	if err != nil {
		return false, err
	}

	return find.MatchString(string(f)), nil
}

// Pointer returns a pointer to any type.
func Pointer[T any](in T) *T {
	return &in
}
