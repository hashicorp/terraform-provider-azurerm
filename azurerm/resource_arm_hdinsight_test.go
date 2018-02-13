package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func Test_HDInsightCreation(t *testing.T) {

	resourceName := "teraform_hdinsight_test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))

	config := testHDInsightClusterCreation(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testCheckAzureRMHDInsigthClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExist(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~1"),
				),
			},
		},
	})
}

func testCheckAzureRMHDInsigthClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).hdiClusterClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hdinsight_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("HDI cluster already  exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMHDInsightClusterExist(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for HDI Cluster : '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).hdiClusterClient

		print(resourceGroup)
		print("name %s", name)
		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on HDI ClusterClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: HDI Cluster '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testHDInsightClusterCreation(rInt int, storage string, location string) string {
	template := fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "terraform-RG-HDInsight"
			location = "%[3]s"
		}

		resource "azurerm_virtual_network" "vnet"{
			name				= 	"vnet_tr"
			location 			=	"${azurerm_resource_group.test.location}"
			resource_group_name	= 	"${azurerm_resource_group.test.name}"
			address_space 		=   ["10.7.0.0/28"] 
		}	

		resource "azurerm_subnet" "subnet"{
			name				 = "subnet_tr"
			virtual_network_name = "${azurerm_virtual_network.vnet.name}"
			resource_group_name	 = "${azurerm_resource_group.test.name}"
			address_prefix		 = "10.7.0.0/29" 
		}

		resource "azurerm_hdinsight_cluster"  "cluster" {
			name 					= "hdinsight-cluster-%[1]d"
			location                = "${azurerm_resource_group.test.location}"
			resource_group_name     = "${azurerm_resource_group.test.name}"
			tier					= "Standard"
			
			kind					= "hadoop"
			cluster_version			= "3.6"
			
			gateway {
				restAuthCredential.isEnabled = "true"
				restAuthCredential.username = "admin"
				restAuthCredential.password = "P&ssw0rd1234"
			} 
			core_site {
				dfs.adls.home.hostname = "sncfpocfbdusl.azuredatalakestore.net"
				dfs.adls.home.mountpoint = "/clusters/HDI-POC-UL/"
			}
			cluster_identity {
				clusterIdentity.applicationId = "a4c29302-0480-4135-b59d-d86d7b6e590a"
				clusterIdentity.certificate = "identityCertificate"
				clusterIdentity.aadTenantId = "https://login.microsoftonline.com/67357eed-9531-498d-9f78-d895d88ed1e1"
				clusterIdentity.resourceUri = "https://management.core.windows.net/"
				clusterIdentity.certificatePassword  = "identityCertificatePassword"
			}

			tags {
				type = "hdisinght"
			}			

			roles {
				name 			= "headnode"
				count 			= "2"
				size 			= "Standard_D3_v2"
				os_profile {
					username	= "adminjf"
					password	= "P&ssw0rd1234"
				}  
				network_profile {
					virtual_network_id	= "${azurerm_virtual_network.vnet.id}"
					subnet_name			= "${azurerm_subnet.subnet.id}"
				}	
			}
			roles {
				name 			= "workernode"
				count 			= "2"
				size 			= "Standard_D3_v2"
				os_profile {
					username	= "adminjf"
					password	= "P&ssw0rd1234"
				}  
				network_profile {
					virtual_network_id	= "${azurerm_virtual_network.vnet.id}"
					subnet_name			= "${azurerm_subnet.subnet.id}"
				}	
			}
		}	
		`, rInt, storage, location)
	print(template)
	return template
}

func Test_HDInsightDelete(t *testing.T) {
	//resourceName := "teraform_hdinsight_test"
}

func Test_HDInsightRead(t *testing.T) {
	//resourceName := "teraform_hdinsight_test"
}
