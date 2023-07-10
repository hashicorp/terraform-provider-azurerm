// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServicePublicCertificateResource struct{}

func TestAccAppServicePublicCertificate_basicAppService(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_public_certificate", "test")
	r := AppServicePublicCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAppService(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").HasValue("8155941A075F972A60AE1C74C749BBDDF82960B5"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServicePublicCertificate_basicFunctionApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_public_certificate", "test")
	r := AppServicePublicCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFunctionApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").HasValue("8155941A075F972A60AE1C74C749BBDDF82960B5"),
			),
		},
		data.ImportStep(),
	})
}

func (r AppServicePublicCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PublicCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.WebAppsClient.GetPublicCertificate(ctx, id.ResourceGroup, id.SiteName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Public Certificate %q (Resource Group %q, App Service %q): %+v", id.Name, id.ResourceGroup, id.SiteName, err)
	}
	return utils.Bool(true), nil
}

func (r AppServicePublicCertificateResource) basicAppService(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestpubcert%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestpubcert-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestpubcert-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

resource "azurerm_app_service_public_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  app_service_name     = azurerm_app_service.test.name
  certificate_name     = "acctestpubcert-%d"
  certificate_location = "Unknown"
  blob                 = filebase64("testdata/app_service_public_certificate.cer")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServicePublicCertificateResource) basicFunctionApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestpubcert-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestpubcert-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctestpubcert-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_app_service_public_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  app_service_name     = azurerm_function_app.test.name
  certificate_name     = "acctestpubcert-%d"
  certificate_location = "Unknown"
  blob                 = filebase64("testdata/app_service_public_certificate.cer")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
