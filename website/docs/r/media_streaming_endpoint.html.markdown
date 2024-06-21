---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_streaming_endpoint"
description: |-
  Manages a Streaming Endpoint.
---

# azurerm_media_streaming_endpoint

Manages a Streaming Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}

resource "azurerm_media_streaming_endpoint" "example" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  media_services_account_name = azurerm_media_services_account.example.name
  scale_units                 = 2
}
```

## Example Usage with Access Control

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}

resource "azurerm_media_streaming_endpoint" "example" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  media_services_account_name = azurerm_media_services_account.example.name
  scale_units                 = 2
  access_control {
    ip_allow {
      name    = "AllowedIP"
      address = "192.168.1.1"
    }

    ip_allow {
      name    = "AnotherIp"
      address = "192.168.1.2"
    }

    akamai_signature_header_authentication_key {
      identifier = "id1"
      expiration = "2030-12-31T16:00:00Z"
      base64_key = "dGVzdGlkMQ=="
    }

    akamai_signature_header_authentication_key {
      identifier = "id2"
      expiration = "2032-01-28T16:00:00Z"
      base64_key = "dGVzdGlkMQ=="
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Streaming Endpoint should exist. Changing this forces a new Streaming Endpoint to be created.

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Streaming Endpoint to be created.

* `name` - (Required) The name which should be used for this Streaming Endpoint maximum length is `24`. Changing this forces a new Streaming Endpoint to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Streaming Endpoint should exist. Changing this forces a new Streaming Endpoint to be created.

* `scale_units` - (Required) The number of scale units. To create a Standard Streaming Endpoint set `0`. For Premium Streaming Endpoint valid values are between `1` and `10`.

---

* `access_control` - (Optional) A `access_control` block as defined below.

* `auto_start_enabled` - (Optional) The flag indicates if the resource should be automatically started on creation.

* `cdn_enabled` - (Optional) The CDN enabled flag.

* `cdn_profile` - (Optional) The CDN profile name.

* `cdn_provider` - (Optional) The CDN provider name. Supported value are `StandardVerizon`,`PremiumVerizon` and `StandardAkamai`

* `cross_site_access_policy` - (Optional) A `cross_site_access_policy` block as defined below.

* `custom_host_names` - (Optional) The custom host names of the streaming endpoint.

* `description` - (Optional) The streaming endpoint description.

* `max_cache_age_seconds` - (Optional) Max cache age in seconds.

* `tags` - (Optional) A mapping of tags which should be assigned to the Streaming Endpoint.

---

A `access_control` block supports the following:

* `akamai_signature_header_authentication_key` - (Optional) One or more `akamai_signature_header_authentication_key` blocks as defined below.

* `ip_allow` - (Optional) A `ip_allow` block as defined below.

---

A `akamai_signature_header_authentication_key` block supports the following:

* `base64_key` - (Optional) Authentication key.

* `expiration` - (Optional) The expiration time of the authentication key.

* `identifier` - (Optional) Identifier of the key.

---

A `ip_allow` block supports the following:

* `address` - (Optional) The IP address to allow.

* `name` - (Optional) The friendly name for the IP address range.

* `subnet_prefix_length` - (Optional) The subnet mask prefix length (see CIDR notation).

---
A `cross_site_access_policy` block supports the following:

* `client_access_policy` - (Optional) The content of `clientaccesspolicy.xml` used by Silverlight.

* `cross_domain_policy` - (Optional) The content of `crossdomain.xml` used by Silverlight.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Streaming Endpoint.

* `host_name` - The host name of the Streaming Endpoint.

* `sku` - A `sku` block defined as below.

---

A `sku` block supports the following:

* `name` - The sku name of Streaming Endpoint.

* `capacity` - The sku capacity of Streaming Endpoint.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Streaming Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Streaming Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Streaming Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Streaming Endpoint.

## Import

Streaming Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_streaming_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaServices/service1/streamingEndpoints/endpoint1
```
