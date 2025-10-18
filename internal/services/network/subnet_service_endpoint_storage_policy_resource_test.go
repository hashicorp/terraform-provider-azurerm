// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubnetServiceEndpointStoragePolicyResource struct{}

func TestAccSubnetServiceEndpointStoragePolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_alias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.alias(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_storage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_update_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_update_alias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.alias(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_update_storage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.storage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetServiceEndpointStoragePolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")
	r := SubnetServiceEndpointStoragePolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t SubnetServiceEndpointStoragePolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.ServiceEndpointPolicies.Get(ctx, *id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r SubnetServiceEndpointStoragePolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetServiceEndpointStoragePolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  definition {
    name        = "resourceid"
    description = "test definition1"
    service     = "Microsoft.Storage"
    service_resources = [
      "/subscriptions/%s",
      azurerm_resource_group.test.id,
      azurerm_storage_account.test.id,
    ]
  }

  definition {
    name        = "alias"
    description = "test definition1"
    service     = "Global"
    service_resources = [
      "/services/Azure",
      "/services/Azure/Batch",
      "/services/Azure/DataFactory",
      "/services/Azure/MachineLearning",
      "/services/Azure/ManagedInstance",
      "/services/Azure/WebPI",
    ]
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.Client().SubscriptionID)
}

func (r SubnetServiceEndpointStoragePolicyResource) alias(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  definition {
    name        = "alias"
    description = "test definition1"
    service     = "Global"
    service_resources = [
      "/services/Azure",
      "/services/Azure/Batch",
      "/services/Azure/DataFactory",
      "/services/Azure/MachineLearning",
      "/services/Azure/ManagedInstance",
      "/services/Azure/WebPI",
    ]
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r SubnetServiceEndpointStoragePolicyResource) storage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  definition {
    name        = "resourceid"
    description = "test definition1"
    service     = "Microsoft.Storage"
    service_resources = [
      "/subscriptions/%s",
      azurerm_resource_group.test.id,
      azurerm_storage_account.test.id,
    ]
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.Client().SubscriptionID)
}

func (r SubnetServiceEndpointStoragePolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "import" {
  name                = azurerm_subnet_service_endpoint_storage_policy.test.name
  resource_group_name = azurerm_subnet_service_endpoint_storage_policy.test.resource_group_name
  location            = azurerm_subnet_service_endpoint_storage_policy.test.location
}
`, r.basic(data))
}

func (SubnetServiceEndpointStoragePolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
