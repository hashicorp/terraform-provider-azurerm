// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesFluxConfigurationResource struct{}

func TestAccKubernetesFluxConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
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

func TestAccKubernetesFluxConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
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

func TestAccKubernetesFluxConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateGitRepositoryWithHttpKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key_base64"),
	})
}

func TestAccKubernetesFluxConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bucket(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("bucket.0.secret_key_base64"),
		{
			Config: r.privateGitRepositoryWithHttpKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key_base64"),
	})
}

func TestAccKubernetesFluxConfiguration_privateRepositoryWithSshKey(t *testing.T) {
	const FluxUrl = "ARM_K8S_FLUX_CONFIG_SSH_URL" // git@github.com:Azure/arc-k8s-demo.git
	const PrivateSshKey = "ARM_K8S_FLUX_CONFIG_SSH_KEY"
	const KnownHosts = "ARM_K8S_FLUX_CONFIG_KNOWN_HOSTS"

	if os.Getenv(FluxUrl) == "" || os.Getenv(PrivateSshKey) == "" {
		t.Skipf("Acceptance test skipped unless env `%s` and `%s` set", FluxUrl, PrivateSshKey)
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateRepositoryWithSshKey(data, os.Getenv(FluxUrl), os.Getenv(PrivateSshKey), os.Getenv(KnownHosts)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.ssh_private_key_base64"),
	})
}

func TestAccKubernetesFluxConfiguration_azureBlobWithAccountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("blob_storage.0.account_key"),
	})
}

func TestAccKubernetesFluxConfiguration_azureBlobWithManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFluxConfiguration_azureBlobWithSasToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithSasToken(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("blob_storage.0.sas_token"),
	})
}

func TestAccKubernetesFluxConfiguration_azureBlobWithServicePrincipalSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithServicePrincipalSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("blob_storage.0.service_principal.0.client_secret"),
	})
}

func TestAccKubernetesFluxConfiguration_azureBlobWithServicePrincipalCertificate(t *testing.T) {
	if os.Getenv("ARM_CLIENT_CERTIFICATE") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE not set")
	}
	if os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE_PASSWORD not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithServicePrincipalCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("blob_storage.0.service_principal.0.client_certificate_base64", "blob_storage.0.service_principal.0.client_certificate_password"),
	})
}

func TestAccKubernetesFluxConfiguration_kustomizationNameDuplicated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_flux_configuration", "test")
	r := KubernetesFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.kustomizationNameDuplicated(data),
			ExpectError: regexp.MustCompile("kustomization name `kustomization-1` is not unique"),
		},
	})
}

func (r KubernetesFluxConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluxconfiguration.ParseScopedFluxConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Containers.KubernetesFluxConfigurationClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r KubernetesFluxConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_extension" "test" {
  name           = "acctest-kce-%d"
  cluster_id     = azurerm_kubernetes_cluster.test.id
  extension_type = "microsoft.flux"
}

