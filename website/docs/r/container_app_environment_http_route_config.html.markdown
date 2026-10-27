---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_http_route_config"
description: |-
  Manages a Container App Environment HTTP Route Config.
---

# azurerm_container_app_environment_http_route_config

Manages a Container App Environment HTTP Route Config.

## Example Usage

```hcl
resource "azurerm_container_app_environment" "example" {
  name                = "example-environment"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_container_app_environment_http_route_config" "example" {
  name                         = "exampleroute"
  container_app_environment_id = azurerm_container_app_environment.example.id

  rules {
    description = "example rule"

    targets {
      container_app = "mycontainerapp"
    }

    routes {
      match {
        prefix         = "/api"
        case_sensitive = true
      }

      action {
        prefix_rewrite = "/v1"
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The ID of the Container App Environment. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Container App Environment HTTP Route Config. Changing this forces a new resource to be created.

* `rules` - (Required) One or more `rules` blocks as defined below.

---

* `custom_domains` - (Optional) One or more `custom_domains` blocks as defined below.

---

A `action` block supports the following:

* `prefix_rewrite` - (Optional) The prefix to rewrite the path with. Defaults to no rewrite.

---

A `custom_domains` block supports the following:

* `name` - (Required) The hostname.

* `binding_type` - (Optional) The Custom Domain binding type. Possible values are `Auto`, `Disabled` and `SniEnabled`.

* `certificate_id` - (Optional) The Resource ID of the Certificate to be bound to this hostname. Must exist in the Managed Environment.

---

A `match` block supports the following:

* `case_sensitive` - (Optional) Whether path matching is case sensitive. Defaults to `true`.

* `path` - (Optional) The exact path to match on.

* `path_separated_prefix` - (Optional) The path separated prefix to match on. Matches on all prefixes, not exact.

* `prefix` - (Optional) The prefix to match on. Matches on all prefixes, not exact.

---

A `routes` block supports the following:

* `match` - (Required) A `match` block as defined above.

* `action` - (Optional) A `action` block as defined above.

---

A `rules` block supports the following:

* `targets` - (Required) One or more `targets` blocks as defined below.

* `description` - (Optional) The description of the rule.

* `routes` - (Optional) One or more `routes` blocks as defined above.

---

A `targets` block supports the following:

* `container_app` - (Required) The Container App name to route requests to.

* `label` - (Optional) The label to route requests to.

* `revision` - (Optional) The revision to route requests to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Container App Environment HTTP Route Config.

* `fqdn` - The FQDN of the HTTP Route Config.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment HTTP Route Config.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment HTTP Route Config.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment HTTP Route Config.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment HTTP Route Config.

## Import

Container App Environment HTTP Route Configs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_http_route_config.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment/httpRouteConfigs/myhttproute
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
