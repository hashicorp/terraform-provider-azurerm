---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via the Azure CLI"
sidebar_current: "docs-azurerm-index-authentication-azure-cli"
description: |-
  This guide will cover how to use the Azure CLI to provide authentication for the Azure Provider.

---

# Azure Provider: Authenticating using the Azure CLI

Terraform supports authenticating to Azure using the Azure CLI, a Service Principal or via Managed Service Identity.

We recommend [using a Service Principal when running in a shared environment](authenticating_via_service_principal.html) (such as within a CI server/automation) and authenticating via the Azure CLI when running Terraform locally.

~> **NOTE:** Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via `az login --service-principal`) you should instead [authenticate via the Service Principal directly](authenticating_via_service_principal.html).

When authenticating via the Azure CLI, Terraform will automatically connect to the Default Subscription - this can be changed by using the Azure CLI - and is documented below.

## Configuring the Azure CLI

~> **Note:** There are multiple versions of the Azure CLI - the latest version is known as [the Azure CLI 2.0 (Python)](https://github.com/Azure/azure-cli), which can be used for authentication. Terraform does not support the older [Azure CLI (Node.JS)](https://github.com/Azure/azure-xplat-cli) or [Azure PowerShell](https://github.com/Azure/azure-powershell).

This guide assumes that you have [the Azure CLI 2.0 (Python)](https://github.com/Azure/azure-cli) installed.

~> **Note:** If you're using the **China**, **German** or **US Government** Azure Clouds, you'll need to first configure the Azure CLI to work with that Cloud. You can do this by running:

```shell
$ az cloud set --name AzureChinaCloud|AzureGermanCloud|AzureUSGovernment
```

---

Firstly, login to the Azure CLI using:

```shell
$ az login
```

~> **NOTE:** Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via `az login --service-principal`) you should instead [authenticate via the Service Principal directly](authenticating_via_service_principal.html).

This will prompt you to open a web browser, as shown below:

```shell
To sign in, use a web browser to open the page https://aka.ms/devicelogin and enter the code XXXXXXXX to authenticate.
```

Once logged in, it's possible to list the Subscriptions associated with the account via:

```shell
$ az account list
```

The output (similar to below) will display one or more Subscriptions:

```json
[
  {
    "cloudName": "AzureCloud",
    "id": "00000000-0000-0000-0000-000000000000",
    "isDefault": true,
    "name": "PAYG Subscription",
    "state": "Enabled",
    "tenantId": "00000000-0000-0000-0000-000000000000",
    "user": {
      "name": "user@example.com",
      "type": "user"
    }
  }
]
```

In the snippet above, `id` refers to the Subscription ID and `isDefault` refers to whether this Subscription is configured as the default.

~> **Note:** When authenticating via the Azure CLI, Terraform will automatically connect to the Default Subscription. Therefore, if you have multiple subscriptions on the account, you may need to set the Default Subscription, via:

```shell
$ az account set --subscription="SUBSCRIPTION_ID"
```

If you're previously authenticated using a Service Principal (configured via Environment Variables) - you must remove the `ARM_*` prefixed Environment Variables in order to be able to authenticate using the Azure CLI.
