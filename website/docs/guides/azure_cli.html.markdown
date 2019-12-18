---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via the Azure CLI"
sidebar_current: "docs-azurerm-guide-authentication-azure-cli"
description: |-
  This guide will cover how to use the Azure CLI as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using the Azure CLI

Terraform supports a number of different methods for authenticating to Azure:

* Authenticating to Azure using the Azure CLI (which is covered in this guide)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Important Notes about Authenticating using the Azure CLI

* Prior to version 1.20 the AzureRM Provider used a different method of authorizing via the Azure CLI where credentials reset after an hour - as such we'd recommend upgrading to version 1.20 or later of the AzureRM Provider.
* Terraform only supports authenticating using the `az` CLI (and this must be available on your PATH) - authenticating using the older `azure` CLI or PowerShell Cmdlets is not supported.
* Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via `az login --service-principal`) you should instead authenticate via the Service Principal directly (either using a [Client Secret](service_principal_client_secret.html) or a [Client Certificate](service_principal_client_certificate.html)).

---

## Logging into the Azure CLI

~> **Note**: If you're using the **China**, **German** or **Government** Azure Clouds - you'll need to first configure the Azure CLI to work with that Cloud.  You can do this by running:

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

The output (similar to below) will display one or more Subscriptions - with the `id` field being the `subscription_id` field referenced above.

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

```bash
$ az account set --subscription="SUBSCRIPTION_ID"
```

---

## Configuring Azure CLI authentication in Terraform

Now that we're logged into the Azure CLI - we can configure Terraform to use these credentials.

To configure Terraform to use the Default Subscription defined in the Azure CLI - we can use the following Provider block:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.38.0"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

It's also possible to configure Terraform to use a specific Subscription - for example:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.38.0"

  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

If you're looking to use Terraform across Tenants - it's possible to do this by configuring the Tenant ID field in the Provider block, as shown below:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.38.0"

  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "11111111-1111-1111-1111-111111111111"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.
