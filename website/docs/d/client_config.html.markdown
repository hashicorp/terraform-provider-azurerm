---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_client_config"
sidebar_current: "docs-azurerm-datasource-client-config"
description: |-
  Gets information about the configuration of the azurerm provider.
---

# Data Source: azurerm_client_config

Use this data source to access the configuration of the AzureRM provider.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

output "account_id" {
  value = "${data.azurerm_client_config.current.service_principal_application_id}"
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `client_id` is set to the Azure Client ID (Application Object ID).
* `tenant_id` is set to the Azure Tenant ID.
* `subscription_id` is set to the Azure Subscription ID.
* `object_id` is set to the Azure Object ID.

---

~> **Note:** the following fields are only available when authenticating via a Service Principal (as opposed to using the Azure CLI) and have been deprecated:

* `service_principal_application_id` is the Service Principal Application ID (same as `client_id`).
* `service_principal_object_id` is the Service Principal Object ID (now available via `object_id`).

~> **Note:** To better understand "application" and "service principal", please read
[Application and service principal objects in Azure Active Directory](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-application-objects).
