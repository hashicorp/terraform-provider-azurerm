package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccDataBoxJob_basic(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBoxJob_basic(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataBoxJob_complete(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBoxJob_complete(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		},
	})
}

func TestAccDataBoxJob_requiresImport(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBoxJob_basic(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccDataBoxJob_requiresImport),
		},
	})
}

func TestAccDataBoxJob_update(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBoxJob_complete(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
			{
				Config: testAccDataBoxJob_update(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		},
	})
}

func TestAccDataBoxJob_withCustomerManaged(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataBoxJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBoxJob_withCustomerManaged(data, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataBoxJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDataBoxJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataBox.JobClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Data Box Job not found: %s", resourceName)
		}

		id, err := parse.DataBoxJobID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Data Box Job %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on DataBox.JobClient: %+v", err)
		}

		return nil
	}
}

func testCheckDataBoxJobDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataBox.JobClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_box_job" {
			continue
		}

		id, err := parse.DataBoxJobID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on DataBox.JobClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testGetLocationFromSubscription() (string, error) {
	subscription := strings.ToLower(os.Getenv("ARM_SUBSCRIPTION_ID"))
	location := ""

	if strings.HasPrefix(subscription, "67a9759d") || strings.HasPrefix(subscription, "85b3dbca") {
		location = "westus"
	} else if strings.HasPrefix(subscription, "1a6092a6") || strings.HasPrefix(subscription, "88720cb0") {
		location = "westcentralus"
	} else {
		return "", fmt.Errorf("Skipping since test is not running as one of the four valid subscriptions allowed to run DataBox tests")
	}

	return location, nil
}

func testAccDataBoxJob_basic(data acceptance.TestData, location string) string {
	template := testAccDataBoxJob_template(data, location)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  preferred_shipment_type = "MicrosoftManaged"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccDataBoxJob_requiresImport(data acceptance.TestData) string {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		return ""
	}

	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "import" {
  name                = azurerm_data_box_job.test.name
  location            = azurerm_data_box_job.test.location
  resource_group_name = azurerm_data_box_job.test.resource_group_name

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  preferred_shipment_type = "MicrosoftManaged"

  sku_name = "DataBox"
}
`, testAccDataBoxJob_basic(data, location))
}

func testAccDataBoxJob_complete(data acceptance.TestData, location string) string {
	template := testAccDataBoxJob_template(data, location)
	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {}

resource "azurerm_storage_account" "test2" {
  name                = "acctestrgstorage2%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
    phone_mobile = "+11234567891"
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

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test2.id
    share_password     = "fddbc123123aa@"
  }

  destination_managed_disk {
    resource_group_id          = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/TestManagedRG%s"
    staging_storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city                  = "San Francisco"
    country               = "US"
    postal_code           = "94107"
    state_or_province     = "CA"
    street_address_1      = "16 TOWNSEND ST"
    address_type          = "Commercial"
    company_name          = "Microsoft"
    street_address_2      = "17 TOWNSEND ST"
    street_address_3      = "18 TOWNSEND ST"
    postal_code_plus_four = "94107"
  }

  expected_data_size_in_tb = 5
  databox_disk_passkey     = "abcabc123123@"

  databox_preferred_disk {
    size_in_tb = 2
    count      = 5
  }

  sku_name                     = "DataBoxDisk"
  datacenter_region_preference = ["westus", "westcentralus"]
  delivery_type                = "NonScheduled"
  preferred_shipment_type      = "MicrosoftManaged"

  tags = {
    env = "TEST"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString)
}

func testAccDataBoxJob_update(data acceptance.TestData, location string) string {
	template := testAccDataBoxJob_template(data, location)
	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {}

resource "azurerm_storage_account" "test2" {
  name                = "acctestrgstorage2%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "DataBoxJobTester2"
    emails       = ["some.user@example.com", "some.user2@example.com"]
    phone_number = "+11234567892"
    phone_mobile = "+11234567892"
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

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test2.id
    share_password     = "fddbc123123aa@"
  }

  destination_managed_disk {
    resource_group_id          = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/TestManagedRG%s"
    staging_storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city                  = "San Diego"
    country               = "US"
    postal_code           = "92111"
    state_or_province     = "CA"
    street_address_1      = "6901 SUN STREET"
    address_type          = "Residential"
    company_name          = "Intel"
    street_address_2      = "6902 SUN STREET"
    street_address_3      = "6903 SUN STREET"
    postal_code_plus_four = "92111"
  }

  expected_data_size_in_tb = 5
  databox_disk_passkey     = "abcabc123123@"

  databox_preferred_disk {
    size_in_tb = 2
    count      = 5
  }

  sku_name                     = "DataBoxDisk"
  datacenter_region_preference = ["westus", "westcentralus"]
  delivery_type                = "NonScheduled"
  preferred_shipment_type      = "MicrosoftManaged"

  tags = {
    env = "TEST2"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomString)
}

func testAccDataBoxJob_withCustomerManaged(data acceptance.TestData, location string) string {
	template := testAccDataBoxJob_template(data, location)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  preferred_shipment_type      = "CustomerManaged"
  delivery_scheduled_date_time = "2020-04-01T05:30:00+05:30"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccDataBoxJob_template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databoxjob-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestrgstorage%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}
`, data.RandomInteger, location, data.RandomString)
}
