package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMServiceFabricCluster_basic(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_basic(ri, location, 3)

	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "client_certificate_thumbprint.#", "0"),
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

func TestAccAzureRMServiceFabricCluster_addOnFeatures(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_addOnFeatures(ri, location)

	resource.Test(t, resource.TestCase{
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
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_certificates(ri, location),
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

func TestAccAzureRMServiceFabricCluster_clientCertificateThumbprint(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(ri, location),
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

func TestAccAzureRMServiceFabricCluster_diagnosticsConfig(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_diagnosticsConfig(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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

func TestAccAzureRMServiceFabricCluster_fabricSettings(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_fabricSettings(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_fabricSettings(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fabric_settings.#", "1"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, location, 3),
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
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_nodeTypeCustomPorts(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_nodeTypeMultiple(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, location, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "3"),
				),
			},
			{
				Config: testAccAzureRMServiceFabricCluster_basic(ri, location, 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceFabricClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "node_type.0.instance_count", "4"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceFabricCluster_tags(t *testing.T) {
	resourceName := "azurerm_service_fabric_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMServiceFabricCluster_tags(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceFabricClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	client := testAccProvider.Meta().(*ArmClient).serviceFabricClustersClient
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

func testCheckAzureRMServiceFabricClusterExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		clusterName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Service Fabric Cluster %q", clusterName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceFabricClustersClient
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
  name = "acctestRG-%d"
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

func testAccAzureRMServiceFabricCluster_addOnFeatures(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
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
  add_on_features     = [ "DnsService", "RepairManager" ]

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
  name = "acctestRG-%d"
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

    parameters {
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

func testAccAzureRMServiceFabricCluster_clientCertificateThumbprint(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
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

    parameters {
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

func testAccAzureRMServiceFabricCluster_diagnosticsConfig(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
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

func testAccAzureRMServiceFabricCluster_fabricSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
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

    parameters {
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
  name = "acctestRG-%d"
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
  name = "acctestRG-%d"
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

func testAccAzureRMServiceFabricCluster_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
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

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt)
}
