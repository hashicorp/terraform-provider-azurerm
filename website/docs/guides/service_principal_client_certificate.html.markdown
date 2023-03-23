---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via a Service Principal and a Client Certificate"
description: |-
  This guide will cover how to use a Service Principal (Shared Account) with a Client Certificate as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using a Service Principal with a Client Certificate

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
* Authenticating to Azure using a Service Principal and a Client Certificate (which is covered in this guide)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)
* [Authenticating to Azure using a Service Principal and OpenID Connect](service_principal_oidc.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

---

## Setting up an Application and Service Principal

A Service Principal is a security principal within Azure Active Directory which can be granted access to resources within Azure Subscriptions. To authenticate with a Service Principal, you will need to create an Application object within Azure Active Directory, which you will use as a means of authentication, either [using a Client Secret](service_principal_client_secret.html), a Client Certificate (which is documented in this guide), or [OpenID Connect](service_principal_oidc.html). This can be done using the Azure Portal.

This guide will cover how to generate a client certificate, how to create an Application and linked Service Principal, and then how to assign the Client Certificate to the Application so that it can be used for authentication. Once that's done finally we're going to grant the Service Principal permission to manage resources in the Subscription - to do this we're going to assign `Contributor` rights to the Subscription - however, [it's possible to assign other permissions](https://docs.microsoft.com/azure/role-based-access-control/built-in-roles) depending on your configuration.

---

## Generating a Client Certificate

Firstly we need to create a certificate which can be used for authentication. To do that we're going to generate a private key and self-signed certificate using OpenSSL or LibreSSL (this can also be achieved using PowerShell, however that's outside the scope of this document):

```shell
$ openssl req -subj '/CN=myclientcertificate/O=MyCompany, Inc./ST=CA/C=US' \
    -new -newkey rsa:4096 -sha256 -days 730 -nodes -x509 -keyout client.key -out client.crt
```

Next we generate a PKCS#12 bundle (.pfx file) which can be used by the AzureRM provider to authenticate with Azure:

```shell
# note: the password is intentionally quoted for shell compatibility, the value does not include the quotes
$ openssl pkcs12 -export -password pass:"Pa55w0rd123" -out client.pfx -inkey client.key -in client.crt
```

Now that we've generated a certificate, we can create the Azure Active Directory Application.

---

## Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this navigate to [the **Azure Active Directory** overview](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview) within the Azure Portal - [then select the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview). Click the **New registration** button at the top to add a new Application within Azure Active Directory. On this page, set the following values then press **Create**:

* **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
* **Supported Account Types** - this should be set to "Accounts in this organizational directory only (single-tenant)"
* **Redirect URI** - you should choose "Web" for the URI type. the actual value can be left blank

At this point the newly created Azure Active Directory application should be visible on-screen - if it's not, navigate to the [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and select the Azure Active Directory application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

### Assigning the Client Certificate to the Azure Active Directory Application

To associate the public portion of the Client Certificate (the `*.crt` file) with the Azure Active Directory Application - to do this select **Certificates & secrets**. This screen displays the Certificates and Client Secrets (i.e. passwords) which are associated with this Azure Active Directory Application.

