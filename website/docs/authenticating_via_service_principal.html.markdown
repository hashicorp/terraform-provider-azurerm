---
layout: "azurerm"
page_title: "AzureRM: Authenticating via a Service Principal"
sidebar_current: "docs-azurerm-index-authentication-service-principal"
description: |-
  The Azure Resource Manager provider supports authenticating via multiple means. This guide will cover creating a Service Principal which can be used to access Azure Resource Manager.

---

# Authenticating to Azure Resource Manager using a Service Principal

Terraform supports authenticating to Azure through a Service Principal or the Azure CLI.

We recommend using a Service Principal when running in a Shared Environment (such as within a CI server/automation) - and [authenticating via the Azure CLI](authenticating_via_azure_cli.html) when you're running Terraform locally.

## Creating a Service Principal

A Service Principal is an application within Azure Active Directory whose authentication tokens can be used as the `client_id`, `client_secret`, and `tenant_id` fields needed by Terraform (`subscription_id` can be independently recovered from your Azure account details).

It's possible to complete this task in either the [Azure CLI](#creating-a-service-principal-using-the-azure-cli) or in the [Azure Portal](#creating-a-service-principal-in-the-azure-portal) - in both we'll create a Service Principal which has `Contributor` rights to the subscription. [It's also possible to assign other rights](https://azure.microsoft.com/en-gb/documentation/articles/role-based-access-built-in-roles/) depending on your configuration.

### Creating a Service Principal using the Azure CLI

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

We can now create the Service Principal which will have permissions to manage resources in the specified Subscription using the following command:

```shell
$ az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID"
```

This command will output 5 values:

```json
{
  "appId": "00000000-0000-0000-0000-000000000000",
  "displayName": "azure-cli-2017-06-05-10-41-15",
  "name": "http://azure-cli-2017-06-05-10-41-15",
  "password": "0000-0000-0000-0000-000000000000",
  "tenant": "00000000-0000-0000-0000-000000000000"
}
```

These values map to the Terraform variables like so:

 - `appId` is the `client_id` defined above.
 - `password` is the `client_secret` defined above.
 - `tenant` is the `tenant_id` defined above.

---

Finally - it's possible to test these values work as expected by first logging in:

```shell
$ az login --service-principal -u CLIENT_ID -p CLIENT_SECRET --tenant TENANT_ID
```

Once logged in as the Service Principal - we should be able to list the VM Sizes by specifying an Azure region, for example here we use the `West US` region:

```shell
$ az vm list-sizes --location westus
```

~> **Note**: If you're using the **China**, **German** or **Government** Azure Clouds - you will need to switch `westus` out for another region. You can find which regions are available by running:

```shell
$ az account list-locations
```

### Creating a Service Principal in the Azure Portal

There are two tasks needed to create a Service Principal via [the Azure Portal](https://portal.azure.com):

 1. Create an Application in Azure Active Directory (which acts as a Service Principal)
 2. Grant the Application access to manage resources in your Azure Subscription

### 1. Creating an Application in Azure Active Directory

Firstly navigate to [the **Azure Active Directory** overview](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview) within the Azure Portal - [then select the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and click **Endpoints** at the top of the **App Registration** blade. A list of URIs will be displayed and you need to located the URI for **OAUTH 2.0 AUTHORIZATION ENDPOINT** which contains a GUID. This is your Tenant ID / the `tenant_id` field mentioned above.

Next, navigate back to [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) - from here we'll create the Application in Azure Active Directory. To do this click **Add** at the top to add a new Application within Azure Active Directory. On this page, set the following values then press **Create**:

- **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
- **Application Type** - this should be set to "Web app / API"
- **Sign-on URL** - this can be anything, providing it's a valid URI (e.g. https://terra.form)

Once that's done - select the Application you just created in [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview). At the top of this page, the "Application ID" GUID is the `client_id` you'll need.

Finally, we can create the `client_secret` by selecting **Keys** and then generating a new key by entering a description, selecting how long the `client_secret` should be valid for - and finally pressing **Save**. This value will only be visible whilst on the page, so be sure to copy it now (otherwise you'll need to regenerate a new key).

### 2. Granting the Application access to manage resources in your Azure Subscription

Once the Application exists in Azure Active Directory - we can grant it permissions to modify resources in the Subscription. To do this, [navigate to the **Subscriptions** blade within the Azure Portal](https://portal.azure.com/#blade/Microsoft_Azure_Billing/SubscriptionsBlade), then select the Subscription you wish to use, then click **Access Control (IAM)**, and finally **Add**.

Firstly specify a Role which grants the appropriate permissions needed for the Service Principal (for example, `Contributor` will grant Read/Write on all resources in the Subscription). There's more information about [the built in roles](https://azure.microsoft.com/en-gb/documentation/articles/role-based-access-built-in-roles/) available here.

Secondly, search for and select the name of the Application created in Azure Active Directory to assign it this role - then press **Save**.

## Creating a Service Principal through the Legacy CLI's

It's also possible to create credentials via [the legacy cross-platform CLI](https://azure.microsoft.com/en-us/documentation/articles/resource-group-authenticate-service-principal-cli/) and the [legacy PowerShell Cmdlets](https://azure.microsoft.com/en-us/documentation/articles/resource-group-authenticate-service-principal/) - however we would highly recommend using the Azure CLI above.
