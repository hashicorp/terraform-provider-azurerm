---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
description: |-
  Gets information about an existing Managed Kubernetes Cluster (AKS)
---

# Data Source: azurerm_kubernetes_cluster

Use this data source to access information about an existing Managed Kubernetes Cluster (AKS).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain text.
[Read more about sensitive data in the state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_kubernetes_cluster" "example" {
  name                = "myakscluster"
  resource_group_name = "my-example-resource-group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the managed Kubernetes Cluster.

* `resource_group_name` - The name of the Resource Group in which the managed Kubernetes Cluster exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes Managed Cluster.

* `api_server_authorized_ip_ranges` - The IP ranges to whitelist for incoming traffic to the primaries.

* `aci_connector_linux` - An `aci_connector_linux` block as documented below.

* `azure_active_directory_role_based_access_control` - An `azure_active_directory_role_based_access_control` block as documented below.

* `azure_policy_enabled` - Is Azure Policy enabled on this managed Kubernetes Cluster?

* `agent_pool_profile` - An `agent_pool_profile` block as documented below.

* `current_kubernetes_version` - Contains the current version of Kubernetes running on the Cluster.

* `dns_prefix` - The DNS Prefix of the managed Kubernetes cluster.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `http_application_routing_enabled` - Is HTTP Application Routing enabled for this managed Kubernetes Cluster?

* `http_application_routing_zone_name` - The Zone Name of the HTTP Application Routing.

* `ingress_application_gateway` - An `ingress_application_gateway` block as documented below.

* `key_management_service` - A `key_management_service` block as documented below.

* `key_vault_secrets_provider` - A `key_vault_secrets_provider` block as documented below.

* `private_fqdn` - The FQDN of this Kubernetes Cluster when private link has been enabled. This name is only resolvable inside the Virtual Network where the Azure Kubernetes Service is located

-> **Note:** At this time Private Link is in Public Preview.

* `kube_admin_config` - A `kube_admin_config` block as defined below. This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts are not disabled.

* `kube_admin_config_raw` - Raw Kubernetes config for the admin account to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools. This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts are not disabled.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Base64 encoded Kubernetes configuration.

* `kubernetes_version` - The version of Kubernetes used on the managed Kubernetes Cluster.

* `private_cluster_enabled` - If the cluster has the Kubernetes API only exposed on internal IP addresses.

* `location` - The Azure Region in which the managed Kubernetes Cluster exists.

* `microsoft_defender` - A `microsoft_defender` block as defined below.

* `oidc_issuer_enabled` - Whether or not the OIDC feature is enabled or disabled.

* `oidc_issuer_url` - The OIDC issuer URL that is associated with the cluster.

* `oms_agent` - An `oms_agent` block as documented below.

* `open_service_mesh_enabled` - Is Open Service Mesh enabled for this managed Kubernetes Cluster?

* `disk_encryption_set_id` - The ID of the Disk Encryption Set used for the Nodes and Volumes.

* `linux_profile` - A `linux_profile` block as documented below.

* `windows_profile` - A `windows_profile` block as documented below.

* `network_profile` - A `network_profile` block as documented below.

* `node_resource_group` - Auto-generated Resource Group containing AKS Cluster resources.

* `node_resource_group_id` - The ID of the Resource Group containing the resources for this Managed Kubernetes Cluster.

* `role_based_access_control_enabled` - Is Role Based Access Control enabled for this managed Kubernetes Cluster?

* `service_principal` - A `service_principal` block as documented below.

* `storage_profile` - A `storage_profile` block as documented below.

* `identity` - An `identity` block as documented below.

* `kubelet_identity` - A `kubelet_identity` block as documented below.

* `tags` - A mapping of tags assigned to this resource.

---

An `aci_connector_linux` block exports the following:

* `subnet_name` - The subnet name for the virtual nodes to run.

---

An `agent_pool_profile` block exports the following:

* `type` - The type of the Agent Pool.

* `count` - The number of Agents (VMs) in the Pool.

* `max_pods` - The maximum number of pods that can run on each agent.

* `auto_scaling_enabled` - If the auto-scaler is enabled.

* `node_public_ip_enabled` - If the Public IPs for the nodes in this Agent Pool are enabled.

* `host_group_id` - The ID of a Dedicated Host Group that this Node Pool should be run on. Changing this forces a new resource to be created.

* `min_count` - Minimum number of nodes for auto-scaling

* `max_count` - Maximum number of nodes for auto-scaling

* `name` - The name assigned to this pool of agents.

* `node_public_ip_prefix_id` - Resource ID for the Public IP Addresses Prefix for the nodes in this Agent Pool.

* `os_disk_size_gb` - The size of the Agent VM's Operating System Disk in GB.

* `os_type` - The Operating System used for the Agents.

* `tags` - A mapping of tags to assign to the resource.

* `orchestrator_version` - Kubernetes version used for the Agents.

* `upgrade_settings` - A `upgrade_settings` block as documented below.

* `vm_size` - The size of each VM in the Agent Pool (e.g. `Standard_F1`).

* `vnet_subnet_id` - The ID of the Subnet where the Agents in the Pool are provisioned.

* `zones` - A list of Availability Zones in which this Kubernetes Cluster is located.

---

An `azure_active_directory_role_based_access_control` block exports the following:

* `tenant_id` - The Tenant ID used for Azure Active Directory Application.

* `admin_group_object_ids` - A list of Object IDs of Azure Active Directory Groups which should have Admin Role on the Cluster.

* `azure_rbac_enabled` - Is Role Based Access Control based on Azure AD enabled?

---

A `upgrade_settings` block exports the following:

* `drain_timeout_in_minutes` - The amount of time in minutes to wait on eviction of pods and graceful termination per node. This eviction wait time honors waiting on pod disruption budgets. If this time is exceeded, the upgrade fails.

* `node_soak_duration_in_minutes` - The amount of time in minutes to wait after draining a node and before reimaging it and moving on to next node.

* `max_surge` - The maximum number or percentage of nodes that will be added to the Node Pool size during an upgrade.

---

A `key_management_service` block supports the following:

* `key_vault_key_id` - Identifier of Azure Key Vault key. See [key identifier format](https://learn.microsoft.com/en-us/azure/key-vault/general/about-keys-secrets-certificates#vault-name-and-object-name) for more details.

* `key_vault_network_access` - Network access of the key vault. The possible values are `Public` and `Private`. `Public` means the key vault allows public access from all networks. `Private` means the key vault disables public access and enables private link.

---

A `key_vault_secrets_provider` block exports the following:

* `secret_rotation_enabled` - Is secret rotation enabled?

* `secret_rotation_interval` - The interval to poll for secret rotation.

* `secret_identity` - A `secret_identity` block as documented below.

---

The `kube_admin_config` and `kube_config` blocks export the following:

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **Note:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```hcl
provider "kubernetes" {
  host                   = data.azurerm_kubernetes_cluster.main.kube_config[0].host
  username               = data.azurerm_kubernetes_cluster.main.kube_config[0].username
  password               = data.azurerm_kubernetes_cluster.main.kube_config[0].password
  client_certificate     = base64decode(data.azurerm_kubernetes_cluster.main.kube_config[0].client_certificate)
  client_key             = base64decode(data.azurerm_kubernetes_cluster.main.kube_config[0].client_key)
  cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.main.kube_config[0].cluster_ca_certificate)
}
```

---

A `linux_profile` block exports the following:

* `admin_username` - The username associated with the administrator account of the managed Kubernetes Cluster.

* `ssh_key` - An `ssh_key` block as defined below.

---

A `microsoft_defender` block exports the following:

* `log_analytics_workspace_id` - The ID of the Log Analytics Workspace which Microsoft Defender uses to send audit logs to.

---

A `windows_profile` block exports the following:

* `admin_username` - The username associated with the administrator account of the Windows VMs.

---

A `network_profile` block exports the following:

* `docker_bridge_cidr` - IP address (in CIDR notation) used as the Docker bridge IP address on nodes.

* `dns_service_ip` - IP address within the Kubernetes service address range used by cluster service discovery (kube-dns).

* `network_plugin` - Network plugin used such as `azure` or `kubenet`.

* `network_policy` - Network policy to be used with Azure CNI. e.g. `calico` or `azure`

* `network_mode` - Network mode to be used with Azure CNI. e.g. `bridge` or `transparent`

* `pod_cidr` - The CIDR used for pod IP addresses.

* `service_cidr` - Network range used by the Kubernetes service.

---

An `oms_agent` block exports the following:

* `log_analytics_workspace_id` - The ID of the Log Analytics Workspace to which the OMS Agent should send data.

* `msi_auth_for_monitoring_enabled` - Is managed identity authentication for monitoring enabled?

* `oms_agent_identity` - An `oms_agent_identity` block as defined below.

---

The `oms_agent_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the OMS Agents.

* `object_id` - The Object ID of the user-defined Managed Identity used by the OMS Agents.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the OMS Agents.

---

An `ingress_application_gateway` block supports the following:

* `effective_gateway_id` - The ID of the Application Gateway associated with the ingress controller deployed to this Kubernetes Cluster.

* `gateway_id` - The ID of the Application Gateway integrated with the ingress controller of this Kubernetes Cluster. This attribute is only set when gateway_id is specified when configuring the `ingress_application_gateway` addon.

* `subnet_cidr` - The subnet CIDR used to create an Application Gateway, which in turn will be integrated with the ingress controller of this Kubernetes Cluster. This attribute is only set when `subnet_cidr` is specified when configuring the `ingress_application_gateway` addon.

* `subnet_id` - The ID of the subnet on which to create an Application Gateway, which in turn will be integrated with the ingress controller of this Kubernetes Cluster. This attribute is only set when `subnet_id` is specified when configuring the `ingress_application_gateway` addon.

* `ingress_application_gateway_identity` - An `ingress_application_gateway_identity` block as defined below.

---

The `ingress_application_gateway_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Application Gateway.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Application Gateway.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Application Gateway.

---

The `secret_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Secret Provider.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Secret Provider.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Secret Provider.

---

A `service_principal` block exports the following:

* `client_id` - The Client ID of the Service Principal used by this Managed Kubernetes Cluster.

---

A `storage_profile` block exports the following:

* `blob_driver_enabled` Is the Blob CSI driver enabled?

* `disk_driver_enabled` Is the Disk CSI driver enabled?

* `file_driver_enabled` Is the File CSI driver enabled?

* `snapshot_controller_enabled` Is the Snapshot Controller enabled?

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Kubernetes Cluster.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Kubernetes Cluster.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Kubernetes Cluster.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Kubernetes Cluster.

---

The `kubelet_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity assigned to the Kubelets.

* `object_id` - The Object ID of the user-defined Managed Identity assigned to the Kubelets.

* `user_assigned_identity_id` - The ID of the User Assigned Identity assigned to the Kubelets.

---

A `ssh_key` block exports the following:

* `key_data` - The Public SSH Key used to access the cluster.

---

A `service_mesh_profile` block exports the following:

* `mode` - The mode of the service mesh.

* `revisions` - List of revisions of the Istio control plane. When an upgrade is not in progress, this holds one value. When canary upgrade is in progress, this can only hold two consecutive values. [Learn More](
  https://learn.microsoft.com/en-us/azure/aks/istio-upgrade).

* `internal_ingress_gateway_enabled` - Is Istio Internal Ingress Gateway enabled?

* `external_ingress_gateway_enabled` - Is Istio External Ingress Gateway enabled?

* `certificate_authority` - A `certificate_authority` block as documented below.

---

A `certificate_authority` block exports the following:

* `key_vault_id` - The resource ID of the Key Vault.

* `root_cert_object_name` - The root certificate object name in Azure Key Vault.

* `cert_chain_object_name` - The certificate chain object name in Azure Key Vault.

* `cert_object_name` - The intermediate certificate object name in Azure Key Vault.

* `key_object_name` - The intermediate certificate private key object name in Azure Key Vault.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Kubernetes Cluster (AKS).
