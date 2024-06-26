---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_load_test"
description: |-
  Gets information about an existing Load Test Service.
---

# Data Source: azurerm_load_test

Use this data source to access information about a Load Test Service.

## Example Usage

```hcl
data "azurerm_load_test" "example" {
  resource_group_name = "example-resources"
  name                = "example-load-test"
}

output "load_test_id" {
  value = data.azurerm_load_test.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Load Test Service.

* `resource_group_name` - The name of the Resource Group in which the Load Test Service exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Load Test Service.

* `data_plane_uri` - Resource data plane URI.

* `description` - Description of the resource.

* `encryption` - An `encryption` block as defined below.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Load Test exists.

* `tags` - A mapping of tags assigned to the Load Test Service.

---

A `identity` block exports the following:

* `identity_ids` - The list of the User Assigned Identity IDs that is assigned to this Load Test Service.

* `principal_id` - The Principal ID for the System-Assigned Managed Identity assigned to this Load Test Service.

* `tenant_id` - The Tenant ID for the System-Assigned Managed Identity assigned to this Load Test Service.

* `type` - Type of Managed Service Identity.

---

A `encryption` block exports the following:

* `identity` - An `identity` block as defined below.

* `key_url` - The URI specifying the Key vault and key to be used to encrypt data in this resource.

---

A `identity` block for `encryption` exports the following:

* `identity_id` - The User Assigned Identity ID that is assigned to this Load Test Encryption.

* `type` - Type of Managed Service Identity that is assigned to this Load Test Encryption.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Load Test Service.
