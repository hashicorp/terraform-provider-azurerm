---
subcategory: "Red Hat Openshift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_openshift_cluster"
description: |-
  Manages fully managed Azure Red Hat Openshift Cluster (also known as ARO)
---

# azurerm_redhatopenshift_cluster

Manages a fully managed Azure Red Hat Openshift Cluster (also known as ARO).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "aro-vnet"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azuread_service_principal" "redhatopenshift" {
  // The is the `Azure Red Hat OpenShift RP` Service Principal ID
  application_id = "00000-0000-000000"
  use_existing   = true
}

resource "azurerm_role_assignment" "redhatopenshift" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.redhatopenshift.id
}

resource "azurerm_subnet" "main_subnet" {
  name                                           = "main-subnet"
  resource_group_name                            = azurerm_resource_group.example.name
  virtual_network_name                           = azurerm_virtual_network.example.name
  address_prefixes                               = ["10.0.0.0/23"]
  service_endpoints                              = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]
  enforce_private_link_service_network_policies  = true
  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]
}

resource "azurerm_redhat_openshift_cluster" "example" {
  name                = "example-redhatopenshift1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  cluster_profile {
    domain = "foo.example.com"
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    subnet_id    = azurerm_subnet.worker_subnet.id
    node_count   = 3
  }
  
  api_server_profile {
    visibility = "Public"
  }

  ingress_profile {
    visibility = "Public"
  }

  tags = {
    Environment = "Production"
  }
  
  depends_on = ["azurerm_role_assignment.redhatopenshift"]
}

output "console_url" {
  value = azurerm_redhatopenshift_cluster.example.console_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Red Hat Openshift Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Azure Red Hat Openshift Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Azure Red Hat Openshift Cluster should exist. Changing this forces a new resource to be created.

* `service_principal` - (Required) A `service_principal` block as defined below.

* `main_profile` - (Required) A `main_profile` block as defined below.

* `worker_profile` - (Required) A `worker_profile` block as defined below.

* `cluster_profile` - (Required) A `cluster_profile` block as defined below.

* `api_server_profile` - (Required) An `api_server_profile` block as defined below.

* `ingress_profile` - (Required) An `ingress_profile` block as defined below.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal. Changing this forces a new resource to be created.

* `client_secret` - (Required) The Client Secret for the Service Principal. Changing this forces a new resource to be created.

~> **Note:** Currently a service principal cannot be associated with more than one ARO clusters on the Azure subscription.

---

A `main_profile` block supports the following:

* `subnet_id` - (Required) The ID of the subnet where main nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the main nodes. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether main virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

---

A `worker_profile` block supports the following:

* `subnet_id` - (Required) The ID of the subnet where worker nodes will be hosted. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machines for the worker nodes. Changing this forces a new resource to be created.

* `disk_size_gb` - (Required) The internal OS disk size of the worker Virtual Machines in GB. Changing this forces a new resource to be created.

* `node_count` - (Required) The initial number of worker nodes which should exist in the cluster. Changing this forces a new resource to be created.

* `encryption_at_host_enabled` - (Optional) Whether worker virtual machines are encrypted at host. Defaults to `false`. Changing this forces a new resource to be created.

* `disk_encryption_set_id` - (Optional) The resource ID of an associated disk encryption set. Changing this forces a new resource to be created.

---

A `cluster_profile` block supports the following:

* `domain` - (Required) The custom domain for the cluster. For more info, see [Prepare a custom domain for your cluster](https://docs.microsoft.com/azure/openshift/tutorial-create-cluster#prepare-a-custom-domain-for-your-cluster-optional).  Changing this forces a new resource to be created.

* `pull_secret` - (Optional) The Red Hat pull secret for the cluster. Changing this forces a new resource to be created.

* `fips_enabled` - (Optional) Whether Federal Information Processing Standard (FIPS) validated cryptographic modules are used. Defaults to `false`. Changing this forces a new resource to be created.

---

A `network_profile` block supports the following:

* `pod_cidr` - (Optional) The CIDR to use for pod IP addresses. Defaults to `10.128.0.0/1`. Changing this forces a new resource to be created.

* `service_cidr` - (Optional) The network range used by the Openshift service. Defaults to `172.30.0.0/16`. Changing this forces a new resource to be created.

---

A `api_server_profile` block supports the following:

* `visibility` - (Optional) Cluster API server visibility. Supported values are `Public` and `Private`. Defaults to `Public`. Changing this forces a new resource to be created.

---

A `ingress_profile` block supports the following:

* `visibility` - (Optional) Cluster Ingress visibility. Supported values are `Public` and `Private`. Defaults to `Public`. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `version` - The Red Hat Openshift cluster version.

* `console_url` - The Red Hat Openshift cluster console URL.

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
* `update` - (Defaults to 90 minutes) Used when updating the Red Hat OpenShift cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Red Hat OpenShift cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Red Hat OpenShift cluster.

## Import

Red Hat Openshift Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhatopenshift_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.RedHatOpenShift/openShiftClusters/cluster1
```
