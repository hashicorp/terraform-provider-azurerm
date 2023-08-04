---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via the Azure Developer CLI (azd)"
description: |-
  This guide will cover how to use the Azure Developer CLI as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using the Azure Developer CLI (azd)

Terraform supports a number of different methods for authenticating to Azure:

* Authenticating to Azure using the Azure Developer CLI (which is covered in this guide)
* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)
* [Authenticating to Azure using a Service Principal and Open ID Connect](service_principal_oidc.html)

---

Use the [Azure Developer CLI](https://learn.microsoft.com/azure/developer/azure-developer-cli/overview) to extend the `azd` authentication for terraform. This is convenience for terraform templates.

---

## Logging into the Azure Developer CLI

Login to the Azure CLI using:

```shell
azd auth login
```
---

## Configuring Azure Developer CLI authentication in Terraform

Now that we're logged into the Azure Developer CLI - we can configure Terraform to use these credentials.

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

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure Developer CLI to authenticate.

---

The Azure Developer CLI uses the selected Azure subscription during `azd provision` or `azd up`. 
