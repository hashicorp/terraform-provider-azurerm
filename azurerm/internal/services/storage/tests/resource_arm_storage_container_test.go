package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"strings"
	"testing"
)

func TestAccAzureRMStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_deleteAndRecreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_template(data),
			},
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basicAzureADAuth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageContainer_requiresImport),
		},
	})
}

func TestAccAzureRMStorageContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_update(data, "private"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStorageContainer_update(data, "container"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_metaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_metaDataUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_metaDataEmpty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_encryptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_encryptionScope(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_immutabilityPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_immutabilityPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_immutabilityPolicyUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_legalHold(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_legalHold(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_legalHoldUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
					testAccARMStorageContainerDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageContainer_root(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_root(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "$root"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageContainer_web(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageContainer_web(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageContainerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "$web"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageContainerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.ContainersClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		resp, err := client.GetProperties(ctx, accountName, containerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container %q (Account %q / Resource Group %q) does not exist", containerName, accountName, account.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on ContainersClient: %+v", err)
		}

		return nil
	}
}

func testAccARMStorageContainerDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.ContainersClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		if _, err := client.Delete(ctx, accountName, containerName); err != nil {
			return fmt.Errorf("Error deleting Container %q (Account %q): %s", containerName, accountName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageContainerDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_container" {
			continue
		}

		containerName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
		}
		if account == nil {
			return nil
		}

		client, err := storageClient.ContainersClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Containers Client: %s", err)
		}

		props, err := client.GetProperties(ctx, accountName, containerName)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Container still exists: %+v", props)
	}

	return nil
}

func testAccAzureRMStorageContainer_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, template)
}

func testAccAzureRMStorageContainer_basicAzureADAuth(data acceptance.TestData) string {
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

func testAccAzureRMStorageContainer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "import" {
  name                  = azurerm_storage_container.test.name
  storage_account_name  = azurerm_storage_container.test.storage_account_name
  container_access_type = azurerm_storage_container.test.container_access_type
}
`, template)
}

func testAccAzureRMStorageContainer_update(data acceptance.TestData, accessType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
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
  allow_blob_public_access = true

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, accessType)
}

func testAccAzureRMStorageContainer_metaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
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

func testAccAzureRMStorageContainer_metaDataUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
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

func testAccAzureRMStorageContainer_metaDataEmpty(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
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

func testAccAzureRMStorageContainer_root(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$root"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, template)
}

func testAccAzureRMStorageContainer_web(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "$web"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, template)
}

func testAccAzureRMStorageContainer_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func TestValidateArmStorageContainerName(t *testing.T) {
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

func testAccAzureRMStorageContainer_encryptionScope(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestES%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_storage_container" "test" {
  name                           = "vhds"
  storage_account_name           = azurerm_storage_account.test.name
  container_access_type          = "private"
  default_encryption_scope       = azurerm_storage_encryption_scope.test.name
  encryption_scope_for_all_blobs = true
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageContainer_immutabilityPolicy(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
  extended_immutability_policy {
    since_creation_in_days        = 7
    allow_protected_append_writes = true
  }
}
`, template)
}

func testAccAzureRMStorageContainer_immutabilityPolicyUpdate(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
  extended_immutability_policy {
    since_creation_in_days        = 3
    allow_protected_append_writes = false
  }
}
`, template)
}

func testAccAzureRMStorageContainer_legalHold(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
  legal_hold {
    tags = ["production"]
  }
}
`, template)
}

func testAccAzureRMStorageContainer_legalHoldUpdate(data acceptance.TestData) string {
	template := testAccAzureRMStorageContainer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"

  legal_hold {
    tags = ["production", "stage", "stage2"]
  }
}
`, template)
}
