---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_principal_attribute_set"
description: |-
  Manages a custom security attribute set assignment on an existing Microsoft Entra service principal.

---

# azurerm_service_principal_attribute_set

Manages a custom security attribute set assignment on an existing Microsoft Entra service principal.

~> **Note:** This resource manages **one** attribute set per resource instance and is designed to avoid overwriting other attribute sets on the same service principal.

## Example Usage (minimal inputs)

```hcl
resource "azurerm_service_principal_attribute_set" "example" {
  service_principal_object_id = var.service_principal_object_id
  attribute_set_name          = var.attribute_set_name

  attributes = {
    Environment = "Production"
    Team        = "Platform"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `service_principal_object_id` - (Required) The Object ID of the target Microsoft Entra service principal. Changing this forces a new resource to be created.

* `attribute_set_name` - (Required) The custom security attribute set name to manage on the service principal. Changing this forces a new resource to be created.

* `attributes` - (Required) A map of attribute names and string values to apply for this attribute set.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The Service Principal Attribute Set resource ID.

## Import

Service principal attribute sets can be imported using the resource ID in the format:

```shell
terraform import azurerm_service_principal_attribute_set.example "<service_principal_object_id>|<attribute_set_name>"
```

## Permissions

The caller requires Microsoft Graph permissions that allow reading and updating service principals and custom security attributes in the tenant.
