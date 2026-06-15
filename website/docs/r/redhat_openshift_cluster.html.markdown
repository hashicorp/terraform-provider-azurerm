---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_openshift_cluster"
description: |-
  Manages a fully managed Azure Red Hat OpenShift Cluster (also known as ARO).
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
    azurerm_role_assignment.role_network1,
    azurerm_role_assignment.role_network2,
  ]
}

output "console_url" {
  value = azurerm_redhat_openshift_cluster.example.console_url
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Red Hat OpenShift Cluster to create. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Azure Red Hat OpenShift Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Azure Red Hat OpenShift Cluster should be created. Changing this forces a new resource to be created.

* `api_server_profile` - (Required) An `api_server_profile` block as defined below. Changing this forces a new resource to be created.

* `cluster_profile` - (Required) A `cluster_profile` block as defined below. Changing this forces a new resource to be created.

* `ingress_profile` - (Required) An `ingress_profile` block as defined below. Changing this forces a new resource to be created.

* `main_profile` - (Required) A `main_profile` block as defined below. Changing this forces a new resource to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `worker_profile` - (Required) A `worker_profile` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below. Exactly one of `identity` or `service_principal` must be specified. Required when `platform_workload_identity_profile` is set.

* `platform_workload_identity_profile` - (Optional) A `platform_workload_identity_profile` block as defined below. Required when `identity` is set.

* `service_principal` - (Optional) A `service_principal` block as defined below. Exactly one of `service_principal` or `identity` must be specified.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `api_server_profile` block supports the following:

* `visibility` - (Required) Cluster API server visibility. Possible values are `Public` and `Private`. Changing this forces a new resource to be created.

---

A `cluster_profile` block supports the following:

* `domain` - (Required) The custom domain for the cluster. For more info, see [Prepare a custom domain for your cluster](https://docs.microsoft.com/azure/openshift/tutorial-create-cluster#prepare-a-custom-domain-for-your-cluster-optional). Changing this forces a new resource to be created.

* `fips_enabled` - (Optional) Whether Federal Information Processing Standard (FIPS) validated cryptographic modules are used. Defaults to `false`. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of a Resource Group which will be created to host VMs of Azure Red Hat OpenShift Cluster. The value cannot contain uppercase characters. Changing this forces a new resource to be created.

* `pull_secret` - (Optional) The Red Hat pull secret for the cluster. For more info, see [Get a Red Hat pull secret](https://learn.microsoft.com/azure/openshift/tutorial-create-cluster#get-a-red-hat-pull-secret-optional). Changing this forces a new resource to be created.

* `version` - (Optional) The version of the OpenShift cluster. When omitted, the Azure Red Hat OpenShift Resource Provider selects its current default supported version. Available versions can be found with the Azure CLI command `az aro get-versions --location <region>`. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `identity_ids` - (Required) A set of User Assigned Managed Identity IDs to assign to the cluster. Exactly one cluster identity is required for now.

* `type` - (Required) The type of Managed Service Identity assigned to the cluster. The only supported value is `UserAssigned`.

---

A `ingress_profile` block supports the following:

* `visibility` - (Required) Cluster Ingress visibility. Possible values are `Public` and `Private`. Changing this forces a new resource to be created.

---

A `load_balancer_profile` block supports the following:

* `managed_outbound_ip_count` - (Optional) The desired number of managed outbound public IP addresses created and assigned to the cluster's outbound load balancer. Possible values range between `1` and `20`. Defaults to `1`.

---

A `main_profile` block supports the following:

* `subnet_id` - (Required) The ID of the subnet where main nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the main nodes. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether main virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** `encryption_at_host_enabled` is only available for certain VM sizes and the `EncryptionAtHost` feature must be enabled for your subscription. Please see the [Azure documentation](https://learn.microsoft.com/azure/virtual-machines/disks-enable-host-based-encryption-portal?tabs=azure-powershell) for more information.

---

A `network_profile` block supports the following:

* `pod_cidr` - (Required) The CIDR to use for pod IP addresses. Changing this forces a new resource to be created.

* `service_cidr` - (Required) The network range used by the OpenShift service. Changing this forces a new resource to be created.

* `load_balancer_profile` - (Optional) A `load_balancer_profile` block as defined below. Only applicable when `outbound_type` is set to `Loadbalancer`.

* `outbound_type` - (Optional) The outbound (egress) routing method. Possible values are `Loadbalancer` and `UserDefinedRouting`. Defaults to `Loadbalancer`. Changing this forces a new resource to be created.

* `preconfigured_network_security_group_enabled` - (Optional) Whether a preconfigured network security group is being used on the subnets. Defaults to `false`. Changing this forces a new resource to be created.

---

A `platform_workload_identity` block supports the following:

* `identity_id` - (Required) The resource ID of the User Assigned Managed Identity to assign to the operator.

* `name` - (Required) The name of the platform workload identity operator.

~> **Note:** The required operator names vary by OpenShift minor version. For 4.19 they are `aro-operator`, `cloud-controller-manager`, `cloud-network-config`, `disk-csi-driver`, `file-csi-driver`, `image-registry`, `ingress`, and `machine-api`. See [Understand managed identities in Azure Red Hat OpenShift](https://learn.microsoft.com/azure/openshift/howto-understand-managed-identities#understand-identity-role-assignment-architecture) for more information.

---

A `platform_workload_identity_profile` block supports the following:

* `platform_workload_identity` - (Required) One or more `platform_workload_identity` blocks as defined below.

* `upgradeable_to` - (Optional) The target OpenShift version (`x.y.z`) the platform workload identities should satisfy.

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.

* `client_secret` - (Required) The Client Secret for the Service Principal.

~> **Note:** Currently a service principal cannot be associated with more than one ARO clusters on the Azure subscription.

---

A `worker_profile` block supports the following:

* `disk_size_gb` - (Required) The internal OS disk size of the worker Virtual Machines in GB. Changing this forces a new resource to be created.

* `node_count` - (Required) The initial number of worker nodes which should exist in the cluster. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet where worker nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the worker nodes. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether worker virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** `encryption_at_host_enabled` is only available for certain VM sizes and the `EncryptionAtHost` feature must be enabled for your subscription. Please see the [Azure documentation](https://learn.microsoft.com/azure/virtual-machines/disks-enable-host-based-encryption-portal?tabs=azure-powershell) for more information.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Red Hat OpenShift Cluster.

* `api_server_profile` - An `api_server_profile` block as defined below.

* `cluster_profile` - A `cluster_profile` block as defined below.

* `console_url` - The console URL of the Azure Red Hat OpenShift Cluster.

* `ingress_profile` - An `ingress_profile` block as defined below.

---

A `api_server_profile` block exports the following:

* `ip_address` - The IP Address the API Server Profile is associated with.

* `url` - The URL the API Server Profile is associated with.

---

A `cluster_profile` block exports the following:

* `resource_group_id` - The Resource Group that the cluster profile is attached to.

---

A `ingress_profile` block exports the following:

* `ip_address` - The IP Address the Ingress Profile is associated with.

* `name` - The name of the Ingress Profile.

---

A `load_balancer_profile` block exports the following:

* `effective_outbound_ips` - The list of effective outbound IP resource IDs of the cluster's public load balancer.

---

A `platform_workload_identity` block exports the following:

* `client_id` - The client ID of the User Assigned Managed Identity assigned to the operator.

* `object_id` - The object (principal) ID of the User Assigned Managed Identity assigned to the operator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Azure Red Hat OpenShift Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Red Hat OpenShift Cluster.
* `update` - (Defaults to 90 minutes) Used when updating the Azure Red Hat OpenShift Cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Azure Red Hat OpenShift Cluster.

## Import

An Azure Red Hat OpenShift Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhat_openshift_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.RedHatOpenShift/openShiftClusters/openShiftCluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.RedHatOpenShift` - 2025-07-25
