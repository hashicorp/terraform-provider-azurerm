---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_standard"
description: |-
  Manages a Logic Application (Standard / Single Tenant).

---

# azurerm_logic_app_standard

Manages a Logic App (Standard / Single Tenant)

## Example Usage (with App Service Plan)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  os_type  = "Windows"
  sku_name = "WS1"
}

resource "azurerm_logic_app_standard" "example" {
  name                       = "example-logic-app"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  app_settings = {
    "FUNCTIONS_WORKER_RUNTIME"     = "node"
    "WEBSITE_NODE_DEFAULT_VERSION" = "~18"
  }
}
```

## Example Usage (for container mode)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "WorkflowStandard"
    size = "WS1"
  }
}

resource "azurerm_logic_app_standard" "example" {
  name                       = "example-logic-app"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    linux_fx_version = "DOCKER|mcr.microsoft.com/azure-functions/dotnet:3.0-appservice"
  }

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"      = "https://<server-name>.azurecr.io"
    "DOCKER_REGISTRY_SERVER_USERNAME" = "username"
    "DOCKER_REGISTRY_SERVER_PASSWORD" = "password"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Logic App. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the resource group in which to create the Logic App. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Logic App.

* `storage_account_name` - (Required) The backend storage account name which will be used by this Logic App (e.g. for Stateful workflows data). Changing this forces a new resource to be created.

* `storage_account_access_key` - (Required) The access key which will be used to access the backend storage account for the Logic App.

---

* `app_settings` - (Optional) A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

~> **Note:** There are a number of application settings that will be managed for you by this resource type and *shouldn't* be configured separately as part of the app_settings you specify.  `AzureWebJobsStorage` is filled based on `storage_account_name` and `storage_account_access_key`. `WEBSITE_CONTENTSHARE` is detailed below. `FUNCTIONS_EXTENSION_VERSION` is filled based on `version`. `APP_KIND` is set to workflowApp and `AzureFunctionsJobHost__extensionBundle__id` and `AzureFunctionsJobHost__extensionBundle__version` are set as detailed below.

* `use_extension_bundle` - (Optional) Should the logic app use the bundled extension package? If true, then application settings for `AzureFunctionsJobHost__extensionBundle__id` and `AzureFunctionsJobHost__extensionBundle__version` will be created. Defaults to `true`.

* `bundle_version` - (Optional) If `use_extension_bundle` is set to `true` this controls the allowed range for bundle versions. Defaults to `[1.*, 2.0.0)`.

* `connection_string` - (Optional) A `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the Logic App send session affinity cookies, which route client requests in the same session to the same instance?

* `client_certificate_mode` - (Optional) The mode of the Logic App's client certificates requirement for incoming requests. Possible values are `Required` and `Optional`.

* `enabled` - (Optional) Is the Logic App enabled? Defaults to `true`.

* `ftp_publish_basic_authentication_enabled` - (Optional) Whether the FTP basic authentication publishing profile is enabled. Defaults to `true`. 

