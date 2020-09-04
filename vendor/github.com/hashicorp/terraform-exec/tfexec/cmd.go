package tfexec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	checkpointDisableEnvVar = "CHECKPOINT_DISABLE"
	logEnvVar               = "TF_LOG"
	inputEnvVar             = "TF_INPUT"
	automationEnvVar        = "TF_IN_AUTOMATION"
	logPathEnvVar           = "TF_LOG_PATH"
	reattachEnvVar          = "TF_REATTACH_PROVIDERS"

	varEnvVarPrefix = "TF_VAR_"
)

var prohibitedEnvVars = []string{
	inputEnvVar,
	automationEnvVar,
	logPathEnvVar,
	logEnvVar,
	reattachEnvVar,
}

func envMap(environ []string) map[string]string {
	env := map[string]string{}
	for _, ev := range environ {
		parts := strings.SplitN(ev, "=", 2)
		if len(parts) == 0 {
			continue
		}
		k := parts[0]
		v := ""
		if len(parts) == 2 {
			v = parts[1]
		}
		env[k] = v
	}
	return env
}

func envSlice(environ map[string]string) []string {
	env := []string{}
	for k, v := range environ {
		env = append(env, k+"="+v)
	}
	return env
}

func (tf *Terraform) buildEnv(mergeEnv map[string]string) []string {
	// set Terraform level env, if env is nil, fall back to os.Environ
	var env map[string]string
	if tf.env == nil {
		env = envMap(os.Environ())
	} else {
		env = make(map[string]string, len(tf.env))
		for k, v := range tf.env {
			env[k] = v
		}
	}

	// override env with any command specific environment
	for k, v := range mergeEnv {
		env[k] = v
	}

	// always propagate CHECKPOINT_DISABLE env var unless it is
	// explicitly overridden with tf.SetEnv or command env
	if _, ok := env[checkpointDisableEnvVar]; !ok {
		env[checkpointDisableEnvVar] = os.Getenv(checkpointDisableEnvVar)
	}

	// always override logging
	if tf.logPath == "" {
		// so logging can't pollute our stderr output
		env[logEnvVar] = ""
		env[logPathEnvVar] = ""
	} else {
		env[logPathEnvVar] = tf.logPath
		// Log levels other than TRACE are currently unreliable, the CLI recommends using TRACE only.
		env[logEnvVar] = "TRACE"
	}

	// constant automation override env vars
	env[automationEnvVar] = "1"

	return envSlice(env)
}

func (tf *Terraform) buildTerraformCmd(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, tf.execPath, args...)
	cmd.Env = tf.buildEnv(nil)
	cmd.Dir = tf.workingDir

	tf.logger.Printf("[INFO] running Terraform command: %s", cmdString(cmd))

	return cmd
}

func (tf *Terraform) runTerraformCmd(cmd *exec.Cmd) error {
	var errBuf strings.Builder

	stdout := tf.stdout
	if cmd.Stdout != nil {
		stdout = io.MultiWriter(stdout, cmd.Stdout)
	}
	cmd.Stdout = stdout

	stderr := io.MultiWriter(&errBuf, tf.stderr)
	if cmd.Stderr != nil {
		stderr = io.MultiWriter(stderr, cmd.Stderr)
	}
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return parseError(err, errBuf.String())
	}
	return nil
}
