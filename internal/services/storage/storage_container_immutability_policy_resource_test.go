// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageContainerImmutabilityPolicyResource struct{}

func TestAccStorageContainerImmutabilityPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container_immutability_policy", "test")
	r := StorageContainerImmutabilityPolicyResource{}

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

func TestAccStorageContainerImmutabilityPolicy_completeUnlocked(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container_immutability_policy", "test")
	r := StorageContainerImmutabilityPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUnlocked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainerImmutabilityPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container_immutability_policy", "test")
	r := StorageContainerImmutabilityPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUnlocked(data),
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

func TestAccStorageContainerImmutabilityPolicy_completeLocked(t *testing.T) {
	// This test has been written for manual testing of the `locked` property. Ordinarily we do not want to test this in automation,
	// since locking an immutability policy renders the container and its storage account **immutable**. This test will always
	// fail during cleanup for this reason. Uncomment the t.Skip() call to continue...
	t.Skip("this test for manual execution only")

	data := acceptance.BuildTestData(t, "azurerm_storage_container_immutability_policy", "test")
	r := StorageContainerImmutabilityPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUnlocked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeLocked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basic(data),
			ExpectError: regexp.MustCompile("unable to set `locked = false` - once an immutability policy locked it cannot be unlocked"),
		},
	})
}

func (r StorageContainerImmutabilityPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageContainerImmutabilityPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	containerId := commonids.NewStorageContainerID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ContainerName)

	resp, err := client.Storage.ResourceManager.BlobContainers.GetImmutabilityPolicy(ctx, containerId, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Properties.State != nil && !strings.EqualFold(string(*resp.Model.Properties.State), "Deleted")), nil
}

func (r StorageContainerImmutabilityPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_container_immutability_policy" "test" {
  storage_container_resource_manager_id = azurerm_storage_container.test.resource_manager_id
  immutability_period_in_days           = 1
}
`, template)
}

func (r StorageContainerImmutabilityPolicyResource) completeUnlocked(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_container_immutability_policy" "test" {
  storage_container_resource_manager_id = azurerm_storage_container.test.resource_manager_id
  immutability_period_in_days           = 2
  protected_append_writes_all_enabled   = false
  protected_append_writes_enabled       = true
}
`, template)
}

func (r StorageContainerImmutabilityPolicyResource) completeLocked(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_container_immutability_policy" "test" {
  storage_container_resource_manager_id = azurerm_storage_container.test.resource_manager_id
  immutability_period_in_days           = 2
  protected_append_writes_all_enabled   = true
  protected_append_writes_enabled       = false

  locked = true
}
`, template)
}

func (r StorageContainerImmutabilityPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "retention"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
