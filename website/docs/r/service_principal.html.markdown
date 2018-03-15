---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_principal"
sidebar_current: "docs-azurerm-resource-authorization-service-principal"
description: |-
  Manage a Service Principal.

---

# azurerm_service_principal

Create a new Service Principal.

## Example Usage

```hcl
resource "azurerm_ad_application" "example" {
  display_name = "example"
}

resource "azurerm_service_principal" "example" {
  app_id = "${azurerm_ad_application.example.app_id}"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the Azure Active Directory Application.

## Attributes Reference

The following attributes are exported:

* `object_id` - The Service Principal Object ID.

* `display_name` - The Service Principal name.

* `service_principal_names` - A list of names used to identify the application.

## Import

Service Principals can be imported using the `object id`, e.g.

```shell
terraform import azurerm_service_principal.test 00000000-0000-0000-0000-000000000000
```
