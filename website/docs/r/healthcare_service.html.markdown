---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_service"
description: |-
  Manages a Healthcare Service.
---

# azurerm_healthcare_service

Manages a Healthcare Service.

## Example Usage

```hcl
resource "azurerm_healthcare_service" "example" {
  name                = "uniquefhirname"
  resource_group_name = "sample-resource-group"
  location            = "westus2"
  kind                = "fhir-R4"
  cosmosdb_throughput = "2000"

  access_policy_object_ids = ["xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]

  tags = {
    "environment" = "testenv"
    "purpose"     = "AcceptanceTests"
  }

  authentication_configuration {
    authority           = "https://login.microsoftonline.com/$%7Bdata.azurerm_client_config.current.tenant_id%7D"
    audience            = "https://azurehealthcareapis.com/"
    smart_proxy_enabled = "true"
  }

  cors_configuration {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["x-tempo-*", "x-tempo2-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = "500"
    allow_credentials  = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service instance. Used for service endpoint, must be unique within the audience.
* `resource_group_name` - (Required) The name of the Resource Group in which to create the Service.
* `location` - (Required) Specifies the supported Azure Region where the Service should be created.

~> **Please Note**: Not all locations support this resource. Some are `West US 2`, `North Central US`, and `UK West`.

* `access_policy_ids` - (Optional) A set of Azure object id's that are allowed to access the Service. If not configured, the default value is the object id of the service principal or user that is running Terraform.
* `authentication_configuration` - (Optional) An `authentication_configuration` block as defined below.
* `cosmosdb_throughput` - (Optional) The provisioned throughput for the backing database. Range of `400`-`1000`. Defaults to `400`.
* `cosmosdb_key_vault_key_versionless_id` - (Optional) A versionless Key Vault Key ID for CMK encryption of the backing database. Changing this forces a new resource to be created.

~> **Please Note** In order to use a `Custom Key` from Key Vault for encryption you must grant Azure Cosmos DB Service access to your key vault. For instructions on how to configure your Key Vault correctly please refer to the [product documentation](https://docs.microsoft.com/en-us/azure/cosmos-db/how-to-setup-cmk#add-an-access-policy-to-your-azure-key-vault-instance)

* `cors_configuration` - (Optional) A `cors_configuration` block as defined below.
* `public_network_access_enabled` - (Optional) Whether public network access is enabled or disabled for this service instance.
* `kind` - (Optional) The type of the service. Values at time of publication are: `fhir`, `fhir-Stu3` and `fhir-R4`. Default value is `fhir`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---
An `authentication_configuration` supports the following:

* `authority` - (Optional) The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
Authority must be registered to Azure AD and in the following format: https://{Azure-AD-endpoint}/{tenant-id}.
* `audience` - (Optional) The intended audience to receive authentication tokens for the service. The default value is https://azurehealthcareapis.com
* `smart_proxy_enabled` - (Boolean) Enables the 'SMART on FHIR' option for mobile and web implementations.

---
A `cors_configuration` block supports the following:

* `allowed_origins` - (Required) A set of origins to be allowed via CORS.
* `allowed_headers` - (Required) A set of headers to be allowed via CORS.
* `allowed_methods` - (Required) The methods to be allowed via CORS.
* `max_age_in_seconds` - (Required) The max age to be allowed via CORS.
* `allow_credentials` - (Boolean) If credentials are allowed via CORS.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Service.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthcare Service.
* `update` - (Defaults to 30 minutes) Used when updating the Healthcare Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare Service.

## Import

Healthcare Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource_group/providers/Microsoft.HealthcareApis/services/service_name
```
