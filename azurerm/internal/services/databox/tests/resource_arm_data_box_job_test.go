package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataBoxJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataBoxJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataBoxJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataBoxJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.contact_name", "DataBoxJobTester"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_number", "+11234567891"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.city", "San Francisco"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code", "94107"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_1", "16 TOWNSEND ST"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataBoxJob_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.contact_name", "DataBoxJobTester2"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_number", "+112345678912"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.city", "San Diego"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code", "92111"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_1", "6901 SUN STREET"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataBoxJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataBoxJob_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep("expected_data_size_in_tb"),
		},
	})
}

func TestAccAzureRMDataBoxJob_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataBoxJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDataBoxJob_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_data_box_job"),
			},
		},
	})
}

func testCheckAzureRMDataBoxJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataBox.JobClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Data Box Job not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Data Box Job %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on dataBox.JobClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDataBoxJobDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataBox.JobClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_box_job" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on dataBox.JobClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDataBoxJob_basic(data acceptance.TestData) string {
	template := testAccAzureRMDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    contact_name = "DataBoxJobTester"
    emails = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_account_details {
    data_destination_type = "StorageAccount"
    storage_account_id    = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city    = "San Francisco"
    country = "US"
    postal_code = "94107"
    state_or_province = "CA"
    street_address_1 = "16 TOWNSEND ST"
  }

  preferences {
    preferred_shipment_type = "CustomerManaged"
  }

  expected_data_size_in_tb = 27

  disk_pass_key = "abcabc123123123@"

  delivery_scheduled_date_time = "2020-03-01T05:30:00+05:30"

  delivery_type = "Scheduled"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccAzureRMDataBoxJob_update(data acceptance.TestData) string {
	template := testAccAzureRMDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    contact_name = "DataBoxJobTester2"
    emails = ["some.user2@example.com"]
    phone_number = "+112345678912"
  }

  destination_account_details {
    data_destination_type = "StorageAccount"
    storage_account_id    = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city    = "San Diego"
    country = "US"
    postal_code = "92111"
    state_or_province = "CA"
    street_address_1 = "6901 SUN STREET"
  }

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccAzureRMDataBoxJob_complete(data acceptance.TestData) string {
	template := testAccAzureRMDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {}

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    contact_name = "DataBoxJobTester"
    emails = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_account_details {
    data_destination_type         = "ManagedDisk"
	resource_group_id             = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/TestManagedRG%s"
    staging_storage_account_id    = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city    = "San Francisco"
    country = "US"
    postal_code = "94107"
    state_or_province = "CA"
    street_address_1 = "16 TOWNSEND ST"
  }
  
  expected_data_size_in_tb = 5

  sku_name = "DataBoxDisk"
  
  tags = {
     env = "TEST"
  }
}
`, template, data.RandomString, data.RandomString)
}

func testAccAzureRMDataBoxJob_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "import" {
  name                = "${azurerm_data_box_job.test.name}"
  location            = "${azurerm_data_box_job.test.location}"
  resource_group_name = "${azurerm_data_box_job.test.resource_group_name}"

  contact_details {
    contact_name = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_account_details {
    data_destination_type = "StorageAccount"
    storage_account_id    = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city    = "San Francisco"
    country = "US"
    postal_code = "94107"
    state_or_province = "CA"
    street_address_1 = "16 TOWNSEND ST"
  }

  sku_name = "DataBox"
}
`, testAccAzureRMDataBoxJob_basic(data))
}

func testAccAzureRMDataBoxJob_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databoxjob-%d"
  location = "westus"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestrgstorage%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}
`, data.RandomInteger, data.RandomString)
}
