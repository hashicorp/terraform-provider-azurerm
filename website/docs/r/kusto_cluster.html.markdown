---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cluster"
description: |-
  Manages Kusto (also known as Azure Data Explorer) Cluster
---

# azurerm_kusto_cluster

Manages a Kusto (also known as Azure Data Explorer) Cluster

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-kusto-cluster-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "kustocluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Cluster to create. Only lowercase Alphanumeric characters allowed, starting with a letter. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Cluster should exist. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `allowed_fqdns` - (Optional) List of allowed FQDNs(Fully Qualified Domain Name) for egress from Cluster.

* `allowed_ip_ranges` - (Optional) The list of ips in the format of CIDR allowed to connect to the cluster.

* `double_encryption_enabled` - (Optional) Is the cluster's double encryption enabled? Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `auto_stop_enabled` - (Optional) Specifies if the cluster could be automatically stopped (due to lack of data or no activity for many days). Defaults to `true`.

* `disk_encryption_enabled` - (Optional) Specifies if the cluster's disks are encrypted.

* `streaming_ingestion_enabled` - (Optional) Specifies if the streaming ingest is enabled.

* `public_ip_type` - (Optional) Indicates what public IP type to create - IPv4 (default), or DualStack (both IPv4 and IPv6). Defaults to `IPv4`.

* `public_network_access_enabled` - (Optional) Is the public network access enabled? Defaults to `true`.

* `outbound_network_access_restricted` - (Optional) Whether to restrict outbound network access. Value is optional but if passed in, must be `true` or `false`, default is `false`.

* `purge_enabled` - (Optional) Specifies if the purge operations are enabled.

* `virtual_network_configuration` - (Optional) A `virtual_network_configuration` block as defined below.

~> **Note:** Currently removing `virtual_network_configuration` sets the `virtual_network_configuration` to `Disabled` state. But any changes to `virtual_network_configuration` in `Disabled` state forces a new resource to be created.

* `language_extensions` - (Optional) An list of `language_extensions` to enable. Valid values are: `PYTHON`, `PYTHON_3.10.8` and `R`. `PYTHON` is used to specify Python 3.6.5 image and `PYTHON_3.10.8` is used to specify Python 3.10.8 image. Note that `PYTHON_3.10.8` is only available in skus which support nested virtualization.

~> **Note:** In `v4.0.0` and later version of the AzureRM Provider, `language_extensions` will be changed to a list of `language_extension` block. In each block, `name` and `image` are required. `name` is the name of the language extension, possible values are `PYTHON`, `R`. `image` is the image of the language extension, possible values are `Python3_6_5`, `Python3_10_8` and `R`.

* `optimized_auto_scale` - (Optional) An `optimized_auto_scale` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `trusted_external_tenants` - (Optional) Specifies a list of tenant IDs that are trusted by the cluster. Default setting trusts all other tenants. Use `trusted_external_tenants = ["*"]` to explicitly allow all other tenants, `trusted_external_tenants = ["MyTenantOnly"]` for only your tenant or `trusted_external_tenants = ["<tenantId1>", "<tenantIdx>"]` to allow specific other tenants.

~> **Note:** In v3.0 of `azurerm` a new or updated Kusto Cluster will only allow your own tenant by default. Explicit configuration of this setting will change from `trusted_external_tenants = ["MyTenantOnly"]` to `trusted_external_tenants = []`.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Kusto Cluster should be located. Changing this forces a new Kusto Cluster to be created.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Possible values are `Dev(No SLA)_Standard_D11_v2`, `Dev(No SLA)_Standard_E2a_v4`, `Standard_D14_v2`, `Standard_D11_v2`, `Standard_D16d_v5`, `Standard_D13_v2`, `Standard_D12_v2`, `Standard_DS14_v2+4TB_PS`, `Standard_DS14_v2+3TB_PS`, `Standard_DS13_v2+1TB_PS`, `Standard_DS13_v2+2TB_PS`, `Standard_D32d_v5`, `Standard_D32d_v4`, `Standard_EC8ads_v5`, `Standard_EC8as_v5+1TB_PS`, `Standard_EC8as_v5+2TB_PS`, `Standard_EC16ads_v5`, `Standard_EC16as_v5+4TB_PS`, `Standard_EC16as_v5+3TB_PS`, `Standard_E80ids_v4`, `Standard_E8a_v4`, `Standard_E8ads_v5`, `Standard_E8as_v5+1TB_PS`, `Standard_E8as_v5+2TB_PS`, `Standard_E8as_v4+1TB_PS`, `Standard_E8as_v4+2TB_PS`, `Standard_E8d_v5`, `Standard_E8d_v4`, `Standard_E8s_v5+1TB_PS`, `Standard_E8s_v5+2TB_PS`, `Standard_E8s_v4+1TB_PS`, `Standard_E8s_v4+2TB_PS`, `Standard_E4a_v4`, `Standard_E4ads_v5`, `Standard_E4d_v5`, `Standard_E4d_v4`, `Standard_E16a_v4`, `Standard_E16ads_v5`, `Standard_E16as_v5+4TB_PS`, `Standard_E16as_v5+3TB_PS`, `Standard_E16as_v4+4TB_PS`, `Standard_E16as_v4+3TB_PS`, `Standard_E16d_v5`, `Standard_E16d_v4`, `Standard_E16s_v5+4TB_PS`, `Standard_E16s_v5+3TB_PS`, `Standard_E16s_v4+4TB_PS`, `Standard_E16s_v4+3TB_PS`, `Standard_E64i_v3`, `Standard_E2a_v4`, `Standard_E2ads_v5`, `Standard_E2d_v5`, `Standard_E2d_v4`, `Standard_L8as_v3`, `Standard_L8s`, `Standard_L8s_v3`, `Standard_L8s_v2`, `Standard_L4s`, `Standard_L16as_v3`, `Standard_L16s`, `Standard_L16s_v3`, `Standard_L16s_v2`, `Standard_L32as_v3` and `Standard_L32s_v3`.
* `capacity` - (Optional) Specifies the node count for the cluster. Boundaries depend on the SKU name.

~> **Note:** If no `optimized_auto_scale` block is defined, then the capacity is required.
~> **Note:** If an `optimized_auto_scale` block is defined and no capacity is set, then the capacity is initially set to the value of `minimum_instances`.

---

A `virtual_network_configuration` block supports the following:

* `subnet_id` - (Required) The subnet resource id.

* `engine_public_ip_id` - (Required) Engine service's public IP address resource id.

* `data_management_public_ip_id` - (Required) Data management's service public IP address resource id.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that is configured on this Kusto Cluster. Possible values are: `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Kusto Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `optimized_auto_scale` block supports the following:

* `minimum_instances` - (Required) The minimum number of allowed instances. Must between `0` and `1000`.

* `maximum_instances` - (Required) The maximum number of allowed instances. Must between `0` and `1000`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Kusto Cluster ID.

* `uri` - The FQDN of the Azure Kusto Cluster.

* `data_ingestion_uri` - The Kusto Cluster URI to be used for data ingestion.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Cluster.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Cluster.

## Import

Kusto Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1
```
