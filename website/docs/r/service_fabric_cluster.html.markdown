---
subcategory: "Service Fabric"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_cluster"
description: |-
  Manages a Service Fabric Cluster.
---

# azurerm_service_fabric_cluster

Manages a Service Fabric Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_fabric_cluster" "example" {
  name                 = "example-servicefabric"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  reliability_level    = "Bronze"
  upgrade_mode         = "Manual"
  cluster_code_version = "7.1.456.959"
  vm_image             = "Windows"
  management_endpoint  = "https://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}

```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Fabric Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Service Fabric Cluster exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Service Fabric Cluster should exist. Changing this forces a new resource to be created.

* `reliability_level` - (Required) Specifies the Reliability Level of the Cluster. Possible values include `None`, `Bronze`, `Silver`, `Gold` and `Platinum`.

-> **NOTE:** The Reliability Level of the Cluster depends on the number of nodes in the Cluster: `Platinum` requires at least 9 VM's, `Gold` requires at least 7 VM's, `Silver` requires at least 5 VM's, `Bronze` requires at least 3 VM's.

* `management_endpoint` - (Required) Specifies the Management Endpoint of the cluster such as `http://example.com`. Changing this forces a new resource to be created.

* `node_type` - (Required) One or more `node_type` blocks as defined below.

* `upgrade_mode` - (Required) Specifies the Upgrade Mode of the cluster. Possible values are `Automatic` or `Manual`.

* `vm_image` - (Required) Specifies the Image expected for the Service Fabric Cluster, such as `Windows`. Changing this forces a new resource to be created.

---

* `cluster_code_version` - (Optional) Required if Upgrade Mode set to `Manual`, Specifies the Version of the Cluster Code of the cluster.

* `add_on_features` - (Optional) A List of one or more features which should be enabled, such as `DnsService`.

* `azure_active_directory` - (Optional) An `azure_active_directory` block as defined below.

* `certificate_common_names` - (Optional) A `certificate_common_names` block as defined below. Conflicts with `certificate`.

* `certificate` - (Optional) A `certificate` block as defined below. Conflicts with `certificate_common_names`.

* `reverse_proxy_certificate` - (Optional) A `reverse_proxy_certificate` block as defined below. Conflicts with `reverse_proxy_certificate_common_names`.

* `reverse_proxy_certificate_common_names` - (Optional) A `reverse_proxy_certificate_common_names` block as defined below. Conflicts with `reverse_proxy_certificate`.

* `client_certificate_thumbprint` - (Optional) One or more `client_certificate_thumbprint` blocks as defined below.

* `client_certificate_common_name` - (Optional) A `client_certificate_common_name` block as defined below.

-> **NOTE:** If Client Certificates are enabled then at a Certificate must be configured on the cluster.

* `diagnostics_config` - (Optional) A `diagnostics_config` block as defined below. Changing this forces a new resource to be created.

* `fabric_settings` - (Optional) One or more `fabric_settings` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `azure_active_directory` block supports the following:

* `tenant_id` - (Required) The Azure Active Directory Tenant ID.

* `cluster_application_id` - (Required) The Azure Active Directory Cluster Application ID.

* `client_application_id` - (Required) The Azure Active Directory Client ID which should be used for the Client Application.

---

A `certificate_common_names` block supports the following:

* `common_names` - (Required) A `common_names` block as defined below.

* `x509_store_name` - (Required) The X509 Store where the Certificate Exists, such as `My`.

---

A `common_names` block supports the following:

* `certificate_common_name` - (Required) The common or subject name of the certificate.

* `certificate_issuer_thumbprint` - (Optional) The Issuer Thumbprint of the Certificate.

-> **NOTE:** Certificate Issuer Thumbprint may become required in the future, `https://docs.microsoft.com/en-us/azure/service-fabric/service-fabric-create-cluster-using-cert-cn#download-and-update-a-sample-template`.

---

A `certificate` block supports the following:

* `thumbprint` - (Required) The Thumbprint of the Certificate.

* `thumbprint_secondary` - (Required) The Secondary Thumbprint of the Certificate.

* `x509_store_name` - (Required) The X509 Store where the Certificate Exists, such as `My`.

