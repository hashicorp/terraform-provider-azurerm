---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
description: |-
  Manages a managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service)
---

# azurerm_kubernetes_cluster

Manages a Managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service)

!> **Note:** As per Microsoft's AKS preview API [deprecation plan](https://learn.microsoft.com/en-us/azure/aks/concepts-preview-api-life-cycle#upcoming-deprecations) several preview APIs have a deprecation schedule and Microsoft recommends performing updates before the deprecation date. Additionally, Microsoft and HashiCorp recommend upgrading to the penultimate 3.x version v3.116.0 to avoid disruption or, ideally, to the latest 4.x provider version to take advantage of the most current API version that the provider supports. Please see [this GitHub issue](https://github.com/hashicorp/terraform-provider-azurerm/issues/28707) for more details.

-> **Note:** Due to the fast-moving nature of AKS, we recommend using the latest version of the Azure Provider when using AKS - you can find [the latest version of the Azure Provider here](https://registry.terraform.io/providers/hashicorp/azurerm/latest).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

This example provisions a basic Managed Kubernetes Cluster. Other examples of the `azurerm_kubernetes_cluster` resource can be found in [the `./examples/kubernetes` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/kubernetes).

An example of how to attach a specific Container Registry to a Managed Kubernetes Cluster can be found in the docs for [azurerm_container_registry](container_registry.html).

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks1"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Production"
  }
}

