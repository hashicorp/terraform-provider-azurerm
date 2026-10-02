---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_automatic_cluster"
description: |-
  Gets information about an existing Managed Kubernetes Automatic Cluster (AKS)
---

# Data Source: azurerm_kubernetes_automatic_cluster

Use this data source to access information about an existing Managed Kubernetes Automatic Cluster (AKS).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain text.
[Read more about sensitive data in the state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_kubernetes_automatic_cluster" "example" {
  name                = "myakscluster"
  resource_group_name = "my-example-resource-group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Kubernetes Automatic Cluster.

* `resource_group_name` - (Required) The name of the Resource Group in which the Managed Kubernetes Automatic Cluster exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Kubernetes Automatic Cluster.

* `location` - The Azure Region in which the Managed Kubernetes Automatic Cluster exists.

* `agent_pools` - One or more `agent_pools` blocks as documented below.

* `api_server_access` - An `api_server_access` block as documented below.

* `azure_active_directory_role_based_access_control` - An `azure_active_directory_role_based_access_control` block as documented below.

* `azure_policy_enabled` - Is Azure Policy enabled on this Managed Kubernetes Automatic Cluster?

* `bootstrap` - A `bootstrap` block as documented below.

* `current_kubernetes_version` - Contains the current version of Kubernetes running on the Cluster.

* `disk_encryption_set_id` - The ID of the Disk Encryption Set used for the Nodes and Volumes.

* `dns_prefix` - The DNS Prefix of the Managed Kubernetes Automatic Cluster.

* `fully_qualified_domain_name` - The FQDN of the Managed Kubernetes Automatic Cluster.

* `hosted_system` - A `hosted_system` block as documented below.

* `identity` - An `identity` block as documented below.

* `key_management_service` - A `key_management_service` block as documented below.

* `key_vault_secrets_provider` - A `key_vault_secrets_provider` block as documented below.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Base64 encoded Kubernetes configuration.

* `kubelet_identity` - A `kubelet_identity` block as documented below.

* `kubernetes_version` - The version of Kubernetes used on the Managed Kubernetes Automatic Cluster.

* `microsoft_defender` - A `microsoft_defender` block as documented below.

* `network` - A `network` block as documented below.

* `node_resource_group` - Auto-generated Resource Group containing AKS Cluster resources.

* `node_resource_group_id` - The ID of the Resource Group containing the resources for this Managed Kubernetes Automatic Cluster.

* `oidc_issuer_enabled` - Whether or not the OIDC feature is enabled or disabled.

* `oidc_issuer_url` - The OIDC issuer URL that is associated with the Managed Kubernetes Automatic Cluster.

* `oms_agent` - An `oms_agent` block as documented below.

* `portal_fully_qualified_domain_name` - The FQDN used by the Azure Portal for this Managed Kubernetes Automatic Cluster.

* `private_cluster` - A `private_cluster` block as documented below.

* `private_fully_qualified_domain_name` - The FQDN of this Managed Kubernetes Automatic Cluster when private link has been enabled. This name is only resolvable inside the Virtual Network where the Azure Kubernetes Service is located.

* `role_based_access_control_enabled` - Is Role Based Access Control enabled for this Managed Kubernetes Automatic Cluster?

* `service_mesh` - A `service_mesh` block as documented below.

* `storage` - A `storage` block as documented below.

* `tags` - A mapping of tags assigned to this resource.

* `web_app_routing_ingress` - A `web_app_routing_ingress` block as documented below.

---

An `agent_pools` block exports the following:

* `name` - The name assigned to this pool of agents.

* `type` - The type of the Agent Pool.

* `count` - The number of Agents (VMs) in the Pool.

* `max_count` - Maximum number of nodes for auto-scaling.

* `min_count` - Minimum number of nodes for auto-scaling.

* `auto_scaling_enabled` - If the auto-scaler is enabled.

* `vm_size` - The size of each VM in the Agent Pool (e.g. `Standard_F1`).

* `tags` - A mapping of tags assigned to the Agent Pool.

* `os_disk_size_gb` - The size of the Agent VM's Operating System Disk in GB.

* `vnet_subnet_id` - The ID of the Subnet where the Agents in the Pool are provisioned.

* `os_type` - The Operating System used for the Agents.

* `orchestrator_version` - Kubernetes version used for the Agents.

* `max_pods` - The maximum number of pods that can run on each agent.

* `node_labels` - A map of Kubernetes labels applied to the nodes in this Agent Pool.

* `node_taints` - A list of Kubernetes taints applied to the nodes in this Agent Pool.

* `node_public_ip_enabled` - If the Public IPs for the nodes in this Agent Pool are enabled.

* `node_public_ip_prefix_id` - Resource ID for the Public IP Addresses Prefix for the nodes in this Agent Pool.

* `upgrade_settings` - An `upgrade_settings` block as documented below.

* `zones` - A list of Availability Zones in which the nodes in this Agent Pool are located.

---

An `upgrade_settings` block exports the following:

* `drain_timeout_in_minutes` - The amount of time in minutes to wait on eviction of pods and graceful termination per node. This eviction wait time honours waiting on pod disruption budgets. If this time is exceeded, the upgrade fails.

* `max_surge` - The maximum number or percentage of nodes that will be added to the Node Pool size during an upgrade.

* `max_unavailable` - The maximum number or percentage of nodes that can be unavailable during an upgrade.

* `node_soak_duration_in_minutes` - The amount of time in minutes to wait after draining a node and before reimaging it and moving on to next node.

* `undrainable_node_behavior` - The action when a node is undrainable during upgrade. Possible values are `Cordon` and `Schedule`.

---

An `api_server_access` block exports the following:

* `authorized_ip_ranges` - A list of IP ranges authorised to access the API server.

* `subnet_id` - The ID of the subnet that the API server is accessible from.

---

An `azure_active_directory_role_based_access_control` block exports the following:

* `tenant_id` - The Tenant ID used for Azure Active Directory Application.

* `admin_group_object_ids` - A list of Object IDs of Azure Active Directory Groups which should have Admin Role on the Cluster.

* `azure_rbac_enabled` - Is Role Based Access Control based on Azure AD enabled?

---

A `bootstrap` block exports the following:

* `artifact_source` - The source from which artifacts are pulled during bootstrap.

* `container_registry_id` - The ID of the Azure Container Registry used for caching artifacts during bootstrap.

---

A `hosted_system` block exports the following:

* `node_subnet_id` - The ID of the subnet used for the hosted system nodes.

* `system_node_subnet_id` - The ID of the subnet used for the system nodes.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Managed Kubernetes Automatic Cluster.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Managed Kubernetes Automatic Cluster.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Managed Kubernetes Automatic Cluster.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Managed Kubernetes Automatic Cluster.

---

A `key_management_service` block exports the following:

* `key_vault_key_id` - Identifier of Azure Key Vault key. See [key identifier format](https://learn.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates#vault-name-and-object-name) for more details.

* `key_vault_network_access` - Network access of the key vault. The possible values are `Public` and `Private`. `Public` means the key vault allows public access from all networks. `Private` means the key vault disables public access and enables private link.

---

A `key_vault_secrets_provider` block exports the following:

* `secret_rotation_enabled` - Is secret rotation enabled?

* `secret_rotation_interval` - The interval to poll for secret rotation.

* `secret_identity` - A `secret_identity` block as documented below.

---

The `secret_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Secret Provider.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Secret Provider.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Secret Provider.

---

The `kube_config` block exports the following:

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **Note:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```hcl
provider "kubernetes" {
  host                   = data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].host
  username               = data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].username
  password               = data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].password
  client_certificate     = base64decode(data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].client_certificate)
  client_key             = base64decode(data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].client_key)
  cluster_ca_certificate = base64decode(data.azurerm_kubernetes_automatic_cluster.example.kube_config[0].cluster_ca_certificate)
}
```

---

The `kubelet_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity assigned to the Kubelets.

