---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_service_versions"
sidebar_current: "docs-azurerm-datasource-kubernetes-service-versions"
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
  value = "${data.azurerm_kubernetes_service_versions.current.versions}"
}

output "latest_version" {
  value = "${data.azurerm_kubernetes_service_versions.current.latest_version}"
}
```

## Argument Reference

* `location` - (Required) Specifies the location in which to query for versions.

* `version_prefix` - (Optional) A prefix filter for the versions of Kubernetes which should be returned; for example `1.` will return `1.9` to `1.14`, whereas `1.12` will return `1.12.2`.

## Attributes Reference

* `versions` - The list of all supported versions.

* `latest_version` - The most recent version available.
