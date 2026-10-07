// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/vcr"
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
	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

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

		newSteps = append(newSteps, step)
	}
	steps = newSteps

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
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
	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
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
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	td.runAcceptanceSequentialTest(t, testCase)
}

// ResourceRegressionTest runs an acceptance test for resource regression.
// It expects one or two steps. If only one step is supplied it will be duplicated for the 2nd verification step resulting in an identical config being set for the locally built resource.
// The first step uses a specified previous provider version constraint. If an empty string is supplied, the test will use the version in ./version/VERSION from the root of the project
// For StateMigration testing, this should be the last version that has the previous `SchemaVersion` value.
// The second step uses the locally built provider code.
func (td TestData) ResourceRegressionTest(t *testing.T, testResource types.TestResource, steps []TestStep, previousVersion string) {
	if len(steps) != 2 {
		if len(steps) == 1 {
			// duplicate step[0] for second stage - this is all that _should_ be required for breaking change testing
			steps = append(steps, steps[0])
			// the duplicated verification step must not re-run step[0]'s PreConfig - that is setup
			// (e.g. unregistering a leftover provider/feature registration) which would tear down
			// the resource the first step just created out from under the verification step
			steps[1].PreConfig = nil
		} else {
			t.Fatal("expected exactly 2 steps for Regression test. Setup and Check")
		}
	}

	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	steps[0].ExternalProviders = td.externalProviders()
	steps[0].ExternalProviders["azurerm"] = resource.ExternalProvider{
		VersionConstraint: providerRelease([]string{previousVersion}...),
		Source:            "hashicorp/azurerm",
	}

	steps[1].ExternalProviders = td.externalProviders()
	steps[1].ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt")
	steps[1].ConfigPlanChecks = resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			helpers.IsNotResourceAction(td.ResourceName, plancheck.ResourceActionReplace),
		},
	}

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	resource.ParallelTest(t, testCase)
}

// ResourceSequentialRegressionTest is ResourceRegressionTest for resources whose tests cannot run in
// parallel - e.g. subscription-level singletons (provider/feature registrations) where concurrent
// tests using the same Azure resource collide.
func (td TestData) ResourceSequentialRegressionTest(t *testing.T, testResource types.TestResource, steps []TestStep, previousVersion string) {
	if len(steps) != 2 {
		if len(steps) == 1 {
			// duplicate step[0] for second stage - this is all that _should_ be required for breaking change testing
			steps = append(steps, steps[0])
			// the duplicated verification step must not re-run step[0]'s PreConfig - that is setup
			// (e.g. unregistering a leftover provider/feature registration) which would tear down
			// the resource the first step just created out from under the verification step
			steps[1].PreConfig = nil
		} else {
			t.Fatal("expected exactly 2 steps for Regression test. Setup and Check")
		}
	}

	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	steps[0].ExternalProviders = td.externalProviders()
	steps[0].ExternalProviders["azurerm"] = resource.ExternalProvider{
		VersionConstraint: providerRelease([]string{previousVersion}...),
		Source:            "hashicorp/azurerm",
	}

	steps[1].ExternalProviders = td.externalProviders()
	steps[1].ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt")
	steps[1].ConfigPlanChecks = resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			helpers.IsNotResourceAction(td.ResourceName, plancheck.ResourceActionReplace),
		},
	}

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	resource.Test(t, testCase)
}

// DataSourceRegressionTest runs an acceptance test for data source regression - the data source equivalent of ResourceRegressionTest.
// It expects one or two steps. If only one step is supplied it will be duplicated for the 2nd verification step resulting in an identical config being set for the locally built provider.
// The first step uses a specified previous provider version constraint. If an empty string is supplied, the test will use the version in ./version/VERSION from the root of the project
// The second step uses the locally built provider code. As data sources create nothing there is no destroy check.
func (td TestData) DataSourceRegressionTest(t *testing.T, steps []TestStep, previousVersion string) {
	if len(steps) != 2 {
		if len(steps) == 1 {
			// duplicate step[0] for second stage - this is all that _should_ be required for breaking change testing
			steps = append(steps, steps[0])
			// the duplicated verification step must not re-run step[0]'s PreConfig - that is setup
			// (e.g. unregistering a leftover provider/feature registration) which would tear down
			// the resource the first step just created out from under the verification step
			steps[1].PreConfig = nil
		} else {
			t.Fatal("expected exactly 2 steps for Regression test. Setup and Check")
		}
	}

	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	steps[0].ExternalProviders = td.externalProviders()
	steps[0].ExternalProviders["azurerm"] = resource.ExternalProvider{
		VersionConstraint: providerRelease([]string{previousVersion}...),
		Source:            "hashicorp/azurerm",
	}

	steps[1].ExternalProviders = td.externalProviders()
	steps[1].ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt")
	steps[1].ConfigPlanChecks = resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			helpers.IsNotResourceAction(td.ResourceName, plancheck.ResourceActionReplace),
		},
	}

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}

	resource.ParallelTest(t, testCase)
}

// ResourceRegressionAdditionalStepsTest runs an acceptance test for resource regression scenarios that require multiple steps to set up.
// It expects multiple steps (3 or more), for regression testing with 1 or 2 steps, use ResourceRegressionTest.
// All steps except the final step use a specified previous provider version constraint. If an empty string is supplied, the test will use the version in ./version/VERSION from the root of the project
// For StateMigration testing, this should be the last version that has the previous `SchemaVersion` value.
// The final step uses the locally built provider code.
func (td TestData) ResourceRegressionAdditionalStepsTest(t *testing.T, testResource types.TestResource, steps []TestStep, previousVersion string) {
	l := len(steps)
	if l < 3 {
		t.Fatalf("expected at least 3 steps, got %d. For tests with less than 3 steps, use `ResourceRegressionTest`", l)
	}

	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	for i := range steps[:l-1] {
		steps[i].ExternalProviders = td.externalProviders()
		steps[i].ExternalProviders["azurerm"] = resource.ExternalProvider{
			VersionConstraint: providerRelease([]string{previousVersion}...),
			Source:            "hashicorp/azurerm",
		}
	}

	steps[l-1].ExternalProviders = td.externalProviders()
	steps[l-1].ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt")
	steps[l-1].ConfigPlanChecks = resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			helpers.IsNotResourceAction(td.ResourceName, plancheck.ResourceActionReplace),
		},
	}

	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	resource.ParallelTest(t, testCase)
}

func RunTestsInSequence(t *testing.T, tests map[string]map[string]func(t *testing.T)) {
	for group, m := range tests {
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	testclient.RegisterTestT(t)
	defer testclient.UnregisterTestT()

	if os.Getenv("TC_TEST_VIA_VCR") != "" {
		defer func(testName string) {
			_ = vcr.StopRecorder(testName)
		}(t.Name())
	}

	testCase.ExternalProviders = td.externalProviders()
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt")

	resource.ParallelTest(t, testCase)
}

func (td TestData) runAcceptanceSequentialTest(t *testing.T, testCase resource.TestCase) {
	if os.Getenv("TC_TEST_VIA_VCR") != "" {
		defer func(testName string) {
			_ = vcr.StopRecorder(testName)
		}(t.Name())
	}

	testCase.ExternalProviders = td.externalProviders()
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm")

	resource.Test(t, testCase)
}

func (td TestData) externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azuread": {
			VersionConstraint: "=3.4.0",
			Source:            "registry.terraform.io/hashicorp/azuread",
		},
		"local": {
			VersionConstraint: "=2.5.2",
			Source:            "registry.terraform.io/hashicorp/local",
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
