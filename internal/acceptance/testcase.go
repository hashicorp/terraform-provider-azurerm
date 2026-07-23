// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

		if !step.ImportState {
			newSteps = append(newSteps, step)
		}
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

// RunTestsInSequence runs all tests in sequence.
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

// getPreviousProviderVersion reads the most recent release version from the repo's version/VERSION file
func getPreviousProviderVersion() string {
	_, b, _, _ := runtime.Caller(0)
	// b is .../terraform-provider-azurerm/internal/acceptance/testcase.go
	repoRoot := filepath.Join(filepath.Dir(b), "..", "..")
	versionFile := filepath.Join(repoRoot, "version", "VERSION")

	data, err := os.ReadFile(versionFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read version file at %s: %s", versionFile, err))
	}

	return strings.TrimSpace(string(data))
}

// ResourceUpgradeTest executes a 2-stage upgrade test. Stage 1 provisions the resource natively using the
// preceding released version downloaded cleanly from the registry. Stage 2 hands over the state cleanly to
// the locally built test daemons.
func (td TestData) ResourceUpgradeTest(t *testing.T, testResource types.TestResource, steps []TestStep) {
	td.ResourceUpgradeWithVersionTest(t, getPreviousProviderVersion(), testResource, steps)
}

// ResourceUpgradeWithVersionTest allows an explicit framework wrapper bridging the state migration checks.
// It executes a 2-stage upgrade test. Stage 1 provisions the resource natively using the
// specified previousVersion. Stage 2 hands over the state to the locally built test daemons.
// If skipStepOneDuringHandover is true, the first step is only executed against the registry version,
// and drops out of the array before local version checks to prevent checking locally invalid legacy config schemas.
func (td TestData) ResourceUpgradeWithVersionTest(t *testing.T, previousVersion string, testResource types.TestResource, steps []TestStep) {
	os.Setenv("TF_ACC_REFRESH_AFTER_APPLY", "true")

	tc := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps:                    steps,
		ExternalProviders:        td.externalProviders(),
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm", "azurerm-alt"),
	}

	upgradeCase := resource.ProviderUpgradeTestCase{
		ProviderName:    "azurerm",
		ProviderSource:  "registry.terraform.io/hashicorp/azurerm",
		ProviderVersion: previousVersion,
		TestCase:        tc, // Hand execution over to local factories
	}

	testclient.RegisterTestT(t)
	defer testclient.UnregisterTestT()

	if os.Getenv("TC_TEST_VIA_VCR") != "" {
		defer func(testName string) {
			_ = vcr.StopRecorder(testName)
		}(t.Name())
	}

	resource.TestProviderUpgrade(t, upgradeCase)
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

func (td TestData) getFeaturesBlock(config string) string {
	return fmt.Sprintf("azurerm_test_data_%s", config)
}
