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
	data := acceptance.BuildTestData(t, "azurerm_databox_job", "test")

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

func TestAccAzureRMDataBoxJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_job", "test")

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
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		},
	})
}

func TestAccAzureRMDataBoxJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataBoxJob_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.name", "DataBoxJobTester"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.emails.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_number", "+11234567891"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.mobile", "+11234567891"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.at_azure_dc", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.data_copied", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.delivered", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.device_prepared", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.dispatched", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.picked_up", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_extension", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.city", "San Francisco"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code", "94107"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_1", "16 TOWNSEND ST"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.address_type", "Commercial"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.company_name", "Microsoft"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_2", "17 TOWNSEND ST"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_3", "18 TOWNSEND ST"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code_ext", "94107"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TEST"),
				),
			},
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
			{
				Config: testAccAzureRMDataBoxJob_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataBoxJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.name", "DataBoxJobTester2"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.emails.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_number", "+11234567892"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.mobile", "+11234567892"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.at_azure_dc", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.data_copied", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.delivered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.device_prepared", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.dispatched", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.notification_preference.0.picked_up", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "contact_details.0.phone_extension", "124"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.city", "San Diego"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code", "92111"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_1", "6901 SUN STREET"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.address_type", "Residential"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.company_name", "Intel"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_2", "6902 SUN STREET"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.street_address_3", "6903 SUN STREET"),
					resource.TestCheckResourceAttr(data.ResourceName, "shipping_address.0.postal_code_ext", "92111"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TEST2"),
				),
			},
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		},
	})
}

func TestAccAzureRMDataBoxJob_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_databox_job", "test")

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
				ExpectError: acceptance.RequiresImportError("azurerm_databox_job"),
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
			return fmt.Errorf("DataBox Job not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: DataBox Job %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on DataBox.JobClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDataBoxJobDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataBox.JobClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databox_job" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on DataBox.JobClient: %+v", err)
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

resource "azurerm_databox_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_account {
    type               = "StorageAccount"
    storage_account_id = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  preferred_shipment_type = "CustomerManaged"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccAzureRMDataBoxJob_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databox_job" "import" {
  name                = "${azurerm_databox_job.test.name}"
  location            = "${azurerm_databox_job.test.location}"
  resource_group_name = "${azurerm_databox_job.test.resource_group_name}"

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_account {
    type               = "StorageAccount"
    storage_account_id = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  preferred_shipment_type = "CustomerManaged"

  sku_name = "DataBox"
}
`, testAccAzureRMDataBoxJob_basic(data))
}

func testAccAzureRMDataBoxJob_complete(data acceptance.TestData) string {
	template := testAccAzureRMDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {}

resource "azurerm_storage_account" "test2" {
  name                = "acctestrgstorage2%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_databox_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
    mobile       = "+11234567891"
    notification_preference {
      at_azure_dc     = false
      data_copied     = false
      delivered       = false
      device_prepared = true
      dispatched      = false
      picked_up       = false
    }
    phone_extension = "123"
  }

  destination_account {
    type               = "StorageAccount"
    storage_account_id = "${azurerm_storage_account.test2.id}"
    share_password     = "fddbc123123aa@"
  }

  destination_account {
    type                       = "ManagedDisk"
    resource_group_id          = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/TestManagedRG%s"
    staging_storage_account_id = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
    address_type      = "Commercial"
    company_name      = "Microsoft"
    street_address_2  = "17 TOWNSEND ST"
    street_address_3  = "18 TOWNSEND ST"
    postal_code_ext   = "94107"
  }

  expected_data_size_in_tb          = 5
  databox_disk_passkey              = "abcabc123123@"
  databox_preferred_disk_count      = 5
  databox_preferred_disk_size_in_tb = 2
  sku_name                          = "DataBoxDisk"
  datacenter_region_preference      = ["westus", "eastus"]
  delivery_type                     = "Scheduled"
  delivery_scheduled_date_time      = "2020-04-01T05:30:00+05:30"
  preferred_shipment_type           = "CustomerManaged"

  tags = {
    env = "TEST"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString)
}

func testAccAzureRMDataBoxJob_update(data acceptance.TestData) string {
	template := testAccAzureRMDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {}

resource "azurerm_storage_account" "test2" {
  name                = "acctestrgstorage2%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_databox_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  contact_details {
    name         = "DataBoxJobTester2"
    emails       = ["some.user@example.com", "some.user2@example.com"]
    phone_number = "+11234567892"
    mobile       = "+11234567892"
    notification_preference {
      at_azure_dc     = true
      data_copied     = true
      delivered       = true
      device_prepared = false
      dispatched      = true
      picked_up       = true
    }
    phone_extension = "124"
  }

  destination_account {
    type               = "StorageAccount"
    storage_account_id = "${azurerm_storage_account.test2.id}"
    share_password     = "fddbc123123aa@"
  }

  destination_account {
    type                       = "ManagedDisk"
    resource_group_id          = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/TestManagedRG%s"
    staging_storage_account_id = "${azurerm_storage_account.test.id}"
  }

  shipping_address {
    city              = "San Diego"
    country           = "US"
    postal_code       = "92111"
    state_or_province = "CA"
    street_address_1  = "6901 SUN STREET"
    address_type      = "Residential"
    company_name      = "Intel"
    street_address_2  = "6902 SUN STREET"
    street_address_3  = "6903 SUN STREET"
    postal_code_ext   = "92111"
  }

  expected_data_size_in_tb          = 6
  databox_disk_passkey              = "abcabc123123@"
  databox_preferred_disk_count      = 5
  databox_preferred_disk_size_in_tb = 2
  sku_name                          = "DataBoxDisk"
  datacenter_region_preference      = ["westus", "eastus"]
  delivery_type                     = "Scheduled"
  delivery_scheduled_date_time      = "2020-04-01T05:30:00+05:30"
  preferred_shipment_type           = "CustomerManaged"

  tags = {
    env = "TEST2"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString)
}

func testAccAzureRMDataBoxJob_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
