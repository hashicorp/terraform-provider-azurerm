package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/importexport/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageExportJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_export_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageExportJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_information.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "shipping_information.0.country_or_region"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "shipping_information.0.recipient_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "shipping_information.0.street_address1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageExportJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_export_job", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageExportJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageExportJob_requiresImport),
		},
	})
}

func TestAccAzureRMStorageExportJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_export_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageExportJob_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageExportJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_export_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageExportJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageExportJob_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageExportJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageExportJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ImportExport.JobClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Azure Export Job not found: %s", resourceName)
		}
		id, err := parse.StorageImportExportJobID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name, id.ResourceGroup); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: azure export job %q does not exist", id.Name)
			}
			return fmt.Errorf("Bad: Get on AzureImportExport JobClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMStorageExportJobDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ImportExport.JobClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_export_job" {
			continue
		}
		id, err := parse.StorageImportExportJobID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name, id.ResourceGroup); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on AzureImportExport jobClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMStorageExportJob_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageExportJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_export_job" "test" {
  name                = "acctest-export-job-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_account_id        = azurerm_storage_account.test.id
  export_blob_path_prefixes = ["/a", "/b"]

  return_shipping {
    carrier_account_number = "123456789"
    carrier_name           = "DHL"
  }

  return_address {
    recipient_name    = "Tets"
    city              = "Redmond"
    country_or_region = "USA"
    email             = "Test@contoso.com"
    phone             = "123456789"
    postal_code       = "98007"
    street_address1   = "Street1"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageExportJob_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMStorageExportJob_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_export_job" "import" {
  name                = azurerm_export_job.test.name
  location            = azurerm_export_job.test.location
  resource_group_name = azurerm_export_job.test.resource_group_name

  storage_account_id        = azurerm_storage_account.test.id
  export_blob_path_prefixes = ["/a", "/b"]

  return_shipping {
    carrier_account_number = "123456789"
    carrier_name           = "DHL"
  }

  return_address {
    recipient_name    = "Tets"
    city              = "Redmond"
    country_or_region = "USA"
    email             = "Test@contoso.com"
    phone             = "123456789"
    postal_code       = "98007"
    street_address1   = "Street1"
  }
}
`, config)
}

func testAccAzureRMStorageExportJob_complete(data acceptance.TestData) string {
	template := testAccAzureRMStorageExportJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_export_job" "test" {
  name                = "acctest-export-job-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_account_id        = azurerm_storage_account.test.id
  export_blob_path_prefixes = ["/a", "/b"]

  return_shipping {
    carrier_account_number = "123456789"
    carrier_name           = "DHL"
  }

  backup_drive_manifest = true
  diagnostics_path      = "waimportexport"
  log_level             = "Verbose"

  return_address {
    recipient_name    = "Tets"
    city              = "Redmond"
    country_or_region = "USA"
    email             = "Test@contoso.com"
    phone             = "123456789"
    postal_code       = "98007"
    street_address1   = "Street1"
    state_or_province = "wa"
    street_address2   = "street2"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageExportJob_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
