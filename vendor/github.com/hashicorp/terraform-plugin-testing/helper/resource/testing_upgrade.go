// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

// ProviderUpgradeTestCase is a specialized test case designed for executing provider
// upgrade and state-migration tests. It first provisions infrastructure using a specified
// prior published version of the provider isolated from the local development daemon.
//
// After confirming the Phase 1 setup successfully applies and settles with an empty plan,
// the test logic transitions the active working directory and its .tfstate to the
// standard local execution context defined in TestCase.
//
// This is intended to be used for two testing patterns:
//  1. Validating backwards compatibility (no unexpected drift or resource replace) when
//     users upgrade their provider without altering their configuration.
//  2. State Upgrades: Asserting that schema transitions perform successfully
//     by pinning the previous release schema version in Phase 1, then upgrading.
//
// The Phase 1 raw configuration is sourced automatically from the definition of the
// first TestStep defined in TestCase (c.TestCase.Steps[0]).
// Note: terraform developer overrides in `.terraformrc` are honoured and can cause unexpected results on Phase 1.
type ProviderUpgradeTestCase struct {
	ProviderName    string
	ProviderSource  string
	ProviderVersion string

	TestCase
}

// TestProviderUpgrade runs the two-phase provider upgrade test workflow.
func TestProviderUpgrade(t *testing.T, c ProviderUpgradeTestCase) {
	t.Helper()

	ctx := context.Background()
	ctx = logging.InitTestContext(ctx, t)

	// Verify TestCase invariants for V2
	err := c.TestCase.validate(ctx, t)
	if err != nil {
		logging.HelperResourceError(ctx, "Test validation error", map[string]any{logging.KeyError: err})
		t.Fatalf("Test validation error: %s", err)
	}

	if len(c.TestCase.Steps) != 2 {
		t.Fatalf("ProviderUpgradeTestCase requires exactly 2 steps, got %d", len(c.TestCase.Steps))
	}

	if os.Getenv(EnvTfAcc) == "" && !c.TestCase.IsUnitTest {
		t.Skipf("Acceptance tests skipped unless env '%s' set", EnvTfAcc)
		return
	}

	// Copy explicitly passed providers to factories for backwards compatibility in V2
	if len(c.TestCase.Providers) > 0 {
		c.TestCase.ProviderFactories = map[string]func() (*schema.Provider, error){}
		for name, p := range c.TestCase.Providers {
			prov := p
			c.TestCase.ProviderFactories[name] = func() (*schema.Provider, error) { //nolint:unparam
				return prov, nil
			}
		}
	}

	if c.TestCase.PreCheck != nil {
		c.TestCase.PreCheck()
	}

	sourceDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working dir: %s", err)
	}
	helper := plugintest.AutoInitProviderHelper(ctx, sourceDir)
	defer func() {
		err := helper.Close()
		if err != nil {
			logging.HelperResourceError(ctx, "Unable to clean up temporary test files", map[string]any{logging.KeyError: err})
		}
	}()

	if c.TestCase.TerraformVersionChecks != nil {
		runTFVersionChecks(ctx, t, helper.TerraformVersion(), c.TestCase.TerraformVersionChecks)
	}

	wd := helper.RequireNewWorkingDir(ctx, t, c.TestCase.WorkingDir)

	ctx = logging.TestTerraformPathContext(ctx, wd.GetHelper().TerraformExecPath())
	ctx = logging.TestWorkingDirectoryContext(ctx, wd.GetHelper().WorkingDirectory())

	logging.HelperResourceDebug(ctx, "Starting Stage 1: V1 Provider Configuration")

	// Read the raw configuration from the first test step natively
	req := teststep.PrepareConfigurationRequest{
		Directory: c.TestCase.Steps[0].ConfigDirectory,
		File:      c.TestCase.Steps[0].ConfigFile,
		Raw:       c.TestCase.Steps[0].Config,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: 1, // 1-based index for logging/hooks
			TestName:   t.Name(),
		},
	}.Exec()

	stepConfig := teststep.Configuration(req)
	if stepConfig == nil {
		t.Fatalf("Stage 1 setup failed: TestCase.Steps[0] returned nil configuration")
	}

	rawProviderConfig := fmt.Sprintf("\nterraform {\n  required_providers {\n    %s = {\n      source = \"%s\"\n      version = \"%s\"\n    }\n  }\n}\n",
		c.ProviderName, c.ProviderSource, c.ProviderVersion)

	// Append the V1 provider version pin natively
	stepConfig = stepConfig.Append(rawProviderConfig)

	err = wd.SetConfig(ctx, stepConfig, nil)
	if err != nil {
		t.Fatalf("Setup phase SetConfig failed: %s", err)
	}

	// Detach from the automatically injected test daemons so Terraform uses the remote registry
	wd.UnsetReattachInfo()

	err = wd.Init(ctx)
	if err != nil {
		t.Fatalf("Setup phase Init failed (could not acquire V1 from public registry?): %s", err)
	}

	// Apply V1 to Prime the State
	err = wd.Apply(ctx)
	if err != nil {
		_ = wd.Destroy(ctx)
		t.Fatalf("Setup phase Apply failed: %s", err)
	}

	// Validate V1 setup stability via CreatePlan & planIsEmpty check
	err = wd.CreatePlan(ctx)
	if err != nil {
		_ = wd.Destroy(ctx)
		t.Fatalf("Setup phase Plan verification failed: %s", err)
	}

	plan, err := wd.SavedPlan(ctx)
	if err != nil {
		_ = wd.Destroy(ctx)
		t.Fatalf("Setup phase Plan retrieval failed: %s", err)
	}

	if !planIsEmpty(plan, wd.GetHelper().TerraformVersion()) {
		_ = wd.Destroy(ctx)
		t.Fatalf("Setup phase produced a non-empty plan after apply, indicating the V1 configuration or provider is inherently unstable before upgrade.")
	}

	logging.HelperResourceDebug(ctx, "Finished Stage 1: V1 State Primed. Transitioning to local V2 test suite.")

	// Stage 2: Handover
	// Wipe local locks so Terraform automatically accommodates migrating registry hashes to a local daemon.
	lockFilePath := filepath.Join(wd.BaseDir(), ".terraform.lock.hcl")
	_ = os.Remove(lockFilePath)

	// Map Step 0's parameters to zero its behaviour without stripping it from the slice to ensure logging / failures
	// output cleanly as: Step 1 to `1/2` and Step 2 to `2/2`. Hacky, but ¯\_(ツ)_/¯
	c.TestCase.Steps[0].SkipFunc = func() (bool, error) { return true, nil }
	c.TestCase.Steps[0].PreConfig = nil
	c.TestCase.Steps[0].Check = nil

	// Execute standard V2 Test suite identically to `resource.Test()` over the existing, populated directory.
	executeTestCaseSteps(ctx, t, c.TestCase, wd, helper)
}
