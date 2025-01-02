// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IoTCentralOrganizationResource struct{}

func TestAccIoTCentralOrganization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_organization", "test")
	r := IoTCentralOrganizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("organization_id").HasValue("org-test-id"),
				check.That(data.ResourceName).Key("display_name").HasValue("Org basic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralOrganization_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_organization", "test")
	r := IoTCentralOrganizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("organization_id").HasValue("org-test-id"),
				check.That(data.ResourceName).Key("display_name").HasValue("Org child"),
				check.That(data.ResourceName).Key("parent_organization_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralOrganization_updateDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_organization", "test")
	r := IoTCentralOrganizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Org basic"),
			),
		},
		{
			Config: r.basicUpdatedDisplayName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Org basic updated"),
			),
		},
		data.ImportStep(),
	})
}

func (IoTCentralOrganizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.OrganizationID(state.ID)
	if err != nil {
		return nil, err
	}

	appId, err := apps.ParseIotAppID(state.Attributes["iotcentral_application_id"])
	if err != nil {
		return nil, err
	}

	app, err := clients.IoTCentral.AppsClient.Get(ctx, *appId)
	if err != nil || app.Model == nil {
		return nil, fmt.Errorf("checking for the presence of existing %q: %+v", appId, err)
	}

	orgClient, err := clients.IoTCentral.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
	if err != nil {
		return nil, fmt.Errorf("creating organization client: %+v", err)
	}

	resp, err := orgClient.Get(ctx, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil || *resp.ID == ""), nil
}

func (r IoTCentralOrganizationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_iotcentral_organization" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  organization_id           = "org-test-id"
  display_name              = "Org basic"
}
`, r.template(data))
}

func (r IoTCentralOrganizationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_iotcentral_organization" "test_parent" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  organization_id           = "org-test-parent-id"
  display_name              = "Org parent"
}
resource "azurerm_iotcentral_organization" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  organization_id           = "org-test-id"
  display_name              = "Org child"

  parent_organization_id = azurerm_iotcentral_organization.test_parent.organization_id
}
`, r.template(data))
}

func (r IoTCentralOrganizationResource) basicUpdatedDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_iotcentral_organization" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  organization_id           = "org-test-id"
  display_name              = "Org basic updated"
}
`, r.template(data))
}

func (IoTCentralOrganizationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  display_name        = "some-display-name"
  sku                 = "ST0"
}
`, data.RandomInteger, data.Locations.Primary)
}
