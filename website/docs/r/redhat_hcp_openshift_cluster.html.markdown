---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_hcp_openshift_cluster"
description: |-
  Manages an Azure Red Hat OpenShift with Hosted Control Planes Cluster.
---

# azurerm_redhat_hcp_openshift_cluster

Manages an Azure Red Hat OpenShift with Hosted Control Planes Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

locals {
  cluster_name        = "example-aro-hcp"
  resource_group_name = "example-resources"
  suffix              = substr(md5(local.cluster_name), 0, 8)
  key_vault_name      = "example-kv-${substr(md5("${data.azurerm_client_config.current.subscription_id}:${local.resource_group_name}"), 0, 13)}"

  control_plane_operator_names = {
    "cluster-api-azure"        = "cp-cluster-api-azure"
    "control-plane"            = "cp-control-plane"
    "cloud-controller-manager" = "cp-cloud-controller-manager"
    "ingress"                  = "cp-ingress"
    "disk-csi-driver"          = "cp-disk-csi-driver"
    "file-csi-driver"          = "cp-file-csi-driver"
    "image-registry"           = "cp-image-registry"
    "cloud-network-config"     = "cp-cloud-network-config"
    "kms"                      = "cp-kms"
  }

  data_plane_operator_names = {
    "disk-csi-driver" = "dp-disk-csi-driver"
    "file-csi-driver" = "dp-file-csi-driver"
    "image-registry"  = "dp-image-registry"
  }

  role_definition_ids = {
    reader                   = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"
    cluster_api_provider     = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/88366f10-ed47-4cc0-9fab-c8a06148393e"
    control_plane_operator   = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/fc0c873f-45e9-4d0d-a7d1-585aab30c6ed"
    cloud_controller_manager = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/a1f96423-95ce-4224-ab27-4e3dc72facd4"
    ingress_operator         = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/0336e1d3-7a87-462b-b6db-342b63f7802c"
    image_registry_operator  = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/8b32b316-c2f5-4ddf-b05b-83dacd2d08b5"
    file_storage_operator    = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/0d7aedc0-15fd-4a67-a412-efad370c947e"
    network_operator         = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/be7a6435-15ae-4171-8f30-4a343eff9e8f"
    key_vault_crypto_user    = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/12338af0-0e69-4776-bea7-57ae8d297424"
    federated_credential     = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/ef318e2a-8334-4a05-9e4a-295a196c6a6e"
    service_managed_identity = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/c0ff367d-66d8-445e-917c-583feb0ef0d4"
  }
}

