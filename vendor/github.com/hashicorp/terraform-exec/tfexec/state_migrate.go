// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"iter"
	"os/exec"
)

type stateMigrateConfig struct {
	srcProviderLockFile string
	dstProviderLockFile string

	upgrade bool

	reattachInfo ReattachInfo
}

var defaultStateMigrateOptions = stateMigrateConfig{}

// StateMigrateOption represents options used in the StateMigrate method.
type StateMigrateOption interface {
	configureStateMigrate(*stateMigrateConfig)
}

func (opt *SourceProviderLockFileOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.srcProviderLockFile = opt.path
}

func (opt *DestinationProviderLockFileOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.dstProviderLockFile = opt.path
}

func (opt *UpgradeOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.upgrade = opt.upgrade
}

func (opt *ReattachOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.reattachInfo = opt.info
}

// StateMigrateJSON executes `terraform state migrate` with `-json` flag
// and any specified options and waits for it to complete.
//
// StateMigrateJSON is likely to be removed in a future major version in favour of
// StateMigrate returning JSON by default.
func (tf *Terraform) StateMigrateJSON(ctx context.Context, opts ...StateMigrateOption) (iter.Seq[NextMessage], error) {
	err := tf.compatible(ctx, tf1_18_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform state migrate -json was added in 1.18.0: %w", err)
	}

	cmd, err := tf.stateMigrateJSONCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return tf.runTerraformCmdJSONLog(ctx, cmd), nil
}

func (tf *Terraform) stateMigrateJSONCmd(ctx context.Context, opts ...StateMigrateOption) (*exec.Cmd, error) {
	c := defaultStateMigrateOptions

	for _, o := range opts {
		o.configureStateMigrate(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	args := tf.buildStateMigrateArgs(c)

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}

func (tf *Terraform) buildStateMigrateArgs(c stateMigrateConfig) []string {
	args := []string{
		"state", "migrate",
		"-no-color",
		"-input=false",
		"-force-copy",
		"-json",
	}

	if c.srcProviderLockFile != "" {
		args = append(args, "-source-provider-lock-file="+c.srcProviderLockFile)
	}
	if c.dstProviderLockFile != "" {
		args = append(args, "-destination-provider-lock-file="+c.dstProviderLockFile)
	}
	if c.upgrade {
		args = append(args, "-upgrade")
	}
	return args
}
