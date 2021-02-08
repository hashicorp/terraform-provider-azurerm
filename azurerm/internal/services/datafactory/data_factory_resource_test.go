package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataFactoryResource struct {
}

func TestAccDataFactory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_tagsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		{
			Config: r.tagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
				check.That(data.ResourceName).Key("tags.updated").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccDataFactory_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.github(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("github_configuration.0.account_name").HasValue(fmt.Sprintf("acctestGH-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("github_configuration.0.git_url").HasValue("https://github.com/terraform-providers/"),
				check.That(data.ResourceName).Key("github_configuration.0.repository_name").HasValue("terraform-provider-azurerm"),
				check.That(data.ResourceName).Key("github_configuration.0.branch_name").HasValue("master"),
				check.That(data.ResourceName).Key("github_configuration.0.root_folder").HasValue("/"),
			),
		},
		{
			Config: r.githubUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("github_configuration.0.account_name").HasValue(fmt.Sprintf("acctestGitHub-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("github_configuration.0.git_url").HasValue("https://github.com/terraform-providers/"),
				check.That(data.ResourceName).Key("github_configuration.0.repository_name").HasValue("terraform-provider-azuread"),
				check.That(data.ResourceName).Key("github_configuration.0.branch_name").HasValue("stable-website"),
				check.That(data.ResourceName).Key("github_configuration.0.root_folder").HasValue("/azuread"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_publicNetworkDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactory_keyVaultKeyEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVaultKeyEncryption(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DataFactoryResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["factories"]

	resp, err := clients.DataFactory.FactoriesClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (DataFactoryResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["factories"]

	resp, err := client.DataFactory.FactoriesClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return nil, fmt.Errorf("delete on dataFactoryClient: %+v", err)
		}
	}

	return utils.Bool(true), nil
}

func (DataFactoryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DataFactoryResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DataFactoryResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
    updated     = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DataFactoryResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DataFactoryResource) github(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  github_configuration {
    git_url         = "https://github.com/terraform-providers/"
    repository_name = "terraform-provider-azurerm"
    branch_name     = "master"
    root_folder     = "/"
    account_name    = "acctestGH-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DataFactoryResource) githubUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  github_configuration {
    git_url         = "https://github.com/terraform-providers/"
    repository_name = "terraform-provider-azuread"
    branch_name     = "stable-website"
    root_folder     = "/azuread"
    account_name    = "acctestGitHub-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DataFactoryResource) publicNetworkDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  public_network_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DataFactoryResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_factory" "test" {
  name                = "acctest%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    user_identity_ids = [
		azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DataFactoryResource) keyVaultKeyEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                     = "acckv%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
      "delete",
      "purge"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "get",
      "unwrapKey",
      "wrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]
}

resource "azurerm_data_factory" "test" {
  name                = "acctest%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    user_identity_ids = [
		azurerm_user_assigned_identity.test.id
    ]
  }

  key_vault_key_id = azurerm_key_vault_key.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
