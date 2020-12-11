package databox_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataBoxJobResource struct {
}

func TestAccDataBoxJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")
	r := DataBoxJobResource{}

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

func TestAccDataBoxJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")
	r := DataBoxJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
	})
}

func TestAccDataBoxJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")
	r := DataBoxJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.requiresImport(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataBoxJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")
	r := DataBoxJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("databox_disk_passkey", "expected_data_size_in_tb"),
	})
}

func TestAccDataBoxJob_withCustomerManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_box_job", "test")
	r := DataBoxJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withCustomerManaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (DataBoxJobResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DataBoxJobID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBox.JobClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving DataBoxJob %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.JobProperties != nil), nil
}

func (DataBoxJobResource) basic(data acceptance.TestData) string {
	template := testAccDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "Terraform Test"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "Redmond"
    country           = "US"
    postal_code       = "98052"
    state_or_province = "WA"
    street_address_1  = "One Microsoft Way"
  }

  preferred_shipment_type = "MicrosoftManaged"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func (r DataBoxJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "import" {
  name                = azurerm_data_box_job.test.name
  location            = azurerm_data_box_job.test.location
  resource_group_name = azurerm_data_box_job.test.resource_group_name

  contact_details {
    name         = "Terraform Test"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "Redmond"
    country           = "US"
    postal_code       = "98052"
    state_or_province = "WA"
    street_address_1  = "One Microsoft Way"
  }

  preferred_shipment_type = "MicrosoftManaged"

  sku_name = "DataBox"
}
`, r.basic(data))
}

func (DataBoxJobResource) complete(data acceptance.TestData) string {
	template := testAccDataBoxJob_template(data)
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
    name         = "Terraform Test"
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
    city                  = "Redmond"
    country               = "US"
    postal_code           = "98052"
    state_or_province     = "WA"
    street_address_1      = "One Microsoft Way"
    address_type          = "Commercial"
		company_name          = "Microsoft"
    street_address_2      = "Two Microsoft Way"
    street_address_3      = "Three Microsoft Way"
    postal_code_plus_four = "6399"
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

func (DataBoxJobResource) update(data acceptance.TestData) string {
	template := testAccDataBoxJob_template(data)
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
    name         = "Terraform Test"
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
    city                  = "Redmond"
    country               = "US"
    postal_code           = "98052"
    state_or_province     = "WA"
    street_address_1      = "One Microsoft Way"
    address_type          = "Residential"
    company_name          = "Microsoft"
    street_address_2      = "Four Microsoft Way"
    street_address_3      = "Five Microsoft Way"
    postal_code_plus_four = "6399"
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

func (DataBoxJobResource) withCustomerManaged(data acceptance.TestData) string {
	template := testAccDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "Terraform Test"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "Redmond"
    country           = "US"
    postal_code       = "98052"
    state_or_province = "WA"
    street_address_1  = "One Microsoft Way"
  }

  preferred_shipment_type      = "CustomerManaged"
  delivery_scheduled_date_time = "2020-04-01T05:30:00+05:30"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}

func testAccDataBoxJob_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
