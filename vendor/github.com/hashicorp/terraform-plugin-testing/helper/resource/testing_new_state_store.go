// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mitchellh/go-testing-interface"
)

// testStepNewStateStore will run a series of Terraform commands with the goal of ensuring that the state store (defined in config):
//   - Can be successfully initialized (validation and configuring)
//   - Can read and write state
//   - Supports workspaces (creating and deleting)
func testStepNewStateStore(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, cfg teststep.Config) error {
	t.Helper()

	err := wd.SetConfig(ctx, cfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// ----- Validate and configure the state store by running init
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	// After initialization the only workspace created should be "default"
	err = assertWorkspaces(ctx, t, wd, providers, []string{"default"})
	if err != nil {
		return fmt.Errorf("After init, expected the \"default\" workspace to be created: %w", err)
	}

	// ----- Create "foo" workspace
	err = createAndAssertEmptyWorkspace(ctx, t, wd, providers, "foo")
	if err != nil {
		return fmt.Errorf("After creating a new workspace, the state should be empty: %w", err)
	}

	// ----- Create "bar" workspace
	err = createAndAssertEmptyWorkspace(ctx, t, wd, providers, "bar")
	if err != nil {
		return fmt.Errorf("After creating a new workspace, the state should be empty: %w", err)
	}

	// ----- Apply test resources to the "bar" workspace and assert they are saved successfully in state
	err = applyTestResources(ctx, t, wd, step, providers, cfg, "bar")
	if err != nil {
		return err
	}

	// ----- Assert the "foo" workspace is still empty
	err = assertEmptyWorkspace(ctx, t, wd, providers, "foo")
	if err != nil {
		return fmt.Errorf("After writing a resource to \"bar\" state, failed assertion: %s", err)
	}

	// ----- Verify workspaces are "default", "foo" (created during this test), and "bar" (created during this test)
	err = assertWorkspaces(ctx, t, wd, providers, []string{"bar", "default", "foo"})
	if err != nil {
		return err
	}

	// ----- Delete "bar" workspace
	err = deleteWorkspace(ctx, t, wd, providers, "bar")
	if err != nil {
		return err
	}

	// ----- Attempting to delete "default" workspace should return an error
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error selecting \"foo\" workspace: %w", err)
	}
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "default", tfexec.Force(true))
	})
	if err == nil {
		return errors.New("Expected error when deleting \"default\" workspace")
	}

	// ----- Recreate the "bar" workspace and assert it is empty (i.e. no left over artifacts)
	err = createAndAssertEmptyWorkspace(ctx, t, wd, providers, "bar")
	if err != nil {
		return fmt.Errorf("After deleting, then recreating a new workspace, the state should be empty: %w", err)
	}

	// ----- Delete "bar" workspace
	err = deleteWorkspace(ctx, t, wd, providers, "bar")
	if err != nil {
		return err
	}

	// ----- List workspaces and verify it's just "foo" and "default"
	err = assertWorkspaces(ctx, t, wd, providers, []string{"default", "foo"})
	if err != nil {
		return err
	}

	// ----- Delete "foo" workspace
	err = deleteWorkspace(ctx, t, wd, providers, "foo")
	if err != nil {
		return err
	}

	// ----- List workspaces and verify it's just "default" (which we did not modify)
	err = assertWorkspaces(ctx, t, wd, providers, []string{"default"})
	if err != nil {
		return err
	}

	return nil
}

func assertWorkspaces(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories, expected []string) error {
	t.Helper()

	var err error
	actualWorkspaces := make([]string, 0)
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		actualWorkspaces, err = wd.Workspaces(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error getting workspaces: %w", err)
	}

	slices.Sort(expected)
	slices.Sort(actualWorkspaces)

	if !slices.Equal(expected, actualWorkspaces) {
		return fmt.Errorf("Expected workspaces to be %#v, got: %#v", expected, actualWorkspaces)
	}

	return nil
}

