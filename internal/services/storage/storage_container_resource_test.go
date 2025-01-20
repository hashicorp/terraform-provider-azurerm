// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

type StorageContainerResource struct{}

func TestAccStorageContainer_basicDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

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

func TestAccStorageContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_deleteAndRecreateDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
		},
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_deleteAndRecreate(t *testing.T) {
	t.Skip("skipping until https://github.com/Azure/azure-rest-api-specs/issues/30456 is resolved")
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_basicAzureADAuthDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAzureADAuthDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAzureADAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_requiresImportDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImportDeprecated),
	})
}

func TestAccStorageContainer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

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

func TestAccStorageContainer_updateDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateDeprecated(data, "private", "yes"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_access_type").HasValue("private"),
			),
		},
		{
			Config: r.updateDeprecated(data, "container", "no"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_access_type").HasValue("container"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data, "private", "yes"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_access_type").HasValue("private"),
			),
		},
		{
			Config: r.update(data, "container", "no"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("container_access_type").HasValue("container"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_encryptionScopeDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionScopeDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_encryptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_metaDataDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metaDataDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataUpdatedDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataEmptyDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metaData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_rootDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rootDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("$root"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_webDeprecated(t *testing.T) {
	if features.FivePointOhBeta() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("$web"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageContainer_web(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.web(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("$web"),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageContainerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	if !features.FivePointOhBeta() && !strings.HasPrefix(state.ID, "/subscriptions") {
		id, err := containers.ParseContainerID(state.ID, client.Storage.StorageDomainSuffix)
		if err != nil {
			return nil, err
		}

		account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Account %q for Container %q: %+v", id.AccountId.AccountName, id.ContainerName, err)
		}
		if account == nil {
			return nil, fmt.Errorf("unable to locate Storage Account %q", id.AccountId.AccountName)
		}

		containersClient, err := client.Storage.ContainersDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return nil, fmt.Errorf("building Containers Client: %+v", err)
		}

		prop, err := containersClient.Get(ctx, id.ContainerName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Container %q in %s: %+v", id.ContainerName, id.AccountId, err)
		}

		return pointer.To(prop != nil), nil
	}

	id, err := commonids.ParseStorageContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	existing, err := client.Storage.ResourceManager.BlobContainers.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(existing.Model != nil), nil
}

func (r StorageContainerResource) basicDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, r.template(data))
}

func (r StorageContainerResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
`, template)
}

func (r StorageContainerResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[3]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                              = "acctest-container-%[2]s"
  storage_account_id                = azurerm_storage_account.test.id
  container_access_type             = "private"
  default_encryption_scope          = azurerm_storage_encryption_scope.test.name
  encryption_scope_override_enabled = true

  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StorageContainerResource) basicAzureADAuthDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
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

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageContainerResource) basicAzureADAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
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

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageContainerResource) requiresImportDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "import" {
  name                  = azurerm_storage_container.test.name
  storage_account_name  = azurerm_storage_container.test.storage_account_name
  container_access_type = azurerm_storage_container.test.container_access_type
}
`, r.basicDeprecated(data))
}

func (r StorageContainerResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "import" {
  name                  = azurerm_storage_container.test.name
  storage_account_id    = azurerm_storage_container.test.storage_account_id
  container_access_type = azurerm_storage_container.test.container_access_type
}
`, template)
}

func (r StorageContainerResource) updateDeprecated(data acceptance.TestData, accessType, metadataVal string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
  metadata = {
    foo  = "bar"
    test = "%s"
  }
}
`, template, accessType, metadataVal)
}

func (r StorageContainerResource) update(data acceptance.TestData, accessType, metadataVal string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "%s"
  metadata = {
    foo  = "bar"
    test = "%s"
  }
}
`, template, accessType, metadataVal)
}

func (r StorageContainerResource) encryptionScopeDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"

  default_encryption_scope = azurerm_storage_encryption_scope.test.name
}
`, template, data.RandomInteger)
}

func (r StorageContainerResource) encryptionScope(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"

  default_encryption_scope = azurerm_storage_encryption_scope.test.name
}
`, template, data.RandomInteger)
}

func (r StorageContainerResource) metaDataDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageContainerResource) metaData(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageContainerResource) metaDataUpdatedDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"

  metadata = {
    hello = "world"
    panda = "pops"
  }
}
`, template)
}

func (r StorageContainerResource) metaDataUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"

  metadata = {
    hello = "world"
    panda = "pops"
  }
}
`, template)
}

func (r StorageContainerResource) metaDataEmptyDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"

  metadata = {}
}
`, template)
}

func (r StorageContainerResource) metaDataEmpty(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"

  metadata = {}
}
`, template)
}

func (r StorageContainerResource) rootDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$root"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, r.template(data))
}

func (r StorageContainerResource) webDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$web"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, template)
}

func (r StorageContainerResource) web(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$web"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
`, template)
}

func (r StorageContainerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func TestValidateStorageContainerName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"$root",
		"$web",
	}
	for _, v := range validNames {
		_, errors := validate.StorageContainerName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Storage Container Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"InvalidName1",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"ww",
		"$notroot",
		"$notweb",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		_, errors := validate.StorageContainerName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Storage Container Name", v)
		}
	}
}
