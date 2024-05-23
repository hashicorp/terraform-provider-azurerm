// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

// expectNonEmptyPlanOutputChangesMinTFVersion is used to keep compatibility for
// Terraform 0.12 and 0.13 after enabling ExpectNonEmptyPlan to check output
// changes. Those older versions will always show outputs being created.
var expectNonEmptyPlanOutputChangesMinTFVersion = tfversion.Version0_14_0

func testStepNewConfig(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, stepIndex int, helper *plugintest.Helper) error {
	t.Helper()

	configRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       step.Config,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepIndex + 1,
			TestName:   t.Name(),
		},
	}.Exec()

	cfg := teststep.Configuration(configRequest)

	var hasTerraformBlock bool
	var hasProviderBlock bool

	if cfg != nil {
		var err error

		hasTerraformBlock, err = cfg.HasTerraformBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains terraform block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains terraform block: %s", err)
		}

		hasProviderBlock, err = cfg.HasProviderBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains provider block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains provider block: %s", err)
		}
	}

	mergedConfig, err := step.mergedConfig(ctx, c, hasTerraformBlock, hasProviderBlock, helper.TerraformVersion())

	if err != nil {
		logging.HelperResourceError(ctx,
			"Error generating merged configuration",
			map[string]interface{}{logging.KeyError: err},
		)
		t.Fatalf("Error generating merged configuration: %s", err)
	}

	confRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       mergedConfig,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepIndex + 1,
			TestName:   t.Name(),
		},
	}.Exec()

	testStepConfig := teststep.Configuration(confRequest)

	err = wd.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// If this step is a PlanOnly step, skip over this first Plan and
	// subsequent Apply, and use the follow-up Plan that checks for
	// permadiffs
	if !step.PlanOnly {
		logging.HelperResourceDebug(ctx, "Running Terraform CLI plan and apply")

		// Plan!
		err := runProviderCommand(ctx, t, func() error {
			var opts []tfexec.PlanOption
			if step.Destroy {
				opts = append(opts, tfexec.Destroy(true))
			}

			if c.AdditionalCLIOptions != nil && c.AdditionalCLIOptions.Plan.AllowDeferral {
				opts = append(opts, tfexec.AllowDeferral(true))
			}

			return wd.CreatePlan(ctx, opts...)
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error running pre-apply plan: %w", err)
		}

		// Run pre-apply plan checks
		if len(step.ConfigPlanChecks.PreApply) > 0 {
			var plan *tfjson.Plan
			err = runProviderCommand(ctx, t, func() error {
				var err error
				plan, err = wd.SavedPlan(ctx)
				return err
			}, wd, providers)
			if err != nil {
				return fmt.Errorf("Error retrieving pre-apply plan: %w", err)
			}

			err = runPlanChecks(ctx, t, plan, step.ConfigPlanChecks.PreApply)
			if err != nil {
				return fmt.Errorf("Pre-apply plan check(s) failed:\n%w", err)
			}
		}

		// We need to keep a copy of the state prior to destroying such
		// that the destroy steps can verify their behavior in the
		// check function
		var stateBeforeApplication *terraform.State

		if step.Check != nil && step.Destroy {
			// Refresh the state before shimming it for destroy checks later.
			// This re-implements previously existing test step logic for the
			// specific situation that a provider developer has applied a
			// resource with a previous schema version and is destroying it with
			// a provider that has a newer schema version. Without this refresh
			// the shim logic will return an error such as:
			//
			//    Failed to marshal state to json: schema version 0 for null_resource.test in state does not match version 1 from the provider
			err := runProviderCommand(ctx, t, func() error {
				return wd.Refresh(ctx)
			}, wd, providers)

			if err != nil {
				return fmt.Errorf("Error running pre-apply refresh: %w", err)
			}

			err = runProviderCommand(ctx, t, func() error {
				stateBeforeApplication, err = getState(ctx, t, wd)
				if err != nil {
					return err
				}
				return nil
			}, wd, providers)

			if err != nil {
				return fmt.Errorf("Error retrieving pre-apply state: %w", err)
			}
		}

		// Apply the diff, creating real resources
		err = runProviderCommand(ctx, t, func() error {
			var opts []tfexec.ApplyOption

			if c.AdditionalCLIOptions != nil && c.AdditionalCLIOptions.Apply.AllowDeferral {
				opts = append(opts, tfexec.AllowDeferral(true))
			}

			return wd.Apply(ctx, opts...)
		}, wd, providers)
		if err != nil {
			if step.Destroy {
				return fmt.Errorf("Error running destroy: %w", err)
			}
			return fmt.Errorf("Error running apply: %w", err)
		}

		// Run any configured checks
		if step.Check != nil {
			logging.HelperResourceTrace(ctx, "Using TestStep Check")

			if step.Destroy {
				if err := step.Check(stateBeforeApplication); err != nil {
					return fmt.Errorf("Check failed: %w", err)
				}
			} else {
				var state *terraform.State

				err := runProviderCommand(ctx, t, func() error {
					state, err = getState(ctx, t, wd)
					if err != nil {
						return err
					}
					return nil
				}, wd, providers)

				if err != nil {
					return fmt.Errorf("Error retrieving state after apply: %w", err)
				}

				if err := step.Check(state); err != nil {
					return fmt.Errorf("Check failed: %w", err)
				}
			}
		}

		// Run state checks
		if len(step.ConfigStateChecks) > 0 {
			var state *tfjson.State

			err = runProviderCommand(ctx, t, func() error {
				var err error
				state, err = wd.State(ctx)
				return err
			}, wd, providers)

			if err != nil {
				return fmt.Errorf("Error retrieving post-apply, post-refresh state: %w", err)
			}

			err = runStateChecks(ctx, t, state, step.ConfigStateChecks)
			if err != nil {
				return fmt.Errorf("Post-apply refresh state check(s) failed:\n%w", err)
			}
		}
	}

	// Test for perpetual diffs by performing a plan, a refresh, and another plan
	logging.HelperResourceDebug(ctx, "Running Terraform CLI plan to check for perpetual differences")

	// do a plan
	err = runProviderCommand(ctx, t, func() error {
		opts := []tfexec.PlanOption{
			tfexec.Refresh(false),
		}
		if step.Destroy {
			opts = append(opts, tfexec.Destroy(true))
		}

		if c.AdditionalCLIOptions != nil && c.AdditionalCLIOptions.Plan.AllowDeferral {
			opts = append(opts, tfexec.AllowDeferral(true))
		}

		return wd.CreatePlan(ctx, opts...)
	}, wd, providers)
	if err != nil {
		if step.PlanOnly {
			return fmt.Errorf("Error running non-refresh plan: %w", err)
		}

		return fmt.Errorf("Error running post-apply non-refresh plan: %w", err)
	}

	var plan *tfjson.Plan
	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		if step.PlanOnly {
			return fmt.Errorf("Error reading saved non-refresh plan: %w", err)
		}

		return fmt.Errorf("Error reading saved post-apply non-refresh plan: %w", err)
	}

	// Run post-apply, pre-refresh plan checks
	if len(step.ConfigPlanChecks.PostApplyPreRefresh) > 0 {
		err = runPlanChecks(ctx, t, plan, step.ConfigPlanChecks.PostApplyPreRefresh)
		if err != nil {
			if step.PlanOnly {
				return fmt.Errorf("Non-refresh plan checks(s) failed:\n%w", err)
			}

			return fmt.Errorf("Post-apply, pre-refresh plan check(s) failed:\n%w", err)
		}
	}

	if !planIsEmpty(plan, helper.TerraformVersion()) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error reading saved human-readable non-refresh plan output: %w", err)
		}

		if step.PlanOnly {
			return fmt.Errorf("The non-refresh plan was not empty.\nstdout:\n\n%s", stdout)
		}

		return fmt.Errorf("After applying this test step, the non-refresh plan was not empty.\nstdout:\n\n%s", stdout)
	}

	// do another plan
	err = runProviderCommand(ctx, t, func() error {
		var opts []tfexec.PlanOption
		if step.Destroy {
			opts = append(opts, tfexec.Destroy(true))

			if step.PreventPostDestroyRefresh {
				opts = append(opts, tfexec.Refresh(false))
			}
		}

		if c.AdditionalCLIOptions != nil && c.AdditionalCLIOptions.Plan.AllowDeferral {
			opts = append(opts, tfexec.AllowDeferral(true))
		}

		return wd.CreatePlan(ctx, opts...)
	}, wd, providers)
	if err != nil {
		if step.PlanOnly {
			return fmt.Errorf("Error running refresh plan: %w", err)
		}

		return fmt.Errorf("Error running post-apply refresh plan: %w", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		if step.PlanOnly {
			return fmt.Errorf("Error reading refresh plan: %w", err)
		}

		return fmt.Errorf("Error reading post-apply refresh plan: %w", err)
	}

	// Run post-apply, post-refresh plan checks
	if len(step.ConfigPlanChecks.PostApplyPostRefresh) > 0 {
		err = runPlanChecks(ctx, t, plan, step.ConfigPlanChecks.PostApplyPostRefresh)
		if err != nil {
			return fmt.Errorf("Post-apply refresh plan check(s) failed:\n%w", err)
		}
	}

	// check if plan is empty
	if !planIsEmpty(plan, helper.TerraformVersion()) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error reading human-readable refresh plan output: %w", err)
		}

		if step.PlanOnly {
			return fmt.Errorf("The refresh plan was not empty.\nstdout\n\n%s", stdout)
		}

		return fmt.Errorf("After applying this test step, the refresh plan was not empty.\nstdout\n\n%s", stdout)
	} else if step.ExpectNonEmptyPlan && planIsEmpty(plan, helper.TerraformVersion()) {
		return errors.New("Expected a non-empty plan, but got an empty refresh plan")
	}

	// ID-ONLY REFRESH
	// If we've never checked an id-only refresh and our state isn't
	// empty, find the first resource and test it.
	if c.IDRefreshName != "" {
		logging.HelperResourceTrace(ctx, "Using TestCase IDRefreshName")

		var state *terraform.State

		err = runProviderCommand(ctx, t, func() error {
			state, err = getState(ctx, t, wd)
			if err != nil {
				return err
			}
			return nil
		}, wd, providers)

		if err != nil {
			return err
		}

		//nolint:staticcheck // legacy usage
		if state.Empty() {
			return nil
		}

		var idRefreshCheck *terraform.ResourceState

		// Find the first non-nil resource in the state
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[c.IDRefreshName]; ok {
					idRefreshCheck = v
				}

				break
			}
		}

		// If we have an instance to check for refreshes, do it
		// immediately. We do it in the middle of another test
		// because it shouldn't affect the overall state (refresh
		// is read-only semantically) and we want to fail early if
		// this fails. If refresh isn't read-only, then this will have
		// caught a different bug.
		if idRefreshCheck != nil {
			if err := testIDRefresh(ctx, t, c, wd, step, idRefreshCheck, providers, stepIndex, helper); err != nil {
				return fmt.Errorf(
					"[ERROR] Test: ID-only test failed: %s", err)
			}
		}
	}

	return nil
}
