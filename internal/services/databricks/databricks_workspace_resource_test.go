// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatabricksWorkspaceResource struct{}

type DatabricksWorkspaceAlternateSubscription struct {
	tenant_id       string
	subscription_id string
}

func TestAccDatabricksWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_defaultStorageFirewall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultStorageFirewall(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
		{
			Config: r.defaultStorageFirewallUpdateToDisabled(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id", "access_connector_id"),
	})
}

func TestAccDatabricksWorkspace_sameName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sameName(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabricksWorkspace_machineLearning(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.machineLearning(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_privateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLink(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_secureClusterConnectivity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.secureClusterConnectivity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_loadBalancerSecureClusterConnectivity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancerSecureClusterConnectivity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_infrastructureEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.infrastructureEncryption(data, "premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_extendedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extendedUpdateCreate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
		{
			Config: r.extendedUpdateUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "trial"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_managedServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedServices(data, databricksPrincipalID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_managedServicesAndRootDbfsCMK(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedServicesAndRootDbfsCMK(data, databricksPrincipalID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_managedDiskCMK(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedDiskCMK(data, databricksPrincipalID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_managedServicesRootDbfsCMKAndPrivateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedServicesRootDbfsCMKAndPrivateLink(data, databricksPrincipalID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id"),
	})
}

func TestAccDatabricksWorkspace_altSubscriptionCmkComplete(t *testing.T) {
	altSubscription := altSubscriptionCheck()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.altSubscriptionCmkComplete(data, databricksPrincipalID, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_services_cmk_key_vault_id", "managed_disk_cmk_key_vault_id"),
	})
}

func TestAccDatabricksWorkspace_altSubscriptionCmkKeysInDifferentKeyVaultsAcrossDifferentSubscriptions(t *testing.T) {
	altSubscription := altSubscriptionCheck()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.altSubscriptionCmkKeysInDifferentKeyVaultsAcrossDifferentSubscriptions(data, databricksPrincipalID, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_services_cmk_key_vault_id", "managed_disk_cmk_key_vault_id"),
	})
}

func TestAccDatabricksWorkspace_altSubscriptionCmkServicesOnly(t *testing.T) {
	altSubscription := altSubscriptionCheck()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.altSubscriptionCmkServicesOnly(data, databricksPrincipalID, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_services_cmk_key_vault_id"),
	})
}

func TestAccDatabricksWorkspace_altSubscriptionCmkDiskOnly(t *testing.T) {
	altSubscription := altSubscriptionCheck()

	if altSubscription == nil {
		t.Skip("Skipping: Test requires `ARM_SUBSCRIPTION_ID_ALT` and `ARM_TENANT_ID` environment variables to be specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	databricksPrincipalID := getDatabricksPrincipalId(data.Client().SubscriptionID)
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.altSubscriptionCmkDiskOnly(data, databricksPrincipalID, altSubscription),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_disk_cmk_key_vault_id"),
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enhancedSecurityCompliance(data, "premium", true, true, []string{"PCI_DSS", "HIPAA"}, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithoutStandards(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enhancedSecurityCompliance(data, "premium", true, true, []string{}, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithInvalidStandards(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, "premium", true, true, []string{"NONE"}, true),
			ExpectError: regexp.MustCompile(`expected .* to be one of .*, got NONE`),
		},
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithInvalidStandardSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, "standard", true, true, []string{"PCI_DSS", "HIPAA"}, true),
			ExpectError: regexp.MustCompile("enhanced_security_compliance.*are only available with a `premium`"),
		},
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithoutEnhancedSecurityMonitoring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, "premium", true, true, []string{"PCI_DSS", "HIPAA"}, false),
			ExpectError: regexp.MustCompile("`enhanced_security_monitoring_enabled` must be set to true when `compliance_security_profile_enabled` is set to true"),
		},
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithoutAutomaticClusterUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, "premium", false, true, []string{"PCI_DSS", "HIPAA"}, true),
			ExpectError: regexp.MustCompile("`automatic_cluster_update_enabled` .* must be set to true when `compliance_security_profile_enabled` is set to true"),
		},
	})
}

func TestAccDatabricksWorkspace_enhancedComplianceSecurityWithInvalidComplianceSecurityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enhancedSecurityCompliance(data, "premium", true, false, []string{"PCI_DSS", "HIPAA"}, true),
			ExpectError: regexp.MustCompile("`compliance_security_profile_standards` cannot be set when `compliance_security_profile_enabled` is false"),
		},
	})
}

