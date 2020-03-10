---
subcategory: "DataBox"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_job"
description: |-
  Get information about an existing DataBox.
---

# Data Source: azurerm_databox_job

Use this data source to access information about an existing DataBox.

## Example Usage

```hcl
data "azurerm_databox_job" "existing" {
  name                = "example-databoxjob"
  resource_group_name = "example-resources"
}

output "id" {
  value = azurerm_databox_job.existing.id
}
```

## Argument Reference

* `name` - Specifies the name of the DataBox.

* `resource_group_name` - Specifies the name of the Resource Group where this DataBox exists.

## Attributes Reference

* `location` - The Azure location where the resource exists.

* `contact_details` - One or more `contact_details` block defined below.

* `destination_account` - One or more `destination_account` block defined below.

* `preferred_shipment_type` - The shipment logistics type that the customer preferred.

* `shipping_address` - One or more `destination_account` block defined below.

* `sku_name` - The sku name.

* `databox_disk_passkey` - The user entered passkey for DataBox Disk job.

* `databox_preferred_disk_count` - The disk count.

* `databox_preferred_disk_size_in_tb` - The disk size in tb.

* `datacenter_region_preference` - The preferred Data Center Region.

* `delivery_scheduled_date_time` - The delivery scheduled date time.

* `delivery_type` - The delivery type of Job.

* `device_password` - The device password for unlocking DataBox Heavy.

* `expected_data_size_in_tb` - The expected size of the data, which needs to be transferred in this job, in tb.

* `tags` - A mapping of tags to assign to the resource.

---

An `contact_details` block supports the following:

* `name` - The contact name of the person.

* `emails` - The list of Email-ids to be notified about job progress.

* `phone_number` - The phone number of the contact person?

* `mobile` - The mobile number of the contact person

* `notification_preference` - One or more `notification_preference` block defined below.

* `phone_extension` - The phone extension number of the contact person.

---

An `notification_preference` block supports the following:

* `at_azure_dc` - (Optional) Is the at_azure_dc allowed?

* `data_copied` - (Optional) Is the data_copied allowed?

* `delivered` - (Optional) Is the delivered allowed?

* `device_prepared` - (Optional) Is the device_prepared allowed?

* `dispatched` - (Optional) Is the dispatched allowed?

* `picked_up` - (Optional) Is the picked_up allowed?

---

An `destination_account` block supports the following:

* `type` -The destination account type.

* `resource_group_id` - The destination Resource Group Id where the Compute disks should be created.

* `share_password` - The share password to be shared by all shares in SA.

* `staging_storage_account_id` - The arm id of the storage account that can be used to copy the vhd for staging.

* `storage_account_id` - The arm id of the destination where the data has to be moved.

---

An `shipping_address` block supports the following:

* `city` - The name of the city.

* `country` - The name of the country.

* `postal_code` - The portal code.

* `state_or_province` - The name of the state or province.

* `street_address_1` - The Street Address line 1.

* `address_type` - The type of address.

* `company_name` - The name of the company.

* `street_address_2` - The Street Address line 2.

* `street_address_3` - The Street Address line 3.

* `postal_code_ext` - The extended portal code.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataBox.
