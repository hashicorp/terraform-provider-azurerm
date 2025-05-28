---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin"
description: |-
  Manages a Front Door (standard/premium) Origin.
---

# azurerm_cdn_frontdoor_origin

Manages a Front Door (standard/premium) Origin.

!> **Note:** If you are attempting to implement an Origin that uses its own Private Link Service with a Load Balancer the Profile resource in your configuration file **must** have a `depends_on` meta-argument which references the `azurerm_private_link_service`, see `Example Usage With Private Link Service` below.

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
  enabled                       = true

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
  enabled                       = true

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

## Example Usage With Private Link Service

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  depends_on = [azurerm_private_link_service.example]

  name                = "profile-example"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin" example {
  name                           = "origin-example"
  cdn_frontdoor_origin_group_id  = azurerm_cdn_frontdoor_origin_group.example.id
  enabled                        = true
  host_name                      = "example.com"
  origin_host_header             = "example.com"
  priority                       = 1
  weight                         = 1000
  certificate_name_check_enabled = false

  private_link {
    request_message        = "Request access for Private Link Origin CDN Frontdoor"
    location               = azurerm_resource_group.example.location
    private_link_target_id = azurerm_private_link_service.example.id
  }
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "group-example"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "vn-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                                          = "sn-example"
  resource_group_name                           = azurerm_resource_group.example.name
  virtual_network_name                          = azurerm_virtual_network.example.name
  address_prefixes                              = ["10.5.1.0/24"]
  private_link_service_network_policies_enabled = false
}

resource "azurerm_public_ip" "example" {
  name                = "ip-example"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "lb-example"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                = "pls-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  visibility_subscription_ids                 = [data.azurerm_client_config.current.subscription_id]
  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.example.frontend_ip_configuration[0].id]

  nat_ip_configuration {
    name                       = "primary"
    private_ip_address         = "10.5.1.17"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = true
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Origin. Changing this forces a new Front Door Origin to be created.

* `cdn_frontdoor_origin_group_id` - (Required) The ID of the Front Door Origin Group within which this Front Door Origin should exist. Changing this forces a new Front Door Origin to be created.

* `host_name` - (Required) The IPv4 address, IPv6 address or Domain name of the Origin.

!> **Note:** This must be unique across all Front Door Origins within a Front Door Endpoint.

* `certificate_name_check_enabled` - (Required) Specifies whether certificate name checks are enabled for this origin.

* `enabled` - (Optional) Should the origin be enabled? Possible values are `true` or `false`. Defaults to `true`.

* `http_port` - (Optional) The value of the HTTP port. Must be between `1` and `65535`. Defaults to `80`.

* `https_port` - (Optional) The value of the HTTPS port. Must be between `1` and `65535`. Defaults to `443`.

* `origin_host_header` - (Optional) The host header value (an IPv4 address, IPv6 address or Domain name) which is sent to the origin with each request. If unspecified the hostname from the request will be used.

-> **Note:** Azure Front Door Origins, such as Web Apps, Blob Storage, and Cloud Services require this host header value to match the origin's hostname. This field's value overrides the host header defined in the Front Door Endpoint. For more information on how to properly set the origin host header value please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/origin?pivots=front-door-standard-premium#origin-host-header).

* `priority` - (Optional) Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy. Must be between `1` and `5` (inclusive). Defaults to `1`.

* `private_link` - (Optional) A `private_link` block as defined below.

-> **Note:** Private Link requires that the Front Door Profile this Origin is hosted within is using the SKU `Premium_AzureFrontDoor` and that the `certificate_name_check_enabled` field is set to `true`.

* `weight` - (Optional) The weight of the origin in a given origin group for load balancing. Must be between `1` and `1000`. Defaults to `500`.

---

A `private_link` block supports the following:

~> **Note:** At this time the Private Link Endpoint **must be approved manually** - for more information and region availability please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/private-link).

!> **Note:** Origin support for direct private endpoint connectivity is limited to `Storage (Azure Blobs)`, `Storage (Static Web Sites)`, `App Services`, `internal load balancers`, `Azure Container Apps (preview)` and `Azure API Management`. The Azure Front Door Private Link feature is region agnostic but for the best latency, you should always pick an Azure region closest to your origin when choosing to enable Azure Front Door Private Link endpoint.

!> **Note:** To associate a Load Balancer with a Front Door Origin via Private Link you must stand up your own `azurerm_private_link_service` - and ensure that a `depends_on` exists on the `azurerm_cdn_frontdoor_origin` resource to ensure it's destroyed before the `azurerm_private_link_service` resource (e.g. `depends_on = [azurerm_private_link_service.example]`) due to the design of the Front Door Service.

* `request_message` - (Optional) Specifies the request message that will be submitted to the `private_link_target_id` when requesting the private link endpoint connection. Values must be between `1` and `140` characters in length. Defaults to `Access request for CDN FrontDoor Private Link Origin`.

* `target_type` - (Optional) Specifies the type of target for this Private Link Endpoint. Possible values are `blob`, `blob_secondary`, `Gateway`, `managedEnvironments`, `sites`, `web` and `web_secondary`.

-> **Note:** `target_type` cannot be specified when using a Load Balancer as an Origin.

* `location` - (Required) Specifies the location where the Private Link resource should exist. Changing this forces a new resource to be created.

* `private_link_target_id` - (Required) The ID of the Azure Resource to connect to via the Private Link.

-> **Note:** the `private_link_target_id` property must specify the Resource ID of the Private Link Service when using Load Balancer as an Origin.

---

## Example HCL Configurations

* [Private Link Origin with Storage Account Blob](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/storage-account-blob)
* [Private Link Origin with Storage Account Static Web Site](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/storage-account-static-site)
* [Private Link Origin with Linux Web Application](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/cdn-frontdoor/linux-web-app)
* [Private Link Origin with Internal Load Balancer](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/private-link-service/cdn-frontdoor/load-balancer)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Origin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Origin.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Origin.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Origin.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Origin.

## Import

Front Door Origins can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_origin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1
```
