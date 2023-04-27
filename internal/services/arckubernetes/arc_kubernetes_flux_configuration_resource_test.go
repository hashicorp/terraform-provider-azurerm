package arckubernetes_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ArcKubernetesFluxConfigurationResource struct{}

func TestAccArcKubernetesFluxConfiguration_basic(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesFluxConfiguration_requiresImport(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, credential, privateKey, publicKey),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccArcKubernetesFluxConfiguration_complete(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateGitRepositoryWithHttpKey(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key"),
	})
}

func TestAccArcKubernetesFluxConfiguration_update(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bucket(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("bucket.0.secret_key"),
		{
			Config: r.privateGitRepositoryWithHttpKey(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key"),
	})
}

func TestAccArcKubernetesFluxConfiguration_privateRepositoryWithSshKey(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	const FluxUrl = "ARM_K8S_FLUX_CONFIG_SSH_URL" // git@github.com:Azure/arc-k8s-demo.git
	const PrivateSshKey = "ARM_K8S_FLUX_CONFIG_SSH_KEY"
	const KnownHosts = "ARM_K8S_FLUX_CONFIG_KNOWN_HOSTS"

	if os.Getenv(FluxUrl) == "" || os.Getenv(PrivateSshKey) == "" {
		t.Skipf("Acceptance test skipped unless env `%s` and `%s` set", FluxUrl, PrivateSshKey)
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateRepositoryWithSshKey(data, os.Getenv(FluxUrl), os.Getenv(PrivateSshKey), os.Getenv(KnownHosts), credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.ssh_private_key"),
	})
}

func TestAccArcKubernetesFluxConfiguration_azureBlobWithAccountKey(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithAccountKey(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azure_blob.0.account_key"),
	})
}

func TestAccArcKubernetesFluxConfiguration_azureBlobWithSasToken(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithSasToken(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azure_blob.0.sas_token"),
	})
}

func TestAccArcKubernetesFluxConfiguration_azureBlobWithServicePrincipalSecret(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithServicePrincipalSecret(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azure_blob.0.service_principal.0.client_secret"),
	})
}

func TestAccArcKubernetesFluxConfiguration_azureBlobWithServicePrincipalCertificate(t *testing.T) {
	// generateKey() is a time-consuming operation, we only run this test if an env var is set.
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	if os.Getenv("ARM_CLIENT_CERTIFICATE") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE not set")
	}
	if os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE_PASSWORD not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_flux_configuration", "test")
	r := ArcKubernetesFluxConfigurationResource{}
	credential := fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
	privateKey, publicKey, err := ArcKubernetesClusterResource{}.generateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %+v", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureBlobWithServicePrincipalCertificate(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azure_blob.0.service_principal.0.client_certificate", "azure_blob.0.service_principal.0.client_certificate_password"),
	})
}

func (r ArcKubernetesFluxConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluxconfiguration.ParseFluxConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ArcKubernetes.FluxConfigurationClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ArcKubernetesFluxConfigurationResource) template(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_arc_kubernetes_cluster_extension" "test" {
  name           = "acctest-kce-%[2]d"
  cluster_id     = azurerm_arc_kubernetes_cluster.test.id
  extension_type = "microsoft.flux"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}
`, ArcKubernetesClusterExtensionResource{}.template(data, credential, privateKey, publicKey), data.RandomInteger, data.RandomInteger)
}

func (r ArcKubernetesFluxConfigurationResource) basic(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
				%s

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  git_repository {
    url             = "https://github.com/Azure/arc-k8s-demo"
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesFluxConfigurationResource) requiresImport(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	config := r.basic(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_flux_configuration" "import" {
  name       = azurerm_arc_kubernetes_flux_configuration.test.name
  cluster_id = azurerm_arc_kubernetes_flux_configuration.test.cluster_id

  git_repository {
    url             = "https://github.com/Azure/arc-k8s-demo"
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, config)
}

func (r ArcKubernetesFluxConfigurationResource) privateGitRepositoryWithHttpKey(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
				%s

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id
  namespace  = "example"
  scope      = "cluster"

  git_repository {
    url                      = "https://github.com/Azure/arc-k8s-demo"
    https_user               = "example"
    https_key                = base64encode("example")
    https_ca_cert            = base64encode("example")
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
    re_creating_enabled        = true
    garbage_collection_enabled = true
  }

  kustomizations {
    name       = "kustomization-2"
    depends_on = ["kustomization-1"]
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesFluxConfigurationResource) bucket(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential, privateKey, publicKey)
	return fmt.Sprintf(`
			%s

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id
  namespace  = "example"
  scope      = "cluster"

  bucket {
    access_key               = "example"
    secret_key               = base64encode("example")
    bucket_name              = "flux"
    sync_interval_in_seconds = 800
    timeout_in_seconds       = 800
    url                      = "https://fluxminiotest.az.minio.io"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesFluxConfigurationResource) privateRepositoryWithSshKey(data acceptance.TestData, url string, sshKey string, knownHosts string, credential string, privateKey string, publicKey string) string {
	template := r.template(data, credential, privateKey, publicKey)
	knownHostsContent := ""
	if knownHosts != "" {
		knownHostsContent = fmt.Sprintf(`ssh_known_hosts = "%s"`, knownHosts)
	}
	return fmt.Sprintf(`
				%s

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  git_repository {
    url             = "%s"
    ssh_private_key = "%s"
    %s
    reference_type  = "branch"
    reference_value = "main"
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, template, data.RandomInteger, url, sshKey, knownHostsContent)
}

func (r ArcKubernetesFluxConfigurationResource) azureBlobWithAccountKey(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
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

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  azure_blob {
    container_name           = azurerm_storage_container.test.name
    url                      = azurerm_storage_account.test.primary_blob_endpoint
    account_key              = azurerm_storage_account.test.primary_access_key
    sync_interval_in_seconds = 800
    timeout_in_seconds       = 800
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger)
}

func (r ArcKubernetesFluxConfigurationResource) azureBlobWithSasToken(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
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

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  azure_blob {
    container_name = azurerm_storage_container.test.name
    url            = azurerm_storage_account.test.primary_blob_endpoint
    sas_token      = data.azurerm_storage_account_sas.test.sas
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger, startDate, endDate)
}

func (r ArcKubernetesFluxConfigurationResource) azureBlobWithServicePrincipalSecret(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
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

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  azure_blob {
    container_name = azurerm_storage_container.test.name
    url            = azurerm_storage_account.test.primary_blob_endpoint
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
    azurerm_arc_kubernetes_cluster_extension.test,
    azurerm_role_assignment.test_queue,
    azurerm_role_assignment.test_blob
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_SECRET"))
}

func (r ArcKubernetesFluxConfigurationResource) azureBlobWithServicePrincipalCertificate(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
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

resource "azurerm_arc_kubernetes_flux_configuration" "test" {
  name       = "acctest-fc-%[2]d"
  cluster_id = azurerm_arc_kubernetes_cluster.test.id

  azure_blob {
    container_name = azurerm_storage_container.test.name
    url            = azurerm_storage_account.test.primary_blob_endpoint
    service_principal {
      client_id                     = "%[3]s"
      tenant_id                     = "%[4]s"
      client_certificate            = "%[5]s"
      client_certificate_password   = "%[6]s"
      client_certificate_send_chain = true
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.test,
    azurerm_role_assignment.test_queue,
    azurerm_role_assignment.test_blob
  ]
}
`, r.template(data, credential, privateKey, publicKey), data.RandomInteger, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_CERTIFICATE"), os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"))
}
