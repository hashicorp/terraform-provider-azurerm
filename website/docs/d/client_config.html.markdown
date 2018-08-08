---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_client_config"
sidebar_current: "docs-azurerm-datasource-client-config"
description: |-
  Get information about the configuration of the azurerm provider.
---

# Data Source: azurerm_client_config

Use this data source to access the configuration of the Azure Resource Manager
provider.

## Example Usage

```hcl
data "azurerm_client_config" "example" {}

output "account_id" {
  value = "${data.azurerm_client_config.example.service_principal_application_id}"
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `client_id` is set to the Azure Client ID (Application Object ID).
* `tenant_id` is set to the Azure Tenant ID.
* `subscription_id` is set to the Azure Subscription ID.

---

~> **Note:** the following fields are only available when authenticating via a Service Principal (as opposed to using the Azure CLI):

* `service_principal_application_id` is the Service Principal Application ID.
* `service_principal_object_id` is the Service Principal Object ID.

~> **Note:** To better understand "application" and "service principal", please read
[Application and service principal objects in Azure Active Directory](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-application-objects).
