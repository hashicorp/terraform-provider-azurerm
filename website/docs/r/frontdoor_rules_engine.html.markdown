---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_rules_engine"
description: |-
  Manages an Azure Front Door Rules Engine configuration and rules.
---

# azurerm_frontdoor_rules_engine

Manages an Azure Front Door Rules Engine configuration and rules.

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

* `rule` - (Required) A `rule` block as defined below.

---

The `rule` block supports the following:

* `name` - (Required) The name of the rule.

* `priority` - (Required) Priority of the rule, must be unique per rules engine definition.

* `action` - (Required) A `rule_action` block as defined below.

* `match_condition` - One or more `match_condition` block as defined below.

---

The `action` block supports the following:

* `request_header` - A `request_header` block as defined below.

* `response_header` - A `response_header` block as defined below.

---

The `request_header` block supports the following:

* `header_action_type` can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` header name (string).

* `value` value name (string).

---

The `response_header` block supports the following:

* `header_action_type` can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` header name (string).

* `value` value name (string).

---

The `match_condition` block supports the following:

* `variable` can be set to `IsMobile`, `RemoteAddr`, `RequestMethod`, `QueryString`, `PostArgs`, `RequestURI`, `RequestPath`, `RequestFilename`, `RequestFilenameExtension`,`RequestHeader`,`RequestBody` or `RequestScheme`.

* `selector` match against a specific key when `variable` is set to `PostArgs` or `RequestHeader`. It cannot be used with `QueryString` and `RequestMethod`. Defaults to `null`.

* `operator` can be set to `Any`, `IPMatch`, `GeoMatch`, `Equal`, `Contains`, `LessThan`, `GreaterThan`, `LessThanOrEqual`, `GreaterThanOrEqual`, `BeginsWith` or `EndsWith`

* `transform` can be set to one or more values out of `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` and `UrlEncode`

* `negate_condition` can be set to `true` or `false` to negate the given condition. Defaults to `true`.

* `value` (array) can contain one or more strings.

## Import

Azure Front Door Rules Engine's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_rules_engine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/frontdoors/frontdoor1/rulesengines/rule1
```
