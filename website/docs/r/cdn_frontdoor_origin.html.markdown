---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin"
description: |-
  Manages a `CDN FrontDoor Origin.`
---

# azurerm_cdn_frontdoor_origin

Manages a CDN FrontDoor Origin.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-origingroup"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  load_balancing {}
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                          = "example-origin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id

  health_probes_enabled          = true
  certificate_name_check_enabled = false

  host_name          = "contoso.com"
  http_port          = 80
  https_port         = 443
  origin_host_header = "www.contoso.com"
  priority           = 1
  weight             = 1
}
```

## Example Usage With Private Link

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Premium"
  account_replication_type = "LRS"

  allow_nested_items_to_be_public = false

  network_rules {
    default_action = "Deny"
  }

  tags = {
    environment = "Example"
  }
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-origin-group"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  load_balancing {}
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                          = "example-origin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id

  health_probes_enabled          = true
  certificate_name_check_enabled = true
  host_name                      = azurerm_storage_account.example.primary_blob_host
  origin_host_header             = azurerm_storage_account.example.primary_blob_host
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for Private Link Origin CDN Frontdoor"
    target_type            = "blob"
    location               = azurerm_storage_account.example.location
    private_link_target_id = azurerm_storage_account.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this CDN FrontDoor Origin. Changing this forces a new CDN FrontDoor Origin to be created.

* `cdn_frontdoor_origin_group_id` - (Required) The ID of the CDN FrontDoor Origin Group within which this CDN FrontDoor Origin should exist. Changing this forces a new CDN FrontDoor Origin to be created.

* `host_name` - (Required) The IPv4 address, IPv6 address or Domain name of the Origin.

!> **IMPORTANT:** This must be unique across all CDN FrontDoor Origins within a CDN FrontDoor Endpoint.

* `certificate_name_check_enabled` - (Required) Specifies whether certificate name checks are enabled for this origin.

* `health_probes_enabled` - (Optional) Should the health probes be enabled against the origins defined within the origin group? Possible values are `true` or `false`. Defaults to `true`. 

-> **NOTE:** Health probes can only be disabled if there is a single enabled origin in single enabled origin group

* `http_port` - (Optional) The value of the HTTP port. Must be between `1` and `65535`. Defaults to `80`.

* `https_port` - (Optional) The value of the HTTPS port. Must be between `1` and `65535`. Defaults to `443`.

* `origin_host_header` - (Optional) The host header value (an IPv4 address, IPv6 address or Domain name) which is sent to the origin with each request. If unspecified the hostname from the request will be used. 

-> Azure FrontDoor Origins, such as Web Apps, Blob Storage, and Cloud Services require this host header value to match the origin's hostname. This field's value overrides the host header defined in the FrontDoor Endpoint. For more information on how to properly set the origin host header value please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/origin?pivots=front-door-standard-premium#origin-host-header).

* `priority` - (Optional) Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy. Must be between `1` and `5` (inclusive). Defaults to `1`.

* `private_link` - (Optional) A `private_link` block as defined below.

-> **NOTE:** Private Link requires that the CDN FrontDoor Profile this Origin is hosted within is using the SKU `Premium_AzureFrontDoor` and that the `certificate_name_check_enabled` field is set to `true`.

* `weight` - (Optional) The weight of the origin in a given origin group for load balancing. Must be between `1` and `1000`. Defaults to `500`.

---

A `private_link` block supports the following:

~> **NOTE:** At this time the Private Link Endpoint **must be approved manually** - for more information and region availability please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/private-link).

!> **IMPORTANT:** Origin support for direct private end point connectivity is limited to `Storage (Azure Blobs)`, `App Services` and `internal load balancers`. The Azure Front Door Private Link feature is region agnostic but for the best latency, you should always pick an Azure region closest to your origin when choosing to enable Azure Front Door Private Link endpoint.

!> **IMPORTANT:** To associate a Load Balancer with a CDN FrontDoor Origin via Private LInk you must stand up your own `azurerm_private_link_service` - and ensure that a `depends_on` exists on the `azurerm_cdn_frontdoor_origin` resource to ensure it's destroyed before the `azurerm_private_link_service` resource (e.g. `depends_on = [azurerm_private_link_service.example]`) due to the design of the CDN FrontDoor Service. 

* `request_message` - (Optional) Specifies the request message that will be submitted to the `private_link_target_id` when requesting the private link endpoint connection. Values must be between `1` and `140` characters in length. Defaults to `Access request for CDN Frontdoor Private Link Origin`.

* `target_type` - (Optional) Specifies the type of target for this Private Link Endpoint. Possible values are `blob`, `blob_secondary`, `web` and `sites`.

-> **NOTE:** `target_type` cannot be specified when using a Load Balancer as an Origin.

* `location` - (Required) Specifies the location where the Private Link resource should exist.

* `private_link_target_id` - (Required) The ID of the Azure Resource to connect to via the Private Link.

---

## Example HCL Configurations

* [Private Link Origin with Storage Account Blob](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/storage-account-blob)
* [Private Link Origin with Storage Account Static Web Site](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/storage-account-static-site)
* [Private Link Origin with Linux Web Application](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/linux-web-app)
* [Private Link Origin with Internal Load Balancer](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/private-link-service/cdn-frontdoor/load-balancer)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN FrontDoor Origin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN FrontDoor Origin.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Origin.
* `update` - (Defaults to 30 minutes) Used when updating the CDN FrontDoor Origin.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN FrontDoor Origin.

## Import

CDN FrontDoor Origin can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_origin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1
```
