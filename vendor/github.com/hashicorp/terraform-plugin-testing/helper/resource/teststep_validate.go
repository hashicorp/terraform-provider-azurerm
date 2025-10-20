// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

// testStepValidateRequest contains data for the (TestStep).validate() method.
type testStepValidateRequest struct {
	// StepConfiguration contains the TestStep configuration derived from
	// TestStep.Config, TestStep.ConfigDirectory, or TestStep.ConfigFile.
	StepConfiguration teststep.Config

	// StepNumber is the index of the TestStep in the TestCase.Steps.
	StepNumber int

	// TestCaseHasExternalProviders is enabled if the TestCase has
	// ExternalProviders.
	TestCaseHasExternalProviders bool

	// TestCaseHasProviders is enabled if the TestCase has set any of
	// ExternalProviders, ProtoV5ProviderFactories, ProtoV6ProviderFactories,
	// or ProviderFactories.
	TestCaseHasProviders bool

	// TestName is the name of the test.
	TestName string
}

// hasExternalProviders returns true if the TestStep has
// ExternalProviders set.
func (s TestStep) hasExternalProviders() bool {
	return len(s.ExternalProviders) > 0
}

// hasProviders returns true if the TestStep has set any of the
// ExternalProviders, ProtoV5ProviderFactories, ProtoV6ProviderFactories, or
// ProviderFactories fields. It will also return true if ConfigDirectory or
// Config contain terraform configuration which specify a provider block.
func (s TestStep) hasProviders(ctx context.Context, stepIndex int, testName string) (bool, error) {
	if len(s.ExternalProviders) > 0 {
		return true, nil
	}

	if len(s.ProtoV5ProviderFactories) > 0 {
		return true, nil
	}

	if len(s.ProtoV6ProviderFactories) > 0 {
		return true, nil
	}

	if len(s.ProviderFactories) > 0 {
		return true, nil
	}

	configRequest := teststep.PrepareConfigurationRequest{
		Directory: s.ConfigDirectory,
		File:      s.ConfigFile,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepIndex + 1,
			TestName:   testName,
		},
	}.Exec()

	cfg := teststep.Configuration(configRequest)

	var cfgHasProviders bool

	if cfg != nil {
		var err error

		cfgHasProviders, err = cfg.HasProviderBlock(ctx)

		if err != nil {
			return false, err
		}
	}

	if cfgHasProviders {
		return true, nil
	}

	return false, nil
}

