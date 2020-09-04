package tfexec

import (
	"context"
	"os/exec"
	"strconv"
)

type refreshConfig struct {
	backup      string
	dir         string
	lock        bool
	lockTimeout string
	state       string
	stateOut    string
	targets     []string
	vars        []string
	varFile     string
}

var defaultRefreshOptions = refreshConfig{
	lock:        true,
	lockTimeout: "0s",
}

type RefreshCmdOption interface {
	configureRefresh(*refreshConfig)
}

func (opt *BackupOption) configureRefresh(conf *refreshConfig) {
	conf.backup = opt.path
}

func (opt *DirOption) configureRefresh(conf *refreshConfig) {
	conf.dir = opt.path
}

func (opt *LockOption) configureRefresh(conf *refreshConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureRefresh(conf *refreshConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureRefresh(conf *refreshConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureRefresh(conf *refreshConfig) {
	conf.stateOut = opt.path
}

func (opt *TargetOption) configureRefresh(conf *refreshConfig) {
	conf.targets = append(conf.targets, opt.target)
}

func (opt *VarOption) configureRefresh(conf *refreshConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

func (opt *VarFileOption) configureRefresh(conf *refreshConfig) {
	conf.varFile = opt.path
}

func (tf *Terraform) Refresh(ctx context.Context, opts ...RefreshCmdOption) error {
	return tf.runTerraformCmd(tf.refreshCmd(ctx, opts...))
}

func (tf *Terraform) refreshCmd(ctx context.Context, opts ...RefreshCmdOption) *exec.Cmd {
	c := defaultRefreshOptions

	for _, o := range opts {
		o.configureRefresh(&c)
	}

	args := []string{"refresh", "-no-color", "-input=false"}

	// string opts: only pass if set
	if c.backup != "" {
		args = append(args, "-backup="+c.backup)
	}
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.stateOut != "" {
		args = append(args, "-state-out="+c.stateOut)
	}
	if c.varFile != "" {
		args = append(args, "-var-file="+c.varFile)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))

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
