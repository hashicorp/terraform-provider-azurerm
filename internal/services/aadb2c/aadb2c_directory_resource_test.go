package aadb2c_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AadB2cDirectoryResource struct{}

func TestAccAadB2cDirectoryResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aadb2c_directory", "test")
	r := AadB2cDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("country_code", "display_name"),
	})
}

func TestAccAadB2cDirectoryResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aadb2c_directory", "test")
	r := AadB2cDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("country_code", "display_name"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("country_code", "display_name"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("country_code", "display_name"),
	})
}

func TestAccAadB2cDirectoryResource_domainNameUnavailable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aadb2c_directory", "test")
	r := AadB2cDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.domainNameUnavailable(data),
			ExpectError: regexp.MustCompile("checking availability of `domain_name`: the specified domain \"[^\"]+\" is unavailable"),
		},
	})
}

func TestAccAadB2cDirectoryResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aadb2c_directory", "test")
	r := AadB2cDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func (r AadB2cDirectoryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tenants.ParseB2CDirectoryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AadB2c.TenantsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse.StatusCode == http.StatusNotFound {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r AadB2cDirectoryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  client_id               = ""
  client_certificate_path = ""
  client_secret           = ""
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AadB2cDirectoryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_aadb2c_directory" "test" {
  country_code            = "US"
  data_residency_location = "United States"
  display_name            = "acctest%[2]d"
  domain_name             = "acctest%[2]d.onmicrosoft.com"
  resource_group_name     = azurerm_resource_group.test.name
  sku_name                = "PremiumP1"
}
`, r.template(data), data.RandomInteger)
}

func (r AadB2cDirectoryResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_aadb2c_directory" "test" {
  country_code            = "US"
  data_residency_location = "United States"
  display_name            = "acctest%[2]d"
  domain_name             = "acctest%[2]d.onmicrosoft.com"
  resource_group_name     = azurerm_resource_group.test.name
  sku_name                = "PremiumP2"

  tags = {
    "Environment" : "Test",
    "Project" : "Locksmith",
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AadB2cDirectoryResource) domainNameUnavailable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "duplicate" {
  name     = "acctestRG-duplicate-%[2]d"
  location = "%[3]s"
}

resource "azurerm_aadb2c_directory" "duplicate" {
  country_code            = azurerm_aadb2c_directory.test.country_code
  data_residency_location = azurerm_aadb2c_directory.test.data_residency_location
  display_name            = "acctest-duplicate-%[2]d"
  domain_name             = azurerm_aadb2c_directory.test.domain_name
  resource_group_name     = azurerm_resource_group.duplicate.name
  sku_name                = "PremiumP1"
}
`, r.basic(data), data.RandomInteger, data.Locations.Secondary)
}

func (r AadB2cDirectoryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_aadb2c_directory" "import" {
  country_code            = azurerm_aadb2c_directory.test.country_code
  data_residency_location = azurerm_aadb2c_directory.test.data_residency_location
  display_name            = azurerm_aadb2c_directory.test.display_name
  domain_name             = azurerm_aadb2c_directory.test.domain_name
  resource_group_name     = azurerm_aadb2c_directory.test.resource_group_name
  sku_name                = azurerm_aadb2c_directory.test.sku_name
}
`, r.basic(data))
}
