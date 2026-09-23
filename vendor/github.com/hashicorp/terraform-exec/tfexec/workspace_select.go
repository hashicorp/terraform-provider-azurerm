// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
)

type workspaceSelectConfig struct {
	reattachInfo ReattachInfo
}

var defaultWorkspaceSelectOptions = workspaceSelectConfig{}

type WorkspaceSelectOption interface {
	configureWorkspaceSelect(*workspaceSelectConfig)
}

func (opt *ReattachOption) configureWorkspaceSelect(conf *workspaceSelectConfig) {
	conf.reattachInfo = opt.info
}

// WorkspaceSelect represents the workspace select subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceSelect(ctx context.Context, workspace string, opts ...WorkspaceSelectOption) error {
	cmd, err := tf.workspaceSelectCmd(ctx, workspace, opts...)
	if err != nil {
		return err
	}

	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) workspaceSelectCmd(ctx context.Context, workspace string, opts ...WorkspaceSelectOption) (*exec.Cmd, error) {
	c := defaultWorkspaceSelectOptions

	for _, o := range opts {
		o.configureWorkspaceSelect(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, "workspace", "select", "-no-color", workspace), nil
}
