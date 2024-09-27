// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OpenShiftClusterResource struct{}

func TestAccOpenShiftCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_private(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.private(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_userDefinedRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userDefinedRouting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_encryptionAtHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionAtHost(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_preconfiguredNetworkSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.preconfiguredNetworkSecurityGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_pullSecret(t *testing.T) {
	// the pull secret can be generated from https://console.redhat.com/openshift/install/pull-secret
	pullSecret := os.Getenv("ARM_TEST_ARO_PULL_SECRET")
	if pullSecret == "" {
		t.Skip("skip the test due to missing environment variable ARM_TEST_ARO_PULL_SECRET")
	}

	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pullSecret(data, pullSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret", "cluster_profile.0.pull_secret"),
	})
}

func TestAccOpenShiftCluster_basicWithFipsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithFipsEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccOpenShiftCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

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

func TestAccOpenShiftCluster_basicResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicResourceGroupName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func (t OpenShiftClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := openshiftclusters.ParseProviderOpenShiftClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RedHatOpenShift.OpenShiftClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Red Hat Openshift Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r OpenShiftClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
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
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %[1]s

resource "azurerm_redhat_openshift_cluster" "import" {
  name                = azurerm_redhat_openshift_cluster.test.name
  resource_group_name = azurerm_redhat_openshift_cluster.test.resource_group_name
  location            = azurerm_redhat_openshift_cluster.test.location

  cluster_profile {
    domain  = azurerm_redhat_openshift_cluster.test.cluster_profile.0.domain
    version = azurerm_redhat_openshift_cluster.test.cluster_profile.0.version
  }

  network_profile {
    pod_cidr     = azurerm_redhat_openshift_cluster.test.network_profile.0.pod_cidr
    service_cidr = azurerm_redhat_openshift_cluster.test.network_profile.0.service_cidr
  }

  main_profile {
    vm_size   = azurerm_redhat_openshift_cluster.test.main_profile.0.vm_size
    subnet_id = azurerm_redhat_openshift_cluster.test.main_profile.0.subnet_id
  }

  api_server_profile {
    visibility = azurerm_redhat_openshift_cluster.test.api_server_profile.0.visibility
  }

  ingress_profile {
    visibility = azurerm_redhat_openshift_cluster.test.ingress_profile.0.visibility
  }

  worker_profile {
    vm_size      = azurerm_redhat_openshift_cluster.test.worker_profile.0.vm_size
    disk_size_gb = azurerm_redhat_openshift_cluster.test.worker_profile.0.disk_size_gb
    node_count   = azurerm_redhat_openshift_cluster.test.worker_profile.0.node_count
    subnet_id    = azurerm_redhat_openshift_cluster.test.worker_profile.0.subnet_id
  }

  service_principal {
    client_id     = azurerm_redhat_openshift_cluster.test.service_principal.0.client_id
    client_secret = azurerm_redhat_openshift_cluster.test.service_principal.0.client_secret
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.basic(data))
}

func (r OpenShiftClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test2" {
  display_name = "acctest-aro-2-%[2]d"
}

resource "azuread_service_principal" "test2" {
  application_id = azuread_application.test2.application_id
}

resource "azuread_service_principal_password" "test2" {
  service_principal_id = azuread_service_principal.test2.object_id
}

resource "azurerm_role_assignment" "role_network3" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.test2.object_id
}

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
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
    client_id     = azuread_application.test2.application_id
    client_secret = azuread_service_principal_password.test2.value
  }

  tags = {
    foo = "bar"
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
    "azurerm_role_assignment.role_network3",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) pullSecret(data acceptance.TestData, pullSecret string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain      = "aro-%[3]s.com"
    version     = "4.14.16"
    pull_secret = <<SECRET
%[4]s
SECRET
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
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString, pullSecret)
}

