---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_orchestrator_versions"
sidebar_current: "docs-azurerm-datasource-container-orchestrator-versions-x"
description: |-
  Gets available versions of container orchestrators.
---

# Data Source: azurerm_container_orchestrator_versions

Use this data source to access information about available versions of orchestrators supported by the Azure Container services.

## Example Usage

```hcl
data "azurerm_container_orchestrator_versions" "current" {
  location = "westeurope"
}

output "versions" {
  value = "${data.azurerm_container_orchestrator_versions.current.versions}"
}

output "latest_version" {
  value = "${data.azurerm_container_orchestrator_versions.current.latest_version}"
}
```

## Argument Reference

* `location` - (Required) Specifies the location in which to query for versions.
* `orchestrator_type` - (Optional) The type of orchestrator for which to fetch versions. Right now the only supported value is `Kubernetes`. If omitted, defaults to `Kubernetes`.
* `version_prefix` - (Optional) A string prefix used to filter the versions retrieved by the query. When this argument is specified, only the versions which match this prefix are returned.

## Attributes Reference

* `versions` - The list of all supported versions.
* `latest_version` - The most recent version available.