`, KubernetesClusterExtensionResource{}.template(data), data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
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
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_flux_configuration" "import" {
  name       = azurerm_kubernetes_flux_configuration.test.name
  cluster_id = azurerm_kubernetes_flux_configuration.test.cluster_id
  namespace  = azurerm_kubernetes_flux_configuration.test.namespace

  git_repository {
    url             = "https://github.com/Azure/arc-k8s-demo"
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, config)
}

func (r KubernetesFluxConfigurationResource) privateGitRepositoryWithHttpKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"
  scope      = "cluster"

  git_repository {
    url                      = "https://github.com/Azure/arc-k8s-demo"
    https_user               = "example"
    https_key_base64         = base64encode("example")
    https_ca_cert_base64     = base64encode("example")
    sync_interval_in_seconds = 800
    timeout_in_seconds       = 800
    reference_type           = "branch"
    reference_value          = "main"
  }

  kustomizations {
    name                       = "kustomization-1"
    path                       = "./test/path"
    timeout_in_seconds         = 800
    sync_interval_in_seconds   = 800
    retry_interval_in_seconds  = 800
    recreating_enabled         = true
    garbage_collection_enabled = true
  }

  kustomizations {
    name       = "kustomization-2"
    depends_on = ["kustomization-1"]
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) bucket(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"
  scope      = "cluster"

  bucket {
    access_key               = "example"
    secret_key_base64        = base64encode("example")
    bucket_name              = "flux"
    sync_interval_in_seconds = 800
    timeout_in_seconds       = 800
    url                      = "https://fluxminiotest.az.minio.io"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) privateRepositoryWithSshKey(data acceptance.TestData, url string, sshKey string, knownHosts string) string {
	template := r.template(data)
	knownHostsContent := ""
	if knownHosts != "" {
		knownHostsContent = fmt.Sprintf(`ssh_known_hosts_base64 = "%s"`, knownHosts)
	}
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  git_repository {
    url                    = "%s"
    ssh_private_key_base64 = "%s"
    %s
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger, url, sshKey, knownHostsContent)
}

func (r KubernetesFluxConfigurationResource) azureBlobWithAccountKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "sa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "asc%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  blob_storage {
    container_id             = azurerm_storage_container.test.id
    account_key              = azurerm_storage_account.test.primary_access_key
    sync_interval_in_seconds = 800
    timeout_in_seconds       = 800
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) azureBlobWithManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "sa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "asc%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_role_assignment" "test_queue" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Queue Data Contributor"
  principal_id         = azurerm_kubernetes_cluster.test.kubelet_identity.0.object_id
}

resource "azurerm_role_assignment" "test_blob" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_kubernetes_cluster.test.kubelet_identity.0.object_id
}

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  blob_storage {
    container_id = azurerm_storage_container.test.id
    managed_identity {
      client_id = azurerm_kubernetes_cluster.test.kubelet_identity.0.client_id
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test,
    azurerm_role_assignment.test_queue,
    azurerm_role_assignment.test_blob
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFluxConfigurationResource) azureBlobWithSasToken(data acceptance.TestData) string {
	utcNow := time.Now().UTC()
	startDate := utcNow.Add(-time.Hour * 24).Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 48).Format(time.RFC3339)

	return fmt.Sprintf(`
				%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "sa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "asc%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  signed_version    = "2019-10-10"

  resource_types {
    service   = true
    container = true
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "%[3]s"
  expiry = "%[4]s"

  permissions {
    read    = true
    write   = true
    delete  = true
    list    = true
    add     = true
    create  = true
    update  = true
    process = true
    tag     = true
    filter  = false
  }
}

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  blob_storage {
    container_id = azurerm_storage_container.test.id
    sas_token    = data.azurerm_storage_account_sas.test.sas
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, r.template(data), data.RandomInteger, startDate, endDate)
}

func (r KubernetesFluxConfigurationResource) azureBlobWithServicePrincipalSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "sa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "asc%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test_queue" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Queue Data Contributor"
  principal_id         = data.azurerm_client_config.test.object_id
}

resource "azurerm_role_assignment" "test_blob" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.test.object_id
}

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  blob_storage {
    container_id = azurerm_storage_container.test.id
    service_principal {
      client_id     = "%[3]s"
      tenant_id     = "%[4]s"
      client_secret = "%[5]s"
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test,
    azurerm_role_assignment.test_queue,
    azurerm_role_assignment.test_blob
  ]
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_SECRET"))
}

func (r KubernetesFluxConfigurationResource) azureBlobWithServicePrincipalCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "sa%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "asc%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test_queue" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Queue Data Contributor"
  principal_id         = data.azurerm_client_config.test.object_id
}

resource "azurerm_role_assignment" "test_blob" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.test.object_id
}

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_kubernetes_cluster.test.id
  namespace  = "flux"

  blob_storage {
    container_id = azurerm_storage_container.test.id
    service_principal {
      client_id                     = "%[3]s"
      tenant_id                     = "%[4]s"
      client_certificate_base64     = "%[5]s"
      client_certificate_password   = "%[6]s"
      client_certificate_send_chain = true
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test,
    azurerm_role_assignment.test_queue,
    azurerm_role_assignment.test_blob
  ]
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_CERTIFICATE"), os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"))
}

func (r KubernetesFluxConfigurationResource) kustomizationNameDuplicated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
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

  kustomizations {
    name = "kustomization-1"
    path = "./test/path"
  }

  depends_on = [
    azurerm_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}
