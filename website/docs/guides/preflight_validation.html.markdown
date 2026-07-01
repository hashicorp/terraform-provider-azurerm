---
layout: "azurerm"
page_title: "Azure Provider: Preflight Validation"
description: |-
  This guide will cover how to use the Preflight Validation feature
---

# Azure Provider: Preflight Validation

The Azure Provider supports an opt-in plan-time validation feature, known as Preflight Validation, which submits the resource payload that would be sent during `Create` and `Update` to the Azure Deployments Validation API during `terraform plan`. Where supported, this surfaces a number of configuration errors (including those raised by Azure Policy assignments) at plan-time, rather than at apply-time.

Preflight Validation is configured within the `enhanced_validation` block, inside the `features` block of the Provider configuration:

```hcl
provider "azurerm" {
  features {
    enhanced_validation {
      preflight_enabled = true
    }
  }
}
```

Preflight Validation can also be enabled by setting the `ARM_PROVIDER_ENHANCED_VALIDATION_PREFLIGHT_ENABLED` environment variable to `true`.

~> **Note:** Preflight Validation does not verify whether dependent resources exist, whether sufficient quota is available, or whether a globally unique name is still available. These checks continue to occur at apply-time.

## Supported Resources

The following resources currently support Preflight Validation:

* `azurerm_eventgrid_namespace`
* `azurerm_managed_redis`
* `azurerm_service_plan`
* `azurerm_dashboard_grafana`
* `azurerm_nginx_deployment`