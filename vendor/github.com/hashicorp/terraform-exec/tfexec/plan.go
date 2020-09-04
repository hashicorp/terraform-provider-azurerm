package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
)

type planConfig struct {
	destroy     bool
	dir         string
	lock        bool
	lockTimeout string
	out         string
	parallelism int
	refresh     bool
	state       string
	targets     []string
	vars        []string
	varFile     string
}

var defaultPlanOptions = planConfig{
	destroy:     false,
	lock:        true,
	lockTimeout: "0s",
	parallelism: 10,
	refresh:     true,
}

type PlanOption interface {
	configurePlan(*planConfig)
}

func (opt *DirOption) configurePlan(conf *planConfig) {
	conf.dir = opt.path
}

func (opt *VarFileOption) configurePlan(conf *planConfig) {
	conf.varFile = opt.path
}

func (opt *VarOption) configurePlan(conf *planConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

func (opt *TargetOption) configurePlan(conf *planConfig) {
	conf.targets = append(conf.targets, opt.target)
}

func (opt *StateOption) configurePlan(conf *planConfig) {
	conf.state = opt.path
}

func (opt *RefreshOption) configurePlan(conf *planConfig) {
	conf.refresh = opt.refresh
}

func (opt *ParallelismOption) configurePlan(conf *planConfig) {
	conf.parallelism = opt.parallelism
}

func (opt *OutOption) configurePlan(conf *planConfig) {
	conf.out = opt.path
}

func (opt *LockTimeoutOption) configurePlan(conf *planConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *LockOption) configurePlan(conf *planConfig) {
	conf.lock = opt.lock
}

func (opt *DestroyFlagOption) configurePlan(conf *planConfig) {
	conf.destroy = opt.destroy
}

func (tf *Terraform) Plan(ctx context.Context, opts ...PlanOption) error {
	return tf.runTerraformCmd(tf.planCmd(ctx, opts...))
}

func (tf *Terraform) planCmd(ctx context.Context, opts ...PlanOption) *exec.Cmd {
	c := defaultPlanOptions

	for _, o := range opts {
		o.configurePlan(&c)
	}

	args := []string{"plan", "-no-color", "-input=false"}

	// string opts: only pass if set
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if c.out != "" {
		args = append(args, "-out="+c.out)
	}
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.varFile != "" {
		args = append(args, "-var-file="+c.varFile)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))
	args = append(args, "-parallelism="+fmt.Sprint(c.parallelism))
	args = append(args, "-refresh="+strconv.FormatBool(c.refresh))

	// unary flags: pass if true
	if c.destroy {
		args = append(args, "-destroy")
	}

	// string slice opts: split into separate args
	if c.targets != nil {
		for _, ta := range c.targets {
			args = append(args, "-target="+ta)
		}
	}
	if c.vars != nil {
		for _, v := range c.vars {
			args = append(args, "-var", v)
		}
	}

	// optional positional argument
	if c.dir != "" {
		args = append(args, c.dir)
	}

	return tf.buildTerraformCmd(ctx, args...)
}