The Public Key associated with the generated Certificate can be uploaded by selecting **Upload Certificate**, selecting the file which should be uploaded (in the example above, that'd be `service-principal.crt`) - and then hit **Add**.

### Allowing the Service Principal to manage the Subscription

Now that we've created the Application within Azure Active Directory and assigned the certificate we're using for authentication, we can now grant the Application permissions to manage the Subscription via its linked Service Principal. To do this, [navigate to the **Subscriptions** blade within the Azure Portal](https://portal.azure.com/#blade/Microsoft_Azure_Billing/SubscriptionsBlade), select the Subscription you wish to use, then click **Access Control (IAM)** and finally **Add** > **Add role assignment**.

Firstly, specify a Role which grants the appropriate permissions needed for the Service Principal (for example, `Contributor` will grant Read/Write on all resources in the Subscription). More information about [the built in roles can be found here](https://docs.microsoft.com/azure/role-based-access-control/built-in-roles).

Secondly, search for and select the name of the Service Principal created in Azure Active Directory to assign it this role - then press **Save**.

At this point the newly created Azure Active Directory Application should be associated with the Certificate that we generated earlier (which can be used as a Client Certificate) - and should have permissions to the Azure Subscription.

---

## Configuring Terraform to use the Client Certificate

Now that we have our Client Certificate uploaded to Azure and ready to use, it's possible to configure Terraform in a few different ways.

The provider can be configured to read the certificate bundle from the .pfx file in your filesystem, or alternatively you can pass a base64-encoded copy of the certificate bundle directly to the provider.

### Environment Variables

Our recommended approach is storing the credentials as Environment Variables, for example:

*Reading the certificate bundle from the filesystem*
```shell-session
# sh
$ export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_CLIENT_CERTIFICATE_PATH="/path/to/my/client/certificate.pfx"
$ export ARM_CLIENT_CERTIFICATE_PASSWORD="Pa55w0rd123"
$ export ARM_TENANT_ID="10000000-0000-0000-0000-000000000000"
$ export ARM_SUBSCRIPTION_ID="20000000-0000-0000-0000-000000000000"
```
```powershell
# PowerShell
> $env:ARM_CLIENT_ID = "00000000-0000-0000-0000-000000000000"
> $env:ARM_CLIENT_CERTIFICATE_PATH = "C:\Users\myusername\Documents\my\client\certificate.pfx"
> $env:ARM_CLIENT_CERTIFICATE_PASSWORD = "Pa55w0rd123"
> $env:ARM_TENANT_ID = "10000000-0000-0000-0000-000000000000"
> $env:ARM_SUBSCRIPTION_ID = "20000000-0000-0000-0000-000000000000"
```

*Passing the encoded certificate bundle directly*
```shell-session
# sh
$ export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_CLIENT_CERTIFICATE="$(base64 /path/to/my/client/certificate.pfx)"
$ export ARM_CLIENT_CERTIFICATE_PASSWORD="Pa55w0rd123"
$ export ARM_TENANT_ID="10000000-0000-0000-0000-000000000000"
$ export ARM_SUBSCRIPTION_ID="20000000-0000-0000-0000-000000000000"
```
```powershell
# PowerShell
> $env:ARM_CLIENT_ID = "00000000-0000-0000-0000-000000000000"
> $env:ARM_CLIENT_CERTIFICATE = [Convert]::ToBase64String([System.IO.File]::ReadAllBytes("C:\Users\myusername\Documents\my\client\certificate.pfx"))
> $env:ARM_CLIENT_CERTIFICATE_PASSWORD = "Pa55w0rd123"
> $env:ARM_TENANT_ID = "10000000-0000-0000-0000-000000000000"
> $env:ARM_SUBSCRIPTION_ID = "20000000-0000-0000-0000-000000000000"
```

The following Terraform and Provider blocks can be specified - where `3.0.0` is the version of the Azure Provider that you'd like to use:

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

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.

### Provider Block

It's also possible to configure these variables either directly, or from variables, in your provider block, like so:

!> **Caution** We recommend not defining these variables in-line since they could easily be checked into Source Control.

*Reading the certificate bundle from the filesystem*
```hcl
variable "client_certificate" {}
variable "client_certificate_password" {}

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

  client_id                   = "00000000-0000-0000-0000-000000000000"
  client_certificate_path     = var.client_certificate_path
  client_certificate_password = var.client_certificate_password
  tenant_id                   = "10000000-0000-0000-0000-000000000000"
  subscription_id             = "20000000-0000-0000-0000-000000000000"
}
```

*Passing the encoded certificate bundle directly*
```hcl
variable "client_certificate" {}
variable "client_certificate_password" {}

# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.43.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}

  client_id                   = "00000000-0000-0000-0000-000000000000"
  client_certificate          = var.client_certificate
  client_certificate_password = var.client_certificate_password
  tenant_id                   = "10000000-0000-0000-0000-000000000000"
  subscription_id             = "20000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.
