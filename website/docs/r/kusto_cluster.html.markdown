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
  name     = "example"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "example"
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

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Cluster to create. Only lowercase Alphanumeric characters allowed, starting with a letter. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Cluster should exist. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

---

* `allowed_fqdns` - (Optional) List of allowed FQDNs (Fully Qualified Domain Name) for egress from Cluster.

* `allowed_ip_ranges` - (Optional) The list of ips in the format of CIDR allowed to connect to the cluster.

* `auto_stop_enabled` - (Optional) Specifies if the cluster could be automatically stopped (due to lack of data or no activity for many days). Defaults to `true`.

* `callout_policy` - (Optional) A `callout_policy` block as defined below.

* `disk_encryption_enabled` - (Optional) Specifies if the cluster's disks are encrypted. Defaults to `false`.

* `double_encryption_enabled` - (Optional) Is the cluster's double encryption enabled? Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `language_extension` - (Optional) A `language_extension` block as defined below.

* `optimized_auto_scale` - (Optional) An `optimized_auto_scale` block as defined below.

* `outbound_network_access_restricted` - (Optional) Whether to restrict outbound network access. Defaults to `false`.

* `public_ip_type` - (Optional) Indicates what public IP type to create - IPv4 (default), or DualStack (both IPv4 and IPv6). Defaults to `IPv4`.

* `public_network_access_enabled` - (Optional) Is the public network access enabled? Defaults to `true`.

* `purge_enabled` - (Optional) Specifies if the purge operations are enabled. Defaults to `false`.

* `streaming_ingestion_enabled` - (Optional) Specifies if the streaming ingest is enabled. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `trusted_external_tenants` - (Optional) Specifies a list of tenant IDs that are trusted by the cluster. Default setting trusts all other tenants. Use `trusted_external_tenants = ["*"]` to explicitly allow all other tenants, `trusted_external_tenants = ["MyTenantOnly"]` for only your tenant or `trusted_external_tenants = ["<tenantId1>", "<tenantIdx>"]` to allow specific other tenants.

~> **Note:** In v3.0 of `azurerm` a new or updated Kusto Cluster will only allow your own tenant by default. Explicit configuration of this setting will change from `trusted_external_tenants = ["MyTenantOnly"]` to `trusted_external_tenants = []`.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Kusto Cluster should be located. Changing this forces a new Kusto Cluster to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that is configured on this Kusto Cluster. Possible values are: `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Kusto Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `callout_policy` block supports the following:

* `callout_type` - (Required) The type of callout service. Possible values are `azure_digital_twins`, `azure_openai`, `cosmosdb`, `external_data`, `genevametrics`, `kusto`, `mysql`, `postgresql`, `sandbox_artifacts`, `sql`, and `webapi`.

* `callout_uri_regex` - (Required) A regular expression or the callout URI.

* `outbound_access` - (Required) Whether outbound access is permitted for the specified service with the URI pattern. Possible values are `Allow` and `Deny`.

---

A `language_extension` block supports the following:

* `name` - (Required) The name of the language extension. Possible values are `PYTHON` and `R`. 

* `image` - (Required) The language extension image. Possible values are `Python3_11_7`, `Python3_11_7_DL`, `Python3_10_8`, `Python3_10_8_DL`, `Python3_6_5`, `PythonCustomImage`, and `R`.

---

An `optimized_auto_scale` block supports the following:

* `minimum_instances` - (Required) The minimum number of allowed instances. Possible values range between `0` and `1000`.

* `maximum_instances` - (Required) The maximum number of allowed instances. Possible values range between `0` and `1000`.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Possible values are `Dev(No SLA)_Standard_D11_v2`, `Dev(No SLA)_Standard_E2a_v4`, `Standard_D14_v2`, `Standard_D11_v2`, `Standard_D16d_v5`, `Standard_D13_v2`, `Standard_D12_v2`, `Standard_DS14_v2+4TB_PS`, `Standard_DS14_v2+3TB_PS`, `Standard_DS13_v2+1TB_PS`, `Standard_DS13_v2+2TB_PS`, `Standard_D32d_v5`, `Standard_D32d_v4`, `Standard_EC8ads_v5`, `Standard_EC8as_v5+1TB_PS`, `Standard_EC8as_v5+2TB_PS`, `Standard_EC16ads_v5`, `Standard_EC16as_v5+4TB_PS`, `Standard_EC16as_v5+3TB_PS`, `Standard_E80ids_v4`, `Standard_E8a_v4`, `Standard_E8ads_v5`, `Standard_E8as_v5+1TB_PS`, `Standard_E8as_v5+2TB_PS`, `Standard_E8as_v4+1TB_PS`, `Standard_E8as_v4+2TB_PS`, `Standard_E8d_v5`, `Standard_E8d_v4`, `Standard_E8s_v5+1TB_PS`, `Standard_E8s_v5+2TB_PS`, `Standard_E8s_v4+1TB_PS`, `Standard_E8s_v4+2TB_PS`, `Standard_E4a_v4`, `Standard_E4ads_v5`, `Standard_E4d_v5`, `Standard_E4d_v4`, `Standard_E16a_v4`, `Standard_E16ads_v5`, `Standard_E16as_v5+4TB_PS`, `Standard_E16as_v5+3TB_PS`, `Standard_E16as_v4+4TB_PS`, `Standard_E16as_v4+3TB_PS`, `Standard_E16d_v5`, `Standard_E16d_v4`, `Standard_E16s_v5+4TB_PS`, `Standard_E16s_v5+3TB_PS`, `Standard_E16s_v4+4TB_PS`, `Standard_E16s_v4+3TB_PS`, `Standard_E64i_v3`, `Standard_E2a_v4`, `Standard_E2ads_v5`, `Standard_E2d_v5`, `Standard_E2d_v4`, `Standard_L8as_v3`, `Standard_L8s`, `Standard_L8s_v3`, `Standard_L8s_v2`, `Standard_L4s`, `Standard_L16as_v3`, `Standard_L16s`, `Standard_L16s_v3`, `Standard_L16s_v2`, `Standard_L32as_v3` and `Standard_L32s_v3`.

* `capacity` - (Optional) Specifies the node count for the cluster. Boundaries depend on the SKU name.

~> **Note:** If no `optimized_auto_scale` block is defined, then the capacity is required.

~> **Note:** If an `optimized_auto_scale` block is defined and no capacity is set, then the capacity is initially set to the value of `minimum_instances`.

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Cluster.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Cluster.

## Import

Kusto Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Kusto` - 2024-04-13
