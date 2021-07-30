package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementProductResource struct {
}

func TestAccApiManagementProduct_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_required").HasValue("false"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Product"),
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("published").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("false"),
				check.That(data.ResourceName).Key("terms").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementProduct_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

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

func TestAccApiManagementProduct_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_required").HasValue("false"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Product"),
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("published").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("false"),
				check.That(data.ResourceName).Key("terms").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_required").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Updated Product"),
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("published").HasValue("true"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
				check.That(data.ResourceName).Key("terms").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Product"),
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("published").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("false"),
				check.That(data.ResourceName).Key("terms").HasValue(""),
			),
		},
	})
}

func TestAccApiManagementProduct_subscriptionsLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subscriptionLimits(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_required").HasValue("true"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
				check.That(data.ResourceName).Key("subscriptions_limit").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementProduct_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_required").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue("This is an example description"),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Product"),
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("published").HasValue("true"),
				check.That(data.ResourceName).Key("subscriptions_limit").HasValue("2"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
				check.That(data.ResourceName).Key("terms").HasValue("These are some example terms and conditions"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementProduct_approvalRequiredError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementProductResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.approvalRequiredError(data),
			ExpectError: regexp.MustCompile("`subscription_required` must be true and `subscriptions_limit` must be greater than 0 to use `approval_required`"),
		},
	})
}

func (ApiManagementProductResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	productId := id.Path["products"]

	resp, err := clients.ApiManagement.ProductsClient.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Product (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementProductResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = false
  published             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementProductResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product" "import" {
  product_id            = azurerm_api_management_product.test.product_id
  api_management_name   = azurerm_api_management_product.test.api_management_name
  resource_group_name   = azurerm_api_management_product.test.resource_group_name
  display_name          = azurerm_api_management_product.test.display_name
  subscription_required = azurerm_api_management_product.test.subscription_required
  approval_required     = azurerm_api_management_product.test.approval_required
  published             = azurerm_api_management_product.test.published
}
`, r.basic(data))
}

func (ApiManagementProductResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Updated Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 1
  published             = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementProductResource) subscriptionLimits(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 2
  published             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementProductResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  published             = true
  subscriptions_limit   = 2
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementProductResource) approvalRequiredError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  approval_required     = true
  subscription_required = false
  published             = true
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