func getDatabricksPrincipalId(subscriptionId string) string {
	databricksPrincipalID := "bb9ef821-a78b-4312-90cc-5ece3fad3430"
	if strings.HasPrefix(strings.ToLower(subscriptionId), "85b3dbca") {
		databricksPrincipalID = "fe597bb2-377c-44f1-8515-82c8a1a62e3d"
	}

	return databricksPrincipalID
}

func altSubscriptionCheck() *DatabricksWorkspaceAlternateSubscription {
	altSubscriptonID := os.Getenv("ARM_SUBSCRIPTION_ID_ALT")
	altTenantID := os.Getenv("ARM_TENANT_ID")

	if altSubscriptonID == "" || altTenantID == "" {
		return nil
	}

	return &DatabricksWorkspaceAlternateSubscription{
		tenant_id:       altTenantID,
		subscription_id: altSubscriptonID,
	}
}

func (DatabricksWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.WorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Analysis Services Server %q (resource group: %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DatabricksWorkspaceResource) basic(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func (DatabricksWorkspaceResource) defaultStorageFirewall(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBWACC%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"

  custom_parameters {
    no_public_ip        = false
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  access_connector_id              = azurerm_databricks_access_connector.test.id
  default_storage_firewall_enabled = true

}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (DatabricksWorkspaceResource) defaultStorageFirewallUpdateToDisabled(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBWACC%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"

  custom_parameters {
    no_public_ip        = false
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  access_connector_id              = azurerm_databricks_access_connector.test.id
  default_storage_firewall_enabled = false

}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (DatabricksWorkspaceResource) sameName(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%s"

  managed_resource_group_name = "acctestDBW-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, data.RandomInteger)
}

func (DatabricksWorkspaceResource) requiresImport(data acceptance.TestData) string {
	template := DatabricksWorkspaceResource{}.basic(data, "standard")
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace" "import" {
  name                = azurerm_databricks_workspace.test.name
  resource_group_name = azurerm_databricks_workspace.test.resource_group_name
  location            = azurerm_databricks_workspace.test.location
  sku                 = azurerm_databricks_workspace.test.sku
}
`, template)
}

func (DatabricksWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  custom_parameters {
    no_public_ip        = false
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) extendedUpdateCreate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                          = "acctestDBW-%[1]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  sku                           = "premium"
  managed_resource_group_name   = "acctestRG-DBW-%[1]d-managed"
  public_network_access_enabled = true
  customer_managed_key_enabled  = true

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) extendedUpdateUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                                  = "acctestDBW-%[1]d"
  resource_group_name                   = azurerm_resource_group.test.name
  location                              = azurerm_resource_group.test.location
  sku                                   = "premium"
  managed_resource_group_name           = "acctestRG-DBW-%[1]d-managed"
  public_network_access_enabled         = false
  network_security_group_rules_required = "NoAzureDatabricksRules"
  customer_managed_key_enabled          = false

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_subnet_network_security_group_association.public, azurerm_subnet_network_security_group_association.private]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Pricing = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) machineLearning(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[4]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {
    error_404_document = "error.html"
    index_document     = "index.html"
  }
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-mlws-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"

  custom_parameters {
    no_public_ip                  = false
    machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  }
}
`, data.RandomInteger, data.Locations.Secondary, sku, data.RandomString)
}

func (DatabricksWorkspaceResource) privateLink(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "privatelink" {
  name                 = "acctest-snpl-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]

  private_endpoint_network_policies = "Enabled"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  public_network_access_enabled         = false
  network_security_group_rules_required = "NoAzureDatabricksRules"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_private_endpoint" "databricks" {
  name                = "acctest-endpoint-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.privatelink.id

  private_service_connection {
    name                           = "acctest-psc-%[1]d"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace.test.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_dns_zone" "test" {
  depends_on = [azurerm_private_endpoint.databricks]

  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = azurerm_databricks_workspace.test.workspace_url
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300
  record              = "eastus2-c2.azuredatabricks.net"
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (DatabricksWorkspaceResource) secureClusterConnectivity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  public_network_access_enabled = true

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) loadBalancerSecureClusterConnectivity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  public_network_access_enabled         = true
  load_balancer_backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  idle_timeout_in_minutes = 4
  allocation_method       = "Static"

  sku = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                 = "acctest-PIP-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_outbound_rule" "test" {
  name                     = "OutboundRule-%[1]d"
  loadbalancer_id          = azurerm_lb.test.id
  protocol                 = "All"
  enable_tcp_reset         = true
  allocated_outbound_ports = 1024
  idle_timeout_in_minutes  = 4

  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = azurerm_lb.test.frontend_ip_configuration.0.name
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (DatabricksWorkspaceResource) infrastructureEncryption(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"

  infrastructure_encryption_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, sku, data.RandomString)
}

func (DatabricksWorkspaceResource) managedServices(data acceptance.TestData, databricksPrincipalID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.test.id

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-certificate"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, databricksPrincipalID)
}

func (DatabricksWorkspaceResource) managedServicesAndRootDbfsCMK(data acceptance.TestData, databricksPrincipalID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  customer_managed_key_enabled          = true
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.test.id

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-certificate"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, databricksPrincipalID)
}

func (DatabricksWorkspaceResource) managedDiskCMK(data acceptance.TestData, databricksPrincipalID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  customer_managed_key_enabled      = true
  managed_disk_cmk_key_vault_key_id = azurerm_key_vault_key.test.id

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-certificate"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "List",
    "Encrypt",
    "Decrypt",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.managed_disk_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.managed_disk_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, databricksPrincipalID)
}

func (DatabricksWorkspaceResource) managedServicesRootDbfsCMKAndPrivateLink(data acceptance.TestData, databricksPrincipalID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}

resource "azurerm_subnet" "privatelink" {
  name                 = "acctest-snpl-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]

  private_endpoint_network_policies = "Enabled"
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  customer_managed_key_enabled          = true
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.test.id
  public_network_access_enabled         = false
  network_security_group_rules_required = "NoAzureDatabricksRules"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-certificate"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_private_endpoint" "databricks" {
  depends_on = [azurerm_databricks_workspace_root_dbfs_customer_managed_key.test]

  name                = "acctest-endpoint-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.privatelink.id

  private_service_connection {
    name                           = "acctest-psc-%[1]d"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace.test.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_dns_zone" "test" {
  depends_on = [azurerm_private_endpoint.databricks]

  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = azurerm_databricks_workspace.test.workspace_url
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300
  record              = "eastus2-c2.azuredatabricks.net"
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, databricksPrincipalID)
}

func (DatabricksWorkspaceResource) altSubscriptionCmkComplete(data acceptance.TestData, databricksPrincipalID string, alt *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[5]s"
  subscription_id = "%[6]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-pri-sub-complete-%[1]d"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-alt-sub-complete-%[1]d"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctest-databricks-pri-sub-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "databricks-pri-sub-managed-rg-%[1]d"

  managed_services_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.services.id

  managed_disk_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_disk_cmk_key_vault_key_id = azurerm_key_vault_key.disk.id

  tags = {
    Environment = "Sandbox"
    Pricing     = "Premium"
  }
}

# Create this in a different subscription...
resource "azurerm_key_vault" "keyVault" {
  provider = azurerm-alt

  name                = "kv-altsub-%[3]s"
  resource_group_name = azurerm_resource_group.keyVault.name
  location            = azurerm_resource_group.keyVault.location
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "services" {
  provider   = azurerm-alt
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-services-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_key" "disk" {
  provider   = azurerm-alt
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-disk-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, databricksPrincipalID, alt.tenant_id, alt.subscription_id)
}

func (DatabricksWorkspaceResource) altSubscriptionCmkKeysInDifferentKeyVaultsAcrossDifferentSubscriptions(data acceptance.TestData, databricksPrincipalID string, alt *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[5]s"
  subscription_id = "%[6]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-pri-sub-diff-%[1]d"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  name     = "acctestRG-databricks-kv-pri-sub-diff-%[1]d"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVaultAlt" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-kv-alt-sub-diff-%[1]d"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctest-databricks-pri-sub-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "databricks-pri-sub-managed-rg-%[1]d"

  managed_services_cmk_key_vault_id     = azurerm_key_vault.keyVaultAlt.id
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.servicesAlt.id

  managed_disk_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_disk_cmk_key_vault_key_id = azurerm_key_vault_key.disk.id

  tags = {
    Environment = "Sandbox"
    Pricing     = "Premium"
  }
}

resource "azurerm_key_vault" "keyVault" {
  name                = "kv-prisub-%[3]s"
  resource_group_name = azurerm_resource_group.keyVault.name
  location            = azurerm_resource_group.keyVault.location
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

# Create this in a different subscription...
resource "azurerm_key_vault" "keyVaultAlt" {
  provider = azurerm-alt

  name                = "kv-altsub-%[3]s"
  resource_group_name = azurerm_resource_group.keyVaultAlt.name
  location            = azurerm_resource_group.keyVaultAlt.location
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "servicesAlt" {
  provider   = azurerm-alt
  depends_on = [azurerm_key_vault_access_policy.terraformAlt]

  name         = "acctest-services-certificate"
  key_vault_id = azurerm_key_vault.keyVaultAlt.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_key" "disk" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-disk-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "terraformAlt" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVaultAlt.id
  tenant_id    = azurerm_key_vault.keyVaultAlt.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managedAlt" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVaultAlt.id
  tenant_id    = azurerm_key_vault.keyVaultAlt.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, databricksPrincipalID, alt.tenant_id, alt.subscription_id)
}

func (DatabricksWorkspaceResource) altSubscriptionCmkServicesOnly(data acceptance.TestData, databricksPrincipalID string, alt *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[5]s"
  subscription_id = "%[6]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-pri-sub-services-%[1]d"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-alt-sub-services-%[1]d"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctest-databricks-pri-sub-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "databricks-pri-sub-managed-rg-%[1]d"

  managed_services_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.services.id

  tags = {
    Environment = "Sandbox"
    Pricing     = "Premium"
  }
}

# Create this in a different subscription...
resource "azurerm_key_vault" "keyVault" {
  provider = azurerm-alt

  name                = "kv-altsub-%[3]s"
  location            = azurerm_resource_group.keyVault.location
  resource_group_name = azurerm_resource_group.keyVault.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "services" {
  provider   = azurerm-alt
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-services-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, databricksPrincipalID, alt.tenant_id, alt.subscription_id)
}

func (DatabricksWorkspaceResource) altSubscriptionCmkDiskOnly(data acceptance.TestData, databricksPrincipalID string, alt *DatabricksWorkspaceAlternateSubscription) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azurerm-alt" {
  features {}

  tenant_id       = "%[5]s"
  subscription_id = "%[6]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-pri-sub-disk-%[1]d"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm-alt

  name     = "acctestRG-databricks-alt-sub-disk-%[1]d"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctest-databricks-pri-sub-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "databricks-pri-sub-managed-rg-%[1]d"

  managed_disk_cmk_key_vault_id     = azurerm_key_vault.keyVault.id
  managed_disk_cmk_key_vault_key_id = azurerm_key_vault_key.disk.id

  tags = {
    Environment = "Sandbox"
    Pricing     = "Premium"
  }
}

# Create this in a different subscription...
resource "azurerm_key_vault" "keyVault" {
  provider = azurerm-alt

  name                = "kv-altsub-%[3]s"
  location            = azurerm_resource_group.keyVault.location
  resource_group_name = azurerm_resource_group.keyVault.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "disk" {
  provider   = azurerm-alt
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-disk-certificate"
  key_vault_id = azurerm_key_vault.keyVault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.keyVault.id
  tenant_id    = azurerm_key_vault.keyVault.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, databricksPrincipalID, alt.tenant_id, alt.subscription_id)
}

func (DatabricksWorkspaceResource) enhancedSecurityCompliance(data acceptance.TestData, sku string, automaticClusterUpdateEnabled bool, complianceSecurityProfileEnabled bool, complianceSecurityProfileStandards []string, enhancedSecurityMonitoringEnabled bool) string {
	complianceSecurityProfileStandardsStr := ""
	if len(complianceSecurityProfileStandards) > 0 {
		complianceSecurityProfileStandardsStr = fmt.Sprintf(`"%s"`, strings.Join(complianceSecurityProfileStandards, `", "`))
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%s"

  enhanced_security_compliance {
    automatic_cluster_update_enabled      = %t
    compliance_security_profile_enabled   = %t
    compliance_security_profile_standards = [%s]
    enhanced_security_monitoring_enabled  = %t
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, automaticClusterUpdateEnabled, complianceSecurityProfileEnabled, complianceSecurityProfileStandardsStr, enhancedSecurityMonitoringEnabled)
}
