// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
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

	// Parse human output into a list of workspaces and the current workspace
	ws, current := parseWorkspaceListHumanOutput(outBuf.String())

	return ws, current, nil
}

func (tf *Terraform) workspaceListCmd(ctx context.Context, opts ...WorkspaceListOption) (*exec.Cmd, error) {
	c := defaultWorkspaceListOptions

	for _, o := range opts {
		o.configureWorkspaceList(&c)
	}

	args := []string{"workspace", "list", "-no-color"}
	return tf.buildWorkspaceListCmd(ctx, c, args)
}

// parseWorkspaceListHumanOutput parses human output from the workspace list command
// to return a slice of workspace names and the current workspace name
func parseWorkspaceListHumanOutput(stdout string) ([]string, string) {
	var currentWorkspacePrefix string = "* "

	lines := strings.Split(stdout, "\n")

	current := ""
	workspaces := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// results from successive newlines in stdout string; ignore
			continue
		}

		if strings.HasPrefix(line, currentWorkspacePrefix) {
			line = strings.TrimPrefix(line, currentWorkspacePrefix)
			current = line
		}

		if line != url.PathEscape(line) {
			// If a line isn't URL-escaped then it's not a valid workspace name.
			// If we're encountering a processed line string that isn't a valid workspace name it's
			// almost certainly due to a warning being present in stdout.
			//
			// See here for prior use of URL-escaping to validate workspace names:
			// https://github.com/hashicorp/terraform/blob/9fb1177c3566c247c72e6c36ca219463c4b6be2b/internal/command/workspace_command.go#L44-L52
			//
			// tl;dr this line isn't a workspace name, so lets skip and see if there is a list of
			// workspaces later in the output, under any warning text.
			continue
		}

		workspaces = append(workspaces, line)
	}

	return workspaces, current
}

// WorkspaceListJSON represents the terraform workspace list subcommand with the `-json` flag.
//
// WorkspaceListJSON is likely to be removed in a future major version, in favour of
// WorkspaceList handling JSON by default instead of human output.
func (tf *Terraform) WorkspaceListJSON(ctx context.Context, opts ...WorkspaceListOption) (*tfjson.WorkspaceListOutput, error) {
	err := tf.compatible(ctx, tf1_16_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform workspace list -json was added in 1.16.0: %w", err)
	}

	cmd, err := tf.workspaceListJSONCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var output tfjson.WorkspaceListOutput

	err = tf.runTerraformCmdJSON(ctx, cmd, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (tf *Terraform) workspaceListJSONCmd(ctx context.Context, opts ...WorkspaceListOption) (*exec.Cmd, error) {
	c := defaultWorkspaceListOptions

	for _, o := range opts {
		o.configureWorkspaceList(&c)
	}

	args := []string{"workspace", "list", "-json"}
	return tf.buildWorkspaceListCmd(ctx, c, args)
}

func (tf *Terraform) buildWorkspaceListCmd(ctx context.Context, c workspaceListConfig, args []string) (*exec.Cmd, error) {
	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
