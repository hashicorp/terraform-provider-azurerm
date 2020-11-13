---
layout: "azurerm"
page_title: "Azure Active Directory: Migrating to the AzureAD Provider"
description: |-
  This page documents how to migrate from using the AzureAD resources within this repository to the resources in the new split-out repository.

---

# Azure Active Directory: Migrating to the AzureAD Provider

In v1.21 of the AzureRM Provider the Azure Active Directory Data Sources and Resources have been split out into a new Provider specifically for Azure Active Directory.

This guide covers how to migrate from using the following Data Sources and Resources in the AzureRM Provider to using them in the new AzureAD Provider:

* Data Source: `azurerm_azuread_application`
* Data Source: `azurerm_azuread_service_principal`
* Resource: `azurerm_azuread_application`
* Resource: `azurerm_azuread_service_principal`
* Resource: `azurerm_azuread_service_principal_password`

## Updating the Provider block

As the AzureAD and AzureRM Provider support the same authentication methods - it's possible to update the Provider block by setting the new Provider name and version, for example:

```hcl
provider "azurerm" {
  version = "=1.44.0"
}
```

can become:

```hcl
provider "azuread" {
  version = "=0.10.0"
}
```

## Updating the Terraform Configurations

The Azure Active Directory Data Sources and Resources have been split out into the new Provider - which means the name of the Data Sources and Resources has changed slightly.

The main difference in naming is that the `azurerm_` prefix has been removed from the names of the Data Sources and Resources - the following table explains the new name for each of the Azure Active Directory resources:


| Type        | Old Name                                   | New Name                           |
| ----------- | ------------------------------------------ | ---------------------------------- |
| Data Source | azurerm_azuread_application                | azuread_application                |
| Data Source | azurerm_azuread_service_principal          | azuread_service_principal          |
| Resource    | azurerm_azuread_application                | azuread_application                |
| Resource    | azurerm_azuread_service_principal          | azuread_service_principal          |
| Resource    | azurerm_azuread_service_principal_password | azuread_service_principal_password |

---

Once the Provider blocks have been updated, it should be possible to replace the `azurerm_` prefix in your Terraform Configuration from each of the AzureAD resources (and any interpolations) so that the new resources in the AzureAD Provider are used instead.

For example the following Terraform Configuration:

```hcl
resource "azurerm_azuread_application" "example" {
  name = "my-application"
}

resource "azurerm_azuread_service_principal" "example" {
  application_id = azurerm_azuread_application.example.application_id
}

resource "azurerm_azuread_service_principal_password" "example" {
  service_principal_id = azurerm_azuread_service_principal.example.id
  value                = "bd018069-622d-4b46-bcb9-2bbee49fe7d9"
  end_date             = "2020-01-01T01:02:03Z"
}
```

we can remove the `azurerm_` prefix from each of the resource names and interpolations to use the `AzureAD` provider instead of by making this:

```hcl
resource "azuread_application" "example" {
  name = "my-application"
}

resource "azuread_service_principal" "example" {
  application_id = azuread_application.example.application_id
}

resource "azuread_service_principal_password" "example" {
  service_principal_id = azuread_service_principal.example.id
  value                = "bd018069-622d-4b46-bcb9-2bbee49fe7d9"
  end_date             = "2020-01-01T01:02:03Z"
}
```

At this point, it should be possible to run `terraform init`, which will download the new AzureAD Provider.


## Migrating Resources in the State

Now that we've updated the Provider Block and the Terraform Configuration we need to update the names of the resources in the state.

The method for performing this differs between Terraform v0.11 and Terraform v0.12, due to improved state handling in v0.12 which protects against moving resources between providers.

### Terraform v0.11

Firstly, it's a good idea to create a backup of your statefile. For a local statefile, simply create a copy. If you are using a Remote State Backend, ensure your backend platform is creating snapshots or backups for rollback purposes.

Let's list the existing items in the state - we can do this by running `terraform state list`, for example:

```shell
$ terraform state list
azurerm_azuread_application.example
azurerm_azuread_service_principal.example
azurerm_azuread_service_principal_password.import
azurerm_azuread_service_principal_password.example
```

As the Terraform Configuration has been updated - we can move each of the resources in the state using the `terraform state mv` command, for example:

```shell
$ terraform state mv azurerm_azuread_application.example azuread_application.example
Moved azurerm_azuread_application.example to azuread_application.example
```

This needs to be repeated for each of the Azure Active Directory resources which exist in the state.

Note that if you encounter any problems with the built-in state management commands, you can also follow the instructions below for Terraform v0.12.

### Terraform v0.12

With Terraform v0.12 (or later), this operation needs to be performed manually. To do this, you will need a local copy of your statefile. If you are using a Remote State Backend, you will first need to download a copy of your statefile.

```shell
$ terraform state pull >current.tfstate
```

Once you have a local copy of your statefile, you can run the following command to replace the necessary values for your resources. This will work on all matching resources in your state, and you will need the [jq](https://stedolan.github.io/jq/download/) tool (version 1.5 or later).

```shell
$ jq '
  def migrateName: sub("^azurerm_azuread_";"azuread_");
  .resources[].type |= migrateName |
  .resources[].instances[].dependencies[]? |= migrateName' \
  <current.tfstate \
  >new.tfstate
```

Inspect the `new.tfstate` file and compare it against your `current.tfstate` file to verify the correct attributes were changed.

```shell
$ diff current.tfstate new.tfstate
10c10
<       "type": "azurerm_azuread_application",
---
>       "type": "azuread_application",
32c32
<       "type": "azurerm_azuread_service_principal",
---
>       "type": "azuread_service_principal",
45c45
<             "azurerm_azuread_application.test"
---
>             "azuread_application.test"
52c52
<       "type": "azurerm_azuread_service_principal_password",
---
>       "type": "azuread_service_principal_password",
68,69c68,69
<             "azurerm_azuread_application.test",
<             "azurerm_azuread_service_principal.test"
---
>             "azuread_application.test",
>             "azuread_service_principal.test"
```

There should be no unexpected or unrelated changes in your diff output.

For a remote state, you will need to push the new statefile to your backend.

```shell
$ terraform state push new.tfstate
```


## Verifying the new State

Once this has been done, running `terraform plan` should show no changes:

```shell
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

At this point, you've switched over to using [the new Azure Active Directory provider](http://terraform.io/docs/providers/azuread/index.html)! You can stay up to date with Releases (and file Feature Requests/Bugs) [on the Github repository](https://github.com/terraform-providers/terraform-provider-azuread).
