---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via Managed Service Identity"
sidebar_current: "docs-azurerm-guide-authentication-managed-service-identity"
description: |-
  This guide will cover how to use Managed Service Identity as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using Managed Service Identity

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* Authenticating to Azure using Managed Service Identity (which is covered in this guide)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## What is Managed Service Identity?

Certain services within Azure (for example Virtual Machines and Virtual Machine Scale Sets) can be assigned an Azure Active Directory identity which can be used to access the Azure Subscription. This identity can then be assigned permissions to a Subscription, Resource Group or other resources using the Azure Identity and Access Management functionality - however by default no permissions are assigned.

Once a resource is configured with an identity, a local metadata service exposes credentials which can be used by applications such as Terraform.

## Configuring Managed Service Identity

The (simplified) Terraform Configuration below configures a Virtual Machine with Managed Service Identity, and then grants it Contributor access to the Subscription:

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_virtual_machine" "test" {
  # ...

  identity = {
    type = "SystemAssigned"
  }
}

data "azurerm_builtin_role_definition" "contributor" {
  name = "Contributor"
}

resource "azurerm_role_assignment" "test" {
  name               = "${azurerm_virtual_machine.test.name}"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.subscription.id}${data.azurerm_builtin_role_definition.contributor.id}"
  principal_id       = "${lookup(azurerm_virtual_machine.test.identity[0], "principal_id")}"
}
```

## Configuring Managed Service Identity in Terraform

At this point we assume that Managed Service Identity is configured on the resource (e.g. Virtual Machine) being used - and that permissions have been assigned via Azure's Identity and Access Management system.

Terraform can be configured to use Managed Service Identity for authentication in one of two ways: using Environment Variables or by defining the fields within the Provider block.

You can configure Terraform to use Managed Service Identity by setting the Environment Variable `ARM_USE_MSI` to `true`; as shown below:

```shell
$ export ARM_USE_MSI=true
```

-> **Using a Custom MSI Endpoint?** In the unlikely event you're using a custom endpoint for Managed Service Identity - this can be configured using the `ARM_MSI_ENDPOINT` Environment Variable - however this shouldn't need to be configured in regular use.

Whilst a Provider block is _technically_ optional when using Environment Variables - we'd strongly recommend defining one to be able to pin the version of the Provider being used:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.22.0"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using Managed Service Identity.

---

It's also possible to configure Managed Service Identity within the Provider Block:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.22.0"

  use_msi = true
}
```

-> **Using a Custom MSI Endpoint?** In the unlikely event you're using a custom endpoint for Managed Service Identity - this can be configured using the `msi_endpoint` field - however this shouldn't need to be configured in regular use.

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using Managed Service Identity.