resource "azurerm_resource_group" "example" {
  name     = local.resource_group_name
  location = "uksouth"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "worker" {
  subnet_id                 = azurerm_subnet.worker_subnet.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet" "vnet_integration_subnet" {
  name                 = "vnet-integration-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "aro-hcp-delegation"

    service_delegation {
      name = "Microsoft.RedHatOpenShift/hcpOpenShiftClusters"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "vnet_integration" {
  subnet_id                 = azurerm_subnet.vnet_integration_subnet.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_key_vault" "example" {
  name                          = local.key_vault_name
  location                      = azurerm_resource_group.example.location
  resource_group_name           = azurerm_resource_group.example.name
  tenant_id                     = data.azurerm_client_config.current.tenant_id
  sku_name                      = "standard"
  rbac_authorization_enabled    = true
  public_network_access_enabled = true
}

resource "azurerm_role_assignment" "key_vault_crypto_officer" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Crypto Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "etcd_encryption_key" {
  name         = "etcd-data-kms-encryption-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [azurerm_role_assignment.key_vault_crypto_officer]
}

resource "azurerm_user_assigned_identity" "control_plane_operator" {
  for_each = local.control_plane_operator_names

  name                = "${local.cluster_name}-${each.value}-${local.suffix}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_user_assigned_identity" "data_plane_operator" {
  for_each = local.data_plane_operator_names

  name                = "${local.cluster_name}-${each.value}-${local.suffix}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_user_assigned_identity" "service" {
  name                = "${local.cluster_name}-service-managed-identity-${local.suffix}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

locals {
  control_plane_operators = {
    for name in keys(local.control_plane_operator_names) :
    name => azurerm_user_assigned_identity.control_plane_operator[name]
  }
  data_plane_operators = {
    for name in keys(local.data_plane_operator_names) :
    name => azurerm_user_assigned_identity.data_plane_operator[name]
  }

  role_assignments = merge(
    {
      for name in keys(local.control_plane_operator_names) :
      "reader-${name}" => {
        scope              = local.control_plane_operators[name].id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.reader
      }
    },
    {
      for name in keys(local.data_plane_operator_names) :
      "federated-credential-${name}" => {
        scope              = local.data_plane_operators[name].id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.federated_credential
      }
    },
    {
      cluster-api-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["cluster-api-azure"].principal_id
        role_definition_id = local.role_definition_ids.cluster_api_provider
      }
      cluster-api-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["cluster-api-azure"].principal_id
        role_definition_id = local.role_definition_ids.cluster_api_provider
      }
      control-plane-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["control-plane"].principal_id
        role_definition_id = local.role_definition_ids.control_plane_operator
      }
      control-plane-nsg = {
        scope              = azurerm_network_security_group.example.id
        principal_id       = local.control_plane_operators["control-plane"].principal_id
        role_definition_id = local.role_definition_ids.control_plane_operator
      }
      cloud-controller-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      cloud-controller-nsg = {
        scope              = azurerm_network_security_group.example.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      cloud-controller-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      ingress-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["ingress"].principal_id
        role_definition_id = local.role_definition_ids.ingress_operator
      }
      ingress-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["ingress"].principal_id
        role_definition_id = local.role_definition_ids.ingress_operator
      }
      file-csi-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      file-csi-nsg = {
        scope              = azurerm_network_security_group.example.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      file-csi-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      image-registry-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["image-registry"].principal_id
        role_definition_id = local.role_definition_ids.image_registry_operator
      }
      network-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["cloud-network-config"].principal_id
        role_definition_id = local.role_definition_ids.network_operator
      }
      network-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.control_plane_operators["cloud-network-config"].principal_id
        role_definition_id = local.role_definition_ids.network_operator
      }
      kms-key-vault = {
        scope              = azurerm_key_vault.example.id
        principal_id       = local.control_plane_operators["kms"].principal_id
        role_definition_id = local.role_definition_ids.key_vault_crypto_user
      }
      data-plane-file-csi-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-file-csi-nsg = {
        scope              = azurerm_network_security_group.example.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-file-csi-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-image-registry-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = local.data_plane_operators["image-registry"].principal_id
        role_definition_id = local.role_definition_ids.image_registry_operator
      }
      service-vnet = {
        scope              = azurerm_virtual_network.example.id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.service_managed_identity
      }
      service-nsg = {
        scope              = azurerm_network_security_group.example.id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.service_managed_identity
      }
    }
  )
}

resource "azurerm_role_assignment" "example" {
  for_each = local.role_assignments

  scope                            = each.value.scope
  principal_id                     = each.value.principal_id
  role_definition_id               = each.value.role_definition_id
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_redhat_hcp_openshift_cluster" "example" {
  name                = local.cluster_name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.example.name
          visibility = "Public"

          active_key {
            name    = azurerm_key_vault_key.etcd_encryption_key.name
            version = azurerm_key_vault_key.etcd_encryption_key.version
          }
        }
      }
    }
  }

  platform {
    subnet_id                  = azurerm_subnet.worker_subnet.id
    vnet_integration_subnet_id = azurerm_subnet.vnet_integration_subnet.id
    network_security_group_id  = azurerm_network_security_group.example.id

    operators_authentication {
      user_assigned_identities {
        control_plane_operators = {
          for name, identity in local.control_plane_operators : name => identity.id
        }
        data_plane_operators = {
          for name, identity in local.data_plane_operators : name => identity.id
        }
        service_managed_identity = azurerm_user_assigned_identity.service.id
      }
    }
  }

  identity {
    type = "UserAssigned"
    identity_ids = concat(
      [azurerm_user_assigned_identity.service.id],
      [for identity in azurerm_user_assigned_identity.control_plane_operator : identity.id],
    )
  }

  depends_on = [
    azurerm_role_assignment.example,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}

