---
subcategory: "Hybrid Kubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_connected_cluster"
description: |-
  Manages a Connected Cluster.
---

# azurerm_connected_cluster

Manages a Connected Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_connected_cluster" "example" {
  name                         = "example-cluster"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = "MIICYzCCAcygAwIBAgIBADANBgkqhkiG9w0BAQUFADAuMQswCQYDVQQGEwJVUzEMMAoGA1UEChMDSUJNMREwDwYDVQQLEwhMb2NhbCBDQTAeFw05OTEyMjIwNTAwMDBaFw0wMDEyMjMwNDU5NTlaMC4xCzAJBgNVBAYTAlVTMQwwCgYDVQQKEwNJQk0xETAPBgNVBAsTCExvY2FsIENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD2bZEo7xGaX2"

  identity {
    type = "SystemAssigned"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Connected Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Connected Cluster exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Connected Cluster exists. Changing this forces a new resource to be created.

* `identity` - (Required) Managed identities which should be assigned to the azure resource. An `identity` block is defined below.

---

* `type` - (Required) The Type of Identity which should be used for this azure resource. Currently, the only possible value is `SystemAssigned`.

---


* `agent_public_key_certificate` - (Required) Base64 encoded public certificate used by the agent to do the initial handshake to the backend services in Azure. Changing this forces a new resource to be created.

* `distribution` - (Optional) The Kubernetes distribution running on this connected cluster. Changing this forces a new resource to be created.

* `infrastructure` - (Optional) The infrastructure on which the Kubernetes cluster represented by this connected cluster is running on. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this Connected Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Connected Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Connected Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Connected Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Connected Cluster.

## Import

Connected Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_connected_cluster.example /subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/example-group/providers/Microsoft.Kubernetes/connectedClusters/example-cluster
```
