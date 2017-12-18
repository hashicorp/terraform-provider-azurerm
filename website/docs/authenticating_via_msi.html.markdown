---
layout: "azurerm"
page_title: "AzureRM: Authenticating via Managed Service Identity"
sidebar_current: "docs-azurerm-index-authentication-msi"
description: |-
  The Azure Resource Manager provider supports authenticating via multiple means. This guide will cover creating a Managed Service Identity which can be used to access Azure Resource Manager.

---

# Authenticating to Azure Resource Manager using Managed Service Identity

Terraform supports authenticating to Azure through Managed Service Identity, Service Principal or the Azure CLI.

We recommend using Managed Service Identity when running in a Shared Environment (such as within a CI server/automation) - and [authenticating via the Azure CLI](authenticating_via_azure_cli.html) when you're running Terraform locally.

## Configuring Managed Service Identity

Managed Service Identity allows an Azure virtual machine to retrieve a token to access the Azure API without needing to pass in credentials. This works by creating a service principal in Azure Active Directory that is associated to a virtual machine. This service principal can then be granted permissions to Azure resources.

### Configuring Managed Service Identity using Terraform

~> **Note**: if you're using the **China**, **German** or **Government** Azure Clouds - you'll need to first configure the Azure CLI to work with that Cloud.  You can do this by running:

```shell
$ az cloud set --name AzureChinaCloud|AzureGermanCloud|AzureUSGovernment
```

---

Firstly, login to the Azure CLI using:

```shell
$ az login
```


Once logged in - it's possible to list the Subscriptions associated with the account via:

```shell
$ az account list
```

The output (similar to below) will display one or more Subscriptions - with the `ID` field being the `subscription_id` field referenced above.

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

Should you have more than one Subscription, you can specify the Subscription to use via the following command:

```shell
$ az account set --subscription="SUBSCRIPTION_ID"
```



### Creating Managed Service Identity in the Azure Portal

