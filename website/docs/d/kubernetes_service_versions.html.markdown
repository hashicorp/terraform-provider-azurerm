---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_service_versions"
description: |-
  Gets the available versions of Kubernetes supported by Azure Kubernetes Service.
---

# Data Source: azurerm_kubernetes_service_versions

Use this data source to retrieve the version of Kubernetes supported by Azure Kubernetes Service.

## Example Usage

```hcl
data "azurerm_kubernetes_service_versions" "current" {
  location = "West Europe"
}

output "versions" {
  value = data.azurerm_kubernetes_service_versions.current.versions
}

output "latest_version" {
  value = data.azurerm_kubernetes_service_versions.current.latest_version
}
```

## Argument Reference

* `location` - Specifies the location in which to query for versions.

* `version_prefix` - (Optional) A prefix filter for the versions of Kubernetes which should be returned; for example `1.` will return `1.9` to `1.14`, whereas `1.12` will return `1.12.2`.

* `include_preview` - (Optional) Should Preview versions of Kubernetes in AKS be included? Defaults to `true`

## Attributes Reference

* `versions` - The list of all supported versions.

* `latest_version` - The most recent version available. If `include_preview == false`, this is the most recent non-preview version available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the versions.
