// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type providersMirrorConfig struct {
	lockFile  bool
	platforms []string
}

var defaultProvidersMirrorOptions = providersMirrorConfig{
	// defaults to true
	// See https://github.com/hashicorp/terraform/blob/v1.14.0/internal/command/providers_mirror.go#L42
	lockFile: true,
}

type ProvidersMirrorOption interface {
	configureProvidersMirror(*providersMirrorConfig)
}

func (opt *LockFileOption) configureProvidersMirror(conf *providersMirrorConfig) {
	conf.lockFile = opt.useLockFile
}

func (opt *PlatformOption) configureProvidersMirror(conf *providersMirrorConfig) {
	conf.platforms = append(conf.platforms, opt.platform)
}

// ProvidersMirror represents the `terraform providers mirror` command
func (tf *Terraform) ProvidersMirror(ctx context.Context, targetDir string, opts ...ProvidersMirrorOption) error {
	err := tf.compatible(ctx, tf0_13_0, nil)
	if err != nil {
		return fmt.Errorf("terraform providers mirror was added in 0.13.0: %w", err)
	}
	if targetDir == "" {
		return fmt.Errorf("targetDir argument needs to be set")
	}

	mirrorCmd, err := tf.providersMirrorCmd(ctx, targetDir, opts...)
	if err != nil {
		return err
	}

	err = tf.runTerraformCmd(ctx, mirrorCmd)
	if err != nil {
		return err
	}

	return err
}

func (tf *Terraform) providersMirrorCmd(ctx context.Context, targetDir string, opts ...ProvidersMirrorOption) (*exec.Cmd, error) {
	c := defaultProvidersMirrorOptions

	args := []string{"providers", "mirror"}

	for _, o := range opts {
		o.configureProvidersMirror(&c)
	}

	for _, p := range c.platforms {
		args = append(args, "-platform="+p)
	}

	// lockFile is true by default, so only pass the flag if the caller has set it
	// to false
	if !c.lockFile {
		err := tf.compatible(ctx, tf1_10_0, nil)
		if err != nil {
			return nil, fmt.Errorf("lock-file option was introduced in Terraform 1.10.0: %w", err)
		}

		args = append(args, "-lock-file=false")
	}

	args = append(args, targetDir)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
