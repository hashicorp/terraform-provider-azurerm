---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_flux_configuration"
description: |-
  Manages a Kubernetes Flux Configuration.
---

# azurerm_kubernetes_flux_configuration

Manages a Kubernetes Flux Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example-aks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_extension" "example" {
  name           = "example-ext"
  cluster_id     = azurerm_kubernetes_cluster.test.id
  extension_type = "microsoft.flux"
}

resource "azurerm_kubernetes_flux_configuration" "example" {
  name       = "example-fc"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  git_repository {
    url             = "https://github.com/Azure/arc-k8s-demo"
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Kubernetes Flux Configuration. Changing this forces a new Kubernetes Flux Configuration to be created.

* `cluster_id` - (Required) Specifies the Cluster ID. Changing this forces a new Kubernetes Cluster Extension to be created.

* `kustomizations` - (Required) A `kustomizations` block as defined below.

* `namespace` - (Required) Specifies the namespace to which this configuration is installed to. Changing this forces a new Kubernetes Flux Configuration to be created.

* `blob_storage` - (Optional) An `blob_storage` block as defined below.

* `bucket` - (Optional) A `bucket` block as defined below.

* `git_repository` - (Optional) A `git_repository` block as defined below.

* `scope` - (Optional) Specifies the scope at which the operator will be installed. Possible values are `cluster` and `namespace`. Defaults to `namespace`. Changing this forces a new Kubernetes Flux Configuration to be created.

* `continuous_reconciliation_enabled` - (Optional) Whether the configuration will keep its reconciliation of its kustomizations and sources with the repository. Defaults to `true`.

---

A `kustomizations` block supports the following:

* `name` - (Required) Specifies the name of the kustomization.

* `path` - (Optional) Specifies the path in the source reference to reconcile on the cluster.

* `timeout_in_seconds` - (Optional) The maximum time to attempt to reconcile the kustomization on the cluster. Defaults to `600`.

* `sync_interval_in_seconds` - (Optional) The interval at which to re-reconcile the kustomization on the cluster. Defaults to `600`.

* `retry_interval_in_seconds` - (Optional) The interval at which to re-reconcile the kustomization on the cluster in the event of failure on reconciliation. Defaults to `600`.

* `recreating_enabled` - (Optional) Whether re-creating Kubernetes resources on the cluster is enabled when patching fails due to an immutable field change. Defaults to `false`.

* `garbage_collection_enabled` - (Optional) Whether garbage collections of Kubernetes objects created by this kustomization is enabled. Defaults to `false`.

* `depends_on` - (Optional) Specifies other kustomizations that this kustomization depends on. This kustomization will not reconcile until all dependencies have completed their reconciliation.

---

An `blob_storage` block supports the following:

* `container_id` - (Required) Specifies the Azure Blob container ID.

* `account_key` - (Optional) Specifies the account key (shared key) to access the storage account.

* `local_auth_reference` - (Optional) Specifies the name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.

* `managed_identity` - (Optional) A `managed_identity` block as defined below.

* `sas_token` - (Optional) Specifies the shared access token to access the storage container.

* `service_principal` - (Optional) A `service_principal` block as defined below.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster Azure Blob source with the remote.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster Azure Blob source with the remote.

---

A `managed_identity` block supports the following:

* `client_id` - (Required) Specifies the client ID for authenticating a Managed Identity.

---

A `service_principal` block supports the following:

* `client_id` - (Required) Specifies the client ID for authenticating a Service Principal.

* `tenant_id` - (Required) Specifies the tenant ID for authenticating a Service Principal.

* `client_certificate_base64` - (Optional) Base64-encoded certificate used to authenticate a Service Principal .

* `client_certificate_password` - (Optional) Specifies the password for the certificate used to authenticate a Service Principal .

* `client_certificate_send_chain` - (Optional) Specifies whether to include x5c header in client claims when acquiring a token to enable subject name / issuer based authentication for the client certificate.

* `client_secret` - (Optional) Specifies the client secret for authenticating a Service Principal.

---

A `bucket` block supports the following:

* `bucket_name` - (Required) Specifies the bucket name to sync from the url endpoint for the flux configuration.

* `url` - (Required) Specifies the URL to sync for the flux configuration S3 bucket. It must start with `http://` or `https://`.

* `access_key` - (Optional) Specifies the plaintext access key used to securely access the S3 bucket.

* `secret_key_base64` - (Optional) Specifies the Base64-encoded secret key used to authenticate with the bucket source.

* `tls_enabled` - (Optional) Specify whether to communicate with a bucket using TLS is enabled. Defaults to `true`.

* `local_auth_reference` - (Optional) Specifies the name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster git repository source with the remote. Defaults to `600`.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster git repository source with the remote. Defaults to `600`.

---

A `git_repository` block supports the following:

* `url` - (Required) Specifies the URL to sync for the flux configuration git repository. It must start with `http://`, `https://`, `git@` or `ssh://`.

* `reference_type` - (Required) Specifies the source reference type for the GitRepository object. Possible values are `branch`, `commit`, `semver` and `tag`.

* `reference_value` - (Required) Specifies the source reference value for the GitRepository object.

* `https_ca_cert_base64` - (Optional) Specifies the Base64-encoded HTTPS certificate authority contents used to access git private git repositories over HTTPS.

* `https_user` - (Optional) Specifies the plaintext HTTPS username used to access private git repositories over HTTPS.

* `https_key_base64` - (Optional) Specifies the Base64-encoded HTTPS personal access token or password that will be used to access the repository.

* `local_auth_reference` - (Optional) Specifies the name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `ssh_private_key_base64` - (Optional) Specifies the Base64-encoded SSH private key in PEM format.

* `ssh_known_hosts_base64` - (Optional) Specifies the Base64-encoded known_hosts value containing public SSH keys required to access private git repositories over SSH.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster git repository source with the remote. Defaults to `600`.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster git repository source with the remote. Defaults to `600`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Flux Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Flux Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Flux Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Flux Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Flux Configuration.

## Import

Kubernetes Flux Configuration can be imported using the `resource id` for different `cluster_resource_name`, e.g.

```shell
terraform import azurerm_kubernetes_flux_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfiguration1
```
