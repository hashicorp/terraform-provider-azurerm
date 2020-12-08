---
subcategory: "DataBox"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_box_job"
description: |-
  Manages a DataBox.
---

# azurerm_data_box_job

Manages a databox.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                = "example-storageaccount"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_data_box_job" "example" {
  name                = "example-databoxjob"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  contact_details {
    name         = "DataBoxJobTester"
    emails       = ["some.user@example.com"]
    phone_number = "+112345678912"
  }

  destination_account {
    type               = "StorageAccount"
    storage_account_id = azurerm_storage_account.example.id
  }

  preferred_shipment_type = "MicrosoftManaged"

  shipping_address {
    city              = "San Francisco"
    country           = "US"
    postal_code       = "94107"
    state_or_province = "CA"
    street_address_1  = "16 TOWNSEND ST"
  }

  sku_name = "DataBox"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DataBox. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the DataBox should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `contact_details` - (Required) One or more `contact_details` block defined below.

* `preferred_shipment_type` - (Required) Specified shipment logistics type that the customer preferred. Possible values include: `CustomerManaged` and `MicrosoftManaged`. Changing this forces a new resource to be created.

* `shipping_address` - (Required) One or more `shipping_address` block defined below.

* `sku_name` - (Required) Specified the sku name. Possible values include: `DataBox`, `DataBoxDisk` and `DataBoxHeavy`. Changing this forces a new resource to be created.

* `databox_disk_passkey` - (Optional) Specified user entered passkey for DataBox Disk job. Changing this forces a new resource to be created.

* `databox_preferred_disk` - (Required) One or more `databox_preferred_disk` block defined below.

* `datacenter_region_preference` - (Optional) Specified the preferred Data Center Region. Changing this forces a new resource to be created.

* `delivery_scheduled_date_time` - (Optional) Specified the delivery scheduled date time. Changing this forces a new resource to be created.

* `delivery_type` - (Optional) Specified the delivery type of Job. Possible values include: `NonScheduled` and `Scheduled`. Changing this forces a new resource to be created.

* `destination_managed_disk` - (Optional) One or more `destination_managed_disk` block defined below.

* `destination_storage_account` - (Optional) One or more `destination_storage_account` block defined below.

* `device_password` - (Optional) Specified the device password for unlocking DataBox Heavy. Changing this forces a new resource to be created.

* `expected_data_size_in_tb` - (Optional) Specified the expected size of the data, which needs to be transferred in this job, in tb. For `DataBoxDisk`, specifying the expected data in terabytes is mandatory. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `contact_details` block supports the following:

* `name` - (Required) Specified the contact name of the person.

* `emails` - (Required) Specified the list of Email-ids to be notified about job progress.

* `phone_number` - (Required) Specified the phone number of the contact person?

* `phone_mobile` - (Optional) Specified the mobile number of the contact person

* `notification_preference` - (Optional) One or more `notification_preference` block defined below.

* `phone_extension` - (Optional) Specified the phone extension number of the contact person.

---

An `databox_preferred_disk` block supports the following:

* `count` - (Optional) Specified the disk count. Changing this forces a new resource to be created.

* `size_in_tb` - (Optional) Specified the disk size in tb. Changing this forces a new resource to be created.

---

An `notification_preference` block supports the following:

* `at_azure_dc` - (Optional) Is the at_azure_dc allowed?

* `data_copied` - (Optional) Is the data_copied allowed?

* `delivered` - (Optional) Is the delivered allowed?

* `device_prepared` - (Optional) Is the device_prepared allowed?

* `dispatched` - (Optional) Is the dispatched allowed?

* `picked_up` - (Optional) Is the picked_up allowed?

---

An `destination_managed_disk` block supports the following:

* `resource_group_id` - (Required) Specified the destination Resource Group Id where the Compute disks should be created. Changing this forces a new resource to be created.

* `staging_storage_account_id` - (Required) Specified arm id of the storage account that can be used to copy the vhd for staging. Changing this forces a new resource to be created.

* `share_password` - (Optional) Specified the share password to be shared by all shares in SA. Changing this forces a new resource to be created.

---

An `destination_storage_account` block supports the following:

* `storage_account_id` - (Required) Specified arm id of the destination where the data has to be moved. Changing this forces a new resource to be created.

* `share_password` - (Optional) Specified the share password to be shared by all shares in SA. Changing this forces a new resource to be created.

---

An `shipping_address` block supports the following:

* `city` - (Required) Specified the name of the city.

* `country` - (Required) Specified the name of the country.

* `postal_code` - (Required) Specified the portal code.

* `state_or_province` - (Required) Specified the name of the state or province.

* `street_address_1` - (Required) Specified the Street Address line 1.

* `address_type` - (Optional) Specified the type of address. Possible values include: `None`, `Residential`, `Commercial`.

* `company_name` - (Optional) Specified the name of the company.

* `street_address_2` - (Optional) Specified the Street Address line 2.

* `street_address_3` - (Optional) Specified the Street Address line 3.

* `postal_code_plus_four` - (Optional) Specified the extended portal code.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DataBox.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DataBox.
* `update` - (Defaults to 30 minutes) Used when updating the DataBox.
* `read` - (Defaults to 5 minutes) Used when retrieving the DataBox.
* `delete` - (Defaults to 30 minutes) Used when deleting the DataBox.

## Import

DataBox can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_box_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataBox/jobs/job1
```
