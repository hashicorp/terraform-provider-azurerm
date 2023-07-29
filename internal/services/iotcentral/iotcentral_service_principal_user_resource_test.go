// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IoTCentralServicePrincipalUserResource struct{}

func TestAccIoTCentralServicePrincipalUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_service_principal_user", "test")
	r := IoTCentralServicePrincipalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("roles.0.role").HasValue(appAdminRoleId),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralServicePrincipalUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_service_principal_user", "test")
	r := IoTCentralServicePrincipalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("roles.0.role").HasValue(orgAdminRoleId),
				check.That(data.ResourceName).Key("roles.0.organization").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func (IoTCentralServicePrincipalUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	userClient, err := clients.IoTCentral.UsersClient(ctx, id.SubDomain)
	if err != nil {
		return nil, fmt.Errorf("creating users client: %+v", err)
	}

	resp, err := userClient.Get(ctx, id.Id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	_, success := resp.Value.AsServicePrincipalUser()

	return &success, nil
}

func (r IoTCentralServicePrincipalUserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  display_name = "acctest-iotcentralsp-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

%s

resource "azurerm_iotcentral_service_principal_user" "test" {
  sub_domain = azurerm_iotcentral_application.test.sub_domain
  object_id  = azuread_service_principal.test.object_id
  tenant_id  = data.azurerm_client_config.current.tenant_id

  roles {
    role = data.azurerm_iotcentral_role.app_admin.id
  }
}
`, data.RandomInteger, r.templateBasic(data))
}

func (r IoTCentralServicePrincipalUserResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  display_name = "acctest-iotcentralsp-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

%s

resource "azurerm_iotcentral_service_principal_user" "test" {
  sub_domain = azurerm_iotcentral_application.test.sub_domain
  object_id  = azuread_service_principal.test.object_id
  tenant_id  = data.azurerm_client_config.current.tenant_id

  roles {
    role         = data.azurerm_iotcentral_role.org_admin.id
    organization = azurerm_iotcentral_organization.test.organization_id
  }
}
`, data.RandomInteger, r.templateComplete(data))
}

func (IoTCentralServicePrincipalUserResource) templateBasic(data acceptance.TestData) string {
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

data "azurerm_iotcentral_role" "app_admin" {
  sub_domain   = azurerm_iotcentral_application.test.sub_domain
  display_name = "%[3]s"
}
`, data.RandomInteger, data.Locations.Primary, appAdminRoleDisplayName)
}

func (IoTCentralServicePrincipalUserResource) templateComplete(data acceptance.TestData) string {
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

resource "azurerm_iotcentral_organization" "test" {
  sub_domain   = azurerm_iotcentral_application.test.sub_domain
  display_name = "Org"
}

data "azurerm_iotcentral_role" "org_admin" {
  sub_domain   = azurerm_iotcentral_application.test.sub_domain
  display_name = "%[3]s"

  depends_on = [azurerm_iotcentral_organization.test]
}
`, data.RandomInteger, data.Locations.Primary, orgAdminRoleDisplayName)
}
