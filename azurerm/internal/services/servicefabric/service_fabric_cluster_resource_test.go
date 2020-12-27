package servicefabric_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabric/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceFabricClusterResource struct{}

func TestAccAzureRMServiceFabricCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.basicNodeTypeUpdate(data, 3, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.1.is_primary").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMServiceFabricCluster_manualClusterCodeVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	codeVersion := "7.2.445.9590"
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.manualClusterCodeVersion(data, codeVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				resource.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
			),
		},
		{
			Config: r.manualClusterCodeVersion(data, codeVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				resource.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_manualLatest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.manualClusterCodeVersion(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_mode").HasValue("Manual"),
				check.That(data.ResourceName).Key("cluster_code_version").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_addOnFeatures(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.addOnFeatures(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.certificates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.reverseProxyCertificates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("19081"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyNotSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.reverseProxyCertificates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("19081"),
			),
		},
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("http://example:80"),
				check.That(data.ResourceName).Key("add_on_features.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("reverse_proxy_certificate.#").HasValue("0"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("0"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.reverse_proxy_endpoint_port").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.clientCertificateThumbprint(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("1"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_withMultipleClientCertificateThumbprints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withMultipleClientCertificateThumbprints(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.clientCertificateCommonNames(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("2"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.common_name").HasValue("firstcertcommonname"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_common_name.0.issuer_thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.common_name").HasValue("secondcertcommonname"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.is_admin").HasValue("false"),
				check.That(data.ResourceName).Key("client_certificate_common_name.1.issuer_thumbprint").IsEmpty(),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("1"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.readerAdminClientCertificateThumbprint(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.#").HasValue("2"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.0.is_admin").HasValue("true"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.1.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("client_certificate_thumbprint.1.is_admin").HasValue("false"),
				check.That(data.ResourceName).Key("client_certificate_common_name.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_certificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.certificateCommonNames(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate_common_names.0.common_names.2962847220.certificate_common_name").HasValue("example"),
				check.That(data.ResourceName).Key("certificate_common_names.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:80"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.azureActiveDirectory(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.cluster_application_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.client_application_id").Exists(),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.azureActiveDirectory(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.cluster_application_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory.0.client_application_id").Exists(),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		{
			Config: r.azureActiveDirectoryDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificate.#").HasValue("1"),
				check.That(data.ResourceName).Key("certificate.0.thumbprint").HasValue("3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
				check.That(data.ResourceName).Key("certificate.0.x509_store_name").HasValue("My"),
				check.That(data.ResourceName).Key("azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("EncryptAndSign"),
				check.That(data.ResourceName).Key("management_endpoint").HasValue("https://example:19080"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.diagnosticsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics_config.0.storage_account_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.protected_account_key_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.blob_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.queue_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.table_endpoint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.diagnosticsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("diagnostics_config.0.storage_account_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.protected_account_key_name").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.blob_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.queue_endpoint").Exists(),
				check.That(data.ResourceName).Key("diagnostics_config.0.table_endpoint").Exists(),
			),
		},
		{
			Config: r.diagnosticsConfigDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostics_config.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fabricSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("fabric_settings.0.name").HasValue("Security"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.%").HasValue("1"),
				check.That(data.ResourceName).Key("fabric_settings.0.parameters.ClusterProtectionLevel").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettingsRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fabricSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("1"),
			),
		},
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fabric_settings.#").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nodeTypeCustomPorts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.0.start_port").HasValue("20000"),
				check.That(data.ResourceName).Key("node_type.0.application_ports.0.end_port").HasValue("29999"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.0.start_port").HasValue("30000"),
				check.That(data.ResourceName).Key("node_type.0.ephemeral_ports.0.end_port").HasValue("39999"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nodeTypeMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.name").HasValue("first"),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("node_type.0.is_primary").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.1.name").HasValue("second"),
				check.That(data.ResourceName).Key("node_type.1.instance_count").HasValue("4"),
				check.That(data.ResourceName).Key("node_type.1.is_primary").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("3"),
			),
		},
		{
			Config: r.basic(data, 4),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.instance_count").HasValue("4"),
			),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nodeTypeProperties(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.0.placement_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.placement_properties.HasSSD").HasValue("true"),
				check.That(data.ResourceName).Key("node_type.0.capacities.%").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.capacities.ClientConnections").HasValue("20000"),
				check.That(data.ResourceName).Key("node_type.0.capacities.MemoryGB").HasValue("8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceFabricCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	r := ServiceFabricClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceFabricClusterResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ServiceFabric.ClustersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Service Fabric Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r ServiceFabricClusterResource) basic(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = %d
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, count)
}

func (r ServiceFabricClusterResource) basicNodeTypeUpdate(data acceptance.TestData, count int, secondary_count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = %d
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = "second"
    instance_count       = %d
    is_primary           = false
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, count, secondary_count)
}

func (r ServiceFabricClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_service_fabric_cluster" "import" {
  name                = azurerm_service_fabric_cluster.test.name
  resource_group_name = azurerm_service_fabric_cluster.test.resource_group_name
  location            = azurerm_service_fabric_cluster.test.location
  reliability_level   = azurerm_service_fabric_cluster.test.reliability_level
  upgrade_mode        = azurerm_service_fabric_cluster.test.upgrade_mode
  vm_image            = azurerm_service_fabric_cluster.test.vm_image
  management_endpoint = azurerm_service_fabric_cluster.test.management_endpoint

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, r.basic(data, 3))
}

func (r ServiceFabricClusterResource) manualClusterCodeVersion(data acceptance.TestData, clusterCodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                 = "acctest-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  reliability_level    = "Bronze"
  upgrade_mode         = "Manual"
  cluster_code_version = "%[3]s"
  vm_image             = "Windows"
  management_endpoint  = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, clusterCodeVersion)
}

func (r ServiceFabricClusterResource) addOnFeatures(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"
  add_on_features     = ["DnsService", "RepairManager"]

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) certificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) reverseProxyCertificates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  reverse_proxy_certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                        = "first"
    instance_count              = 3
    is_primary                  = true
    client_endpoint_port        = 2020
    http_endpoint_port          = 80
    reverse_proxy_endpoint_port = 19081
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clientCertificateThumbprint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) withMultipleClientCertificateThumbprints(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cluster-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-cluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "1341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  client_certificate_thumbprint {
    thumbprint = "2341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = false
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) clientCertificateCommonNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_common_name {
    common_name       = "firstcertcommonname"
    issuer_thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin          = true
  }

  client_certificate_common_name {
    common_name = "secondcertcommonname"
    is_admin    = false
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) readerAdminClientCertificateThumbprint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = true
  }

  client_certificate_thumbprint {
    thumbprint = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    is_admin   = false
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) certificateCommonNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate_common_names {
    common_names {
      certificate_common_name = "example"
    }

    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) azureActiveDirectory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

resource "azuread_application" "cluster_explorer" {
  name                       = "${azurerm_resource_group.test.name}-explorer-AAD"
  homepage                   = "https://example:19080/Explorer/index.html"
  identifier_uris            = ["https://example:19080/Explorer/index.html"]
  reply_urls                 = ["https://example:19080/Explorer/index.html"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true

  # https://blogs.msdn.microsoft.com/aaddevsup/2018/06/06/guid-table-for-windows-azure-active-directory-permissions/
  # https://shawntabrizi.com/aad/common-microsoft-resources-azure-active-directory/
  required_resource_access {
    resource_app_id = "00000002-0000-0000-c000-000000000000"

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "cluster_explorer" {
  application_id = azuread_application.cluster_explorer.application_id
}

resource "azuread_application" "cluster_console" {
  name                       = "${azurerm_resource_group.test.name}-console-AAD"
  type                       = "native"
  reply_urls                 = ["urn:ietf:wg:oauth:2.0:oob"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true

  # https://blogs.msdn.microsoft.com/aaddevsup/2018/06/06/guid-table-for-windows-azure-active-directory-permissions/
  # https://shawntabrizi.com/aad/common-microsoft-resources-azure-active-directory/
  required_resource_access {
    resource_app_id = "00000002-0000-0000-c000-000000000000"

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
      type = "Scope"
    }
  }

  required_resource_access {
    resource_app_id = azuread_application.cluster_explorer.application_id

    resource_access {
      id   = azuread_application.cluster_explorer.oauth2_permissions[0].id
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "cluster_console" {
  application_id = azuread_application.cluster_console.application_id
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  azure_active_directory {
    tenant_id              = data.azurerm_client_config.current.tenant_id
    cluster_application_id = azuread_application.cluster_explorer.application_id
    client_application_id  = azuread_application.cluster_console.application_id
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "system"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) azureActiveDirectoryDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"
    x509_store_name = "My"
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  node_type {
    name                 = "system"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) diagnosticsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  diagnostics_config {
    storage_account_name       = azurerm_storage_account.test.name
    protected_account_key_name = "StorageAccountKey1"
    blob_endpoint              = azurerm_storage_account.test.primary_blob_endpoint
    queue_endpoint             = azurerm_storage_account.test.primary_queue_endpoint
    table_endpoint             = azurerm_storage_account.test.primary_table_endpoint
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) diagnosticsConfigDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r ServiceFabricClusterResource) fabricSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "None"
    }
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeCustomPorts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80

    application_ports {
      start_port = 20000
      end_port   = 29999
    }

    ephemeral_ports {
      start_port = 30000
      end_port   = 39999
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = "second"
    instance_count       = 4
    is_primary           = false
    client_endpoint_port = 2121
    http_endpoint_port   = 81
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) nodeTypeProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name = "first"

    placement_properties = {
      "HasSSD" = "true"
    }

    capacities = {
      "ClientConnections" = "20000"
      "MemoryGB"          = "8"
    }

    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ServiceFabricClusterResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