* `https_only` - (Optional) Can the Logic App only be accessed via HTTPS? Defaults to `false`.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access` - (Optional) Whether Public Network Access should be enabled or not. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

~> **Note:** Setting this property will also set it in the Site Config.

* `scm_publish_basic_authentication_enabled` - (Optional) Whether the default SCM basic authentication publishing profile is enabled. Defaults to `true`.

* `site_config` - (Optional) A `site_config` object as defined below.

* `storage_account_share_name` - (Optional) The name of the share used by the logic app, if you want to use a custom name. This corresponds to the WEBSITE_CONTENTSHARE appsetting, which this resource will create for you. If you don't specify a name, then this resource will generate a dynamic name. This setting is useful if you want to provision a storage account and create a share using `azurerm_storage_share`.

~> **Note:** When integrating a `CI/CD pipeline` and expecting to run from a deployed package in `Azure` you must seed your `app settings` as part of terraform code for Logic App to be successfully deployed. `Important Default key pairs`: (`"WEBSITE_RUN_FROM_PACKAGE" = ""`, `"FUNCTIONS_WORKER_RUNTIME" = "node"` (or Python, etc.), `"WEBSITE_NODE_DEFAULT_VERSION" = "10.14.1"`, `"APPINSIGHTS_INSTRUMENTATIONKEY" = ""`).

~> **Note:** When using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `version` - (Optional) The runtime version associated with the Logic App. Defaults to `~4`.

* `virtual_network_subnet_id` - (Optional) The subnet ID which will be used by this resource for [regional virtual network integration](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).

~> **Note:** The AzureRM Terraform provider provides regional virtual network integration via the standalone resource [app_service_virtual_network_swift_connection](app_service_virtual_network_swift_connection.html) and in-line within this resource using the `virtual_network_subnet_id` property. You cannot use both methods simultaneously.

~> **Note:** Assigning the `virtual_network_subnet_id` property requires [RBAC permissions on the subnet](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#permissions)

* `vnet_content_share_enabled` - (Optional) Specifies whether allow routing traffic between the Logic App and Storage Account content share through a virtual network. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `connection_string` block supports the following:

* `name` - (Required) The name of the Connection String.

* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and `SQLServer`.

* `value` - (Required) The value for the Connection String.

---

The `site_config` block supports the following:

* `always_on` - (Optional) Should the Logic App be loaded at all times? Defaults to `false`.

* `app_scale_limit` - (Optional) The number of workers this Logic App can scale out to. Only applicable to apps on the Consumption and Premium plan.

* `auto_swap_slot_name` - (Optional) The Auto-swap slot name.

* `cors` - (Optional) A `cors` block as defined below.

* `dotnet_framework_version` - (Optional) The version of the .NET framework's CLR used in this Logic App Possible values are `v4.0` (including .NET Core 2.1 and 3.1), `v5.0`, `v6.0` and `v8.0`. [For more information on which .NET Framework version to use based on the runtime version you're targeting - please see this table](https://docs.microsoft.com/azure/azure-functions/functions-dotnet-class-library#supported-versions). Defaults to `v4.0`.

* `elastic_instance_minimum` - (Optional) The number of minimum instances for this Logic App Only affects apps on the Premium plan.

* `ftps_state` - (Optional) State of FTP / FTPS service for this Logic App. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `AllAllowed`.

* `health_check_path` - (Optional) Path which will be checked for this Logic App health.

* `http2_enabled` - (Optional) Specifies whether the HTTP2 protocol should be enabled. Defaults to `false`.

* `ip_restriction` - (Optional) A list of `ip_restriction` objects representing IP restrictions as defined below.

-> **Note:** User has to explicitly set `ip_restriction` to empty slice (`[]`) to remove it.

* `scm_ip_restriction` - (Optional) A list of `scm_ip_restriction` objects representing SCM IP restrictions as defined below.

-> **Note:** User has to explicitly set `scm_ip_restriction` to empty slice (`[]`) to remove it.

* `scm_use_main_ip_restriction` - (Optional) Should the Logic App `ip_restriction` configuration be used for the SCM too. Defaults to `false`.

* `scm_min_tls_version` - (Optional) Configures the minimum version of TLS required for SSL requests to the SCM site. Possible values are `1.0`, `1.1` and `1.2`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `scm_type` - (Optional) The type of Source Control used by the Logic App in use by the Windows Function App. Defaults to `None`. Possible values are: `BitbucketGit`, `BitbucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None`, `OneDrive`, `Tfs`, `VSO`, and `VSTSRM`

* `linux_fx_version` - (Optional) Linux App Framework and version for the App Service, e.g. `DOCKER|(golang:latest)`. Setting this value will also set the `kind` of application deployed to `functionapp,linux,container,workflowapp`.

~> **Note:** You must set `os_type` in `azurerm_service_plan` to `Linux` when this property is set.

* `min_tls_version` - (Optional) The minimum supported TLS version for the Logic App. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new Logic Apps.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `pre_warmed_instance_count` - (Optional) The number of pre-warmed instances for this Logic App Only affects apps on the Premium plan.

* `runtime_scale_monitoring_enabled` - (Optional) Should Runtime Scale Monitoring be enabled?. Only applicable to apps on the Premium plan. Defaults to `false`.

* `use_32_bit_worker_process` - (Optional) Should the Logic App run in 32 bit mode, rather than 64 bit mode? Defaults to `true`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `vnet_route_all_enabled` - (Optional) Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Logic App Standard. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Logic App Standard.

~> **Note:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Logic App has been created. More details are available below.

~> **Note:** The `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, the priority is set to 65000 if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`. 

* `headers` - (Optional) The `headers` block for this specific as a `ip_restriction` block as defined below.

---

A `scm_ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, the priority is set to `65000` if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`.

* `headers` - (Optional) The `headers` block for this specific `ip_restriction` as defined below.

---

A `headers` block supports the following:

* `x_azure_fdid` - (Optional) A list of allowed Azure FrontDoor IDs in UUID notation with a maximum of 8.

* `x_fd_health_probe` - (Optional) A list to allow the Azure FrontDoor health probe header. Only allowed value is `1`.

* `x_forwarded_for` - (Optional) A list of allowed 'X-Forwarded-For' IPs in CIDR notation with a maximum of 8.

* `x_forwarded_host` - (Optional) A list of allowed 'X-Forwarded-Host' domains with a maximum of 8.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App.

* `custom_domain_verification_id` - An identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `default_hostname` - The default hostname associated with the Logic App - such as `mysite.azurewebsites.net`.

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`.

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `kind` - The Logic App kind.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service.

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service.

* `password` - The password associated with the username, which can be used to publish to this App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App
* `update` - (Defaults to 30 minutes) Used when updating the Logic App
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App

## Import

Logic Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_standard.logicapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/logicapp1
```