* `object_id` - The Object ID of the user-defined Managed Identity assigned to the Kubelets.

* `user_assigned_identity_id` - The ID of the User Assigned Identity assigned to the Kubelets.

---

A `microsoft_defender` block exports the following:

* `log_analytics_workspace_id` - The ID of the Log Analytics Workspace which Microsoft Defender uses to send audit logs to.

---

A `network` block exports the following:

* `dns_service_ip` - IP address within the Kubernetes service address range used by cluster service discovery (kube-dns).

* `load_balancer_sku` - The SKU of the Load Balancer used for this Managed Kubernetes Automatic Cluster.

* `network_plugin` - Network plugin used such as `azure` or `kubenet`.

* `network_policy` - Network policy to be used with Azure CNI. e.g. `calico`, `azure` or `cilium`.

* `outbound_type` - The outbound (egress) routing method which is used for cluster egress traffic.

* `pod_cidr` - The CIDR used for pod IP addresses.

* `service_cidr` - Network range used by the Kubernetes service.

---

An `oms_agent` block exports the following:

* `log_analytics_workspace_id` - The ID of the Log Analytics Workspace to which the OMS Agent should send data.

* `msi_auth_for_monitoring_enabled` - Is managed identity authentication for monitoring enabled?

* `retina_flow_logs_enabled` - Is Retina Flow Logs collection enabled?