func createAndAssertEmptyWorkspace(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories, workspace string) error {
	t.Helper()

	err := runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.CreateWorkspace(ctx, workspace)
	})
	if err != nil {
		return fmt.Errorf("Error creating %q workspace: %w", workspace, err)
	}

	return assertEmptyWorkspace(ctx, t, wd, providers, workspace)
}

func assertEmptyWorkspace(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories, workspace string) error {
	t.Helper()

	err := runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, workspace)
	})
	if err != nil {
		return fmt.Errorf("Error selecting %q workspace: %w", workspace, err)
	}

	var stateObj *tfjson.State
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		stateObj, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving %q state: %w", workspace, err)
	}

	if stateObj.Values != nil && stateObj.Values.RootModule != nil && len(stateObj.Values.RootModule.Resources) > 0 {
		return fmt.Errorf("Expected %q state to be empty. Found %d resources.", workspace, len(stateObj.Values.RootModule.Resources))
	}

	return nil
}

func deleteWorkspace(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories, workspace string) error {
	t.Helper()

	// Select "default" workspace so we can delete the requested workspace
	err := runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "default")
	})
	if err != nil {
		return fmt.Errorf("Error selecting \"default\" workspace: %w", err)
	}
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, workspace, tfexec.Force(true))
	})
	if err != nil {
		return fmt.Errorf("Error deleting %q workspace: %w", workspace, err)
	}

	return nil
}

// This is the primary place that state storage is being tested, we create two test resources (using the built-in terraform_data resource) with
// two different "terraform apply" commands, then check the state for their presence.
func applyTestResources(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, cfg teststep.Config, workspace string) error {
	t.Helper()

	err := runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, workspace)
	})
	if err != nil {
		return fmt.Errorf("Error selecting %q workspace: %w", workspace, err)
	}

	// ----- Apply test resource 1 to workspace
	expectedOutput := "this resource was injected by terraform-plugin-testing"
	testResourceCfg1 := fmt.Sprintf(`
		resource "terraform_data" "tf_plugin_testing_resource_1" {
			input = %q
		}`, expectedOutput)

	cfg = cfg.Append(testResourceCfg1)
	err = wd.SetConfig(ctx, cfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error creating test resource in %q workspace: %w", workspace, err)
	}

	// ----- Apply test resource 2 to workspace
	testResourceCfg2 := fmt.Sprintf(`
		resource "terraform_data" "tf_plugin_testing_resource_2" {
			input = %q
		}`, expectedOutput)

	cfg = cfg.Append(testResourceCfg2)
	err = wd.SetConfig(ctx, cfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error creating test resource in %q workspace: %w", workspace, err)
	}

	// ----- Check if the resources exist in the state
	var stateObj *tfjson.State
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		stateObj, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving %q state: %w", workspace, err)
	}

	checkOutput := statecheck.ExpectKnownValue("terraform_data.tf_plugin_testing_resource_1", tfjsonpath.New("output"), knownvalue.StringExact(expectedOutput))
	checkResp := statecheck.CheckStateResponse{}

	checkOutput.CheckState(ctx, statecheck.CheckStateRequest{State: stateObj}, &checkResp)
	if checkResp.Error != nil {
		return fmt.Errorf("After writing a test resource instance object to %q state and re-reading it, the object has vanished: %w", workspace, checkResp.Error)
	}

	checkOutput = statecheck.ExpectKnownValue("terraform_data.tf_plugin_testing_resource_2", tfjsonpath.New("output"), knownvalue.StringExact(expectedOutput))
	checkResp = statecheck.CheckStateResponse{}

	checkOutput.CheckState(ctx, statecheck.CheckStateRequest{State: stateObj}, &checkResp)
	if checkResp.Error != nil {
		return fmt.Errorf("After writing a test resource instance object to %q state and re-reading it, the object has vanished: %w", workspace, checkResp.Error)
	}

	return nil
}
