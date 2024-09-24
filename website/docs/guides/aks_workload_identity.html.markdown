---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via AKS Workload Identity"
description: |-
  This guide will cover how to use AKS Workload Identity for pods in Azure AKS clusters as authentication for the Azure Provider.
---

# Azure Provider: Authenticating using managed identities for Azure Kubernetes Service with Workload Identity

Terraform supports a number of different methods for authenticating to Azure:

- [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
- [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
- [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)
- [Authenticating to Azure using OpenID Connect](service_principal_oidc.html)
- Authenticating to Azure using AKS Workload Identity (covered in this guide)

---

We recommend using a service principal or a managed identity when running Terraform non-interactively (such as when running Terraform in a CI/CD pipeline), and authenticating using the Azure CLI when running Terraform locally.

## What is AKS Workload Identity?

[AKS Workload Identity](https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview) can be used to authenticate to services that support Azure Active Directory (Azure AD) authentication when running in Azure Kubernetes Service clusters.

When a service account and pod are configured to use AKS Workload Identity, a federated identity token is injected into the pod at run-time, along with environment variables to use that identity.

## Configuring a workload to use an AKS Workload Identity

The (simplified) Terraform configuration below provisions a cluster with workload identity enabled, creates an identity and federated identity credential suitable for a workload identity, and then grants the Contributor role to the identity.

```hcl
data "azurerm_subscription" "current" {}

variable "workload_sa_name" {
  type        = string
  description = "Kubernetes service account to permit"
}

variable "workload_sa_namespace" {
  type        = string
  description = "Kubernetes service account namespace to permit"
}

resource "azurerm_kubernetes_cluster" "mycluster" {
  # ...
  workload_identity_enabled = true
}

resource "azurerm_user_assigned_identity" "myworkload_identity" {
  # ...
  name = "myworkloadidentity"
}

resource "azurerm_federated_identity_credential" "myworkload_identity" {
  name                = azurerm_user_assigned_identity.myworkload_identity.name
  resource_group_name = azurerm_user_assigned_identity.myworkload_identity.resource_group_name
  parent_id           = azurerm_user_assigned_identity.myworkload_identity.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = azurerm_kubernetes_cluster.mycluster.oidc_issuer_url
  subject             = "system:serviceaccount:${workload_sa_namespace}:${workload_sa_name}"
}

data "azurerm_role_definition" "contributor" {
  name = "Contributor"
}

resource "azurerm_role_assignment" "example" {
  scope              = data.azurerm_subscription.current.id
  role_definition_id = "${data.azurerm_subscription.current.id}${data.azurerm_role_definition.contributor.id}"
  principal_id       = azurerm_user_assigned_identity.wayfinder_main.principal_id
}

output "myworkload_identity_client_id" {
  description = "The client ID of the created managed identity to use for the annotation 'azure.workload.identity/client-id' on your service account"
  value       = azurerm_user_assigned_identity.myworkload_identity.client_id
}
```

## Configuring Terraform to use an AKS workload identity

At this point we assume that workload identity is configured on the AKS cluster being used and that permissions have been assigned via Azure's Identity and Access Management system.

Terraform can be configured to use AKS workload identity for authentication in one of two ways: using environment variables, or by defining the field within the provider block.

### Configuring with environment variables

Setting the `ARM_USE_AKS_WORKLOAD_IDENTITY` environment variable (equivalent to provider block argument [`use_aks_workload_identity`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#use_aks_workload_identity)) to `true` tells Terraform to use an AKS workload identity. It is also suggested to disable Azure CLI authentication by setting the `ARM_USE_CLI` environment variable (equivalent to provider block argument [`use_cli`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#use_cli)) to `false`.

If you have not annotated your Kubernetes service account with `azure.workload.identity/client-id`, you will need to specify the `ARM_CLIENT_ID` environment variable (equivalent to provider block argument [`client_id`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#client_id)) to the [client id](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/user_assigned_identity#client_id) of the identity.

In addition to a properly-configured managed identity, Terraform needs to know the subscription ID to fully configure the AzureRM provider. The tenant ID will be detected from the environment provided by AKS Workload Identity.

```shell
export ARM_USE_AKS_WORKLOAD_IDENTITY=true
export ARM_USE_CLI=false
export ARM_SUBSCRIPTION_ID=159f2485-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export ARM_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx # only necessary if the service account is not annotated with the relevant client ID
```

A provider block is _technically_ optional when using environment variables. Even so, we recommend defining provider blocks so that you can pin or constrain the version of the provider being used, and configure other optional settings:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.1.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
}
```

### Configuring with the provider block

It's also possible to configure an AKS workload identity within the provider block:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.1.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}

  use_aks_workload_identity = true
  use_cli                   = false
  #...
}
```

More information on [the fields supported in the provider block can be found here](../index.html#argument-reference).
