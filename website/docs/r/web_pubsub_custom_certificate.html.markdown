---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_custom_certificate"
description: |-
  Manages an Azure Web PubSub Custom Certificate.
---

# azurerm_web_pubsub_custom_certificate

Manages an Azure Web PubSub Custom Certificate.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_web_pubsub_service" "example" {
  name                = "example-webpubsub"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Premium_P1"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Get",
      "List",
    ]

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_web_pubsub_service.test.identity[0].principal_id

    certificate_permissions = [
      "Create",
      "Get",
      "List",
    ]

    secret_permissions = [
      "Get",
      "List",
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "imported-cert"
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("certificate-to-import.pfx")
    password = ""
  }
}

resource "azurerm_web_pubsub_custom_certificate" "test" {
  name                  = "example-cert"
  web_pubsub_id         = azurerm_web_pubsub_service.example.id
  custom_certificate_id = azurerm_key_vault_certificate.example.id

  depends_on = [azurerm_key_vault_access_policy.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web PubSub Custom Certificate. Changing this forces a new resource to be created.

* `web_pubsub_id` - (Required) The Web PubSub ID of the Web PubSub Custom Certificate. Changing this forces a new resource to be created.

-> **Note:** custom certificate is only available for Web PubSub Premium tier. Please enable managed identity in the corresponding Web PubSub Service and give the managed identity access to the key vault, the required permission is Get Certificate and Secret.

* `custom_certificate_id` - (Required) The certificate ID of the Web PubSub Custom Certificate. Changing this forces a new resource to be created.

-> **Note:** Self assigned certificate is not supported and the provisioning status will fail.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web PubSub Custom Certificate.

* `certificate_version` - The certificate version of the Web PubSub Custom Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Custom Certificate of the Web PubSub service
* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Certificate of the Web PubSub service
* `delete` - (Defaults to 30 minutes) Used when deleting the Custom Certificate of the Web PubSub service

## Import

Custom Certificate for a Web PubSub service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_custom_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubSub/WebPubsub1/customCertificates/cert1
```