---

A `reverse_proxy_certificate` block supports the following:

* `thumbprint` - (Required) The Thumbprint of the Certificate.

* `thumbprint_secondary` - (Required) The Secondary Thumbprint of the Certificate.

* `x509_store_name` - (Required) The X509 Store where the Certificate Exists, such as `My`.

---

A `reverse_proxy_certificate_common_names` block supports the following:

* `common_names` - (Required) A `common_names` block as defined below.

* `x509_store_name` - (Required) The X509 Store where the Certificate Exists, such as `My`.

---

A `client_certificate_thumbprint` block supports the following:

* `thumbprint` - (Required) The Thumbprint associated with the Client Certificate.

* `is_admin` - (Required) Does the Client Certificate have Admin Access to the cluster? Non-admin clients can only perform read only operations on the cluster.

---

A `client_certificate_common_name` block supports the following:

* `common_name` - (Required) The common or subject name of the certificate.

* `certificate_issuer_thumbprint` - (Optional) The Issuer Thumbprint of the Certificate.

-> **NOTE:** Certificate Issuer Thumbprint may become required in the future, `https://docs.microsoft.com/en-us/azure/service-fabric/service-fabric-create-cluster-using-cert-cn#download-and-update-a-sample-template`.

* `is_admin` - (Required) Does the Client Certificate have Admin Access to the cluster? Non-admin clients can only perform read only operations on the cluster.

---

A `diagnostics_config` block supports the following:

* `storage_account_name` - (Required) The name of the Storage Account where the Diagnostics should be sent to.

* `protected_account_key_name` - (Required) The protected diagnostics storage key name, such as `StorageAccountKey1`.

* `blob_endpoint` - (Required) The Blob Endpoint of the Storage Account.

* `queue_endpoint` - (Required) The Queue Endpoint of the Storage Account.

* `table_endpoint` - (Required) The Table Endpoint of the Storage Account.

---

A `fabric_settings` block supports the following:

* `name` - (Required) The name of the Fabric Setting, such as `Security` or `Federation`.

* `parameters` - (Optional) A map containing settings for the specified Fabric Setting.

---

A `node_type` block supports the following:

* `name` - (Required) The name of the Node Type. Changing this forces a new resource to be created.

* `placement_properties` - (Optional) The placement tags applied to nodes in the node type, which can be used to indicate where certain services (workload) should run.

* `capacities` - (Optional) The capacity tags applied to the nodes in the node type, the cluster resource manager uses these tags to understand how much resource a node has.

* `instance_count` - (Required) The number of nodes for this Node Type.

* `is_primary` - (Required) Is this the Primary Node Type? Changing this forces a new resource to be created.

* `client_endpoint_port` - (Required) The Port used for the Client Endpoint for this Node Type. Changing this forces a new resource to be created.

* `http_endpoint_port` - (Required) The Port used for the HTTP Endpoint for this Node Type. Changing this forces a new resource to be created.

* `durability_level` - (Optional) The Durability Level for this Node Type. Possible values include `Bronze`, `Gold` and `Silver`. Defaults to `Bronze`. Changing this forces a new resource to be created.

* `application_ports` - (Optional) A `application_ports` block as defined below.

* `ephemeral_ports` - (Optional) A `ephemeral_ports` block as defined below.

* `reverse_proxy_endpoint_port` - (Optional) The Port used for the Reverse Proxy Endpoint  for this Node Type. Changing this will upgrade the cluster.

---

A `application_ports` block supports the following:

* `start_port` - (Required) The start of the Application Port Range on this Node Type.

* `end_port` - (Required) The end of the Application Port Range on this Node Type.

---

A `ephemeral_ports` block supports the following:

* `start_port` - (Required) The start of the Ephemeral Port Range on this Node Type.

* `end_port` - (Required) The end of the Ephemeral Port Range on this Node Type.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Service Fabric Cluster.

* `cluster_endpoint` - The Cluster Endpoint for this Service Fabric Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Fabric Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Service Fabric Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Fabric Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Fabric Cluster.

## Import

Service Fabric Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceFabric/clusters/cluster1
```
