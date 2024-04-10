---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via the Azure CLI"
description: |-
  This guide will cover how to use the Azure CLI as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using the Azure CLI

Terraform supports a number of different methods for authenticating to Azure:

* Authenticating to Azure using the Azure CLI (which is covered in this guide)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)
* [Authenticating to Azure using a Service Principal and Open ID Connect](service_principal_oidc.html)
* [Authenticating to Azure using AKS Workload Identity](aks_workload_identity.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Important Notes about Authenticating using the Azure CLI

* Prior to version 1.20, the AzureRM Provider used a different method of authorizing via the Azure CLI where credentials reset after an hour - as such, we'd recommend upgrading to version 1.20 or later of the AzureRM Provider.
* Terraform only supports authenticating using the `az` CLI (and this must be available on your PATH) - authenticating using the older `azure` CLI or PowerShell Cmdlets are not supported.
* Prior to version 3.44, authenticating via the Azure CLI was only supported when using a User Account. For example `az login --service-principal` was not supported and you had to use either a [Client Secret](service_principal_client_secret.html) or a [Client Certificate](service_principal_client_certificate.html). From 3.44 upwards, authenticating via the Azure CLI is supported when using a Service Principal or Managed Identity.

---

## Logging into the Azure CLI

~> **Note:** If you're using the **China** or **Government** Azure Clouds - you'll need to first configure the Azure CLI to work with that Cloud.  You can do this by running:

```shell
az cloud set --name AzureChinaCloud|AzureUSGovernment
```

---

Firstly, login to the Azure CLI using a User, Service Principal or Managed Identity.

User Account:

```shell
az login
```

Service Principal with a Secret:

```shell
az login --service-principal -u "CLIENT_ID" -p "CLIENT_SECRET" --tenant "TENANT_ID"
```

Service Principal with a Certificate:

```shell
az login --service-principal -u "CLIENT_ID" -p "CERTIFICATE_PEM" --tenant "TENANT_ID"
```

Service Principal with Open ID Connect (for use in CI / CD):

```shell
az login --service-principal -u "CLIENT_ID" --tenant "TENANT_ID"
```

Managed Identity:

```shell
az login --identity

or

az login --identity --username "CLIENT_ID"
```

---

Once logged in - it's possible to list the Subscriptions associated with the account via:

```shell
az account list
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
az account set --subscription="SUBSCRIPTION_ID"
```

---

## Configuring Azure CLI authentication in Terraform

Now that we're logged into the Azure CLI - we can configure Terraform to use these credentials.

To configure Terraform to use the Default Subscription defined in the Azure CLI - we can use the following Provider block:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

It's also possible to configure Terraform to use a specific Subscription - for example:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}

  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

If you're looking to use Terraform across Tenants - it's possible to do this by configuring the Tenant ID field in the Provider block, as shown below:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}

  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "11111111-1111-1111-1111-111111111111"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.
