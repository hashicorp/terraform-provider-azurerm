// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIStoragePathResource struct{}

func TestAccStackHCIStoragePath_basic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_storage_path", "test")
	r := StackHCIStoragePathResource{}

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

func TestAccStackHCIStoragePath_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_storage_path", "test")
	r := StackHCIStoragePathResource{}

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

func TestAccStackHCIStoragePath_update(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_storage_path", "test")
	r := StackHCIStoragePathResource{}

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
			Config: r.update(data),
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

func TestAccStackHCIStoragePath_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_storage_path", "test")
	r := StackHCIStoragePathResource{}

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

func (r StackHCIStoragePathResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.StorageContainers
	id, err := storagecontainers.ParseStorageContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIStoragePathResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-${var.random_string}"
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIStoragePathResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_storage_path" "import" {
  name                = azurerm_stack_hci_storage_path.test.name
  resource_group_name = azurerm_stack_hci_storage_path.test.resource_group_name
  location            = azurerm_stack_hci_storage_path.test.location
  custom_location_id  = azurerm_stack_hci_storage_path.test.custom_location_id
  path                = azurerm_stack_hci_storage_path.test.path
}
`, config)
}

func (r StackHCIStoragePathResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-${var.random_string}"
  tags = {
    foo = "bar"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIStoragePathResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-${var.random_string}"
  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIStoragePathResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_string" {
  default = %q
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-sp-${var.random_string}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomString)
}
