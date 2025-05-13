---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain_association"
description: |-
  Manages the association between a Front Door (standard/premium) Custom Domain and one or more Front Door (standard/premium) Routes.
---

# azurerm_cdn_frontdoor_custom_domain_association

Manages the association between a Front Door (standard/premium) Custom Domain and one or more Front Door (standard/premium) Routes.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "domain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-origin-group"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  session_affinity_enabled = true

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "HEAD"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
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

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name                     = "ExampleRuleSet"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_route" "example" {
  name                          = "example-route"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.example.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.example.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_rule_set.example.id]
  enabled                       = true

  forwarding_protocol    = "HttpsOnly"
  https_redirect_enabled = true
  patterns_to_match      = ["/*"]
  supported_protocols    = ["Http", "Https"]

  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.example.id]
  link_to_default_domain          = false
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example-customDomain"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  dns_zone_id              = azurerm_dns_zone.example.id
  host_name                = join(".", ["contoso", azurerm_dns_zone.example.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_custom_domain_association" "example" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.example.id
  cdn_frontdoor_route_ids        = [azurerm_cdn_frontdoor_route.example.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_custom_domain_id` - (Required) The ID of the Front Door Custom Domain that should be managed by the association resource. Changing this forces a new association resource to be created.

* `cdn_frontdoor_route_ids` - (Required) One or more IDs of the Front Door Route to which the Front Door Custom Domain is associated with.

-> **Note:** This should include all of the Front Door Route resources that the Front Door Custom Domain is associated with. If the list of Front Door Routes is not complete you will receive the service side error `This resource is still associated with a route. Please delete the association with the route first before deleting this resource` when you attempt to `destroy`/`delete` your Front Door Custom Domain.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Custom Domain Association.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Custom Domain Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Custom Domain Association.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Custom Domain Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Custom Domain Association.

## Import

Front Door Custom Domain Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_custom_domain_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/associations/assoc1
```
