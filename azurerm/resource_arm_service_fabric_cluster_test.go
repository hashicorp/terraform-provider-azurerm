package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMServiceFabricCluster_basic(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(ri, testLocation(), 3, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "node_type.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config:      testAccAzureRMServiceFabricCluster_requiresImport(ri, testLocation(), 3),
				ExpectError: testRequiresImportError("azurerm_service_fabric_cluster"),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_manualClusterCodeVersion(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	codeVersion := "6.4.637.9590"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(ri, testLocation(), codeVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttr(resourceName, "cluster_code_version", codeVersion),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(ri, testLocation(), codeVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttr(resourceName, "cluster_code_version", codeVersion),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_manualLatest(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(ri, testLocation(), ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "upgrade_mode", "Manual"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_code_version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_addOnFeatures(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceFabricCluster_addOnFeatures(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_certificate(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_certificates(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyCertificate(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_reverseProxyCertificates(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.reverse_proxy_endpoint_port", "19081"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyNotSet(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_basic(ri, location, 3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.reverse_proxy_endpoint_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_reverseProxyUpdate(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	configBasic := testAccAzureRMServiceFabricCluster_basic(ri, location, 3)
	configProxy := testAccAzureRMServiceFabricCluster_reverseProxyCertificates(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: configProxy,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.reverse_proxy_endpoint_port", "19081"),
				),
			},
			{
				Config: configBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "http://example:80"),
					resource.TestCheckResourceAttr(resourceName, "add_on_features.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reverse_proxy_certificate.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.reverse_proxy_endpoint_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_clientCertificateThumbprint(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.0.is_admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.0.is_admin", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.1.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.1.is_admin", "false"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_certificateCommonNames(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_certificateCommonNames(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate_common_names.0.common_names.2962847220.certificate_common_name", "example"),
					resource.TestCheckResourceAttr(resourceName, "certificate_common_names.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectory(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectory(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.cluster_application_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.client_application_id"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:19080"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectory(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.cluster_application_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_active_directory.0.client_application_id"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:19080"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.thumbprint", "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"),
					resource.TestCheckResourceAttr(resourceName, "certificate.0.x509_store_name", "My"),
					resource.TestCheckResourceAttr(resourceName, "azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "EncryptAndSign"),
					resource.TestCheckResourceAttr(resourceName, "management_endpoint", "https://example:19080"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfig(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.protected_account_key_name"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.blob_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.queue_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.table_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfig(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.storage_account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.protected_account_key_name"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.blob_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.queue_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnostics_config.0.table_endpoint"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "diagnostics_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettings(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_fabricSettings(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.name", "Security"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.0.parameters.ClusterProtectionLevel", "None"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_fabricSettingsRemove(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_fabricSettings(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.#", "1"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.application_ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.application_ports.0.start_port", "20000"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.application_ports.0.end_port", "29999"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.ephemeral_ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.ephemeral_ports.0.start_port", "30000"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.ephemeral_ports.0.end_port", "39999"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesMultiple(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeMultiple(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.name", "first"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.is_primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "node_type.1.name", "second"),
					resource.TestCheckResourceAttr(resourceName, "node_type.1.instance_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "node_type.1.is_primary", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypesUpdate(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, testLocation(), 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "4"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_nodeTypeProperties(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_nodeTypeProperties(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.placement_properties.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.placement_properties.HasSSD", "true"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.capacities.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.capacities.ClientConnections", "20000"),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.capacities.MemoryGB", "8"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_tags(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_tags(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMServiceFabricClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceFabric.ClustersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_fabric_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

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
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		clusterName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Service Fabric Cluster %q", clusterName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceFabric.ClustersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, clusterName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceFabricClustersClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Fabric Cluster %q (Resource Group: %q) does not exist", clusterName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceFabricCluster_basic(rInt int, location string, count int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt, count)
}

func testAccAzureRMServiceFabricCluster_basicNodeTypeUpdate(rInt int, location string, count int, secondary_count int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt, count, secondary_count)
}

func testAccAzureRMServiceFabricCluster_requiresImport(rInt int, location string, count int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_service_fabric_cluster" "import" {
  name                = "${azurerm_service_fabric_cluster.test.name}"
  resource_group_name = "${azurerm_service_fabric_cluster.test.resource_group_name}"
  location            = "${azurerm_service_fabric_cluster.test.location}"
  reliability_level   = "${azurerm_service_fabric_cluster.test.reliability_level}"
  upgrade_mode        = "${azurerm_service_fabric_cluster.test.upgrade_mode}"
  vm_image            = "${azurerm_service_fabric_cluster.test.vm_image}"
  management_endpoint = "${azurerm_service_fabric_cluster.test.management_endpoint}"

  node_type {
    name                 = "first"
    instance_count       = %d
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, testAccAzureRMServiceFabricCluster_basic(rInt, location, count), count)
}

func testAccAzureRMServiceFabricCluster_manualClusterCodeVersion(rInt int, location, clusterCodeVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                 = "acctest-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
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
`, rInt, location, clusterCodeVersion)
}

func testAccAzureRMServiceFabricCluster_addOnFeatures(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_certificates(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_reverseProxyCertificates(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
    x509_store_name = "My"
  }

  reverse_proxy_certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_readerAdminClientCertificateThumbprint(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:80"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
    is_admin   = true
  }

  client_certificate_thumbprint {
    thumbprint = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_certificateCommonNames(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_azureActiveDirectory(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

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
  application_id = "${azuread_application.cluster_explorer.application_id}"
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
    resource_app_id = "${azuread_application.cluster_explorer.application_id}"

    resource_access {
      id   = "${azuread_application.cluster_explorer.oauth2_permissions.0.id}"
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "cluster_console" {
  application_id = "${azuread_application.cluster_console.application_id}"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
    x509_store_name = "My"
  }

  azure_active_directory {
    tenant_id              = "${data.azurerm_client_config.current.tenant_id}"
    cluster_application_id = "${azuread_application.cluster_explorer.application_id}"
    client_application_id  = "${azuread_application.cluster_console.application_id}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_azureActiveDirectoryDelete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://example:19080"

  certificate {
    thumbprint      = "33:41:DB:6C:F2:AF:72:C6:11:DF:3B:E3:72:1A:65:3A:F1:D4:3E:CD:50:F5:84:F8:28:79:3D:BE:91:03:C3:EE"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_diagnosticsConfig(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "http://example:80"

  diagnostics_config {
    storage_account_name       = "${azurerm_storage_account.test.name}"
    protected_account_key_name = "StorageAccountKey1"
    blob_endpoint              = "${azurerm_storage_account.test.primary_blob_endpoint}"
    queue_endpoint             = "${azurerm_storage_account.test.primary_queue_endpoint}"
    table_endpoint             = "${azurerm_storage_account.test.primary_table_endpoint}"
  }

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
`, rInt, location, rString, rInt)
}

func testAccAzureRMServiceFabricCluster_diagnosticsConfigDelete(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rString, rInt)
}

func testAccAzureRMServiceFabricCluster_fabricSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_nodeTypeMultiple(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_nodeTypeProperties(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}

func testAccAzureRMServiceFabricCluster_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt)
}
