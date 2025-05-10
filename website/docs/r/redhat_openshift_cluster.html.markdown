---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_openshift_cluster"
description: |-
  Manages fully managed Azure Red Hat OpenShift Cluster (also known as ARO)
---

# azurerm_redhat_openshift_cluster

Manages a fully managed Azure Red Hat OpenShift Cluster (also known as ARO).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_client_config" "example" {}

data "azuread_client_config" "example" {}

resource "azuread_application" "example" {
  display_name = "example-aro"
}

resource "azuread_service_principal" "example" {
  client_id = azuread_application.example.client_id
}

resource "azuread_service_principal_password" "example" {
  service_principal_id = azuread_service_principal.example.object_id
}

data "azuread_service_principal" "redhatopenshift" {
  // This is the Azure Red Hat OpenShift RP service principal id, do NOT delete it
  client_id = "f1dd0a37-89c6-4e07-bcd1-ffd3d43d8875"
}

resource "azurerm_role_assignment" "role_network1" {
  scope                = azurerm_virtual_network.example.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "role_network2" {
  scope                = azurerm_virtual_network.example.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.redhatopenshift.object_id
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "main_subnet" {
  name                 = "main-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/23"]
  service_endpoints    = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]
}

resource "azurerm_redhat_openshift_cluster" "example" {
  name                = "examplearo"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  cluster_profile {
    domain  = "aro-example.com"
    version = "4.13.23"
  }

  network_profile {
    pod_cidr     = "10.128.0.0/14"
    service_cidr = "172.30.0.0/16"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  api_server_profile {
    visibility = "Public"
  }

  ingress_profile {
    visibility = "Public"
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    node_count   = 3
    subnet_id    = azurerm_subnet.worker_subnet.id
  }

  service_principal {
    client_id     = azuread_application.example.client_id
    client_secret = azuread_service_principal_password.example.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}

output "console_url" {
  value = azurerm_redhat_openshift_cluster.example.console_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Red Hat OpenShift Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Azure Red Hat OpenShift Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Azure Red Hat OpenShift Cluster should exist. Changing this forces a new resource to be created.

* `service_principal` - (Required) A `service_principal` block as defined below.

* `main_profile` - (Required) A `main_profile` block as defined below. Changing this forces a new resource to be created.

* `worker_profile` - (Required) A `worker_profile` block as defined below. Changing this forces a new resource to be created.

* `cluster_profile` - (Required) A `cluster_profile` block as defined below. Changing this forces a new resource to be created.

* `api_server_profile` - (Required) An `api_server_profile` block as defined below. Changing this forces a new resource to be created.

* `ingress_profile` - (Required) An `ingress_profile` block as defined below. Changing this forces a new resource to be created.

* `network_profile` - (Required) A `network_profile` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.

* `client_secret` - (Required) The Client Secret for the Service Principal.

~> **Note:** Currently a service principal cannot be associated with more than one ARO clusters on the Azure subscription.

---

A `main_profile` block supports the following:

* `subnet_id` - (Required) The ID of the subnet where main nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the main nodes. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether main virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** `encryption_at_host_enabled` is only available for certain VM sizes and the `EncryptionAtHost` feature must be enabled for your subscription. Please see the [Azure documentation](https://learn.microsoft.com/azure/virtual-machines/disks-enable-host-based-encryption-portal?tabs=azure-powershell) for more information.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

---

A `worker_profile` block supports the following:

* `subnet_id` - (Required) The ID of the subnet where worker nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the worker nodes. Changing this forces a new resource to be created.

* `disk_size_gb` - (Required) The internal OS disk size of the worker Virtual Machines in GB. Changing this forces a new resource to be created.

* `node_count` - (Required) The initial number of worker nodes which should exist in the cluster. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether worker virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** `encryption_at_host_enabled` is only available for certain VM sizes and the `EncryptionAtHost` feature must be enabled for your subscription. Please see the [Azure documentation](https://learn.microsoft.com/azure/virtual-machines/disks-enable-host-based-encryption-portal?tabs=azure-powershell) for more information.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

---

A `cluster_profile` block supports the following:

* `version` - (Required) The version of the OpenShift cluster. Available versions can be found with the Azure CLI command `az aro get-versions --location <region>`. Changing this forces a new resource to be created.

* `domain` - (Required) The custom domain for the cluster. For more info, see [Prepare a custom domain for your cluster](https://docs.microsoft.com/azure/openshift/tutorial-create-cluster#prepare-a-custom-domain-for-your-cluster-optional). Changing this forces a new resource to be created.

* `pull_secret` - (Optional) The Red Hat pull secret for the cluster. For more info, see [Get a Red Hat pull secret](https://learn.microsoft.com/azure/openshift/tutorial-create-cluster#get-a-red-hat-pull-secret-optional). Changing this forces a new resource to be created.

* `fips_enabled` - (Optional) Whether Federal Information Processing Standard (FIPS) validated cryptographic modules are used. Defaults to `false`. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of a Resource Group which will be created to host VMs of Azure Red Hat OpenShift Cluster. The value cannot contain uppercase characters. Defaults to `aro-{domain}`. Changing this forces a new resource to be created.

---

A `network_profile` block supports the following:

* `pod_cidr` - (Required) The CIDR to use for pod IP addresses. Changing this forces a new resource to be created.

* `service_cidr` - (Required) The network range used by the OpenShift service. Changing this forces a new resource to be created.

* `outbound_type` - (Optional) The outbound (egress) routing method. Possible values are `Loadbalancer` and `UserDefinedRouting`. Defaults to `Loadbalancer`. Changing this forces a new resource to be created.

* `preconfigured_network_security_group_enabled` - (Optional) Whether a preconfigured network security group is being used on the subnets.  Defaults to `false`.  Changing this forces a new resource to be created.

---

A `api_server_profile` block supports the following:

* `visibility` - (Required) Cluster API server visibility. Supported values are `Public` and `Private`. Changing this forces a new resource to be created.

---

A `ingress_profile` block supports the following:

* `visibility` - (Required) Cluster Ingress visibility. Supported values are `Public` and `Private`. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `console_url` - The Red Hat OpenShift cluster console URL.

* `cluster_profile` - A `cluster_profile` block as defined below.

* `api_server_profile` - An `api_server_profile` block as defined below.

* `ingress_profile` - An `ingress_profile` block as defined below.

---

A `cluster_profile` block exports the following:

* `resource_group_id` - The resource group that the cluster profile is attached to.

---

A `api_server_profile` block exports the following:

* `ip_address` - The IP Address the API Server Profile is associated with.

* `url` - The URL the API Server Profile is associated with.

---

A `ingress_profile` block exports the following:

* `name` - The name of the Ingress Profile.

* `ip_address` - The IP Address the Ingress Profile is associated with.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Red Hat OpenShift cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Red Hat OpenShift cluster.
* `update` - (Defaults to 90 minutes) Used when updating the Red Hat OpenShift cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Red Hat OpenShift cluster.

## Import

Red Hat OpenShift Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhat_openshift_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.RedHatOpenShift/openShiftClusters/cluster1
```
