---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_automatic_cluster"
description: |-
  Manages a managed Kubernetes Automatic Cluster (a special SKU of AKS / Azure Kubernetes Service)
---

# azurerm_kubernetes_automatic_cluster

Manages a Managed Kubernetes Automatic Cluster (a special SKU of AKS / Azure Kubernetes Service)

-> **Note:** Due to the fast-moving nature of AKS, we recommend using the latest version of the Azure Provider when using AKS - you can find [the latest version of the Azure Provider here](https://registry.terraform.io/providers/hashicorp/azurerm/latest).

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_automatic_cluster" "example" {
  name                = "example-aks1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Production"
  }
}

output "client_certificate" {
  value     = azurerm_kubernetes_automatic_cluster.example.kube_config[0].client_certificate
  sensitive = true
}

output "kube_config" {
  value = azurerm_kubernetes_automatic_cluster.example.kube_config_raw

  sensitive = true
}
```

## Bring your own networking example

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "node" {
  name                 = "example-node-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "api" {
  name                 = "example-api-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.1.1.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_subnet" "systemnode" {
  name                 = "example-systemnode-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.1.2.0/24"]

  lifecycle {
    ignore_changes = [
      delegation
    ]
  }
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  name                = "example_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.example.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "example" {
  name                = "example-aks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name


  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Kubernetes Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Managed Kubernetes Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Managed Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

---

* `api_server_access` - (Optional) An `api_server_access` block as defined below.

* `hosted_system` - (Optional) A `hosted_system` block as defined below.

* `private_cluster` - (Optional) A `private_cluster` block as defined below.

* `web_app_routing_ingress` - (Optional) A `web_app_routing_ingress` block as defined below.

* `service_mesh` - (Optional) A `service_mesh` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `api_server_access` block supports the following:

* `authorized_ip_ranges` - (Optional) Set of authorized IP ranges to allow access to API server, e.g. ["198.51.100.0/24"].

* `subnet_id` - (Optional) The ID of the Subnet where the API server endpoint is delegated to. Is required for bring your own networking.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Kubernetes Cluster. Possible values are `SystemAssigned` or `UserAssigned`.  `UserAssigned` is required for bring your own networking

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Kubernetes Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `service_mesh` block supports the following:

* `revisions` - (Required) Specify `1` or `2` Istio control plane revisions for managing minor upgrades using the canary upgrade process. For example, create the resource with `revisions` set to `["asm-1-27"]`. To start the canary upgrade, change `revisions` to `["asm-1-27", "asm-1-28"]`. To roll back the canary upgrade, revert to `["asm-1-27"]`. To confirm the upgrade, change to `["asm-1-28"]`.

-> **Note:** Upgrading to a new (canary) revision does not affect existing sidecar proxies. You need to apply the canary revision label to selected namespaces and restart pods with kubectl to inject the new sidecar proxy. [Learn more](https://istio.io/latest/docs/setup/upgrade/canary/#data-plane).

* `internal_ingress_gateway_enabled` - (Optional) Enables Istio Internal Ingress Gateway. Defaults to `false`.

* `external_ingress_gateway_enabled` - (Optional) Enables Istio External Ingress Gateway. Defaults to `false`.

-> **Note:** Currently only one Internal Ingress Gateway and one External Ingress Gateway are allowed per cluster

* `proxy_redirect_mechanism` - (Optional) The mechanism used to redirect application traffic to the Istio sidecar proxy. Possible values are `CNIChaining` and `InitContainers`. Defaults to `InitContainers`.

* `certificate_authority` - (Optional) A `certificate_authority` block as defined below. This configuration allows you to bring your own root certificate and keys for Istio CA in the Istio-based service mesh add-on for Azure Kubernetes Service.

---

A `certificate_authority` block supports the following:

* `key_vault_id` - (Required) The resource ID of the Key Vault.

* `root_certificate_object_name` - (Required) The root certificate object name in Azure Key Vault.

* `certificate_chain_object_name` - (Required) The certificate chain object name in Azure Key Vault.

* `certificate_object_name` - (Required) The intermediate certificate object name in Azure Key Vault.

* `key_object_name` - (Required) The intermediate certificate private key object name in Azure Key Vault.

-> **Note:** For more information on [Istio-based service mesh add-on with plug-in CA certificates and how to generate these certificates](https://learn.microsoft.com/en-us/azure/aks/istio-plugin-ca),

---

A `web_app_routing_ingress` block supports the following:

* `dns_zone_ids` - (Optional) Resource IDs of the DNS zones to be associated with the Application Routing add-on. Public and private DNS zones can be in different resource groups, but all public DNS zones must be in the same resource group and all private DNS zones must be in the same resource group.

* `default_nginx_controller` - (Optional) Specifies the ingress type for the default `NginxIngressController` custom resource. The allowed values are `Internal`, `External` and `AnnotationControlled`. At least one of `default_nginx_controller` or `istio_enabled` must be specified.

* `istio_enabled` - (Optional) Enables Istio as a Gateway API implementation. Defaults to `false`. At least one of `default_nginx_controller` or `istio_enabled` must be specified.

---

A `hosted_system` block supports the following:

* `node_subnet_id` - (Required) The ID of the Subnet where the user nodes are hosted. Is required for bring your own networking

* `system_node_subnet_id` - (Required) The ID of the Subnet where the system nodes are hosted. Changing this forces a new resource to be created. Is required for bring your own networking

---

A `private_cluster` block supports the following:

* `public_fully_qualified_domain_name_enabled` - (Optional) Provisions a Public FQDN for the private cluster. Defaults to `false`.

* `private_dns_zone_id` - (Optional) The ID of the Private DNS Zone which should be used for this Kubernetes Cluster. Possible values are `System`, `None` or the ID of a Private DNS Zone. Defaults to `System`. Changing this forces a new resource to be created.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `current_kubernetes_version` - The current version running on the Azure Kubernetes Managed Cluster.

* `fully_qualified_domain_name` - The FQDN of the Azure Kubernetes Managed Cluster.

* `private_fully_qualified_domain_name` - The FQDN for the Kubernetes Cluster when private link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.

* `portal_fully_qualified_domain_name` - The FQDN for the Azure Portal resources when private link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools.

* `oidc_issuer_url` - The OIDC issuer URL that is associated with the cluster.

* `node_resource_group_id` - The ID of the Resource Group containing the resources for this Managed Kubernetes Cluster.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `kube_config` blocks export the following:

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **Note:** It's possible to use these credentials with [the Kubernetes Provider](/providers/hashicorp/kubernetes/latest/docs) like so:

```hcl
provider "kubernetes" {
  host                   = azurerm_kubernetes_automatic_cluster.main.kube_config[0].host
  username               = azurerm_kubernetes_automatic_cluster.main.kube_config[0].username
  password               = azurerm_kubernetes_automatic_cluster.main.kube_config[0].password
  client_certificate     = base64decode(azurerm_kubernetes_automatic_cluster.main.kube_config[0].client_certificate)
  client_key             = base64decode(azurerm_kubernetes_automatic_cluster.main.kube_config[0].client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_automatic_cluster.main.kube_config[0].cluster_ca_certificate)
}
```

---

The `web_app_routing_identity` block exports the following:

* `client_id` - The Client ID of the user-defined Managed Identity used for Web App Routing.

* `object_id` - The Object ID of the user-defined Managed Identity used for Web App Routing

* `user_assigned_identity_id` - The ID of the User Assigned Identity used for Web App Routing.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster.
* `update` - (Defaults to 90 minutes) Used when updating the Kubernetes Cluster.
* `delete` - (Defaults to 90 minutes) Used when deleting the Kubernetes Cluster.

## Import

Managed Kubernetes Automatic Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_automatic_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2026-04-01
