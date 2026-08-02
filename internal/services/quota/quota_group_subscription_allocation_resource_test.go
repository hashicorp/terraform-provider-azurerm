// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type QuotaGroupSubscriptionAllocationResource struct{}

func TestAccQuotaGroupSubscriptionAllocation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group_subscription_allocation", "test")
	r := QuotaGroupSubscriptionAllocationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccQuotaGroupSubscriptionAllocation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group_subscription_allocation", "test")
	r := QuotaGroupSubscriptionAllocationResource{}

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

func TestAccQuotaGroupSubscriptionAllocation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group_subscription_allocation", "test")
	r := QuotaGroupSubscriptionAllocationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r QuotaGroupSubscriptionAllocationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := subscriptionquotaallocation.ParseQuotaAllocationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Quota.SubscriptionQuotaAllocationClient.GroupQuotaSubscriptionAllocationList(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	// The resource exists if the list returns any non-zero allocations.
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Value != nil {
		for _, item := range *resp.Model.Properties.Value {
			if item.Properties != nil && pointer.From(item.Properties.Limit) > 0 {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}

func (r QuotaGroupSubscriptionAllocationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "test" {
  name                = "acctestqg%d"
  management_group_id = data.azurerm_management_group.test.id

  associated_subscription_ids = [data.azurerm_client_config.current.subscription_id]

  quota_request {
    resource_name = "standardddv4family"
    location      = "%s"
    limit         = 20
  }
}

resource "azurerm_quota_group_subscription_allocation" "test" {
  quota_group_id  = azurerm_quota_group.test.id
  subscription_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  location        = "%s"

  allocation {
    resource_name = "standardddv4family"
    limit         = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Primary)
}

func (r QuotaGroupSubscriptionAllocationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_quota_group_subscription_allocation" "import" {
  quota_group_id  = azurerm_quota_group_subscription_allocation.test.quota_group_id
  subscription_id = azurerm_quota_group_subscription_allocation.test.subscription_id
  location        = azurerm_quota_group_subscription_allocation.test.location

  allocation {
    resource_name = "standardddv4family"
    limit         = 10
  }
}
`, r.basic(data))
}

func (r QuotaGroupSubscriptionAllocationResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "test" {
  name                = "acctestqg%d"
  management_group_id = data.azurerm_management_group.test.id

  associated_subscription_ids = [data.azurerm_client_config.current.subscription_id]

  quota_request {
    resource_name = "standardddv4family"
    location      = "%s"
    limit         = 30
  }

  quota_request {
    resource_name = "standardDSv3Family"
    location      = "%s"
    limit         = 20
  }
}

resource "azurerm_quota_group_subscription_allocation" "test" {
  quota_group_id  = azurerm_quota_group.test.id
  subscription_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  location        = "%s"

  allocation {
    resource_name = "standardddv4family"
    limit         = 15
  }

  allocation {
  resource_name = "standardDSv3Family"
  limit         = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Primary, data.Locations.Primary)
}