output "client_certificate" {
  value     = azurerm_kubernetes_cluster.example.kube_config[0].client_certificate
  sensitive = true
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.example.kube_config_raw

  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Kubernetes Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Managed Kubernetes Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Managed Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `default_node_pool` - (Required) Specifies configuration for "System" mode node pool. A `default_node_pool` block as defined below.

* `dns_prefix` - (Optional) DNS prefix specified when creating the managed cluster. Possible values must begin and end with a letter or number, contain only letters, numbers, and hyphens and be between 1 and 54 characters in length. Changing this forces a new resource to be created.

* `dns_prefix_private_cluster` - (Optional) Specifies the DNS prefix to use with private clusters. Changing this forces a new resource to be created.

-> **Note:** You must define either a `dns_prefix` or a `dns_prefix_private_cluster` field.

In addition, one of either `identity` or `service_principal` blocks must be specified.

---

* `aci_connector_linux` - (Optional) A `aci_connector_linux` block as defined below. For more details, please visit [Create and configure an AKS cluster to use virtual nodes](https://docs.microsoft.com/azure/aks/virtual-nodes-portal).

* `automatic_upgrade_channel` - (Optional) The upgrade channel for this Kubernetes Cluster. Possible values are `patch`, `rapid`, `node-image` and `stable`. Omitting this field sets this value to `none`.

!> **Note:** Cluster Auto-Upgrade will update the Kubernetes Cluster (and its Node Pools) to the latest GA version of Kubernetes automatically - please [see the Azure documentation for more information](https://docs.microsoft.com/azure/aks/upgrade-cluster#set-auto-upgrade-channel).

-> **Note:** Cluster Auto-Upgrade only updates to GA versions of Kubernetes and will not update to Preview versions.

* `api_server_access_profile` - (Optional) An `api_server_access_profile` block as defined below.

* `auto_scaler_profile` - (Optional) A `auto_scaler_profile` block as defined below.

* `azure_active_directory_role_based_access_control` - (Optional) A `azure_active_directory_role_based_access_control` block as defined below.

* `azure_policy_enabled` - (Optional) Should the Azure Policy Add-On be enabled? For more details please visit [Understand Azure Policy for Azure Kubernetes Service](https://docs.microsoft.com/en-ie/azure/governance/policy/concepts/rego-for-aks)

* `confidential_computing` - (Optional) A `confidential_computing` block as defined below. For more details please [the documentation](https://learn.microsoft.com/en-us/azure/confidential-computing/confidential-nodes-aks-overview)

* `cost_analysis_enabled` - (Optional) Should cost analysis be enabled for this Kubernetes Cluster? Defaults to `false`. The `sku_tier` must be set to `Standard` or `Premium` to enable this feature. Enabling this will add Kubernetes Namespace and Deployment details to the Cost Analysis views in the Azure portal.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used for the Nodes and Volumes. More information [can be found in the documentation](https://docs.microsoft.com/azure/aks/azure-disk-customer-managed-keys). Changing this forces a new resource to be created.

* `edge_zone` - (Optional) Specifies the Extended Zone (formerly called Edge Zone) within the Azure Region where this Managed Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `http_application_routing_enabled` - (Optional) Should HTTP Application Routing be enabled?

-> **Note:** At this time HTTP Application Routing is not supported in Azure China or Azure US Government.

* `http_proxy_config` - (Optional) A `http_proxy_config` block as defined below.

* `identity` - (Optional) An `identity` block as defined below. One of either `identity` or `service_principal` must be specified.

!> **Note:** A migration scenario from `service_principal` to `identity` is supported. When upgrading `service_principal` to `identity`, your cluster's control plane and addon pods will switch to use managed identity, but the kubelets will keep using your configured `service_principal` until you upgrade your Node Pool.

* `image_cleaner_enabled` - (Optional) Specifies whether Image Cleaner is enabled.

* `image_cleaner_interval_hours` - (Optional) Specifies the interval in hours when images should be cleaned up. Defaults to `0`.

* `ingress_application_gateway` - (Optional) A `ingress_application_gateway` block as defined below.

-> **Note:** Since the Application Gateway is deployed inside a Virtual Network, users (and Service Principals) that are operating the Application Gateway must have the `Microsoft.Network/virtualNetworks/subnets/join/action` permission on the Virtual Network or Subnet. For more details, please visit [Virtual Network Permission](https://learn.microsoft.com/en-us/azure/application-gateway/configuration-infrastructure#virtual-network-permission).

* `key_management_service` - (Optional) A `key_management_service` block as defined below. For more details, please visit [Key Management Service (KMS) etcd encryption to an AKS cluster](https://learn.microsoft.com/en-us/azure/aks/use-kms-etcd-encryption).

* `key_vault_secrets_provider` - (Optional) A `key_vault_secrets_provider` block as defined below. For more details, please visit [Azure Keyvault Secrets Provider for AKS](https://docs.microsoft.com/azure/aks/csi-secrets-store-driver).

* `kubelet_identity` - (Optional) A `kubelet_identity` block as defined below.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade). AKS does not require an exact patch version to be specified, minor version aliases such as `1.22` are also supported. - The minor version's latest GA patch is automatically chosen in that case. More details can be found in [the documentation](https://docs.microsoft.com/en-us/azure/aks/supported-kubernetes-versions?tabs=azure-cli#alias-minor-version).

-> **Note:** Upgrading your cluster may take up to 10 minutes per node.

* `linux_profile` - (Optional) A `linux_profile` block as defined below.

* `local_account_disabled` - (Optional) If `true` local accounts will be disabled. See [the documentation](https://docs.microsoft.com/azure/aks/managed-aad#disable-local-accounts) for more information.

-> **Note:** If `local_account_disabled` is set to `true`, it is required to enable Kubernetes RBAC and AKS-managed Azure AD integration. See [the documentation](https://docs.microsoft.com/azure/aks/managed-aad#azure-ad-authentication-overview) for more information.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `maintenance_window_auto_upgrade` - (Optional) A `maintenance_window_auto_upgrade` block as defined below.

* `maintenance_window_node_os` - (Optional) A `maintenance_window_node_os` block as defined below.

* `microsoft_defender` - (Optional) A `microsoft_defender` block as defined below.

* `monitor_metrics` - (Optional) Specifies a Prometheus add-on profile for the Kubernetes Cluster. A `monitor_metrics` block as defined below.

-> **Note:** If deploying Managed Prometheus, the `monitor_metrics` properties are required to configure the cluster for metrics collection. If no value is needed, set properties to `null`.

* `network_profile` - (Optional) A `network_profile` block as defined below. Changing this forces a new resource to be created.

-> **Note:** If `network_profile` is not defined, `kubenet` profile will be used by default.

* `node_os_upgrade_channel` - (Optional) The upgrade channel for this Kubernetes Cluster Nodes' OS Image. Possible values are `Unmanaged`, `SecurityPatch`, `NodeImage` and `None`. Defaults to `NodeImage`.

-> **Note:** `node_os_upgrade_channel` must be set to `NodeImage` if `automatic_upgrade_channel` has been set to `node-image`

* `node_resource_group` - (Optional) The name of the Resource Group where the Kubernetes Nodes should exist. Changing this forces a new resource to be created.

-> **Note:** Azure requires that a new, non-existent Resource Group is used, as otherwise, the provisioning of the Kubernetes Service will fail.

* `oidc_issuer_enabled` - (Optional) Enable or Disable the [OIDC issuer URL](https://learn.microsoft.com/en-gb/azure/aks/use-oidc-issuer)

* `oms_agent` - (Optional) A `oms_agent` block as defined below.

* `open_service_mesh_enabled` - (Optional) Is Open Service Mesh enabled? For more details, please visit [Open Service Mesh for AKS](https://docs.microsoft.com/azure/aks/open-service-mesh-about).

* `private_cluster_enabled` - (Optional) Should this Kubernetes Cluster have its API server only exposed on internal IP addresses? This provides a Private IP Address for the Kubernetes API on the Virtual Network where the Kubernetes Cluster is located. Defaults to `false`. Changing this forces a new resource to be created.

* `private_dns_zone_id` - (Optional) Either the ID of Private DNS Zone which should be delegated to this Cluster, `System` to have AKS manage this or `None`. In case of `None` you will need to bring your own DNS server and set up resolving, otherwise, the cluster will have issues after provisioning. Changing this forces a new resource to be created.

* `private_cluster_public_fqdn_enabled` - (Optional) Specifies whether a Public FQDN for this Private Cluster should be added. Defaults to `false`.

-> **Note:** If you use BYO DNS Zone, the AKS cluster should either use a User Assigned Identity or a service principal (which is deprecated) with the `Private DNS Zone Contributor` role and access to this Private DNS Zone. If `UserAssigned` identity is used - to prevent improper resource order destruction - the cluster should depend on the role assignment, like in this example:

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "privatelink.eastus2.azmk8s.io"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "aks-example-identity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_private_dns_zone.example.id
  role_definition_name = "Private DNS Zone Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_kubernetes_cluster" "example" {
  name                    = "aksexamplewithprivatednszone1"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  dns_prefix              = "aksexamplednsprefix1"
  private_cluster_enabled = true
  private_dns_zone_id     = azurerm_private_dns_zone.example.id

  # rest of configuration omitted for brevity

  depends_on = [
    azurerm_role_assignment.example,
  ]
}

```

* `service_mesh_profile` - (Optional) A `service_mesh_profile` block as defined below.

* `workload_autoscaler_profile` - (Optional) A `workload_autoscaler_profile` block defined below.

* `workload_identity_enabled` - (Optional) Specifies whether Azure AD Workload Identity should be enabled for the Cluster. Defaults to `false`.

-> **Note:** To enable Azure AD Workload Identity `oidc_issuer_enabled` must be set to `true`.

-> **Note:** Enabling this option will allocate Workload Identity resources to the `kube-system` namespace in Kubernetes. If you wish to customize the deployment of Workload Identity, you can refer to [the documentation on Azure AD Workload Identity.](https://azure.github.io/azure-workload-identity/docs/installation/mutating-admission-webhook.html) The documentation provides guidance on how to install the mutating admission webhook, which allows for the customization of Workload Identity deployment.

* `role_based_access_control_enabled` - (Optional) Whether Role Based Access Control for the Kubernetes Cluster should be enabled. Defaults to `true`. Changing this forces a new resource to be created.

* `run_command_enabled` - (Optional) Whether to enable run command for the cluster or not. Defaults to `true`.

* `service_principal` - (Optional) A `service_principal` block as documented below. One of either `identity` or `service_principal` must be specified.

!> **Note:** A migration scenario from `service_principal` to `identity` is supported. When upgrading `service_principal` to `identity`, your cluster's control plane and addon pods will switch to use managed identity, but the kubelets will keep using your configured `service_principal` until you upgrade your Node Pool.

* `sku_tier` - (Optional) The SKU Tier that should be used for this Kubernetes Cluster. Possible values are `Free`, `Standard` (which includes the Uptime SLA) and `Premium`. Defaults to `Free`.

-> **Note:** Whilst the AKS API previously supported the `Paid` SKU - the AKS API introduced a breaking change in API Version `2023-02-01` (used in v3.51.0 and later) where the value `Paid` must now be set to `Standard`.

* `storage_profile` - (Optional) A `storage_profile` block as defined below.

* `support_plan` - (Optional) Specifies the support plan which should be used for this Kubernetes Cluster. Possible values are `KubernetesOfficial` and `AKSLongTermSupport`. Defaults to `KubernetesOfficial`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `upgrade_override` - (Optional) A `upgrade_override` block as defined below.

* `web_app_routing` - (Optional) A `web_app_routing` block as defined below.

* `windows_profile` - (Optional) A `windows_profile` block as defined below.

---

An `aci_connector_linux` block supports the following:

* `subnet_name` - (Required) The subnet name for the virtual nodes to run.

-> **Note:** At this time ACI Connectors are not supported in Azure China.

-> **Note:** AKS will add a delegation to the subnet named here. To prevent further runs from failing you should make sure that the subnet you create for virtual nodes has a delegation, like so.

```hcl
resource "azurerm_subnet" "virtual" {

  #...

  delegation {
    name = "aciDelegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
```

---

An `api_server_access_profile` block supports the following:

* `authorized_ip_ranges` - (Optional) Set of authorized IP ranges to allow access to API server, e.g. ["198.51.100.0/24"].

---

An `auto_scaler_profile` block supports the following:

* `balance_similar_node_groups` - (Optional) Detect similar node groups and balance the number of nodes between them. Defaults to `false`.

* `daemonset_eviction_for_empty_nodes_enabled` - (Optional) Whether DaemonSet pods will be gracefully terminated from empty nodes. Defaults to `false`.

* `daemonset_eviction_for_occupied_nodes_enabled` - (Optional) Whether DaemonSet pods will be gracefully terminated from non-empty nodes. Defaults to `true`.

* `expander` - (Optional) Expander to use. Possible values are `least-waste`, `priority`, `most-pods` and `random`. Defaults to `random`.

* `ignore_daemonsets_utilization_enabled` - (Optional) Whether DaemonSet pods will be ignored when calculating resource utilization for scale down. Defaults to `false`.

* `max_graceful_termination_sec` - (Optional) Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node. Defaults to `600`.

* `max_node_provisioning_time` - (Optional) Maximum time the autoscaler waits for a node to be provisioned. Defaults to `15m`.

* `max_unready_nodes` - (Optional) Maximum Number of allowed unready nodes. Defaults to `3`.

* `max_unready_percentage` - (Optional) Maximum percentage of unready nodes the cluster autoscaler will stop if the percentage is exceeded. Defaults to `45`.

* `new_pod_scale_up_delay` - (Optional) For scenarios like burst/batch scale where you don't want CA to act before the kubernetes scheduler could schedule all the pods, you can tell CA to ignore unscheduled pods before they're a certain age. Defaults to `10s`.

* `scale_down_delay_after_add` - (Optional) How long after the scale up of AKS nodes the scale down evaluation resumes. Defaults to `10m`.

* `scale_down_delay_after_delete` - (Optional) How long after node deletion that scale down evaluation resumes. Defaults to the value used for `scan_interval`.

* `scale_down_delay_after_failure` - (Optional) How long after scale down failure that scale down evaluation resumes. Defaults to `3m`.

* `scan_interval` - (Optional) How often the AKS Cluster should be re-evaluated for scale up/down. Defaults to `10s`.

* `scale_down_unneeded` - (Optional) How long a node should be unneeded before it is eligible for scale down. Defaults to `10m`.

* `scale_down_unready` - (Optional) How long an unready node should be unneeded before it is eligible for scale down. Defaults to `20m`.

* `scale_down_utilization_threshold` - (Optional) Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down. Defaults to `0.5`.

* `empty_bulk_delete_max` - (Optional) Maximum number of empty nodes that can be deleted at the same time. Defaults to `10`.

* `skip_nodes_with_local_storage` - (Optional) If `true` cluster autoscaler will never delete nodes with pods with local storage, for example, EmptyDir or HostPath. Defaults to `true`.

* `skip_nodes_with_system_pods` - (Optional) If `true` cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet or mirror pods). Defaults to `true`.

---

An `azure_active_directory_role_based_access_control` block supports the following:

* `tenant_id` - (Optional) The Tenant ID used for Azure Active Directory Application. If this isn't specified the Tenant ID of the current Subscription is used.

* `admin_group_object_ids` - (Optional) A list of Object IDs of Azure Active Directory Groups which should have Admin Role on the Cluster.

* `azure_rbac_enabled` - (Optional) Is Role Based Access Control based on Azure AD enabled?

---

A `confidential_computing` block supports the following:

* `sgx_quote_helper_enabled` - (Required) Should the SGX quote helper be enabled?

---

An `monitor_metrics` block supports the following:

* `annotations_allowed` - (Optional) Specifies a comma-separated list of Kubernetes annotation keys that will be used in the resource's labels metric.

* `labels_allowed` - (Optional) Specifies a Comma-separated list of additional Kubernetes label keys that will be used in the resource's labels metric.

-> **Note:** Both properties `annotations_allowed` and `labels_allowed` are required if you are enabling Managed Prometheus with an existing Azure Monitor Workspace.

---

A `default_node_pool` block supports the following:

-> **Note:** Changing certain properties of the `default_node_pool` is done by cycling the system node pool of the cluster. When cycling the system node pool, it doesn't perform cordon and drain, and it will disrupt rescheduling pods currently running on the previous system node pool.`temporary_name_for_rotation` must be specified when changing any of the following properties: `host_encryption_enabled`, `node_public_ip_enabled`, `fips_enabled`, `kubelet_config`, `kubelet_disk_type`, `linux_os_config`, `max_pods`, `only_critical_addons_enabled`, `os_disk_size_gb`, `os_disk_type`, `os_sku`, `pod_subnet_id`, `snapshot_id`, `ultra_ssd_enabled`, `vnet_subnet_id`, `vm_size`, `zones`.

* `name` - (Required) The name which should be used for the default Kubernetes Node Pool.

* `vm_size` - (Required) The size of the Virtual Machine, such as `Standard_DS2_v2`. `temporary_name_for_rotation` must be specified when attempting a resize.

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the Capacity Reservation Group within which this AKS Cluster should be created. Changing this forces a new resource to be created.

* `auto_scaling_enabled` - (Optional) Should [the Kubernetes Auto Scaler](https://docs.microsoft.com/azure/aks/cluster-autoscaler) be enabled for this Node Pool?

-> **Note:** This requires that the `type` is set to `VirtualMachineScaleSets`.

-> **Note:** If you're using AutoScaling, you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore changes to the `node_count` field.

* `host_encryption_enabled` - (Optional) Should the nodes in the Default Node Pool have host encryption enabled? `temporary_name_for_rotation` must be specified when changing this property.

-> **Note:** This requires that the  Feature `Microsoft.ContainerService/EnableEncryptionAtHost` is enabled and the Resource Provider is registered.

* `node_public_ip_enabled` - (Optional) Should nodes in this Node Pool have a Public IP Address? `temporary_name_for_rotation` must be specified when changing this property.

* `gpu_instance` - (Optional) Specifies the GPU MIG instance profile for supported GPU VM SKU. The allowed values are `MIG1g`, `MIG2g`, `MIG3g`, `MIG4g` and `MIG7g`. Changing this forces a new resource to be created.

* `host_group_id` - (Optional) Specifies the ID of the Host Group within which this AKS Cluster should be created. Changing this forces a new resource to be created.

* `kubelet_config` - (Optional) A `kubelet_config` block as defined below. `temporary_name_for_rotation` must be specified when changing this block.

* `linux_os_config` - (Optional) A `linux_os_config` block as defined below. `temporary_name_for_rotation` must be specified when changing this block.

* `fips_enabled` - (Optional) Should the nodes in this Node Pool have Federal Information Processing Standard enabled? `temporary_name_for_rotation` must be specified when changing this block. Changing this forces a new resource to be created.

* `kubelet_disk_type` - (Optional) The type of disk used by kubelet. Possible values are `OS` and `Temporary`. `temporary_name_for_rotation` must be specified when changing this block.

* `max_pods` - (Optional) The maximum number of pods that can run on each agent. `temporary_name_for_rotation` must be specified when changing this property.

* `node_network_profile` - (Optional) A `node_network_profile` block as documented below.

* `node_public_ip_prefix_id` - (Optional) Resource ID for the Public IP Addresses Prefix for the nodes in this Node Pool. `node_public_ip_enabled` should be `true`. Changing this forces a new resource to be created.

* `node_labels` - (Optional) A map of Kubernetes labels which should be applied to nodes in the Default Node Pool.

* `only_critical_addons_enabled` - (Optional) Enabling this option will taint default node pool with `CriticalAddonsOnly=true:NoSchedule` taint. `temporary_name_for_rotation` must be specified when changing this property.

* `orchestrator_version` - (Optional) Version of Kubernetes used for the Agents. If not specified, the default node pool will be created with the version specified by `kubernetes_version`. If both are unspecified, the latest recommended version will be used at provisioning time (but won't auto-upgrade). AKS does not require an exact patch version to be specified, minor version aliases such as `1.22` are also supported. - The minor version's latest GA patch is automatically chosen in that case. More details can be found in [the documentation](https://docs.microsoft.com/en-us/azure/aks/supported-kubernetes-versions?tabs=azure-cli#alias-minor-version).

-> **Note:** This version must be supported by the Kubernetes Cluster - as such the version of Kubernetes used on the Cluster/Control Plane may need to be upgraded first.

* `os_disk_size_gb` - (Optional) The size of the OS Disk which should be used for each agent in the Node Pool. `temporary_name_for_rotation` must be specified when attempting a change.

* `os_disk_type` - (Optional) The type of disk which should be used for the Operating System. Possible values are `Ephemeral` and `Managed`. Defaults to `Managed`. `temporary_name_for_rotation` must be specified when attempting a change.

* `os_sku` - (Optional) Specifies the OS SKU used by the agent pool. Possible values are `AzureLinux`, `Ubuntu`, `Windows2019` and `Windows2022`. If not specified, the default is `Ubuntu` if OSType=Linux or `Windows2019` if OSType=Windows. And the default Windows OSSKU will be changed to `Windows2022` after Windows2019 is deprecated. Changing this from `AzureLinux` or `Ubuntu` to `AzureLinux` or `Ubuntu` will not replace the resource, otherwise `temporary_name_for_rotation` must be specified when attempting a change.

* `pod_subnet_id` - (Optional) The ID of the Subnet where the pods in the default Node Pool should exist.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group. Changing this forces a new resource to be created.

* `scale_down_mode` - (Optional) Specifies the autoscaling behaviour of the Kubernetes Cluster. Allowed values are `Delete` and `Deallocate`. Defaults to `Delete`.

* `snapshot_id` - (Optional) The ID of the Snapshot which should be used to create this default Node Pool. `temporary_name_for_rotation` must be specified when changing this property.

* `temporary_name_for_rotation` - (Optional) Specifies the name of the temporary node pool used to cycle the default node pool for VM resizing.

* `type` - (Optional) The type of Node Pool which should be created. Possible values are `VirtualMachineScaleSets`. Defaults to `VirtualMachineScaleSets`. Changing this forces a new resource to be created.

-> **Note:** When creating a cluster that supports multiple node pools, the cluster must use `VirtualMachineScaleSets`. For more information on the limitations of clusters using multiple node pools see [the documentation](https://learn.microsoft.com/en-us/azure/aks/use-multiple-node-pools#limitations).

* `tags` - (Optional) A mapping of tags to assign to the Node Pool.

~> **Note:** At this time there's a bug in the AKS API where Tags for a Node Pool are not stored in the correct case - you [may wish to use Terraform's `ignore_changes` functionality to ignore changes to the casing](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) until this is fixed in the AKS API.

* `ultra_ssd_enabled` - (Optional) Used to specify whether the UltraSSD is enabled in the Default Node Pool. Defaults to `false`. See [the documentation](https://docs.microsoft.com/azure/aks/use-ultra-disks) for more information. `temporary_name_for_rotation` must be specified when attempting a change.

* `upgrade_settings` - (Optional) A `upgrade_settings` block as documented below.

* `vnet_subnet_id` - (Optional) The ID of a Subnet where the Kubernetes Node Pool should exist.

~> **Note:** A Route Table must be configured on this Subnet.

* `workload_runtime` - (Optional) Specifies the workload runtime used by the node pool. Possible value is `OCIContainer`.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Kubernetes Cluster should be located. `temporary_name_for_rotation` must be specified when changing this property.

-> **Note:** This requires that the `type` is set to `VirtualMachineScaleSets` and that `load_balancer_sku` is set to `standard`.

If `auto_scaling_enabled` is set to `true`, then the following fields can also be configured:

* `max_count` - (Optional) The maximum number of nodes which should exist in this Node Pool. If specified this must be between `1` and `1000`.

* `min_count` - (Optional) The minimum number of nodes which should exist in this Node Pool. If specified this must be between `1` and `1000`.

* `node_count` - (Optional) The initial number of nodes which should exist in this Node Pool. If specified this must be between `1` and `1000` and between `min_count` and `max_count`.

-> **Note:** If specified you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) to ignore changes to this field.

-> **Note:** If `auto_scaling_enabled` is set to `false` both `min_count` and `max_count` fields need to be set to `null` or omitted from the configuration.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Kubernetes Cluster. Possible values are `SystemAssigned` or `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Kubernetes Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned`. Currently only one User Assigned Identity is supported.

---

A `key_management_service` block supports the following:

* `key_vault_key_id` - (Required) Identifier of Azure Key Vault key. See [key identifier format](https://learn.microsoft.com/en-us/azure/key-vault/general/about-keys-secrets-certificates#vault-name-and-object-name) for more details.

* `key_vault_network_access` - (Optional) Network access of the key vault Network access of key vault. The possible values are `Public` and `Private`. `Public` means the key vault allows public access from all networks. `Private` means the key vault disables public access and enables private link. Defaults to `Public`.

---

A `key_vault_secrets_provider` block supports the following:

* `secret_rotation_enabled` - (Optional) Should the secret store CSI driver on the AKS cluster be enabled?

* `secret_rotation_interval` - (Optional) The interval to poll for secret rotation. This attribute is only set when `secret_rotation_enabled` is true. Defaults to `2m`.

-> **Note:** To enable`key_vault_secrets_provider` either `secret_rotation_enabled` or `secret_rotation_interval` must be specified.

---

A `kubelet_config` block supports the following:

* `allowed_unsafe_sysctls` - (Optional) Specifies the allow list of unsafe sysctls command or patterns (ending in `*`).

* `container_log_max_line` - (Optional) Specifies the maximum number of container log files that can be present for a container. must be at least 2.

* `container_log_max_size_mb` - (Optional) Specifies the maximum size (e.g. 10MB) of container log file before it is rotated.

* `cpu_cfs_quota_enabled` - (Optional) Is CPU CFS quota enforcement for containers enabled? Defaults to `true`.

* `cpu_cfs_quota_period` - (Optional) Specifies the CPU CFS quota period value.

* `cpu_manager_policy` - (Optional) Specifies the CPU Manager policy to use. Possible values are `none` and `static`,.

* `image_gc_high_threshold` - (Optional) Specifies the percent of disk usage above which image garbage collection is always run. Must be between `0` and `100`.

* `image_gc_low_threshold` - (Optional) Specifies the percent of disk usage lower than which image garbage collection is never run. Must be between `0` and `100`.

* `pod_max_pid` - (Optional) Specifies the maximum number of processes per pod.

* `topology_manager_policy` - (Optional) Specifies the Topology Manager policy to use. Possible values are `none`, `best-effort`, `restricted` or `single-numa-node`.

---

The `kubelet_identity` block supports the following:

* `client_id` - (Optional) The Client ID of the user-defined Managed Identity to be assigned to the Kubelets. If not specified a Managed Identity is created automatically. Changing this forces a new resource to be created.

* `object_id` - (Optional) The Object ID of the user-defined Managed Identity assigned to the Kubelets.If not specified a Managed Identity is created automatically. Changing this forces a new resource to be created.

* `user_assigned_identity_id` - (Optional) The ID of the User Assigned Identity assigned to the Kubelets. If not specified a Managed Identity is created automatically. Changing this forces a new resource to be created.

-> **Note:** When `kubelet_identity` is enabled - The `type` field in the `identity` block must be set to `UserAssigned` and `identity_ids` must be set.

---

A `linux_os_config` block supports the following:

* `swap_file_size_mb` - (Optional) Specifies the size of the swap file on each node in MB.

* `sysctl_config` - (Optional) A `sysctl_config` block as defined below.

* `transparent_huge_page_defrag` - (Optional) specifies the defrag configuration for Transparent Huge Page. Possible values are `always`, `defer`, `defer+madvise`, `madvise` and `never`.

* `transparent_huge_page_enabled` - (Optional) Specifies the Transparent Huge Page enabled configuration. Possible values are `always`, `madvise` and `never`.

---

A `node_network_profile` block supports the following:

* `allowed_host_ports` - (Optional) One or more `allowed_host_ports` blocks as defined below.

* `application_security_group_ids` - (Optional) A list of Application Security Group IDs which should be associated with this Node Pool.

* `node_public_ip_tags` - (Optional) Specifies a mapping of tags to the instance-level public IPs. Changing this forces a new resource to be created.

---

An `allowed_host_ports` block supports the following:

* `port_start` - (Optional) Specifies the start of the port range.

* `port_end` - (Optional) Specifies the end of the port range.

* `protocol` - (Optional) Specifies the protocol of the port range. Possible values are `TCP` and `UDP`.

---

A `linux_profile` block supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster. Changing this forces a new resource to be created.

* `ssh_key` - (Required) An `ssh_key` block as defined below. Only one is currently allowed. Changing this will update the key on all node pools. More information can be found in [the documentation](https://learn.microsoft.com/en-us/azure/aks/node-access#update-ssh-key-on-an-existing-aks-cluster-preview).

---

A `maintenance_window` block supports the following:

* `allowed` - (Optional) One or more `allowed` blocks as defined below.

* `not_allowed` - (Optional) One or more `not_allowed` block as defined below.

---

A `maintenance_window_auto_upgrade` block supports the following:

* `frequency` - (Required) Frequency of maintenance. Possible options are `Weekly`, `AbsoluteMonthly` and `RelativeMonthly`.

* `interval` - (Required) The interval for maintenance runs. Depending on the frequency this interval is week or month based.

* `duration` - (Required) The duration of the window for maintenance to run in hours. Possible options are between `4` to `24`.

* `day_of_week` - (Optional) The day of the week for the maintenance run. Required in combination with weekly frequency. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday` and `Wednesday`.

* `day_of_month` - (Optional) The day of the month for the maintenance run. Required in combination with AbsoluteMonthly frequency. Value between 0 and 31 (inclusive).

* `week_index` - (Optional) Specifies on which instance of the allowed days specified in `day_of_week` the maintenance occurs. Options are `First`, `Second`, `Third`, `Fourth`, and `Last`.
 Required in combination with relative monthly frequency.

* `start_time` - (Optional) The time for maintenance to begin, based on the timezone determined by `utc_offset`. Format is `HH:mm`.

* `utc_offset` - (Optional) Used to determine the timezone for cluster maintenance.

* `start_date` - (Optional) The date on which the maintenance window begins to take effect.

* `not_allowed` - (Optional) One or more `not_allowed` block as defined below.

---

A `maintenance_window_node_os` block supports the following:

* `frequency` - (Required) Frequency of maintenance. Possible options are `Daily`, `Weekly`, `AbsoluteMonthly` and `RelativeMonthly`.

* `interval` - (Required) The interval for maintenance runs. Depending on the frequency this interval is week or month based.

* `duration` - (Required) The duration of the window for maintenance to run in hours. Possible options are between `4` to `24`.

* `day_of_week` - (Optional) The day of the week for the maintenance run. Required in combination with weekly frequency. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday` and `Wednesday`.

* `day_of_month` - (Optional) The day of the month for the maintenance run. Required in combination with AbsoluteMonthly frequency. Value between 0 and 31 (inclusive).

* `week_index` - (Optional) The week in the month used for the maintenance run. Options are `First`, `Second`, `Third`, `Fourth`, and `Last`.

* `start_time` - (Optional) The time for maintenance to begin, based on the timezone determined by `utc_offset`. Format is `HH:mm`.

* `utc_offset` - (Optional) Used to determine the timezone for cluster maintenance.

* `start_date` - (Optional) The date on which the maintenance window begins to take effect.

* `not_allowed` - (Optional) One or more `not_allowed` block as defined below.

---

An `allowed` block supports the following:

* `day` - (Required) A day in a week. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

* `hours` - (Required) An array of hour slots in a day. For example, specifying `1` will allow maintenance from 1:00am to 2:00am. Specifying `1`, `2` will allow maintenance from 1:00am to 3:00m. Possible values are between `0` and `23`.

---

A `not_allowed` block supports the following:

* `end` - (Required) The end of a time span, formatted as an RFC3339 string.

* `start` - (Required) The start of a time span, formatted as an RFC3339 string.

---

A `microsoft_defender` block supports the following:

* `log_analytics_workspace_id` - (Required) Specifies the ID of the Log Analytics Workspace where the audit logs collected by Microsoft Defender should be sent to.

---

A `network_profile` block supports the following:

* `network_plugin` - (Required) Network plugin to use for networking. Currently supported values are `azure`, `kubenet` and `none`. Changing this forces a new resource to be created.

-> **Note:** When `network_plugin` is set to `azure` - the `pod_cidr` field must not be set, unless specifying `network_plugin_mode` to `overlay`.

* `network_mode` - (Optional) Network mode to be used with Azure CNI. Possible values are `bridge` and `transparent`. Changing this forces a new resource to be created.

~> **Note:** `network_mode` can only be set to `bridge` for existing Kubernetes Clusters and cannot be used to provision new Clusters - this will be removed by Azure in the future.

~> **Note:** This property can only be set when `network_plugin` is set to `azure`.

* `network_policy` - (Optional) Sets up network policy to be used with Azure CNI. [Network policy allows us to control the traffic flow between pods](https://docs.microsoft.com/azure/aks/use-network-policies). Currently supported values are `calico`, `azure` and `cilium`.

~> **Note:** When `network_policy` is set to `azure`, the `network_plugin` field can only be set to `azure`.

~> **Note:** When `network_policy` is set to `cilium`, the `network_data_plane` field must be set to `cilium`.

* `dns_service_ip` - (Optional) IP address within the Kubernetes service address range that will be used by cluster service discovery (kube-dns). Changing this forces a new resource to be created.

* `network_data_plane` - (Optional) Specifies the data plane used for building the Kubernetes network. Possible values are `azure` and `cilium`. Defaults to `azure`. Disabling this forces a new resource to be created.

~> **Note:** When `network_data_plane` is set to `cilium`, the `network_plugin` field can only be set to `azure`.

~> **Note:** When `network_data_plane` is set to `cilium`, one of either `network_plugin_mode = "overlay"` or `pod_subnet_id` must be specified.

* `network_plugin_mode` - (Optional) Specifies the network plugin mode used for building the Kubernetes network. Possible value is `overlay`.

~> **Note:** When `network_plugin_mode` is set to `overlay`, the `network_plugin` field can only be set to `azure`. When upgrading from Azure CNI without overlay, `pod_subnet_id` must be specified.

* `outbound_type` - (Optional) The outbound (egress) routing method which should be used for this Kubernetes Cluster. Possible values are `loadBalancer`, `userDefinedRouting`, `managedNATGateway` and `userAssignedNATGateway`. Defaults to `loadBalancer`. More information on supported migration paths for `outbound_type` can be found in [this documentation](https://learn.microsoft.com/azure/aks/egress-outboundtype#updating-outboundtype-after-cluster-creation).

* `pod_cidr` - (Optional) The CIDR to use for pod IP addresses. This field can only be set when `network_plugin` is set to `kubenet` or `network_plugin_mode` is set to `overlay`. Changing this forces a new resource to be created.

* `pod_cidrs` - (Optional) A list of CIDRs to use for pod IP addresses. For single-stack networking a single IPv4 CIDR is expected. For dual-stack networking an IPv4 and IPv6 CIDR are expected. Changing this forces a new resource to be created.

* `service_cidr` - (Optional) The Network Range used by the Kubernetes service. Changing this forces a new resource to be created.

* `service_cidrs` - (Optional) A list of CIDRs to use for Kubernetes services. For single-stack networking a single IPv4 CIDR is expected. For dual-stack networking an IPv4 and IPv6 CIDR are expected. Changing this forces a new resource to be created.

~> **Note:** This range should not be used by any network element on or connected to this VNet. Service address CIDR must be smaller than /12. `docker_bridge_cidr`, `dns_service_ip` and `service_cidr` should all be empty or all should be set.

Examples of how to use [AKS with Advanced Networking](https://docs.microsoft.com/azure/aks/networking-overview#advanced-networking) can be [found in the `./examples/kubernetes/` directory in the GitHub repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/kubernetes).

* `ip_versions` - (Optional) Specifies a list of IP versions the Kubernetes Cluster will use to assign IP addresses to its nodes and pods. Possible values are `IPv4` and/or `IPv6`. `IPv4` must always be specified. Changing this forces a new resource to be created.

-> **Note:** To configure dual-stack networking `ip_versions` should be set to `["IPv4", "IPv6"]`.

-> **Note:** Dual-stack networking requires that the Preview Feature `Microsoft.ContainerService/AKS-EnableDualStack` is enabled and the Resource Provider is re-registered, see [the documentation](https://docs.microsoft.com/azure/aks/configure-kubenet-dual-stack?tabs=azure-cli%2Ckubectl#register-the-aks-enabledualstack-preview-feature) for more information.

* `load_balancer_sku` - (Optional) Specifies the SKU of the Load Balancer used for this Kubernetes Cluster. Possible values are `basic` and `standard`. Defaults to `standard`. Changing this forces a new resource to be created.

* `load_balancer_profile` - (Optional) A `load_balancer_profile` block as defined below. This can only be specified when `load_balancer_sku` is set to `standard`. Changing this forces a new resource to be created.

* `nat_gateway_profile` - (Optional) A `nat_gateway_profile` block as defined below. This can only be specified when `load_balancer_sku` is set to `standard` and `outbound_type` is set to `managedNATGateway` or `userAssignedNATGateway`. Changing this forces a new resource to be created.

---

A `load_balancer_profile` block supports the following:

~> **Note:** The fields `managed_outbound_ip_count`, `outbound_ip_address_ids` and `outbound_ip_prefix_ids` are mutually exclusive. Note that when specifying `outbound_ip_address_ids` ([azurerm_public_ip](/docs/providers/azurerm/r/public_ip.html)) the SKU must be `standard`.

* `backend_pool_type` - (Optional) The type of the managed inbound Load Balancer Backend Pool. Possible values are `NodeIP` and `NodeIPConfiguration`. Defaults to `NodeIPConfiguration`. See [the documentation](https://learn.microsoft.com/en-us/azure/aks/load-balancer-standard#change-the-inbound-pool-type) for more information.

* `idle_timeout_in_minutes` - (Optional) Desired outbound flow idle timeout in minutes for the cluster load balancer. Must be between `4` and `100` inclusive. Defaults to `30`.

* `managed_outbound_ip_count` - (Optional) Count of desired managed outbound IPs for the cluster load balancer. Must be between `1` and `100` inclusive.

* `managed_outbound_ipv6_count` - (Optional) The desired number of IPv6 outbound IPs created and managed by Azure for the cluster load balancer. Must be in the range of 1 to 100 (inclusive). The default value is 0 for single-stack and 1 for dual-stack.

~> **Note:** `managed_outbound_ipv6_count` requires dual-stack networking. To enable dual-stack networking the Preview Feature `Microsoft.ContainerService/AKS-EnableDualStack` needs to be enabled and the Resource Provider re-registered, see [the documentation](https://docs.microsoft.com/azure/aks/configure-kubenet-dual-stack?tabs=azure-cli%2Ckubectl#register-the-aks-enabledualstack-preview-feature) for more information.

* `outbound_ip_address_ids` - (Optional) The ID of the Public IP Addresses which should be used for outbound communication for the cluster load balancer.

-> **Note:** Set `outbound_ip_address_ids` to an empty slice `[]` in order to unlink it from the cluster. Unlinking a `outbound_ip_address_ids` will revert the load balancing for the cluster back to a managed one.

* `outbound_ip_prefix_ids` - (Optional) The ID of the outbound Public IP Address Prefixes which should be used for the cluster load balancer.

-> **Note:** Set `outbound_ip_prefix_ids` to an empty slice `[]` in order to unlink it from the cluster. Unlinking a `outbound_ip_prefix_ids` will revert the load balancing for the cluster back to a managed one.

* `outbound_ports_allocated` - (Optional) Number of desired SNAT port for each VM in the clusters load balancer. Must be between `0` and `64000` inclusive. Defaults to `0`.

---

A `nat_gateway_profile` block supports the following:

* `idle_timeout_in_minutes` - (Optional) Desired outbound flow idle timeout in minutes for the managed nat gateway. Must be between `4` and `120` inclusive. Defaults to `4`.

* `managed_outbound_ip_count` - (Optional) Count of desired managed outbound IPs for the managed nat gateway. Must be between `1` and `16` inclusive.

---

An `oms_agent` block supports the following:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace which the OMS Agent should send data to.

* `msi_auth_for_monitoring_enabled` - (Optional) Is managed identity authentication for monitoring enabled?

---

An `ingress_application_gateway` block supports the following:

* `gateway_id` - (Optional) The ID of the Application Gateway to integrate with the ingress controller of this Kubernetes Cluster. See [this](https://docs.microsoft.com/azure/application-gateway/tutorial-ingress-controller-add-on-existing) page for further details.

* `gateway_name` - (Optional) The name of the Application Gateway to be used or created in the Nodepool Resource Group, which in turn will be integrated with the ingress controller of this Kubernetes Cluster. See [this](https://docs.microsoft.com/azure/application-gateway/tutorial-ingress-controller-add-on-new) page for further details.

* `subnet_cidr` - (Optional) The subnet CIDR to be used to create an Application Gateway, which in turn will be integrated with the ingress controller of this Kubernetes Cluster. See [this](https://docs.microsoft.com/azure/application-gateway/tutorial-ingress-controller-add-on-new) page for further details.

* `subnet_id` - (Optional) The ID of the subnet on which to create an Application Gateway, which in turn will be integrated with the ingress controller of this Kubernetes Cluster. See [this](https://docs.microsoft.com/azure/application-gateway/tutorial-ingress-controller-add-on-new) page for further details.

-> **Note:** Exactly one of `gateway_id`, `gateway_name`, `subnet_id`, or `subnet_cidr` must be specified.

-> **Note:** If specifying `ingress_application_gateway` in conjunction with `only_critical_addons_enabled`, the AGIC pod will fail to start. A separate `azurerm_kubernetes_cluster_node_pool` is required to run the AGIC pod successfully. This is because AGIC is classed as a "non-critical addon".

---

A `service_mesh_profile` block supports the following:

* `mode` - (Required) The mode of the service mesh. Possible value is `Istio`.

* `revisions` - (Required) Specify 1 or 2 Istio control plane revisions for managing minor upgrades using the canary upgrade process. For example, create the resource with `revisions` set to `["asm-1-20"]`, or leave it empty (the `revisions` will only be known after apply). To start the canary upgrade, change `revisions` to `["asm-1-20", "asm-1-21"]`. To roll back the canary upgrade, revert to `["asm-1-20"]`. To confirm the upgrade, change to `["asm-1-21"]`.

-> **Note:** Upgrading to a new (canary) revision does not affect existing sidecar proxies. You need to apply the canary revision label to selected namespaces and restart pods with kubectl to inject the new sidecar proxy. [Learn more](https://istio.io/latest/docs/setup/upgrade/canary/#data-plane).

* `internal_ingress_gateway_enabled` - (Optional) Is Istio Internal Ingress Gateway enabled?

* `external_ingress_gateway_enabled` - (Optional) Is Istio External Ingress Gateway enabled?

-> **Note:** Currently only one Internal Ingress Gateway and one External Ingress Gateway are allowed per cluster

* `certificate_authority` - (Optional) A `certificate_authority` block as defined below. When this property is specified, `key_vault_secrets_provider` is also required to be set. This configuration allows you to bring your own root certificate and keys for Istio CA in the Istio-based service mesh add-on for Azure Kubernetes Service.

---

A `certificate_authority` block supports the following:

* `key_vault_id` - (Required) The resource ID of the Key Vault.

* `root_cert_object_name` - (Required) The root certificate object name in Azure Key Vault.

* `cert_chain_object_name` - (Required) The certificate chain object name in Azure Key Vault.

* `cert_object_name` - (Required) The intermediate certificate object name in Azure Key Vault.

* `key_object_name` - (Required) The intermediate certificate private key object name in Azure Key Vault.

-> **Note:** For more information on [Istio-based service mesh add-on with plug-in CA certificates and how to generate these certificates](https://learn.microsoft.com/en-us/azure/aks/istio-plugin-ca),

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal.

* `client_secret` - (Required) The Client Secret for the Service Principal.

---

A `ssh_key` block supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster. Changing this forces a new resource to be created.

---

A `storage_profile` block supports the following:

* `blob_driver_enabled` - (Optional) Is the Blob CSI driver enabled? Defaults to `false`.

* `disk_driver_enabled` - (Optional) Is the Disk CSI driver enabled? Defaults to `true`.

* `file_driver_enabled` - (Optional) Is the File CSI driver enabled? Defaults to `true`.

* `snapshot_controller_enabled` - (Optional) Is the Snapshot Controller enabled? Defaults to `true`.

---

A `sysctl_config` block supports the following:

~> **Note:** For more information, please refer to [Linux Kernel Doc](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/index.html).

* `fs_aio_max_nr` - (Optional) The sysctl setting fs.aio-max-nr. Must be between `65536` and `6553500`.

* `fs_file_max` - (Optional) The sysctl setting fs.file-max. Must be between `8192` and `12000500`.

* `fs_inotify_max_user_watches` - (Optional) The sysctl setting fs.inotify.max_user_watches. Must be between `781250` and `2097152`.

* `fs_nr_open` - (Optional) The sysctl setting fs.nr_open. Must be between `8192` and `20000500`.

* `kernel_threads_max` - (Optional) The sysctl setting kernel.threads-max. Must be between `20` and `513785`.

* `net_core_netdev_max_backlog` - (Optional) The sysctl setting net.core.netdev_max_backlog. Must be between `1000` and `3240000`.

* `net_core_optmem_max` - (Optional) The sysctl setting net.core.optmem_max. Must be between `20480` and `4194304`.

* `net_core_rmem_default` - (Optional) The sysctl setting net.core.rmem_default. Must be between `212992` and `134217728`.

* `net_core_rmem_max` - (Optional) The sysctl setting net.core.rmem_max. Must be between `212992` and `134217728`.

* `net_core_somaxconn` - (Optional) The sysctl setting net.core.somaxconn. Must be between `4096` and `3240000`.

* `net_core_wmem_default` - (Optional) The sysctl setting net.core.wmem_default. Must be between `212992` and `134217728`.

* `net_core_wmem_max` - (Optional) The sysctl setting net.core.wmem_max. Must be between `212992` and `134217728`.

* `net_ipv4_ip_local_port_range_max` - (Optional) The sysctl setting net.ipv4.ip_local_port_range max value. Must be between `32768` and `65535`.

* `net_ipv4_ip_local_port_range_min` - (Optional) The sysctl setting net.ipv4.ip_local_port_range min value. Must be between `1024` and `60999`.

* `net_ipv4_neigh_default_gc_thresh1` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh1. Must be between `128` and `80000`.

* `net_ipv4_neigh_default_gc_thresh2` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh2. Must be between `512` and `90000`.

* `net_ipv4_neigh_default_gc_thresh3` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh3. Must be between `1024` and `100000`.

* `net_ipv4_tcp_fin_timeout` - (Optional) The sysctl setting net.ipv4.tcp_fin_timeout. Must be between `5` and `120`.

* `net_ipv4_tcp_keepalive_intvl` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_intvl. Must be between `10` and `90`.

* `net_ipv4_tcp_keepalive_probes` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_probes. Must be between `1` and `15`.

* `net_ipv4_tcp_keepalive_time` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_time. Must be between `30` and `432000`.

* `net_ipv4_tcp_max_syn_backlog` - (Optional) The sysctl setting net.ipv4.tcp_max_syn_backlog. Must be between `128` and `3240000`.

* `net_ipv4_tcp_max_tw_buckets` - (Optional) The sysctl setting net.ipv4.tcp_max_tw_buckets. Must be between `8000` and `1440000`.

* `net_ipv4_tcp_tw_reuse` - (Optional) The sysctl setting net.ipv4.tcp_tw_reuse.

* `net_netfilter_nf_conntrack_buckets` - (Optional) The sysctl setting net.netfilter.nf_conntrack_buckets. Must be between `65536` and `524288`.

* `net_netfilter_nf_conntrack_max` - (Optional) The sysctl setting net.netfilter.nf_conntrack_max. Must be between `131072` and `2097152`.

* `vm_max_map_count` - (Optional) The sysctl setting vm.max_map_count. Must be between `65530` and `262144`.

* `vm_swappiness` - (Optional) The sysctl setting vm.swappiness. Must be between `0` and `100`.

* `vm_vfs_cache_pressure` - (Optional) The sysctl setting vm.vfs_cache_pressure. Must be between `0` and `100`.

---

The `upgrade_override` block supports the following:

-> **Note:** Once set, the `upgrade_override` block cannot be removed from the configuration.

* `force_upgrade_enabled` - (Required) Whether to force upgrade the cluster. Possible values are `true` or `false`.

!> **Note:** The `force_upgrade_enabled` field instructs the upgrade operation to bypass upgrade protections (e.g. checking for deprecated API usage) which may render the cluster inoperative after the upgrade process has completed. Use the `force_upgrade_enabled` option with extreme caution only.

* `effective_until` - (Optional) Specifies the duration, in RFC 3339 format (e.g., `2025-10-01T13:00:00Z`), the `upgrade_override` values are effective. This field must be set for the `upgrade_override` values to take effect. The date-time must be within the next 30 days.

-> **Note:** This only matches the start time of an upgrade, and the effectiveness won't change once an upgrade starts even if the `effective_until` value expires as the upgrade proceeds.

---

A `web_app_routing` block supports the following:

* `dns_zone_ids` - (Required) Specifies the list of the DNS Zone IDs in which DNS entries are created for applications deployed to the cluster when Web App Routing is enabled. If not using Bring-Your-Own DNS zones this property should be set to an empty list.

---

A `windows_profile` block supports the following:

* `admin_username` - (Required) The Admin Username for Windows VMs. Changing this forces a new resource to be created.

* `admin_password` - (Optional) The Admin Password for Windows VMs. Length must be between 14 and 123 characters.

* `license` - (Optional) Specifies the type of on-premise license which should be used for Node Pool Windows Virtual Machine. At this time the only possible value is `Windows_Server`.

* `gmsa` - (Optional) A `gmsa` block as defined below.

---

A `gmsa` block supports the following:

* `dns_server` - (Required) Specifies the DNS server for Windows gMSA. Set this to an empty string if you have configured the DNS server in the VNet which was used to create the managed cluster.

* `root_domain` - (Required) Specifies the root domain name for Windows gMSA. Set this to an empty string if you have configured the DNS server in the VNet which was used to create the managed cluster.

-> **Note:** The properties `dns_server` and `root_domain` must both either be set or unset, i.e. empty.

---

A `workload_autoscaler_profile` block supports the following:

* `keda_enabled` - (Optional) Specifies whether KEDA Autoscaler can be used for workloads.

* `vertical_pod_autoscaler_enabled` - (Optional) Specifies whether Vertical Pod Autoscaler should be enabled.

-> **Note:** This requires that the Preview Feature `Microsoft.ContainerService/AKS-VPAPreview` is enabled and the Resource Provider is re-registered, see [the documentation]([Microsoft.ContainerService/AKS-VPAPreview](https://learn.microsoft.com/en-us/azure/aks/vertical-pod-autoscaler#register-the-aks-vpapreview-feature-flag) for more information.

---

A `http_proxy_config` block supports the following:

* `http_proxy` - (Optional) The proxy address to be used when communicating over HTTP.

* `https_proxy` - (Optional) The proxy address to be used when communicating over HTTPS.

* `no_proxy` - (Optional) The list of domains that will not use the proxy for communication.

-> **Note:** If you specify the `default_node_pool[0].vnet_subnet_id`, be sure to include the Subnet CIDR in the `no_proxy` list.

-> **Note:** You may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore the changes to this field.

* `trusted_ca` - (Optional) The base64 encoded alternative CA certificate content in PEM format.

---

A `upgrade_settings` block supports the following:

* `drain_timeout_in_minutes` - (Optional) The amount of time in minutes to wait on eviction of pods and graceful termination per node. This eviction wait time honors pod disruption budgets for upgrades. If this time is exceeded, the upgrade fails. Unsetting this after configuring it will force a new resource to be created.

* `node_soak_duration_in_minutes` - (Optional) The amount of time in minutes to wait after draining a node and before reimaging and moving on to next node. Defaults to `0`.

* `max_surge` - (Required) The maximum number or percentage of nodes which will be added to the Node Pool size during an upgrade.

-> **Note:** If a percentage is provided, the number of surge nodes is calculated from the `node_count` value on the current cluster. Node surge can allow a cluster to have more nodes than `max_count` during an upgrade. Ensure that your cluster has enough [IP space](https://docs.microsoft.com/azure/aks/upgrade-cluster#customize-node-surge-upgrade) during an upgrade.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `current_kubernetes_version` - The current version running on the Azure Kubernetes Managed Cluster.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `private_fqdn` - The FQDN for the Kubernetes Cluster when private link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.

* `portal_fqdn` - The FQDN for the Azure Portal resources when private link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.

* `kube_admin_config` - A `kube_admin_config` block as defined below. This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts enabled.

* `kube_admin_config_raw` - Raw Kubernetes config for the admin account to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools. This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts enabled.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools.

* `http_application_routing_zone_name` - The Zone Name of the HTTP Application Routing.

* `oidc_issuer_url` - The OIDC issuer URL that is associated with the cluster.

* `node_resource_group` - The auto-generated Resource Group which contains the resources for this Managed Kubernetes Cluster.

* `node_resource_group_id` - The ID of the Resource Group containing the resources for this Managed Kubernetes Cluster.

* `network_profile` - A `network_profile` block as defined below.

* `ingress_application_gateway` - An `ingress_application_gateway` block as defined below.

* `oms_agent` - An `oms_agent` block as defined below.

* `key_vault_secrets_provider` - A `key_vault_secrets_provider` block as defined below.

---

The `aci_connector_linux` block exports the following:

* `connector_identity` - A `connector_identity` block is exported. The exported attributes are defined below.

---

The `connector_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the ACI Connector.

* `object_id` - The Object ID of the user-defined Managed Identity used by the ACI Connector.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the ACI Connector.

---

The `kubelet_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity to be assigned to the Kubelets.

* `object_id` - The Object ID of the user-defined Managed Identity assigned to the Kubelets.

* `user_assigned_identity_id` - The ID of the User Assigned Identity assigned to the Kubelets.

---

A `network_profile` block exports the following:

* `load_balancer_profile` - A `load_balancer_profile` block as defined below.

* `nat_gateway_profile` - A `nat_gateway_profile` block as defined below.

---

A `load_balancer_profile` block exports the following:

* `effective_outbound_ips` - The outcome (resource IDs) of the specified arguments.

---

A `nat_gateway_profile` block exports the following:

* `effective_outbound_ips` - The outcome (resource IDs) of the specified arguments.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

The `kube_admin_config` and `kube_config` blocks export the following:

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **Note:** It's possible to use these credentials with [the Kubernetes Provider](/providers/hashicorp/kubernetes/latest/docs) like so:

```hcl
provider "kubernetes" {
  host                   = azurerm_kubernetes_cluster.main.kube_config[0].host
  username               = azurerm_kubernetes_cluster.main.kube_config[0].username
  password               = azurerm_kubernetes_cluster.main.kube_config[0].password
  client_certificate     = base64decode(azurerm_kubernetes_cluster.main.kube_config[0].client_certificate)
  client_key             = base64decode(azurerm_kubernetes_cluster.main.kube_config[0].client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.main.kube_config[0].cluster_ca_certificate)
}
```

---

The `ingress_application_gateway` block exports the following:

* `effective_gateway_id` - The ID of the Application Gateway associated with the ingress controller deployed to this Kubernetes Cluster.

* `ingress_application_gateway_identity` - An `ingress_application_gateway_identity` block is exported. The exported attributes are defined below.

---

The `ingress_application_gateway_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Application Gateway.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Application Gateway.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Application Gateway.

---

The `oms_agent` block exports the following:

* `oms_agent_identity` - An `oms_agent_identity` block is exported. The exported attributes are defined below.

---

The `oms_agent_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the OMS Agents.

* `object_id` - The Object ID of the user-defined Managed Identity used by the OMS Agents.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the OMS Agents.

---

The `key_vault_secrets_provider` block exports the following:

* `secret_identity` - An `secret_identity` block is exported. The exported attributes are defined below.

---

The `secret_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Secret Provider.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Secret Provider.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Secret Provider.

---

A `web_app_routing` block exports the following:

* `web_app_routing_identity` - A `web_app_routing_identity` block is exported. The exported attributes are defined below.

---

The `web_app_routing_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used for Web App Routing.

* `object_id` - The Object ID of the user-defined Managed Identity used for Web App Routing

* `user_assigned_identity_id` - The ID of the User Assigned Identity used for Web App Routing.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster.
* `update` - (Defaults to 90 minutes) Used when updating the Kubernetes Cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Kubernetes Cluster.

## Import

Managed Kubernetes Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