output "cluster_id" {
  value = azurerm_redhat_hcp_openshift_cluster.example.id
}

output "console_url" {
  value = azurerm_redhat_hcp_openshift_cluster.example.console_url
}
```

## Arguments Reference

The following arguments are supported:

* `etcd` - (Required) A `etcd` block as defined below.

* `identity` - (Required) A `identity` block as defined below.

* `location` - (Required) The Azure Region where the Azure Red Hat OpenShift with Hosted Control Planes Cluster should exist. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Azure Red Hat OpenShift with Hosted Control Planes Cluster. Must be between 1 and 54 characters, begin with a letter, end with a letter or number, and contain only letters, numbers and hyphens. Changing this forces a new resource to be created.

* `platform` - (Required) A `platform` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Red Hat OpenShift with Hosted Control Planes Cluster should exist. Changing this forces a new resource to be created.

* `version` - (Required) A `version` block as defined below.

---

* `api` - (Optional) A `api` block as defined below.

* `autoscaling` - (Optional) A `autoscaling` block as defined below.

* `cluster_image_registry` - (Optional) A `cluster_image_registry` block as defined below. Changing this forces a new resource to be created.

* `crypto_restrictions` - (Optional) Cryptographic restrictions for kernel and userspace libraries. Possible values are `None` and `FIPS`. `FIPS` uses only FIPS-validated algorithms and implementations. Defaults to `None`. Changing this forces a new resource to be created.

* `dns` - (Optional) A `dns` block as defined below. Changing this forces a new resource to be created.

* `image_digest_mirrors` - (Optional) One or more `image_digest_mirrors` blocks as defined below.

* `ingress` - (Optional) A `ingress` block as defined below. Changing this forces a new resource to be created.

* `network` - (Optional) A `network` block as defined below. Changing this forces a new resource to be created.

* `node_drain_timeout_minutes` - (Optional) The grace period for how long Pod Disruption Budget-protected workloads will be respected during any node draining operation. After this grace period, any workloads protected by Pod Disruption Budgets that have not been successfully drained from a node will be forcibly evicted. This is especially relevant to cluster upgrades. Valid values are in minutes and from `0` to `10080` minutes (1 week). `0` means that the NodePool can be drained without any time limitation. This value is used as a default for all NodePools. It can be overridden by specifying `node_drain_timeout_minutes` for a given NodePool. Defaults to `0`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Red Hat OpenShift with Hosted Control Planes Cluster.

---

A `active_key` block supports the following:

* `name` - (Required) The name of the keyvault key used for encryption/decryption. Changing this forces a new resource to be created.

* `version` - (Required) The version of the key to use.

---

A `api` block supports the following:

* `authorized_cidrs` - (Optional) The list of authorized IPv4 CIDR blocks allowed to access the API server. Maximum 500 entries.

* `visibility` - (Optional) The internet visibility of the OpenShift API server. Possible values are `Public` and `Private`. Defaults to `Public`. Changing this forces a new resource to be created.

---

A `autoscaling` block supports the following:

* `max_node_provision_time_seconds` - (Optional) The maximum time to wait for node provisioning before considering the provisioning to be unsuccessful. Defaults to `900` seconds, or 15 minutes.

* `max_nodes_total` - (Optional) The maximum allowable number of nodes. The autoscaler will not grow the cluster beyond this number. If omitted, the autoscaler will not have a maximum limit.

* `max_pod_grace_period_seconds` - (Optional) The maximum amount of time in seconds to wait for graceful pod termination before scaling down a NodePool. Defaults to `600`.

* `pod_priority_threshold` - (Optional) Enables users to schedule "best-effort" pods, which shouldn't trigger autoscaler actions, but only run when there are spare resources available. Defaults to `-10`.

---

A `cluster_image_registry` block supports the following:

* `state` - (Optional) The desired ImageStream-backed cluster image registry installation mode. Possible values are `Enabled` and `Disabled`. `Enabled` means the ImageStream-backed image registry will be run as pods on worker nodes in the cluster. `Disabled` means the ImageStream-backed image registry will not be present in the cluster. This can only be set during cluster creation and cannot be changed after cluster creation. Defaults to `Enabled`.

---

A `container_registry` block supports the following:

* `managed_identity` - (Required) The Azure resource ID of the user-assigned managed identity used for container registry image pulls.

---

A `customer_managed` block supports the following:

* `kms` - (Required) A `kms` block as defined below.

* `encryption_type` - (Optional) The encryption type used. The only possible value at this time is `KMS`. Defaults to `KMS`. Changing this forces a new resource to be created.

---

A `data_encryption` block supports the following:

* `key_management_mode` - (Required) Specify the key management strategy used for the encryption key that encrypts the ETCD data. The only possible value is `CustomerManaged`. Changing this forces a new resource to be created.

* `customer_managed` - (Optional) A `customer_managed` block as defined above. Required when `key_management_mode` is `CustomerManaged` and must not be set for any other mode.

---

A `dns` block supports the following:

* `base_domain_prefix` - (Optional) The base domain prefix for the cluster ingress. Must be between 1 and 15 characters, begin with a lowercase letter, end with a lowercase letter or number, and contain only lowercase letters, numbers and hyphens. It will be used to configure ingress in the cluster through the subdomain "BaseDomainPrefix.BaseDomain". Changing this forces a new resource to be created.

---

A `etcd` block supports the following:

* `data_encryption` - (Required) A `data_encryption` block as defined above.

---

A `identity` block supports the following:

* `type` - (Required) The type of managed identity assigned to the cluster. Must be `UserAssigned`, as the service currently only supports User Assigned Managed Identities.

* `identity_ids` - (Optional) Specifies a set of User Assigned Managed Identity IDs to be assigned to the cluster.

---

A `image_digest_mirrors` block supports the following:

* `mirrors` - (Required) A list of locations that may also contain the same images. Images can be pulled from these mirrors only if they are referenced by their digests. The mirrored location is obtained by replacing the part of the input reference that matches `source` by the mirrors entry. The order of mirrors in this list is treated as the user's desired priority, while `source` is by default considered lower priority than all mirrors. If all image pulls from the mirror list fail, the image will continue to be pulled from the repository in the pull spec. Other cluster configuration may impact the exact order mirrors are contacted in, or some mirrors may be contacted in parallel, so this should be considered a preference rather than a guarantee of ordering. Each entry uses one of the following formats: `host[:port]`, `host[:port]/namespace[/namespace...]` or `host[:port]/namespace[/namespace...]/repo`.

* `source` - (Required) The repository that users refer to, e.g. in image pull specifications. Setting `source` to a registry hostname, e.g. `docker.io`, `quay.io`, or `registry.redhat.io`, will match the image pull specification of the corresponding registry. Uses one of the following formats: `host[:port]`, `host[:port]/namespace[/namespace...]`, `host[:port]/namespace[/namespace...]/repo` or `[*.]host`.

---

A `ingress` block supports the following:

* `type` - (Optional) The type of the default cluster ingress. Possible values are `Public`, `Private` and `Disabled`. Defaults to `Public`. Changing this forces a new resource to be created.

---

A `kms` block supports the following:

* `active_key` - (Required) A `active_key` block as defined above.

* `vault_name` - (Required) The name of the keyvault that contains the secret. Changing this forces a new resource to be created.

* `visibility` - (Required) The visibility of the keyvault that contains the secret. Possible values are `Public` and `Private`. Changing this forces a new resource to be created.

* `key_vault_type` - (Optional) The type of keyvault used for KMS encryption. Possible values are `KeyVault` (Standard Azure Key Vault) and `ManagedHSM` (Azure Managed Hardware Security Module). Defaults to `KeyVault`. Changing this forces a new resource to be created.

---

A `network` block supports the following:

* `host_prefix` - (Optional) Network host prefix. Must be between `0` and `32`. Defaults to `23`. Changing this forces a new resource to be created.

* `machine_cidr` - (Optional) The IPv4 CIDR block from which to assign machine IP addresses. Defaults to `10.0.0.0/16`. Changing this forces a new resource to be created.

* `network_type` - (Optional) The main controller responsible for rendering the core networking components. Possible values are `OVNKubernetes` and `Other`. Defaults to `OVNKubernetes`. Changing this forces a new resource to be created.

* `pod_cidr` - (Optional) The IPv4 CIDR of the pod IP addresses. Defaults to `10.128.0.0/14`. Changing this forces a new resource to be created.

* `service_cidr` - (Optional) The IPv4 CIDR block for assigned service IPs. Defaults to `172.30.0.0/16`. Changing this forces a new resource to be created.

---

A `operators_authentication` block supports the following:

* `user_assigned_identities` - (Required) A `user_assigned_identities` block as defined below.

---

A `platform` block supports the following:

* `network_security_group_id` - (Required) Resource Id for the NSG (network security group) attached to the cluster subnet. Note that NSGs cannot be reused for other ARO-HCP clusters. Changing this forces a new resource to be created.

* `operators_authentication` - (Required) A `operators_authentication` block as defined above.

* `subnet_id` - (Required) The Azure resource ID of the worker subnet. Note that a subnet cannot be reused between ARO-HCP Clusters. Changing this forces a new resource to be created.

* `vnet_integration_subnet_id` - (Required) The Azure resource ID of a subnet that enables direct, private network connectivity between the hosted control plane and your cluster's nodes. This subnet must be dedicated to ARO HCP and cannot be shared with the cluster subnet or any node pool subnets. Changing this forces a new resource to be created.

* `container_registry` - (Optional) A `container_registry` block as defined above.

* `managed_resource_group` - (Optional) The name of the resource group to put cluster resources in. If not specified then a unique name is generated from the pattern `aro-hcp-<clusterName>-<UUID>`, where `clusterName` means the `hcpOpenShiftClusters` resource name (up to 45 characters) followed by a 16-byte universally unique identifier per RFC 4122. Changing this forces a new resource to be created.

* `outbound_type` - (Optional) The core outgoing configuration. The only possible value is `LoadBalancer`. Defaults to `LoadBalancer`. Changing this forces a new resource to be created.

---

A `user_assigned_identities` block supports the following:

* `control_plane_operators` - (Required) A mapping of operator names to the Azure User-Assigned Managed Identity resource IDs leveraged for the Control Plane operators of the cluster. The set of required managed identities is dependent on the Cluster's OpenShift version.

* `data_plane_operators` - (Required) A mapping of operator names to the Azure User-Assigned Managed Identity resource IDs leveraged for the Data Plane operators of the cluster. The set of required managed identities is dependent on the Cluster's OpenShift version.

* `service_managed_identity` - (Required) The Azure resource ID of an Azure User-Assigned Managed Identity whose purpose is to perform service level actions.

---

A `version` block supports the following:

* `id` - (Required) The desired X.Y version of the cluster control plane.

* `channel_group` - (Optional) The name of the set to which this version belongs. Each version belongs to only a single set. Possible values are `nightly`, `candidate`, `fast` and `stable`. Defaults to `stable`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Red Hat OpenShift with Hosted Control Planes Cluster.

* `console_url` - The cluster web console URL endpoint.

---

An `api` block exports the following:

* `url` - The URL endpoint of the OpenShift API server.

---

A `dns` block exports the following:

* `base_domain` - The base domain of the cluster.

---

A `platform` block exports the following:

* `issuer_url` - The OIDC issuer URL of the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Azure Red Hat OpenShift with Hosted Control Planes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Red Hat OpenShift with Hosted Control Planes Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the Azure Red Hat OpenShift with Hosted Control Planes Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the Azure Red Hat OpenShift with Hosted Control Planes Cluster.

## Import

Azure Red Hat OpenShift with Hosted Control Planes Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhat_hcp_openshift_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.RedHatOpenShift` - 2026-10-01
