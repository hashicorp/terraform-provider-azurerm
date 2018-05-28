---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via the Azure CLI"
sidebar_current: "docs-azurerm-index-authentication-azure-cli"
description: |-
  This guide will cover how to use the Azure CLI to provide authentication for the Azure Provider.

---

# Azure Provider: Authenticating using the Azure CLI

Terraform supports authenticating to Azure through a Service Principal or the Azure CLI.

We recommend [using a Service Principal when running in a shared environment](authenticating_via_service_principal.html) (such as within a CI server/automation) - and authenticating via the Azure CLI when you're running Terraform locally.

~> **NOTE:** Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (e.g. via `az login --service-principal`) you should instead [authenticate via the Service Principal directly](authenticating_via_service_principal.html).

When authenticating via the Azure CLI, Terraform will automatically connect to the Default Subscription - this can be changed by using the Azure CLI - and is documented below.

## Configuring the Azure CLI

~> **Note:** There are multiple versions of the Azure CLI - the latest version is known as [the Azure CLI 2.0 (Python)](https://github.com/Azure/azure-cli) and [the older Azure CLI (Node.JS)](https://github.com/Azure/azure-xplat-cli). While Terraform currently supports both - we highly recommend users upgrade to the Azure CLI 2.0 (Python) if possible.

This guide assumes that you have [the Azure CLI 2.0 (Python)](https://github.com/Azure/azure-cli) installed.

~> **Note:** If you're using the **China**, **German** or **Government** Azure Clouds - you'll need to first configure the Azure CLI to work with that Cloud.  You can do this by running:

```shell
$ az cloud set --name AzureChinaCloud|AzureGermanCloud|AzureUSGovernment
```

---

Firstly, login to the Azure CLI using:

```shell
$ az login
```

~> **NOTE:** Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (e.g. via `az login --service-principal`) you should instead [authenticate via the Service Principal directly](authenticating_via_service_principal.html).

This will prompt you to open a web browser, as shown below:

```shell
To sign in, use a web browser to open the page https://aka.ms/devicelogin and enter the code XXXXXXXX to authenticate.
```

Once logged in - it's possible to list the Subscriptions associated with the account via:

```shell
$ az account list
```

The output (similar to below) will display one or more Subscriptions - with the `id` field being the Subscription ID.

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

~> **Note:** When authenticating via the Azure CLI, Terraform will automatically connect to the Default Subscription. As such if you have multiple subscriptions on the account, you may need to set the Default Subscription, via:

```shell
$ az account set --subscription="SUBSCRIPTION_ID"
```

Also, if you have been authenticating with a service principal and you switch to Azure CLI, you must null out the ARM_* environment variables. Failure to do so causes errors to be thrown.
