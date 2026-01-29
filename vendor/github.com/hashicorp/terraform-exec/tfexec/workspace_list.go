// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
	"strings"
)

type workspaceListConfig struct {
	reattachInfo ReattachInfo
}

var defaultWorkspaceListOptions = workspaceListConfig{}

type WorkspaceListOption interface {
	configureWorkspaceList(*workspaceListConfig)
}

func (opt *ReattachOption) configureWorkspaceList(conf *workspaceListConfig) {
	conf.reattachInfo = opt.info
}

// WorkspaceList represents the workspace list subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceList(ctx context.Context, opts ...WorkspaceListOption) ([]string, string, error) {
	wlCmd, err := tf.workspaceListCmd(ctx, opts...)
	if err != nil {
		return nil, "", err
	}

	var outBuf strings.Builder
	wlCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(ctx, wlCmd)
	if err != nil {
		return nil, "", err
	}

	ws, current := parseWorkspaceList(outBuf.String())

	return ws, current, nil
}

const currentWorkspacePrefix = "* "

func (tf *Terraform) workspaceListCmd(ctx context.Context, opts ...WorkspaceListOption) (*exec.Cmd, error) {
	c := defaultWorkspaceListOptions

	for _, o := range opts {
		o.configureWorkspaceList(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, "workspace", "list", "-no-color"), nil
}

func parseWorkspaceList(stdout string) ([]string, string) {
	lines := strings.Split(stdout, "\n")

	current := ""
	workspaces := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, currentWorkspacePrefix) {
			line = strings.TrimPrefix(line, currentWorkspacePrefix)
			current = line
		}
		workspaces = append(workspaces, line)
	}

	return workspaces, current
}
