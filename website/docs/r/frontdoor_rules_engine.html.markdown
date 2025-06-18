---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_rules_engine"
description: |-
  Manages an Azure Front Door (classic) Rules Engine configuration and rules.
---

# azurerm_frontdoor_rules_engine

!> **Note:** This deploys an Azure Front Door (classic) resource which has been deprecated and will receive security updates only. Please migrate your existing Azure Front Door (classic) deployments to the new [Azure Front Door (standard/premium) resources](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cdn_frontdoor_custom_domain). For your convenience, the service team has exposed a `Front Door Classic` to `Front Door Standard/Premium` [migration tool](https://learn.microsoft.com/azure/frontdoor/tier-migration) to allow you to migrate your existing `Front Door Classic` instances to the new `Front Door Standard/Premium` product tiers.

!> **Note:** On `1 April 2025`, Azure Front Door (classic) will be retired for the public cloud, existing Azure Front Door (classic) resources must be migrated out of Azure Front Door (classic) to Azure Front Door Standard/Premium before `1 October 2025` to avoid potential disruptions in service.

Manages an Azure Front Door (classic) Rules Engine configuration and rules.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_frontdoor" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name

  backend_pool {
    name                = "exampleBackendBing"
    load_balancing_name = "exampleLoadBalancingSettings1"
    health_probe_name   = "exampleHealthProbeSetting1"

    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }
  }

  backend_pool_health_probe {
    name = "exampleHealthProbeSetting1"
  }

  backend_pool_load_balancing {
    name = "exampleLoadBalancingSettings1"
  }

  frontend_endpoint {
    name      = "exampleFrontendEndpoint1"
    host_name = "example-FrontDoor.azurefd.net"
  }

  routing_rule {
    name               = "exampleRoutingRule1"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = ["exampleFrontendEndpoint1"]
  }
}

resource "azurerm_frontdoor_rules_engine" "example_rules_engine" {
  name                = "exampleRulesEngineConfig1"
  frontdoor_name      = azurerm_frontdoor.example.name
  resource_group_name = azurerm_frontdoor.example.resource_group_name

  rule {
    name     = "debuggingoutput"
    priority = 1

    action {
      response_header {
        header_action_type = "Append"
        header_name        = "X-TEST-HEADER"
        value              = "Append Header Rule"
      }
    }
  }

  rule {
    name     = "overwriteorigin"
    priority = 2

    match_condition {
      variable = "RequestMethod"
      operator = "Equal"
      value    = ["GET", "POST"]
    }

    action {

      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Origin"
        value              = "*"
      }

      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Credentials"
        value              = "true"
      }

    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Rules engine configuration. Changing this forces a new resource to be created. 

* `frontdoor_name` - (Required) The name of the Front Door instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `enabled` - (Optional) Whether this Rules engine configuration is enabled? Defaults to `true`.

* `rule` - (Optional) A `rule` block as defined below.

---

The `rule` block supports the following:

* `name` - (Required) The name of the rule.

* `priority` - (Required) Priority of the rule, must be unique per rules engine definition.

* `action` - (Optional) An `action` block as defined below.

* `match_condition` - (Optional) One or more `match_condition` block as defined below.

---

The `action` block supports the following:

* `request_header` - (Optional) A `request_header` block as defined below.

* `response_header` - (Optional) A `response_header` block as defined below.

---

The `request_header` block supports the following:

* `header_action_type` - (Optional) can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` - (Optional) header name (string).

* `value` - (Optional) value name (string).

---

The `response_header` block supports the following:

* `header_action_type` - (Optional) can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` - (Optional) header name (string).

* `value` - (Optional) value name (string).

---

The `match_condition` block supports the following:

* `variable` - (Optional) can be set to `IsMobile`, `RemoteAddr`, `RequestMethod`, `QueryString`, `PostArgs`, `RequestURI`, `RequestPath`, `RequestFilename`, `RequestFilenameExtension`,`RequestHeader`,`RequestBody` or `RequestScheme`.

* `selector` - (Optional) match against a specific key when `variable` is set to `PostArgs` or `RequestHeader`. It cannot be used with `QueryString` and `RequestMethod`.

* `operator` - (Required) can be set to `Any`, `IPMatch`, `GeoMatch`, `Equal`, `Contains`, `LessThan`, `GreaterThan`, `LessThanOrEqual`, `GreaterThanOrEqual`, `BeginsWith` or `EndsWith`

* `transform` - (Optional) can be set to one or more values out of `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` and `UrlEncode`

* `negate_condition` - (Optional) can be set to `true` or `false` to negate the given condition. Defaults to `false`.

* `value` - (Optional) (array) can contain one or more strings.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Frontdoor Rules Engine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Rules Engine.
* `update` - (Defaults to 6 hours) Used when updating the Frontdoor Rules Engine.
* `delete` - (Defaults to 6 hours) Used when deleting the Frontdoor Rules Engine.

## Import

Azure Front Door Rules Engine's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_rules_engine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/frontdoors/frontdoor1/rulesEngines/rule1
```
