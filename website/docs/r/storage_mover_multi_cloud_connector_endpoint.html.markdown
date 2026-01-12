---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_multi_cloud_connector_endpoint"
description: |-
  Manages a Storage Mover Multi-Cloud Connector Endpoint.
---

# azurerm_storage_mover_multi_cloud_connector_endpoint

Manages a Storage Mover Multi-Cloud Connector Endpoint for migrating data from AWS S3 to Azure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_multi_cloud_connector_endpoint" "example" {
  name                    = "example-mcce"
  storage_mover_id        = azurerm_storage_mover.example.id
  multi_cloud_connector_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.HybridConnectivity/publicCloudConnectors/example-connector"
  aws_s3_bucket_id        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/aws-rg/providers/Microsoft.AWSConnector/s3Buckets/example-bucket"
  description             = "Example Multi-Cloud Connector Endpoint"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Multi-Cloud Connector Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover for this Multi-Cloud Connector Endpoint. Changing this forces a new resource to be created.

* `multi_cloud_connector_id` - (Required) Specifies the resource ID of the Multi-Cloud Connector. Changing this forces a new resource to be created.

* `aws_s3_bucket_id` - (Required) Specifies the resource ID of the AWS S3 bucket. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for the Storage Mover Multi-Cloud Connector Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Multi-Cloud Connector Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Multi-Cloud Connector Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Multi-Cloud Connector Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Multi-Cloud Connector Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Multi-Cloud Connector Endpoint.

## Import

Storage Mover Multi-Cloud Connector Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_multi_cloud_connector_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01