* `oms_agent_identity` - An `oms_agent_identity` block as documented below.

---

The `oms_agent_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the OMS Agents.

* `object_id` - The Object ID of the user-defined Managed Identity used by the OMS Agents.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the OMS Agents.

---

A `private_cluster` block exports the following:

* `public_fully_qualified_domain_name_enabled` - If the public FQDN for this Managed Kubernetes Automatic Cluster is enabled.

* `private_dns_zone_id` - The ID of the Private DNS Zone used by this Managed Kubernetes Automatic Cluster.

---

A `service_mesh` block exports the following:

* `revisions` - List of revisions of the Istio control plane. When an upgrade is not in progress, this holds one value. When a canary upgrade is in progress, this can hold two consecutive values. [Learn More](https://learn.microsoft.com/azure/aks/istio-upgrade).

* `internal_ingress_gateway_enabled` - If the Istio Internal Ingress Gateway is enabled.

* `external_ingress_gateway_enabled` - If the Istio External Ingress Gateway is enabled.

* `proxy_redirect_mechanism` - The proxy redirect mechanism configured for the Istio service mesh.

* `certificate_authority` - A `certificate_authority` block as documented below.

---

A `certificate_authority` block exports the following:

* `key_vault_id` - The resource ID of the Key Vault.

* `root_certificate_object_name` - The root certificate object name in Azure Key Vault.

* `certificate_chain_object_name` - The certificate chain object name in Azure Key Vault.

* `certificate_object_name` - The intermediate certificate object name in Azure Key Vault.

* `key_object_name` - The intermediate certificate private key object name in Azure Key Vault.

---

A `storage` block exports the following:

* `blob_driver_enabled` - Is the Blob CSI driver enabled?

* `disk_driver_enabled` - Is the Disk CSI driver enabled?

* `file_driver_enabled` - Is the File CSI driver enabled?

* `snapshot_controller_enabled` - Is the Snapshot Controller enabled?

---

A `web_app_routing_ingress` block exports the following:

* `dns_zone_ids` - A list of DNS Zone IDs associated with the web app routing ingress.

* `default_nginx_controller` - The default Nginx controller for the web app routing ingress.

* `istio_enabled` - If Istio is enabled for the web app routing ingress.

* `web_app_routing_identity` - A `web_app_routing_identity` block as documented below.

---

The `web_app_routing_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used by the Web App Routing ingress controller.

* `object_id` - The Object ID of the user-defined Managed Identity used by the Web App Routing ingress controller.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by the Web App Routing ingress controller.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Kubernetes Automatic Cluster (AKS).

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2026-05-01
