// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type QuotaGroupResource struct{}

func TestAccQuotaGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group", "test")
	r := QuotaGroupResource{}

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

func TestAccQuotaGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group", "test")
	r := QuotaGroupResource{}

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

func TestAccQuotaGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group", "test")
	r := QuotaGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("Acceptance Test Quota Group %d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccQuotaGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_quota_group", "test")
	r := QuotaGroupResource{}

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
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("Acceptance Test Quota Group %d", data.RandomInteger)),
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

func (r QuotaGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := groupquotas.ParseGroupQuotaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Quota.GroupQuotasClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r QuotaGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "test" {
  name                = "acctestqg-%d"
  management_group_id = data.azurerm_management_group.test.id
}
`, data.RandomInteger)
}

func (r QuotaGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_quota_group" "import" {
  name                = azurerm_quota_group.test.name
  management_group_id = azurerm_quota_group.test.management_group_id
}
`, r.basic(data))
}

func (r QuotaGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "test" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "test" {
  name                = "acctestqg-%d"
  management_group_id = data.azurerm_management_group.test.id
  display_name        = "Acceptance Test Quota Group %d"

  quota_request {
    resource_name = "standardddv4family"
    location      = "%s"
    limit         = 10
    comment       = "Acceptance test quota request"
  }
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}
