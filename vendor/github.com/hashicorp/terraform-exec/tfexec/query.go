// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"iter"
	"os/exec"
)

type queryConfig struct {
	dir            string
	generateConfig string
	reattachInfo   ReattachInfo
	vars           []string
	varFiles       []string
}

var defaultQueryOptions = queryConfig{}

// QueryOption represents options used in the Query method.
type QueryOption interface {
	configureQuery(*queryConfig)
}

func (opt *DirOption) configureQuery(conf *queryConfig) {
	conf.dir = opt.path
}

func (opt *GenerateConfigOutOption) configureQuery(conf *queryConfig) {
	conf.generateConfig = opt.path
}

func (opt *ReattachOption) configureQuery(conf *queryConfig) {
	conf.reattachInfo = opt.info
}

func (opt *VarFileOption) configureQuery(conf *queryConfig) {
	conf.varFiles = append(conf.varFiles, opt.path)
}

func (opt *VarOption) configureQuery(conf *queryConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

// QueryJSON executes `terraform query` with the specified options as well as the
// `-json` flag and waits for it to complete.
//
// Using the `-json` flag will result in
// [machine-readable](https://developer.hashicorp.com/terraform/internals/machine-readable-ui)
// JSON being written to the supplied `io.Writer`.
//
// The returned error is nil if `terraform query` has been executed and exits
// with 0.
//
// QueryJSON is likely to be removed in a future major version in favour of
// query returning JSON by default.
func (tf *Terraform) QueryJSON(ctx context.Context, opts ...QueryOption) (iter.Seq[NextMessage], error) {
	err := tf.compatible(ctx, tf1_14_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform query -json was added in 1.14.0: %w", err)
	}

	queryCmd, err := tf.queryJSONCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return tf.runTerraformCmdJSONLog(ctx, queryCmd), nil
}

func (tf *Terraform) queryJSONCmd(ctx context.Context, opts ...QueryOption) (*exec.Cmd, error) {
	c := defaultQueryOptions

	for _, o := range opts {
		o.configureQuery(&c)
	}

	args, err := tf.buildQueryArgs(ctx, c)
	if err != nil {
		return nil, err
	}

	args = append(args, "-json")

	return tf.buildQueryCmd(ctx, c, args)
}

func (tf *Terraform) buildQueryArgs(ctx context.Context, c queryConfig) ([]string, error) {
	args := []string{"query", "-no-color"}

	if c.generateConfig != "" {
		args = append(args, "-generate-config-out="+c.generateConfig)
	}

	for _, vf := range c.varFiles {
		args = append(args, "-var-file="+vf)
	}

	if c.vars != nil {
		for _, v := range c.vars {
			args = append(args, "-var", v)
		}
	}

	return args, nil
}

func (tf *Terraform) buildQueryCmd(ctx context.Context, c queryConfig, args []string) (*exec.Cmd, error) {
	// optional positional argument
	if c.dir != "" {
		args = append(args, c.dir)
	}

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
