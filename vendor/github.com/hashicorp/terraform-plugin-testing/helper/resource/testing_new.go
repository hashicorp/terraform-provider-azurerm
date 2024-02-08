// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func runPostTestDestroy(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, providers *providerFactories, statePreDestroy *terraform.State) error {
	t.Helper()

	err := runProviderCommand(ctx, t, func() error {
		return wd.Destroy(ctx)
	}, wd, providers)
	if err != nil {
		return err
	}

	if c.CheckDestroy != nil {
		logging.HelperResourceTrace(ctx, "Using TestCase CheckDestroy")
		logging.HelperResourceDebug(ctx, "Calling TestCase CheckDestroy")

		if err := c.CheckDestroy(statePreDestroy); err != nil {
			return err
		}

		logging.HelperResourceDebug(ctx, "Called TestCase CheckDestroy")
	}

	return nil
}

func runNewTest(ctx context.Context, t testing.T, c TestCase, helper *plugintest.Helper) {
	t.Helper()

	wd := helper.RequireNewWorkingDir(ctx, t, c.WorkingDir)

	ctx = logging.TestTerraformPathContext(ctx, wd.GetHelper().TerraformExecPath())
	ctx = logging.TestWorkingDirectoryContext(ctx, wd.GetHelper().WorkingDirectory())

	providers := &providerFactories{
		legacy:  c.ProviderFactories,
		protov5: c.ProtoV5ProviderFactories,
		protov6: c.ProtoV6ProviderFactories,
	}

	defer func() {
		var statePreDestroy *terraform.State
		var err error
		err = runProviderCommand(ctx, t, func() error {
			statePreDestroy, err = getState(ctx, t, wd)
			if err != nil {
				return err
			}
			return nil
		}, wd, providers)
		if err != nil {
			logging.HelperResourceError(ctx,
				"Error retrieving state, there may be dangling resources",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error retrieving state, there may be dangling resources: %s", err.Error())
			return
		}

		if !stateIsEmpty(statePreDestroy) {
			err := runPostTestDestroy(ctx, t, c, wd, providers, statePreDestroy)
			if err != nil {
				logging.HelperResourceError(ctx,
					"Error running post-test destroy, there may be dangling resources",
					map[string]interface{}{logging.KeyError: err},
				)
				t.Fatalf("Error running post-test destroy, there may be dangling resources: %s", err.Error())
			}
		}

		wd.Close()
	}()

	// Return value from c.ProviderConfig() is assigned to Raw as this was previously being
	// passed to wd.SetConfig() when the second argument accept a configuration string.
	if c.hasProviders(ctx) {
		config := teststep.Configuration(
			teststep.ConfigurationRequest{
				Raw: teststep.Pointer(c.providerConfig(ctx, false)),
			},
		)

		err := wd.SetConfig(ctx, config, nil)

		if err != nil {
			logging.HelperResourceError(ctx,
				"TestCase error setting provider configuration",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("TestCase error setting provider configuration: %s", err)
		}

		err = runProviderCommand(ctx, t, func() error {
			return wd.Init(ctx)
		}, wd, providers)

		if err != nil {
			logging.HelperResourceError(ctx,
				"TestCase error running init",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("TestCase error running init: %s", err.Error())
		}
	}

	logging.HelperResourceDebug(ctx, "Starting TestSteps")

	// use this to track last step successfully applied
	// acts as default for import tests
	var appliedCfg teststep.Config
	var stepNumber int

	for stepIndex, step := range c.Steps {
		if stepNumber > 0 {
			copyWorkingDir(ctx, t, stepNumber, wd)
		}

		stepNumber = stepIndex + 1 // 1-based indexing for humans

		configRequest := teststep.PrepareConfigurationRequest{
			Directory: step.ConfigDirectory,
			File:      step.ConfigFile,
			Raw:       step.Config,
			TestStepConfigRequest: config.TestStepConfigRequest{
				StepNumber: stepNumber,
				TestName:   t.Name(),
			},
		}.Exec()

		cfg := teststep.Configuration(configRequest)

		ctx = logging.TestStepNumberContext(ctx, stepNumber)

		logging.HelperResourceDebug(ctx, "Starting TestStep")

		if step.PreConfig != nil {
			logging.HelperResourceDebug(ctx, "Calling TestStep PreConfig")
			step.PreConfig()
			logging.HelperResourceDebug(ctx, "Called TestStep PreConfig")
		}

		if step.SkipFunc != nil {
			logging.HelperResourceDebug(ctx, "Calling TestStep SkipFunc")

			skip, err := step.SkipFunc()
			if err != nil {
				logging.HelperResourceError(ctx,
					"Error calling TestStep SkipFunc",
					map[string]interface{}{logging.KeyError: err},
				)
				t.Fatalf("Error calling TestStep SkipFunc: %s", err.Error())
			}

			logging.HelperResourceDebug(ctx, "Called TestStep SkipFunc")

			if skip {
				t.Logf("Skipping step %d/%d due to SkipFunc", stepNumber, len(c.Steps))
				logging.HelperResourceWarn(ctx, "Skipping TestStep due to SkipFunc")
				continue
			}
		}

		if cfg != nil && !step.Destroy && len(step.Taint) > 0 {
			err := testStepTaint(ctx, step, wd)

			if err != nil {
				logging.HelperResourceError(ctx,
					"TestStep error tainting resources",
					map[string]interface{}{logging.KeyError: err},
				)
				t.Fatalf("TestStep %d/%d error tainting resources: %s", stepNumber, len(c.Steps), err)
			}
		}

		hasProviders, err := step.hasProviders(ctx, stepIndex, t.Name())

		if err != nil {
			logging.HelperResourceError(ctx,
				"TestStep error checking for providers",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("TestStep %d/%d error checking for providers: %s", stepNumber, len(c.Steps), err)
		}

		if hasProviders {
			providers = &providerFactories{
				legacy:  sdkProviderFactories(c.ProviderFactories).merge(step.ProviderFactories),
				protov5: protov5ProviderFactories(c.ProtoV5ProviderFactories).merge(step.ProtoV5ProviderFactories),
				protov6: protov6ProviderFactories(c.ProtoV6ProviderFactories).merge(step.ProtoV6ProviderFactories),
			}

			var hasProviderBlock bool

			if cfg != nil {
				hasProviderBlock, err = cfg.HasProviderBlock(ctx)

				if err != nil {
					logging.HelperResourceError(ctx,
						"TestStep error determining whether configuration contains provider block",
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("TestStep %d/%d error determining whether configuration contains provider block: %s", stepNumber, len(c.Steps), err)
				}
			}

			var testStepConfig teststep.Config

			// Return value from step.providerConfig() is assigned to Raw as this was previously being
			// passed to wd.SetConfig() directly when the second argument to wd.SetConfig() accepted a
			// configuration string.
			confRequest := teststep.PrepareConfigurationRequest{
				Directory: step.ConfigDirectory,
				File:      step.ConfigFile,
				Raw:       step.providerConfig(ctx, hasProviderBlock),
				TestStepConfigRequest: config.TestStepConfigRequest{
					StepNumber: stepIndex + 1,
					TestName:   t.Name(),
				},
			}.Exec()

			testStepConfig = teststep.Configuration(confRequest)

			err = wd.SetConfig(ctx, testStepConfig, step.ConfigVariables)

			if err != nil {
				logging.HelperResourceError(ctx,
					"TestStep error setting provider configuration",
					map[string]interface{}{logging.KeyError: err},
				)
				t.Fatalf("TestStep %d/%d error setting test provider configuration: %s", stepNumber, len(c.Steps), err)
			}

			err = runProviderCommand(
				ctx,
				t,
				func() error {
					return wd.Init(ctx)
				},
				wd,
				providers,
			)

			if err != nil {
				logging.HelperResourceError(ctx,
					"TestStep error running init",
					map[string]interface{}{logging.KeyError: err},
				)
				t.Fatalf("TestStep %d/%d running init: %s", stepNumber, len(c.Steps), err.Error())
				return
			}
		}

		if step.ImportState {
			logging.HelperResourceTrace(ctx, "TestStep is ImportState mode")

			err := testStepNewImportState(ctx, t, helper, wd, step, appliedCfg, providers, stepIndex)
			if step.ExpectError != nil {
				logging.HelperResourceDebug(ctx, "Checking TestStep ExpectError")
				if err == nil {
					logging.HelperResourceError(ctx,
						"Error running import: expected an error but got none",
					)
					t.Fatalf("Step %d/%d error running import: expected an error but got none", stepNumber, len(c.Steps))
				}
				if !step.ExpectError.MatchString(err.Error()) {
					logging.HelperResourceError(ctx,
						fmt.Sprintf("Error running import: expected an error with pattern (%s)", step.ExpectError.String()),
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d error running import, expected an error with pattern (%s), no match on: %s", stepNumber, len(c.Steps), step.ExpectError.String(), err)
				}
			} else {
				if err != nil && c.ErrorCheck != nil {
					logging.HelperResourceDebug(ctx, "Calling TestCase ErrorCheck")
					err = c.ErrorCheck(err)
					logging.HelperResourceDebug(ctx, "Called TestCase ErrorCheck")
				}
				if err != nil {
					logging.HelperResourceError(ctx,
						"Error running import",
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d error running import: %s", stepNumber, len(c.Steps), err)
				}
			}

			logging.HelperResourceDebug(ctx, "Finished TestStep")

			continue
		}

		if step.RefreshState {
			logging.HelperResourceTrace(ctx, "TestStep is RefreshState mode")

			err := testStepNewRefreshState(ctx, t, wd, step, providers)
			if step.ExpectError != nil {
				logging.HelperResourceDebug(ctx, "Checking TestStep ExpectError")
				if err == nil {
					logging.HelperResourceError(ctx,
						"Error running refresh: expected an error but got none",
					)
					t.Fatalf("Step %d/%d error running refresh: expected an error but got none", stepNumber, len(c.Steps))
				}
				if !step.ExpectError.MatchString(err.Error()) {
					logging.HelperResourceError(ctx,
						fmt.Sprintf("Error running refresh: expected an error with pattern (%s)", step.ExpectError.String()),
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d error running refresh, expected an error with pattern (%s), no match on: %s", stepNumber, len(c.Steps), step.ExpectError.String(), err)
				}
			} else {
				if err != nil && c.ErrorCheck != nil {
					logging.HelperResourceDebug(ctx, "Calling TestCase ErrorCheck")
					err = c.ErrorCheck(err)
					logging.HelperResourceDebug(ctx, "Called TestCase ErrorCheck")
				}
				if err != nil {
					logging.HelperResourceError(ctx,
						"Error running refresh",
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d error running refresh: %s", stepNumber, len(c.Steps), err)
				}
			}

			logging.HelperResourceDebug(ctx, "Finished TestStep")

			continue
		}

		if cfg != nil {
			logging.HelperResourceTrace(ctx, "TestStep is Config mode")

			err := testStepNewConfig(ctx, t, c, wd, step, providers, stepIndex)
			if step.ExpectError != nil {
				logging.HelperResourceDebug(ctx, "Checking TestStep ExpectError")

				if err == nil {
					logging.HelperResourceError(ctx,
						"Expected an error but got none",
					)
					t.Fatalf("Step %d/%d, expected an error but got none", stepNumber, len(c.Steps))
				}
				if !step.ExpectError.MatchString(err.Error()) {
					logging.HelperResourceError(ctx,
						fmt.Sprintf("Expected an error with pattern (%s)", step.ExpectError.String()),
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d, expected an error with pattern, no match on: %s", stepNumber, len(c.Steps), err)
				}
			} else {
				if err != nil && c.ErrorCheck != nil {
					logging.HelperResourceDebug(ctx, "Calling TestCase ErrorCheck")

					err = c.ErrorCheck(err)

					logging.HelperResourceDebug(ctx, "Called TestCase ErrorCheck")
				}
				if err != nil {
					logging.HelperResourceError(ctx,
						"Unexpected error",
						map[string]interface{}{logging.KeyError: err},
					)
					t.Fatalf("Step %d/%d error: %s", stepNumber, len(c.Steps), err)
				}
			}

			var hasTerraformBlock bool
			var hasProviderBlock bool

			if cfg != nil {
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

			mergedConfig := step.mergedConfig(ctx, c, hasTerraformBlock, hasProviderBlock)

			confRequest := teststep.PrepareConfigurationRequest{
				Directory: step.ConfigDirectory,
				File:      step.ConfigFile,
				Raw:       mergedConfig,
				TestStepConfigRequest: config.TestStepConfigRequest{
					StepNumber: stepIndex + 1,
					TestName:   t.Name(),
				},
			}.Exec()

			appliedCfg = teststep.Configuration(confRequest)

			logging.HelperResourceDebug(ctx, "Finished TestStep")

			continue
		}

		t.Fatalf("Step %d/%d, unsupported test mode", stepNumber, len(c.Steps))
	}

	if stepNumber > 0 {
		copyWorkingDir(ctx, t, stepNumber, wd)
	}
}

func getState(ctx context.Context, t testing.T, wd *plugintest.WorkingDir) (*terraform.State, error) {
	t.Helper()

	jsonState, err := wd.State(ctx)
	if err != nil {
		return nil, err
	}
	state, err := shimStateFromJson(jsonState)
	if err != nil {
		t.Fatal(err)
	}
	return state, nil
}

func stateIsEmpty(state *terraform.State) bool {
	return state.Empty() || !state.HasResources() //nolint:staticcheck // legacy usage
}

func planIsEmpty(plan *tfjson.Plan) bool {
	for _, rc := range plan.ResourceChanges {
		for _, a := range rc.Change.Actions {
			if a != tfjson.ActionNoop {
				return false
			}
		}
	}
	return true
}

func testIDRefresh(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, step TestStep, r *terraform.ResourceState, providers *providerFactories, stepIndex int) error {
	t.Helper()

	// Build the state. The state is just the resource with an ID. There
	// are no attributes. We only set what is needed to perform a refresh.
	state := terraform.NewState() //nolint:staticcheck // legacy usage
	state.RootModule().Resources = make(map[string]*terraform.ResourceState)
	state.RootModule().Resources[c.IDRefreshName] = &terraform.ResourceState{}

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

	var hasProviderBlock bool

	if cfg != nil {
		var err error

		hasProviderBlock, err = cfg.HasProviderBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains provider block for import test config",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains provider block for import test config: %s", err)
		}
	}

	// Return value from c.ProviderConfig() is assigned to Raw as this was previously being
	// passed to wd.SetConfig() when the second argument accept a configuration string.
	testStepConfig := teststep.Configuration(
		teststep.ConfigurationRequest{
			Raw: teststep.Pointer(c.providerConfig(ctx, hasProviderBlock)),
		},
	)

	// Temporarily set the config to a minimal provider config for the refresh
	// test. After the refresh we can reset it.
	err := wd.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		t.Fatalf("Error setting import test config: %s", err)
	}

	defer func() {
		confRequest := teststep.PrepareConfigurationRequest{
			Directory: step.ConfigDirectory,
			File:      step.ConfigFile,
			Raw:       step.providerConfig(ctx, hasProviderBlock),
			TestStepConfigRequest: config.TestStepConfigRequest{
				StepNumber: stepIndex + 1,
				TestName:   t.Name(),
			},
		}.Exec()

		testStepConfigDefer := teststep.Configuration(confRequest)

		err = wd.SetConfig(ctx, testStepConfigDefer, step.ConfigVariables)

		if err != nil {
			t.Fatalf("Error resetting test config: %s", err)
		}
	}()

	// Refresh!
	err = runProviderCommand(ctx, t, func() error {
		err = wd.Refresh(ctx)
		if err != nil {
			t.Fatalf("Error running terraform refresh: %s", err)
		}
		state, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		return err
	}

	// Verify attribute equivalence.
	actualR := state.RootModule().Resources[c.IDRefreshName]
	if actualR == nil {
		return fmt.Errorf("Resource gone!")
	}
	if actualR.Primary == nil {
		return fmt.Errorf("Resource has no primary instance")
	}
	actual := actualR.Primary.Attributes
	expected := r.Primary.Attributes

	if len(c.IDRefreshIgnore) > 0 {
		logging.HelperResourceTrace(ctx, fmt.Sprintf("Using TestCase IDRefreshIgnore: %v", c.IDRefreshIgnore))
	}

	// Remove fields we're ignoring
	for _, v := range c.IDRefreshIgnore {
		for k := range actual {
			if strings.HasPrefix(k, v) {
				delete(actual, k)
			}
		}
		for k := range expected {
			if strings.HasPrefix(k, v) {
				delete(expected, k)
			}
		}
	}

	if !reflect.DeepEqual(actual, expected) {
		// Determine only the different attributes
		for k, v := range expected {
			if av, ok := actual[k]; ok && v == av {
				delete(expected, k)
				delete(actual, k)
			}
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			return fmt.Errorf("IDRefreshName attributes not equivalent. Difference is shown below. The - symbol indicates attributes missing after refresh.\n\n%s", diff)
		}
	}

	return nil
}

func copyWorkingDir(ctx context.Context, t testing.T, stepNumber int, wd *plugintest.WorkingDir) {
	if os.Getenv(plugintest.EnvTfAccPersistWorkingDir) == "" {
		return
	}

	workingDir := wd.GetHelper().WorkingDirectory()

	dest := filepath.Join(workingDir, fmt.Sprintf("%s%s", "step_", strconv.Itoa(stepNumber)))

	baseDir := wd.BaseDir()
	rootBaseDir := strings.TrimLeft(baseDir, workingDir)

	err := plugintest.CopyDir(workingDir, dest, rootBaseDir)
	if err != nil {
		logging.HelperResourceError(ctx,
			"Unexpected error copying working directory files",
			map[string]interface{}{logging.KeyError: err},
		)
		t.Fatalf("TestStep %d/%d error copying working directory files: %s", stepNumber, err)
	}

	t.Logf("Working directory and files have been copied to: %s", dest)
}
