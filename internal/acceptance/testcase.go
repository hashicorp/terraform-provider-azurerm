// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
>>>>>>> 5dbe4f4bb6 (go-vcr dummy chnages with Meta)
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/vcr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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

func (td TestData) ResourceIdentityTest(t *testing.T, steps []TestStep, sequential bool) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0"))),
		},
		Steps: steps,
	}

	if sequential {
		td.runAcceptanceSequentialTest(t, testCase)
		return
	}

	td.runAcceptanceTest(t, testCase)
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

// ResourceTestWithVCR is an opt-in test method that uses VCR for HTTP recording/playback.
// Tests using this method will use a VCR-wrapped HTTP client for all Azure API calls.
// This enables faster test execution by replaying recorded HTTP interactions.
func (td TestData) ResourceTestWithVCR(t *testing.T, testResource types.TestResource, steps []TestStep) {
	refreshStep := TestStep{
		RefreshState: true,
	}

	newSteps := make([]TestStep, 0)
	for _, step := range steps {
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

	td.runAcceptanceTestWithVCR(t, testCase)
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
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm", "azurerm-alt")

	resource.ParallelTest(t, testCase)
}

// runAcceptanceTestWithVCR runs acceptance test with a VCR-wrapped HTTP client for recording/playback.
func (td TestData) runAcceptanceTestWithVCR(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()

	// Create provider factories with VCR HTTP client injection
	testCase.ProtoV5ProviderFactories = make(map[string]func() (tfprotov5.ProviderServer, error))
	for _, name := range []string{"azurerm", "azurerm-alt"} {
		testCase.ProtoV5ProviderFactories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, v2Provider, err := framework.ProtoV5ProviderServerFactory(context.Background())
			if err != nil {
				return nil, err
			}

			// Wrap the original ConfigureContextFunc to inject HTTPClient into Meta
			configureContextFunc := v2Provider.ConfigureContextFunc
			v2Provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				// Get VCR HTTP client
				httpClient := vcr.GetHTTPClient(t)
				var meta *clients.Client
				if v, ok := v2Provider.Meta().(*clients.Client); ok && v != nil {
					meta = v
				} else {
					meta = new(clients.Client)
				}
				meta.HTTPClient = httpClient
				v2Provider.SetMeta(meta)

				// Call the underline ConfigureContextFunc
				return configureContextFunc(ctx, d)
			}

			return providerServerFactory(), nil
		}
	}

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
			VersionConstraint: "=3.4.0",
			Source:            "registry.terraform.io/hashicorp/azuread",
		},
		"random": {
			VersionConstraint: "=3.7.2",
			Source:            "registry.terraform.io/hashicorp/random",
		},
		"time": {
			VersionConstraint: "=0.13.1",
			Source:            "registry.terraform.io/hashicorp/time",
		},
		"tls": {
			VersionConstraint: "=4.1.0",
			Source:            "registry.terraform.io/hashicorp/tls",
		},
	}
}