func (r OpenShiftClusterResource) userDefinedRouting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
  }

  network_profile {
    pod_cidr      = "10.128.0.0/14"
    service_cidr  = "172.30.0.0/16"
    outbound_type = "UserDefinedRouting"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    node_count   = 3
    subnet_id    = azurerm_subnet.worker_subnet.id
  }

  api_server_profile {
    visibility = "Private"
  }

  ingress_profile {
    visibility = "Private"
  }

  service_principal {
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) private(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
  }

  network_profile {
    pod_cidr     = "10.128.0.0/14"
    service_cidr = "172.30.0.0/16"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    node_count   = 3
    subnet_id    = azurerm_subnet.worker_subnet.id
  }

  api_server_profile {
    visibility = "Private"
  }

  ingress_profile {
    visibility = "Private"
  }

  service_principal {
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) basicWithFipsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain       = "aro-%[3]s.com"
    version      = "4.14.16"
    fips_enabled = true
  }

  network_profile {
    pod_cidr     = "10.128.0.0/14"
    service_cidr = "172.30.0.0/16"
  }

  api_server_profile {
    visibility = "Public"
  }

  ingress_profile {
    visibility = "Public"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    node_count   = 3
    subnet_id    = azurerm_subnet.worker_subnet.id
  }

  service_principal {
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) preconfiguredNetworkSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_network_security_group" "test" {
  name                = "test-network-security-group"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_all_inbound" {
  name                        = "test_allow_all_inbound"
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_all_outbound" {
  name                        = "test_allow_all_outbound"
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test_main" {
  subnet_id                 = azurerm_subnet.main_subnet.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_subnet_network_security_group_association" "test_worker" {
  subnet_id                 = azurerm_subnet.worker_subnet.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_role_assignment" "role_network3" {
  scope                = azurerm_network_security_group.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "role_network4" {
  scope                = azurerm_network_security_group.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.redhatopenshift.object_id
}

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
  }

  network_profile {
    pod_cidr                                     = "10.128.0.0/14"
    service_cidr                                 = "172.30.0.0/16"
    preconfigured_network_security_group_enabled = true
  }

  api_server_profile {
    visibility = "Public"
  }

  ingress_profile {
    visibility = "Public"
  }

  main_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.main_subnet.id
  }

  worker_profile {
    vm_size      = "Standard_D4s_v3"
    disk_size_gb = 128
    node_count   = 3
    subnet_id    = azurerm_subnet.worker_subnet.id
  }

  service_principal {
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
    "azurerm_role_assignment.role_network3",
    "azurerm_role_assignment.role_network4",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) encryptionAtHost(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                        = "acctestKV-%[3]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.test.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.test.tenant_id
  object_id    = data.azurerm_client_config.test.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "GetRotationPolicy",
    "Purge",
    "Update",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.service-principal
  ]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk_encryption" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id    = azurerm_disk_encryption_set.test.identity.0.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey"
  ]
}

resource "azurerm_role_assignment" "disk_encryption_reader1" {
  scope                = azurerm_disk_encryption_set.test.id
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "disk_encryption_reader2" {
  scope                = azurerm_disk_encryption_set.test.id
  role_definition_name = "Reader"
  principal_id         = data.azuread_service_principal.redhatopenshift.object_id
}

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain  = "aro-%[3]s.com"
    version = "4.14.16"
  }

  network_profile {
    pod_cidr     = "10.128.0.0/14"
    service_cidr = "172.30.0.0/16"
  }

  api_server_profile {
    visibility = "Public"
  }

  ingress_profile {
    visibility = "Public"
  }

  main_profile {
    vm_size                    = "Standard_D8s_v3"
    subnet_id                  = azurerm_subnet.main_subnet.id
    encryption_at_host_enabled = true
    disk_encryption_set_id     = azurerm_disk_encryption_set.test.id
  }

  worker_profile {
    vm_size                    = "Standard_D4s_v3"
    disk_size_gb               = 128
    node_count                 = 3
    subnet_id                  = azurerm_subnet.worker_subnet.id
    encryption_at_host_enabled = true
    disk_encryption_set_id     = azurerm_disk_encryption_set.test.id
  }

  service_principal {
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_key_vault_access_policy.disk_encryption",
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
    "azurerm_role_assignment.disk_encryption_reader1",
    "azurerm_role_assignment.disk_encryption_reader2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (r OpenShiftClusterResource) basicResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_openshift_cluster" "test" {
  name                = "acctestaro%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cluster_profile {
    domain                      = "aro-%[3]s.com"
    version                     = "4.14.16"
    managed_resource_group_name = "acctestrg-aro-infra-%[3]s"
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
    client_id     = azuread_application.test.application_id
    client_secret = azuread_service_principal_password.test.value
  }

  depends_on = [
    "azurerm_role_assignment.role_network1",
    "azurerm_role_assignment.role_network2",
  ]
}
  `, r.template(data), data.RandomInteger, data.RandomString)
}

func (OpenShiftClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  skip_provider_registration = true
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

data "azurerm_client_config" "test" {}

data "azuread_client_config" "test" {}

resource "azuread_application" "test" {
  display_name = "acctest-aro-%[1]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.object_id
}

data "azuread_service_principal" "redhatopenshift" {
  // This is the Azure Red Hat OpenShift RP service principal id, do NOT delete it
  application_id = "f1dd0a37-89c6-4e07-bcd1-ffd3d43d8875"
}

resource "azurerm_role_assignment" "role_network1" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "role_network2" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.redhatopenshift.object_id
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aro-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "main_subnet" {
  name                 = "main-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/23"]
  service_endpoints    = ["Microsoft.Storage", "Microsoft.ContainerRegistry"]
}
 `, data.RandomInteger, data.Locations.Primary)
}
