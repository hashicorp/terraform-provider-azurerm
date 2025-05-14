---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_custom_domain"
description: |-
  Manages an Azure Web PubSub Custom Domain.
---

# azurerm_web_pubsub_custom_domain

Manages an Azure Web PubSub Custom Domain.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_web_pubsub" "example" {
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
    object_id = azurerm_web_pubsub.test.identity[0].principal_id
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
  web_pubsub_id         = azurerm_web_pubsub.example.id
  custom_certificate_id = azurerm_key_vault_certificate.example.id

  depends_on = [azurerm_key_vault_access_policy.example]
}

resource "azurerm_web_pubsub_custom_domain" "test" {
  name                             = "example-domain"
  domain_name                      = "tftest.com"
  web_pubsub_id                    = azurerm_web_pubsub.test.id
  web_pubsub_custom_certificate_id = azurerm_web_pubsub_custom_certificate.test.id
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Web PubSub Custom Domain. Changing this forces a new resource to be created.

* `domain_name` - (Required) Specifies the custom domain name of the Web PubSub Custom Domain. Changing this forces a new resource to be created.

-> **Note:** Please ensure the custom domain name is included in the Subject Alternative Names of the selected Web PubSub Custom Certificate.

* `web_pubsub_id` - (Required) Specifies the Web PubSub ID of the Web PubSub Custom Domain. Changing this forces a new resource to be created.

* `web_pubsub_custom_certificate_id` - (Required) Specifies the Web PubSub Custom Certificate ID of the Web PubSub Custom Domain. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web PubSub Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the custom domain of the Web PubSub service
* `read` - (Defaults to 5 minutes) Used when retrieving the custom domain of the Web PubSub service
* `delete` - (Defaults to 30 minutes) Used when deleting the custom domain of the Web PubSub service

## Import

Custom Domain for a Web PubSub service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/webPubSub/webpubsub1/customDomains/customDomain1
```
