package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

// Show reads the default state path and outputs the state.
// To read a state or plan file, ShowState or ShowPlan must be used instead.
func (tf *Terraform) Show(ctx context.Context) (*tfjson.State, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	showCmd := tf.showCmd(ctx)

	var ret tfjson.State
	var outBuf bytes.Buffer
	showCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(showCmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowStateFile reads a given state file and outputs the state.
func (tf *Terraform) ShowStateFile(ctx context.Context, statePath string) (*tfjson.State, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	if statePath == "" {
		return nil, fmt.Errorf("statePath cannot be blank: use Show() if not passing statePath")
	}

	showCmd := tf.showCmd(ctx, statePath)

	var ret tfjson.State
	var outBuf bytes.Buffer
	showCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(showCmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowPlanFile reads a given plan file and outputs the plan.
func (tf *Terraform) ShowPlanFile(ctx context.Context, planPath string) (*tfjson.Plan, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	if planPath == "" {
		return nil, fmt.Errorf("planPath cannot be blank: use Show() if not passing planPath")
	}

	showCmd := tf.showCmd(ctx, planPath)

	var ret tfjson.Plan
	var outBuf bytes.Buffer
	showCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(showCmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

func (tf *Terraform) showCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"show", "-json", "-no-color"}
	allArgs = append(allArgs, args...)

	return tf.buildTerraformCmd(ctx, allArgs...)
}
