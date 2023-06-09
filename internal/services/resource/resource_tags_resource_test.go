package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceTagsResource struct{}

func TestAccResourceTags_ReadOnlyBasicKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tags", "test")
	r := ResourceTagsResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.ReadOnlyBasicKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceTags_ReadOnlyBasicStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tags", "test")
	r := ResourceTagsResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.ReadOnlyBasicStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceTags_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tags", "test")
	r := ResourceTagsResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.ReadOnlyBasicKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t ResourceTagsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ResourceTagsID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.TagsClient.GetAtScope(ctx, id.ParentResourceID())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Properties.Tags != nil), nil
}

func (ResourceTagsResource) ReadOnlyBasicKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"

  enable_rbac_authorization  = true

  contact {
	email = "test@example.com"
  }

  lifecycle {
	ignore_changes = [tags]
  }
}

resource "azurerm_resource_tags" "test" {
  resource_id = azurerm_key_vault.test.id
  tags = {
    "test" = "test",
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ResourceTagsResource) ReadOnlyBasicStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "tagstest%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  lifecycle {
	ignore_changes = [tags]
  }
}

resource "azurerm_resource_tags" "test" {
  resource_id = azurerm_storage_account.test.id
  tags = {
    "test" = "test",
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ResourceTagsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_tags" "import" {
  resource_id       = azurerm_resource_tags.test.resource_id
  tags      		= azurerm_resource_tags.test.tags
}
`, r.ReadOnlyBasicKeyVault(data))
}
