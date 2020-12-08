package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabric/parse"
)

func TestAccAzureRMServiceFabricCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(data, 3, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.1.is_primary", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceFabricCluster_requiresImport),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_manualClusterCodeVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	codeVersion := "6.5.676.9590"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(data, codeVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(data, codeVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttr(data.ResourceName, "cluster_code_version", codeVersion),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_manualLatest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "cluster_code_version"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_addOnFeatures(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_addOnFeatures(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_certificates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_reverseProxyCertificates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.reverse_proxy_endpoint_port", "19081"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyNotSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.reverse_proxy_endpoint_port", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")
	configBasic := testAccAzureRMServiceFabricCluster_basic(data, 3)
	configProxy := testAccAzureRMServiceFabricCluster_reverseProxyCertificates(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: configProxy,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.reverse_proxy_endpoint_port", "19081"),
				),
			},
			{
				Config: configBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(data.ResourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.reverse_proxy_endpoint_port", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.is_admin", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_withMultipleClientCertificateThumbprints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_withMultipleClientCertificateThumbprints(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_clientCertificateCommonNames(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.0.common_name", "firstcertcommonname"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.0.is_admin", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.0.issuer_thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.1.common_name", "secondcertcommonname"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.1.is_admin", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.1.issuer_thumbprint", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.is_admin", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.0.is_admin", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.1.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_thumbprint.1.is_admin", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "client_certificate_common_name.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_certificateCommonNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_certificateCommonNames(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_common_names.0.common_names.2962847220.certificate_common_name", "example"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_common_names.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:80"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectory(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.cluster_application_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.client_application_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:19080"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectory(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.cluster_application_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "azure_active_directory.0.client_application_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:19080"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.thumbprint", "3341DB6CF2AF72C611DF3BE3721A653AF1D43ECD50F584F828793DBE9103C3EE"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(data.ResourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(data.ResourceName, "management_endpoint", "https://example:19080"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.protected_account_key_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.blob_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.queue_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.table_endpoint"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.protected_account_key_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.blob_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.queue_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "diagnostics_config.0.table_endpoint"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "diagnostics_config.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_fabricSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "None"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettingsRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_fabricSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.#", "1"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "fabric_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.application_ports.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.application_ports.0.start_port", "20000"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.application_ports.0.end_port", "29999"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.ephemeral_ports.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.ephemeral_ports.0.start_port", "30000"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.ephemeral_ports.0.end_port", "39999"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.name", "first"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.1.name", "second"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.1.instance_count", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.1.is_primary", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "3"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(data, 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.instance_count", "4"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.placement_properties.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.placement_properties.HasSSD", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.capacities.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.capacities.ClientConnections", "20000"),
					resource.TestCheckResourceAttr(data.ResourceName, "node_type.0.capacities.MemoryGB", "8"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceFabricCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMServiceFabricClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabric.ClustersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_fabric_cluster" {
			continue
		}

		id, err := parse.ClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Service Fabric Cluster still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceFabricClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceFabric.ClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceFabricClustersClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Fabric Cluster %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceFabricCluster_basic(data acceptance.TestData, count int) string {
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

func testAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(data acceptance.TestData, count int, secondary_count int) string {
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

func testAccAzureRMServiceFabricCluster_requiresImport(data acceptance.TestData) string {
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
`, testAccAzureRMServiceFabricCluster_basic(data, 3))
}

func testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(data acceptance.TestData, clusterCodeVersion string) string {
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

func testAccAzureRMServiceFabricCluster_addOnFeatures(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_certificates(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_reverseProxyCertificates(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_withMultipleClientCertificateThumbprints(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_clientCertificateCommonNames(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_certificateCommonNames(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_azureActiveDirectory(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_diagnosticsConfig(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_fabricSettings(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_nodeTypeMultiple(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_nodeTypeProperties(data acceptance.TestData) string {
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

func testAccAzureRMServiceFabricCluster_tags(data acceptance.TestData) string {
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