// validate ensures the TestStep is valid based on the following criteria:
//
//   - Config or ImportState or RefreshState is set.
//   - Config and RefreshState are not both set.
//   - RefreshState and Destroy are not both set.
//   - RefreshState is not the first TestStep.
//   - Providers are not specified (ExternalProviders,
//     ProtoV5ProviderFactories, ProtoV6ProviderFactories, ProviderFactories)
//     if specified at the TestCase level.
//   - Providers are specified (ExternalProviders, ProtoV5ProviderFactories,
//     ProtoV6ProviderFactories, ProviderFactories) if not specified at the
//     TestCase level.
//   - No overlapping ExternalProviders and ProviderFactories entries
//   - ResourceName is not empty when ImportState is true, ImportStateIdFunc
//     is not set, and ImportStateId is not set.
//   - ConfigPlanChecks (PreApply, PostApplyPreRefresh, PostApplyPostRefresh) are only set when Config is set.
//   - ConfigPlanChecks.PreApply are only set when PlanOnly is false.
//   - RefreshPlanChecks (PostRefresh) are only set when RefreshState is set.
func (s TestStep) validate(ctx context.Context, req testStepValidateRequest) error {
	ctx = logging.TestStepNumberContext(ctx, req.StepNumber)

	logging.HelperResourceTrace(ctx, "Validating TestStep")

	if req.StepConfiguration == nil && !s.ImportState && !s.RefreshState {
		err := fmt.Errorf("TestStep missing Config or ConfigDirectory or ConfigFile or ImportState or RefreshState")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if req.StepConfiguration != nil && s.RefreshState {
		err := fmt.Errorf("TestStep cannot have Config or ConfigDirectory or ConfigFile and RefreshState")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if s.RefreshState && s.Destroy {
		err := fmt.Errorf("TestStep cannot have RefreshState and Destroy")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if s.RefreshState && req.StepNumber == 1 {
		err := fmt.Errorf("TestStep cannot have RefreshState as first step")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if s.ImportState && s.RefreshState {
		err := fmt.Errorf("TestStep cannot have ImportState and RefreshState in same step")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	for name := range s.ExternalProviders {
		if _, ok := s.ProviderFactories[name]; ok {
			err := fmt.Errorf("TestStep provider %q set in both ExternalProviders and ProviderFactories", name)
			logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
			return err
		}
	}

	if req.TestCaseHasExternalProviders && req.StepConfiguration != nil && req.StepConfiguration.HasConfigurationFiles() {
		err := fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if s.hasExternalProviders() && req.StepConfiguration != nil && req.StepConfiguration.HasConfigurationFiles() {
		err := fmt.Errorf("Providers must only be specified within the terraform configuration files when using TestStep.Config")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	// We need a 0-based step index for consistency
	hasProviders, err := s.hasProviders(ctx, req.StepNumber-1, req.TestName)

	if err != nil {
		logging.HelperResourceError(ctx, "TestStep error checking for providers", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if req.TestCaseHasProviders && hasProviders {
		err := fmt.Errorf("Providers must only be specified either at the TestCase or TestStep level")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	var cfgHasProviderBlock bool

	if req.StepConfiguration != nil {
		cfgHasProviderBlock, err = req.StepConfiguration.HasProviderBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx, "TestStep error checking for if configuration has provider block", map[string]interface{}{logging.KeyError: err})
			return err
		}
	}

	if !req.TestCaseHasProviders && !hasProviders && !cfgHasProviderBlock {
		err := fmt.Errorf("Providers must be specified at the TestCase level, or in all TestStep, or in TestStep.ConfigDirectory or TestStep.ConfigFile")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if s.ImportState {
		if s.ImportStateId == "" && s.ImportStateIdFunc == nil && s.ResourceName == "" {
			err := fmt.Errorf("TestStep ImportState must be specified with ImportStateId, ImportStateIdFunc, or ResourceName")
			logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
			return err
		}
	}

	if len(s.ConfigPlanChecks.PreApply) > 0 {
		if req.StepConfiguration == nil {
			err := fmt.Errorf("TestStep ConfigPlanChecks.PreApply must only be specified with Config, ConfigDirectory or ConfigFile")
			logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
			return err
		}

		if s.PlanOnly {
			err := fmt.Errorf("TestStep ConfigPlanChecks.PreApply cannot be run with PlanOnly")
			logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
			return err
		}
	}

	if len(s.ConfigPlanChecks.PostApplyPreRefresh) > 0 && req.StepConfiguration == nil {
		err := fmt.Errorf("TestStep ConfigPlanChecks.PostApplyPreRefresh must only be specified with Config, ConfigDirectory or ConfigFile")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if len(s.ConfigPlanChecks.PostApplyPostRefresh) > 0 && req.StepConfiguration == nil {
		err := fmt.Errorf("TestStep ConfigPlanChecks.PostApplyPostRefresh must only be specified with Config, ConfigDirectory or ConfigFile")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if len(s.RefreshPlanChecks.PostRefresh) > 0 && !s.RefreshState {
		err := fmt.Errorf("TestStep RefreshPlanChecks.PostRefresh must only be specified with RefreshState")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	if len(s.ConfigStateChecks) > 0 && req.StepConfiguration == nil {
		err := fmt.Errorf("TestStep ConfigStateChecks must only be specified with Config, ConfigDirectory or ConfigFile")
		logging.HelperResourceError(ctx, "TestStep validation error", map[string]interface{}{logging.KeyError: err})
		return err
	}

	return nil
}
