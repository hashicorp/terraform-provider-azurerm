// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type workspaceShowConfig struct {
	reattachInfo ReattachInfo
}

var defaultWorkspaceShowOptions = workspaceShowConfig{}

type WorkspaceShowOption interface {
	configureWorkspaceShow(*workspaceShowConfig)
}

func (opt *ReattachOption) configureWorkspaceShow(conf *workspaceShowConfig) {
	conf.reattachInfo = opt.info
}

// WorkspaceShow represents the workspace show subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceShow(ctx context.Context, opts ...WorkspaceShowOption) (string, error) {
	workspaceShowCmd, err := tf.workspaceShowCmd(ctx, opts...)
	if err != nil {
		return "", err
	}

	var outBuffer strings.Builder
	workspaceShowCmd.Stdout = &outBuffer

	err = tf.runTerraformCmd(ctx, workspaceShowCmd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(outBuffer.String()), nil
}

func (tf *Terraform) workspaceShowCmd(ctx context.Context, opts ...WorkspaceShowOption) (*exec.Cmd, error) {
	err := tf.compatible(ctx, tf0_10_0, nil)
	if err != nil {
		return nil, fmt.Errorf("workspace show was first introduced in Terraform 0.10.0: %w", err)
	}

	c := defaultWorkspaceShowOptions

	for _, o := range opts {
		o.configureWorkspaceShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, "workspace", "show", "-no-color"), nil
}
