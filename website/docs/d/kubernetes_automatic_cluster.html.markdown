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

* `api_server_access` - An `api_server_access` block as documented below.

* `current_kubernetes_version` - Contains the current version of Kubernetes running on the Cluster.

* `dns_prefix` - The DNS Prefix of the Managed Kubernetes Automatic Cluster.

* `fully_qualified_domain_name` - The FQDN of the Managed Kubernetes Automatic Cluster.

* `hosted_system` - A `hosted_system` block as documented below.

* `identity` - An `identity` block as documented below.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Base64 encoded Kubernetes configuration.

* `kubelet_identity` - A `kubelet_identity` block as documented below.

* `kubernetes_version` - The version of Kubernetes used on the Managed Kubernetes Automatic Cluster.

* `node_resource_group` - Auto-generated Resource Group containing AKS Cluster resources.

* `node_resource_group_id` - The ID of the Resource Group containing the resources for this Managed Kubernetes Automatic Cluster.

* `portal_fully_qualified_domain_name` - The FQDN used by the Azure Portal for this Managed Kubernetes Automatic Cluster.

* `private_cluster` - A `private_cluster` block as documented below.

* `private_fully_qualified_domain_name` - The FQDN of this Managed Kubernetes Automatic Cluster when private link has been enabled. This name is only resolvable inside the Virtual Network where the Azure Kubernetes Service is located.

* `service_mesh` - A `service_mesh` block as documented below.

* `tags` - A mapping of tags assigned to this resource.

* `web_app_routing_ingress` - A `web_app_routing_ingress` block as documented below.

---

An `api_server_access` block exports the following:

* `authorized_ip_ranges` - A list of IP ranges authorised to access the API server.

* `subnet_id` - The ID of the subnet that the API server is accessible from.

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

A `private_cluster` block exports the following:

* `public_fully_qualified_domain_name_enabled` - If the public FQDN for this Managed Kubernetes Automatic Cluster is enabled.

* `private_dns_zone_id` - The ID of the Private DNS Zone used by this Managed Kubernetes Automatic Cluster.

---

A `service_mesh` block exports the following:

* `revisions` - List of revisions of the Istio control plane. When an upgrade is not in progress, this holds one value. When a canary upgrade is in progress, this can hold two consecutive values. [Learn More](https://learn.microsoft.com/en-us/azure/aks/istio-upgrade).

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

* `Microsoft.ContainerService` - 2026-04-01
