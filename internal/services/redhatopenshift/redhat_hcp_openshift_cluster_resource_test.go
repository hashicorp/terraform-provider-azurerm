// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedhatHcpOpenshiftClusterResource struct{}

func TestAccRedhatHcpOpenshiftCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccRedhatHcpOpenshiftCluster_private(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.private(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.visibility").HasValue("Private"),
				check.That(data.ResourceName).Key("ingress.0.type").HasValue("Private"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_fips(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fips(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("crypto_restrictions").HasValue("FIPS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_customManagedResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customManagedResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.managed_resource_group").HasValue(fmt.Sprintf("acctestManagedRG-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_clusterImageRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clusterImageRegistry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cluster_image_registry.0.state").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_dnsBaseDomainPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsBaseDomainPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dns.0.base_domain_prefix").HasValue(fmt.Sprintf("hcp-%s", data.RandomString)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_containerRegistryIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerRegistryIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.container_registry.0.managed_identity").MatchesOtherKey(check.That("azurerm_user_assigned_identity.container_registry").Key("id")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_autoscaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaling(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("autoscaling.0.max_nodes_total").HasValue("100"),
				check.That(data.ResourceName).Key("autoscaling.0.max_pod_grace_period_seconds").HasValue("300"),
				check.That(data.ResourceName).Key("autoscaling.0.max_node_provision_time_seconds").HasValue("1200"),
				check.That(data.ResourceName).Key("autoscaling.0.pod_priority_threshold").HasValue("-5"),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("version.0.channel_group").HasValue("fast"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_network(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.network(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.authorized_cidrs.#").HasValue("2"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "api.0.authorized_cidrs.*", "192.0.2.0/24"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "api.0.authorized_cidrs.*", "198.51.100.0/24"),
				check.That(data.ResourceName).Key("network.0.network_type").HasValue("OVNKubernetes"),
				check.That(data.ResourceName).Key("network.0.pod_cidr").HasValue("10.132.0.0/14"),
				check.That(data.ResourceName).Key("network.0.service_cidr").HasValue("172.31.0.0/16"),
				check.That(data.ResourceName).Key("network.0.machine_cidr").HasValue("10.0.0.0/16"),
				check.That(data.ResourceName).Key("network.0.host_prefix").HasValue("24"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_imageDigestMirrors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageDigestMirrors(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("image_digest_mirrors.#").HasValue("1"),
				check.That(data.ResourceName).Key("image_digest_mirrors.0.source").HasValue("quay.io/example/app"),
				check.That(data.ResourceName).Key("image_digest_mirrors.0.mirrors.#").HasValue("2"),
				check.That(data.ResourceName).Key("image_digest_mirrors.0.mirrors.0").HasValue("mirror1.example.com/example/app"),
				check.That(data.ResourceName).Key("image_digest_mirrors.0.mirrors.1").HasValue("mirror2.example.com/example/app"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedhatHcpOpenshiftCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_cluster", "test")
	r := RedhatHcpOpenshiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerRegistryIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.container_registry.0.managed_identity").MatchesOtherKey(check.That("azurerm_user_assigned_identity.container_registry").Key("id")),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.containerRegistryIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.container_registry.0.managed_identity").MatchesOtherKey(check.That("azurerm_user_assigned_identity.container_registry_updated").Key("id")),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("30"),
			),
		},
		data.ImportStep(),
		{
			Config: r.containerRegistryIdentityUpdateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.container_registry.0.managed_identity").MatchesOtherKey(check.That("azurerm_user_assigned_identity.container_registry_updated").Key("id")),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("updated"),
			),
		},
		data.ImportStep(),
		{
			Config: r.containerRegistryIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.container_registry.0.managed_identity").MatchesOtherKey(check.That("azurerm_user_assigned_identity.container_registry").Key("id")),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (t RedhatHcpOpenshiftClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RedHatOpenShift.HcpOpenShiftClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Red Hat HCP OpenShift Cluster (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r RedhatHcpOpenshiftClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "import" {
  name                = azurerm_redhat_hcp_openshift_cluster.test.name
  location            = azurerm_redhat_hcp_openshift_cluster.test.location
  resource_group_name = azurerm_redhat_hcp_openshift_cluster.test.resource_group_name

  version {
    id = azurerm_redhat_hcp_openshift_cluster.test.version.0.id
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    subnet_id                  = azurerm_redhat_hcp_openshift_cluster.test.platform.0.subnet_id
    vnet_integration_subnet_id = azurerm_redhat_hcp_openshift_cluster.test.platform.0.vnet_integration_subnet_id
    network_security_group_id  = azurerm_redhat_hcp_openshift_cluster.test.platform.0.network_security_group_id

    operators_authentication {
      user_assigned_identities {
        control_plane_operators  = azurerm_redhat_hcp_openshift_cluster.test.platform.0.operators_authentication.0.user_assigned_identities.0.control_plane_operators
        data_plane_operators     = azurerm_redhat_hcp_openshift_cluster.test.platform.0.operators_authentication.0.user_assigned_identities.0.data_plane_operators
        service_managed_identity = azurerm_redhat_hcp_openshift_cluster.test.platform.0.operators_authentication.0.user_assigned_identities.0.service_managed_identity
      }
    }
  }

  identity {
    type         = azurerm_redhat_hcp_openshift_cluster.test.identity.0.type
    identity_ids = azurerm_redhat_hcp_openshift_cluster.test.identity.0.identity_ids
  }
}
  `, r.basic(data))
}

func (r RedhatHcpOpenshiftClusterResource) private(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  api {
    visibility = "Private"
  }

  ingress {
    type = "Private"
  }

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) fips(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  crypto_restrictions = "FIPS"

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) customManagedResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    managed_resource_group     = "acctestManagedRG-%[2]d"
    subnet_id                  = azurerm_subnet.worker_subnet.id
    vnet_integration_subnet_id = azurerm_subnet.vnet_integration_subnet.id
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) clusterImageRegistry(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_image_registry {
    state = "Disabled"
  }

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) dnsBaseDomainPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns {
    base_domain_prefix = "hcp-%[3]s"
  }

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r RedhatHcpOpenshiftClusterResource) containerRegistryIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "container_registry" {
  name                = "acctest-aro-hcp-%[2]d-container-registry"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "container_registry_reader" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "container_registry_managed_identity_operator" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830"
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

    container_registry {
      managed_identity = azurerm_user_assigned_identity.container_registry.id
    }

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
    azurerm_role_assignment.test,
    azurerm_role_assignment.container_registry_reader,
    azurerm_role_assignment.container_registry_managed_identity_operator,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) autoscaling(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                       = "acctestarohcp%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  node_drain_timeout_minutes = 30

  autoscaling {
    max_nodes_total                 = 100
    max_pod_grace_period_seconds    = 300
    max_node_provision_time_seconds = 1200
    pod_priority_threshold          = -5
  }

  version {
    id            = "4.22"
    channel_group = "fast"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) containerRegistryIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "container_registry" {
  name                = "acctest-aro-hcp-%[2]d-container-registry"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "container_registry_reader" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "container_registry_managed_identity_operator" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830"
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_user_assigned_identity" "container_registry_updated" {
  name                = "acctest-aro-hcp-%[2]d-container-registry-updated"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "container_registry_updated_reader" {
  scope                            = azurerm_user_assigned_identity.container_registry_updated.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "container_registry_updated_managed_identity_operator" {
  scope                            = azurerm_user_assigned_identity.container_registry_updated.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830"
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                       = "acctestarohcp%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  node_drain_timeout_minutes = 30

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

    container_registry {
      managed_identity = azurerm_user_assigned_identity.container_registry_updated.id
    }

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
    azurerm_role_assignment.test,
    azurerm_role_assignment.container_registry_reader,
    azurerm_role_assignment.container_registry_managed_identity_operator,
    azurerm_role_assignment.container_registry_updated_reader,
    azurerm_role_assignment.container_registry_updated_managed_identity_operator,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) containerRegistryIdentityUpdateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "container_registry" {
  name                = "acctest-aro-hcp-%[2]d-container-registry"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "container_registry_reader" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "container_registry_managed_identity_operator" {
  scope                            = azurerm_user_assigned_identity.container_registry.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830"
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_user_assigned_identity" "container_registry_updated" {
  name                = "acctest-aro-hcp-%[2]d-container-registry-updated"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "container_registry_updated_reader" {
  scope                            = azurerm_user_assigned_identity.container_registry_updated.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "container_registry_updated_managed_identity_operator" {
  scope                            = azurerm_user_assigned_identity.container_registry_updated.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830"
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                       = "acctestarohcp%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  node_drain_timeout_minutes = 30

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

    container_registry {
      managed_identity = azurerm_user_assigned_identity.container_registry_updated.id
    }

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

  tags = {
    environment = "updated"
  }

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_role_assignment.container_registry_reader,
    azurerm_role_assignment.container_registry_managed_identity_operator,
    azurerm_role_assignment.container_registry_updated_reader,
    azurerm_role_assignment.container_registry_updated_managed_identity_operator,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) network(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  api {
    authorized_cidrs = ["192.0.2.0/24", "198.51.100.0/24"]
  }

  network {
    network_type = "OVNKubernetes"
    pod_cidr     = "10.132.0.0/14"
    service_cidr = "172.31.0.0/16"
    machine_cidr = "10.0.0.0/16"
    host_prefix  = 24
  }

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r RedhatHcpOpenshiftClusterResource) imageDigestMirrors(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_cluster" "test" {
  name                = "acctestarohcp%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  image_digest_mirrors {
    source = "quay.io/example/app"
    mirrors = [
      "mirror1.example.com/example/app",
      "mirror2.example.com/example/app",
    ]
  }

  version {
    id = "4.22"
  }

  etcd {
    data_encryption {
      key_management_mode = "CustomerManaged"

      customer_managed {
        kms {
          vault_name = azurerm_key_vault.test.name
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
    network_security_group_id  = azurerm_network_security_group.test.id

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
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.worker,
    azurerm_subnet_network_security_group_association.vnet_integration,
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (RedhatHcpOpenshiftClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

locals {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aro-hcp-%[2]d"
  location = "%[1]s"

  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "test-nsg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "worker" {
  subnet_id                 = azurerm_subnet.worker_subnet.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_subnet" "vnet_integration_subnet" {
  name                 = "vnet-integration-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
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
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_key_vault" "test" {
  name                          = "acctestkv-%[3]s"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  tenant_id                     = data.azurerm_client_config.current.tenant_id
  sku_name                      = "standard"
  rbac_authorization_enabled    = true
  public_network_access_enabled = true
}

resource "azurerm_role_assignment" "key_vault_crypto_officer" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Key Vault Crypto Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "etcd_encryption_key" {
  name         = "etcd-data-kms-encryption-key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [azurerm_role_assignment.key_vault_crypto_officer]
}

resource "azurerm_user_assigned_identity" "control_plane_operator" {
  count = length(local.control_plane_operator_names)

  name                = "acctest-aro-hcp-%[2]d-${values(local.control_plane_operator_names)[count.index]}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "data_plane_operator" {
  count = length(local.data_plane_operator_names)

  name                = "acctest-aro-hcp-%[2]d-${values(local.data_plane_operator_names)[count.index]}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "service" {
  name                = "acctest-aro-hcp-%[2]d-service-managed-identity"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

locals {
  control_plane_operators = {
    for index, name in keys(local.control_plane_operator_names) :
    name => azurerm_user_assigned_identity.control_plane_operator[index]
  }
  data_plane_operators = {
    for index, name in keys(local.data_plane_operator_names) :
    name => azurerm_user_assigned_identity.data_plane_operator[index]
  }

  role_assignments = merge(
    {
      for name, identity in local.control_plane_operators :
      "reader-${name}" => {
        scope              = identity.id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.reader
      }
    },
    {
      for name, identity in local.data_plane_operators :
      "federated-credential-${name}" => {
        scope              = identity.id
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
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["cluster-api-azure"].principal_id
        role_definition_id = local.role_definition_ids.cluster_api_provider
      }
      control-plane-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["control-plane"].principal_id
        role_definition_id = local.role_definition_ids.control_plane_operator
      }
      control-plane-nsg = {
        scope              = azurerm_network_security_group.test.id
        principal_id       = local.control_plane_operators["control-plane"].principal_id
        role_definition_id = local.role_definition_ids.control_plane_operator
      }
      cloud-controller-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      cloud-controller-nsg = {
        scope              = azurerm_network_security_group.test.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      cloud-controller-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["cloud-controller-manager"].principal_id
        role_definition_id = local.role_definition_ids.cloud_controller_manager
      }
      ingress-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["ingress"].principal_id
        role_definition_id = local.role_definition_ids.ingress_operator
      }
      ingress-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["ingress"].principal_id
        role_definition_id = local.role_definition_ids.ingress_operator
      }
      file-csi-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      file-csi-nsg = {
        scope              = azurerm_network_security_group.test.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      file-csi-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      image-registry-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["image-registry"].principal_id
        role_definition_id = local.role_definition_ids.image_registry_operator
      }
      network-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.control_plane_operators["cloud-network-config"].principal_id
        role_definition_id = local.role_definition_ids.network_operator
      }
      network-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.control_plane_operators["cloud-network-config"].principal_id
        role_definition_id = local.role_definition_ids.network_operator
      }
      kms-key-vault = {
        scope              = azurerm_key_vault.test.id
        principal_id       = local.control_plane_operators["kms"].principal_id
        role_definition_id = local.role_definition_ids.key_vault_crypto_user
      }
      data-plane-file-csi-subnet = {
        scope              = azurerm_subnet.worker_subnet.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-file-csi-nsg = {
        scope              = azurerm_network_security_group.test.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-file-csi-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.data_plane_operators["file-csi-driver"].principal_id
        role_definition_id = local.role_definition_ids.file_storage_operator
      }
      data-plane-image-registry-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = local.data_plane_operators["image-registry"].principal_id
        role_definition_id = local.role_definition_ids.image_registry_operator
      }
      service-vnet = {
        scope              = azurerm_virtual_network.test.id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.service_managed_identity
      }
      service-nsg = {
        scope              = azurerm_network_security_group.test.id
        principal_id       = azurerm_user_assigned_identity.service.principal_id
        role_definition_id = local.role_definition_ids.service_managed_identity
      }
    }
  )
}

resource "azurerm_role_assignment" "test" {
  count = length(local.role_assignments)

  scope                            = values(local.role_assignments)[count.index].scope
  principal_id                     = values(local.role_assignments)[count.index].principal_id
  role_definition_id               = values(local.role_assignments)[count.index].role_definition_id
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}
 `, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
