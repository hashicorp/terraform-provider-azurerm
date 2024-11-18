// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func (td TestData) DataSourceTest(t *testing.T, steps []TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out

	// lintignore:AT001
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) DataSourceTestInSequence(t *testing.T, steps []TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out

	// lintignore:AT001
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}

	td.runAcceptanceSequentialTest(t, testCase)
}

func (td TestData) ResourceTest(t *testing.T, testResource types.TestResource, steps []TestStep) {
	// Testing framework as of 1.6.0 no longer auto-refreshes state, so adding it back in here for all steps that update
	// the config rather than having to modify 1000's of tests individually to add a refresh-only step
	refreshStep := TestStep{
		RefreshState: true,
	}

	newSteps := make([]TestStep, 0)
	for _, step := range steps {
		// This block adds a check to make sure tests aren't recreating a resource
		if (step.Config != "" || step.ConfigDirectory != nil || step.ConfigFile != nil) && !step.PlanOnly {
			step.ConfigPlanChecks = resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					helpers.IsNotResourceAction(td.ResourceName, plancheck.ResourceActionReplace),
				},
			}
		}

		if !step.ImportState {
			newSteps = append(newSteps, step)
		} else {
			newSteps = append(newSteps, refreshStep)
			newSteps = append(newSteps, step)
		}
	}
	steps = newSteps

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}
	td.runAcceptanceTest(t, testCase)
}

// ResourceTestIgnoreRecreate should be used when checking that a resource should be recreated during a test.
func (td TestData) ResourceTestIgnoreRecreate(t *testing.T, testResource types.TestResource, steps []TestStep) {
	// Testing framework as of 1.6.0 no longer auto-refreshes state, so adding it back in here for all steps that update
	// the config rather than having to modify 1000's of tests individually to add a refresh-only step
	refreshStep := TestStep{
		RefreshState: true,
	}

	newSteps := make([]TestStep, 0)
	for _, step := range steps {
		if !step.ImportState {
			newSteps = append(newSteps, step)
		} else {
			newSteps = append(newSteps, refreshStep)
			newSteps = append(newSteps, step)
		}
	}
	steps = newSteps

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}
	td.runAcceptanceTest(t, testCase)
}

// ResourceTestIgnoreCheckDestroyed skips the check to confirm the resource test has been destroyed.
// This is done because certain resources can't actually be deleted.
func (td TestData) ResourceTestSkipCheckDestroyed(t *testing.T, steps []TestStep) {
	// lintignore:AT001
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) ResourceSequentialTestSkipCheckDestroyed(t *testing.T, steps []TestStep) {
	// lintignore:AT001
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	td.runAcceptanceSequentialTest(t, testCase)
}

func (td TestData) ResourceSequentialTest(t *testing.T, testResource types.TestResource, steps []TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	td.runAcceptanceSequentialTest(t, testCase)
}

func RunTestsInSequence(t *testing.T, tests map[string]map[string]func(t *testing.T)) {
	for group, m := range tests {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm")

	resource.ParallelTest(t, testCase)
}

func (td TestData) runAcceptanceSequentialTest(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm")

	resource.Test(t, testCase)
}

func (td TestData) externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azuread": {
			VersionConstraint: "=2.47.0",
			Source:            "registry.terraform.io/hashicorp/azuread",
		},
		"time": {
			VersionConstraint: "=0.9.1",
			Source:            "registry.terraform.io/hashicorp/time",
		},
		"tls": {
			VersionConstraint: "=4.0.4",
			Source:            "registry.terraform.io/hashicorp/tls",
		},
	}
}
