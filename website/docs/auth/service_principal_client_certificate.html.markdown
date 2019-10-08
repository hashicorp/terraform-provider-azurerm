---
layout: "azurerm"
page_title: "Azure Provider: Authenticating via a Service Principal and a Client Certificate"
sidebar_current: "docs-azurerm-guide-authentication-service-principal-client-certificate"
description: |-
  This guide will cover how to use a Service Principal (Shared Account) with a Client Certificate as authentication for the Azure Provider.

---

# Azure Provider: Authenticating using a Service Principal with a Client Certificate

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](managed_service_identity.html)
* Authenticating to Azure using a Service Principal and a Client Certificate (which is covered in this guide)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

---

## Creating a Service Principal

A Service Principal is an application within Azure Active Directory which can be used as a means of authentication, either [using a Client Secret](service_principal_client_secret.html) or a Client Certificate (which is documented in this guide) and can be created though the Azure Portal.

This guide will cover how to generate a client certificate, how to create a Service Principal and then how to assign the Client Certificate to the Service Principal so that it can be used for authentication. Once that's done finally we're going to grant the Service Principal permission to manage resources in the Subscription - to do this we're going to assign `Contributor` rights to the Subscription - however [it's possible to assign other permissions](https://azure.microsoft.com/en-gb/documentation/articles/role-based-access-built-in-roles/) depending on your configuration.

---

## Generating a Client Certificate

Firstly we need to create a certificate which can be used for authentication. To do that we're going to generate a Certificate Signing Request (also known as a CSR) using `openssl` (this can also be achieved using PowerShell, however that's outside the scope of this document):

```shell
$ openssl req -newkey rsa:4096 -nodes -keyout "service-principal.key" -out "service-principal.csr"
```

-> During the generation of the certificate you'll be prompted for various bits of information required for the certificate signing request - at least one item has to be specified for this to complete.

We can now sign that Certificate Signing Request, in this example we're going to self-sign this certificate using the Key we just generated; however it's also possible to do this using a Certificate Authority. In order to do that we're again going to use `openssl`:

```shell
$ openssl x509 -signkey "service-principal.key" -in "service-principal.csr" -req -days 365 -out "service-principal.crt"
```

Finally we can generate a PFX file which can be used to authenticate with Azure:

```shell
$ openssl pkcs12 -export -out "service-principal.pfx" -inkey "service-principal.key" -in "service-principal.crt"
```

Now that we've generated a certificate, we can create the Azure Active Directory application.

---

### Creating the Service Principal

We're going to create the Service Principal in the Azure Portal - to do this navigate to [the **Azure Active Directory** overview](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview) within the Azure Portal - [then select the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and click **Endpoints** at the top of the **App Registration** blade. A list of URIs will be displayed and you need to locate the URI for **OAUTH 2.0 AUTHORIZATION ENDPOINT** which contains a GUID. This GUID is your Tenant ID (the `tenant_id` field mentioned above).

Next, navigate back to [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) - from here we'll create the Application in Azure Active Directory. To do this click **New application registration** at the top to add a new Application within Azure Active Directory. On this page, set the following values then press **Create**:

- **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
- **Application Type** - this should be set to "Web app / API"
- **Sign-on URL** - this can be anything, providing it's a valid URI (e.g. https://terra.form)

At this point the newly created Azure Active Directory application should be visible on-screen - if it's not, navigate to the [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and select the Azure Active Directory application. At the top of this page, the "Application ID" GUID is the `client_id` you'll need.

### Assigning the Client Certificate to the Service Principal

To associate the public portion of the Client Certificate (the `*.crt` file) with the Azure Active Directory Application - to do this select **Settings** and then **Keys**. This screen displays the Passwords (Client Secrets) and Public Keys (Client Certificates) which are associated with this Azure Active Directory Application.

The Public Key associated with the generated Certificate can be uploaded by selecting **Upload Public Key**, selecting the file which should be uploaded (in the example above, this'd be `service-principal.crt`) - and then hitting **Save**.

### Allowing the Service Principal to manage the Subscription

Now that we've created the Application within Azure Active Directory and assigned the certificate we're using for authentication, we can now grant the Application permissions to manage the Subscription. To do this, [navigate to the **Subscriptions** blade within the Azure Portal](https://portal.azure.com/#blade/Microsoft_Azure_Billing/SubscriptionsBlade), select the Subscription you wish to use, then click **Access Control (IAM)** and finally **Add role assignment**.

Firstly, specify a Role which grants the appropriate permissions needed for the Service Principal (for example, `Contributor` will grant Read/Write on all resources in the Subscription). More information about [the built in roles can be found here](https://azure.microsoft.com/en-gb/documentation/articles/role-based-access-built-in-roles/).

Secondly, search for and select the name of the Application created in Azure Active Directory to assign it this role - then press **Save**.

At this point the newly created Azure Active Directory Application should be associated with the Certificate that we generated earlier (which can be used as a Client Certificate) - and should have permissions to the Azure Subscription.

---

### Configuring the Service Principal in Terraform

As we've obtained the credentials for this Service Principal - it's possible to configure them in a few different ways.

When storing the credentials as Environment Variables, for example:

```bash
$ export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_CLIENT_CERTIFICATE_PATH="/path/to/my/client/certificate.pfx"
$ export ARM_CLIENT_CERTIFICATE_PASSWORD="Pa55w0rd123"
$ export ARM_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_TENANT_ID="00000000-0000-0000-0000-000000000000"
```

The following Provider block can be specified - where `1.21.0` is the version of the Azure Provider that you'd like to use:

```hcl
provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.34.0"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.

---

It's also possible to configure these variables either in-line or from using variables in Terraform (as the `client_certificate_path` and `client_certificate_password` are in this example), like so:

~> **NOTE:** We'd recommend not defining these variables in-line since they could easily be checked into Source Control.

```hcl
variable "client_certificate_path" {}
variable "client_certificate_password" {}

provider "azurerm" {
  # Whilst version is optional, we /strongly recommend/ using it to pin the version of the Provider being used
  version = "=1.34.0"

  subscription_id             = "00000000-0000-0000-0000-000000000000"
  client_id                   = "00000000-0000-0000-0000-000000000000"
  client_certificate_path     = "${var.client_certificate_path}"
  client_certificate_password = "${var.client_certificate_password}"
  tenant_id                   = "00000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.
