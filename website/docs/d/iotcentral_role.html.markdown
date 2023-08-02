---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_role"
description: |-
  Gets information about an existing IotCentral Role.
---

# Data Source: azurerm_iotcentral_role

Use this data source to access information about an existing IotCentral Role.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource"
  location = "West Europe"
}

resource "azurerm_iotcentral_application" "example" {
  name                = "example-iotcentral-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sub_domain          = "example-iotcentral-app-subdomain"

  display_name = "example-iotcentral-app-display-name"
  sku          = "ST1"
  template     = "iotc-default@1.0.0"

  tags = {
    Foo = "Bar"
  }
}

data "azurerm_iotcentral_role" "example" {
  sub_domain   = azurerm_iotcentral_application.example.sub_domain
  display_name = "App Administrator"
}
```

## Argument Reference

The following arguments are supported:

* `sub_domain` - (Required) The application `sub_domain`.

* `display_name` - (Required) The `display_name` of the role.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTCentral Role.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the IoTCentral Role.
