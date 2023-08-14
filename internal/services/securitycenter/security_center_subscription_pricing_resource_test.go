// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	pricings_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterSubscriptionPricingResource struct{}

func TestAccServerVulnerabilityAssessment(t *testing.T) {
	// these tests need to change `azurerm_security_center_subscription_pricing` of `VirtualMachines` in their test configs, so we need to run them serially.
	// `securityCenterAssessmentPolicy` is included because it's using same `azurerm_security_center_assessment_policy` with other tests
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"securityCenterAssessment": {
			"basic":          testAccSecurityCenterAssessment_basic,
			"complete":       testAccSecurityCenterAssessment_complete,
			"update":         testAccSecurityCenterAssessment_update,
			"requiresImport": testAccSecurityCenterAssessment_requiresImport,
		},
		"securityCenterAssessmentPolicy": {
			"basic":    testAccSecurityCenterAssessmentPolicy_basic,
			"complete": testAccSecurityCenterAssessmentPolicy_complete,
			"update":   testAccSecurityCenterAssessmentPolicy_update,
		},
		"serverVulnerabilityAssessment": {
			"basic":          testAccServerVulnerabilityAssessment_basic,
			"requiresImport": testAccServerVulnerabilityAssessment_requiresImport,
		},
		"serverVulnerabilityAssessmentVirtualMachine": {
			"basic":          testAccServerVulnerabilityAssessmentVirtualMachine_basic,
			"requiresImport": testAccServerVulnerabilityAssessmentVirtualMachine_requiresImport,
		},
		"workSpace": {
			"basic":          testAccSecurityCenterWorkspace_basic,
			"update":         testAccSecurityCenterWorkspace_update,
			"requiresImport": testAccSecurityCenterWorkspace_requiresImport,
		},
	})
}

func TestAccSecurityCenterSubscriptionPricing_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")
	r := SecurityCenterSubscriptionPricingResource{}

	data.ResourceSequentialTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.tier("Standard", "AppServices"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterSubscriptionPricing_cosmosDbs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")
	r := SecurityCenterSubscriptionPricingResource{}

	data.ResourceSequentialTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.tier("Standard", "CosmosDbs"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterSubscriptionPricing_storageAccountSubplan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")
	r := SecurityCenterSubscriptionPricingResource{}

	data.ResourceSequentialTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.storageAccountSubplan(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("subplan").HasValue("PerStorageAccount"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterSubscriptionPricing_cloudPostureExtension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_subscription_pricing", "test")
	r := SecurityCenterSubscriptionPricingResource{}

	data.ResourceSequentialTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.cloudPostureExtension(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("resource_type").HasValue("CloudPosture"),
				check.That(data.ResourceName).Key("extension.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cloudPostureExtensionUpdated(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("resource_type").HasValue("CloudPosture"),
				check.That(data.ResourceName).Key("extension.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cloudPostureExtension(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("resource_type").HasValue("CloudPosture"),
				check.That(data.ResourceName).Key("extension.#").HasValue("2"),
			),
		},
	})
}

func (SecurityCenterSubscriptionPricingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pricings_v2023_01_01.ParsePricingIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.PricingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (SecurityCenterSubscriptionPricingResource) tier(tier string, resource_type string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "%s"
  resource_type = "%s"
}
`, tier, resource_type)
}

func (SecurityCenterSubscriptionPricingResource) storageAccountSubplan() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Standard"
  resource_type = "StorageAccounts"
  subplan       = "PerStorageAccount"
}
`
}

func (SecurityCenterSubscriptionPricingResource) cloudPostureExtension() string {
	return `
provider "azurerm" {
  features {

  }
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Standard"
  resource_type = "CloudPosture"

  extension {
    name = "SensitiveDataDiscovery"
  }

  extension {
    name = "AgentlessVmScanning"
    additional_extension_properties = {
      ExclusionTags = "[]"
    }
  }
}
`
}

func (SecurityCenterSubscriptionPricingResource) cloudPostureExtensionUpdated() string {
	return `
provider "azurerm" {
  features {

  }
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Standard"
  resource_type = "CloudPosture"

  extension {
    name = "ContainerRegistriesVulnerabilityAssessments"
  }

  extension {
    name = "AgentlessVmScanning"
    additional_extension_properties = {
      ExclusionTags = "[]"
    }
  }

  extension {
    name = "AgentlessDiscoveryForKubernetes"
  }
}
`
}
