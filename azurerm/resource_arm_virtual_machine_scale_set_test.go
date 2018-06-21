package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMVirtualMachineScaleSet_basic(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					// testing default scaleset values
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(resourceName, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicPublicIP(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicPublicIP(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetIsPrimary(resourceName, true),
					testCheckAzureRMVirtualMachineScaleSetPublicIPName(resourceName, "TestPublicIPConfiguration"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetAcceleratedNetworking(resourceName, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicIPForwarding(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	networkProfileName := fmt.Sprintf("TestNetworkProfile-%d", ri)
	networkProfile := map[string]interface{}{"name": networkProfileName, "primary": true}
	networkProfileHash := fmt.Sprintf("%d", resourceArmVirtualMachineScaleSetNetworkConfigurationHash(networkProfile))
	config := testAccAzureRMVirtualMachineScaleSet_basicIPForwarding(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_profile."+networkProfileHash+".ip_forwarding", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicDNSSettings(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	networkProfileName := fmt.Sprintf("TestNetworkProfile-%d", ri)
	networkProfile := map[string]interface{}{"name": networkProfileName, "primary": true}
	networkProfileHash := fmt.Sprintf("%d", resourceArmVirtualMachineScaleSetNetworkConfigurationHash(networkProfile))
	config := testAccAzureRMVirtualMachineScaleSet_basicDNSSettings(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_profile."+networkProfileHash+".dns_settings.0.dns_servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "network_profile."+networkProfileHash+".dns_settings.0.dns_servers.1", "8.8.4.4"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_bootDiagnostic(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "boot_diagnostics.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicWindows(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),

					// single placement group should default to true
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(resourceName, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(resourceName, false),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_linuxUpdated(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSet_linux(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSet_linuxUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customDataUpdated(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSet_linux(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D1_v2")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk_resize(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D1_v2")
	postConfig := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D2_v2")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_disappears(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_planManagedDisk(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_planManagedDisk(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customImage(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := fmt.Sprintf("testadmin%d", ri)
	password := fmt.Sprintf("Password1234!%d", ri)
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"
	config := testAccAzureRMVirtualMachineScaleSet_customImage(ri, testLocation(), userName, password, hostName)
	preConfig := testAccAzureRMImage_standaloneImage_setup(ri, userName, password, hostName, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				//need to create a vm and then reference it in the image creation
				Config:  preConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, testLocation()),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_applicationGateway(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetApplicationGatewayTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetHasApplicationGateway(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancer(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetHasLoadbalancer(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancerManagedDataDisks(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetHasDataDisks(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_overprovision(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetOverprovision(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_priority(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetPriorityTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "priority", "Low"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_MSI(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetMSITemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttrSet(resourceName, "identity.0.principal_id"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extension(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetExtensionTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extensionUpdate(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSetExtensionTemplate(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleExtensions(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Conflict between `vhd_containers`"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_NonStandardCasing(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetNonStandardCasing(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func testGetAzureRMVirtualMachineScaleSet(s *terraform.State, resourceName string) (result *compute.VirtualMachineScaleSet, err error) {
	// Ensure we have enough information in state to look up in API
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", resourceName)
	}

	// Name of the actual scale set
	name := rs.Primary.Attributes["name"]

	resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
	if !hasResourceGroup {
		return nil, fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
	}

	client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	vmss, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
	}

	if vmss.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Bad: VirtualMachineScaleSet %q (resource group: %q) does not exist", name, resourceGroup)
	}

	return &vmss, err
}

func testCheckAzureRMVirtualMachineScaleSetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		return err
	}
}

func testCheckAzureRMVirtualMachineScaleSetDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_scale_set" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Machine Scale Set still exists:\n%#v", resp.VirtualMachineScaleSetProperties)
		}
	}

	return nil
}

func testCheckAzureRMVirtualMachineScaleSetHasLoadbalancer(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		pools := (*ip)[0].LoadBalancerBackendAddressPools
		if pools == nil || len(*pools) == 0 {
			return fmt.Errorf("Bad: Load balancer backend pools is empty for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetHasApplicationGateway(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		pools := (*ip)[0].ApplicationGatewayBackendAddressPools
		if pools == nil || len(*pools) == 0 {
			return fmt.Errorf("Bad: Application gateway backend pools is empty for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetIsPrimary(name string, boolean bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		primary := *(*ip)[0].Primary
		if primary != boolean {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, boolean, primary)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetPublicIPName(name, publicIPName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		publicIPConfigName := *(*ip)[0].PublicIPAddressConfiguration.Name
		if publicIPConfigName != publicIPName {
			return fmt.Errorf("Bad: Public IP Config Name set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, publicIPName, publicIPConfigName)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetAcceleratedNetworking(name string, boolean bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		acceleratedNetworking := *(*n)[0].EnableAcceleratedNetworking
		if acceleratedNetworking != boolean {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, boolean, acceleratedNetworking)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetOverprovision(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		if *resp.Overprovision {
			return fmt.Errorf("Bad: Overprovision should have been false for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(name string, expectedSinglePlacementGroup bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		if *resp.SinglePlacementGroup != expectedSinglePlacementGroup {
			return fmt.Errorf("Bad: SinglePlacementGroup should have been %t for scale set %v", expectedSinglePlacementGroup, name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetMSI(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		identityType := resp.Identity.Type
		if identityType != "systemAssigned" {
			return fmt.Errorf("Bad: Identity Type is not systemAssigned for scale set %v", name)
		}

		principalID := *resp.Identity.PrincipalID
		if len(principalID) == 0 {
			return fmt.Errorf("Bad: Could not get principal_id for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetExtension(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.ExtensionProfile.Extensions
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get extensions for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetHasDataDisks(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachineScaleSet %q (resource group: %q) does not exist", name, resourceGroup)
		}

		storageProfile := resp.VirtualMachineProfile.StorageProfile.DataDisks
		if storageProfile == nil || len(*storageProfile) == 0 {
			return fmt.Errorf("Bad: Could not get data disks configurations for scale set %v", name)
		}

		return nil
	}
}

func testAccAzureRMVirtualMachineScaleSet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicPublicIP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
				primary = true
				public_ip_address_configuration {
					name = "TestPublicIPConfiguration"
					domain_name_label = "test-domain-label-%[1]d"
					idle_timeout = 4
				}
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D4_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
			accelerated_networking = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicIPForwarding(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D4_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
	  ip_forwarding = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicDNSSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D4_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
	  dns_settings {
		  dns_servers = ["8.8.8.8", "8.8.4.4"]
	  }
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

	boot_diagnostics {
			storage_uri = "${azurerm_storage_account.test.primary_blob_endpoint}"
	}

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
			network_security_group_id = "${azurerm_network_security_group.test.id}"
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
				primary = true
				public_ip_address_configuration {
					name = "TestPublicIPConfiguration"
					domain_name_label = "test-domain-label-%[1]d"
					idle_timeout = 4
				}
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%[1]d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  os_profile_windows_config {
    enable_automatic_upgrades = false
    provision_vm_agent = true

    winrm {
      protocol = "http"
    }
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  single_placement_group = false

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_linux(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "custom data!"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_linuxUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "custom data!"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  tags {
    ThisIs = "a test"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "updated custom data!"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
    name = "TestNetworkProfile-%[1]d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name 		  = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk_withZones(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
    name = "acctvmss-%[1]d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    upgrade_policy_mode = "Manual"
    zones = ["1", "2"]

    sku {
        name = "Standard_D1_v2"
        tier = "Standard"
        capacity = 2
    }

    os_profile {
        computer_name_prefix = "testvm-%[1]d"
        admin_username = "myadmin"
        admin_password = "Passwword1234"
    }

    network_profile {
        name = "TestNetworkProfile-%[1]d"
        primary = true
        ip_configuration {
            name = "TestIPConfiguration"
            subnet_id = "${azurerm_subnet.test.id}"
        }
    }

    storage_profile_os_disk {
        caching = "ReadWrite"
        create_option = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    storage_profile_image_reference {
        publisher = "Canonical"
        offer = "UbuntuServer"
        sku = "16.04-LTS"
        version = "latest"
    }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(rInt int, location string, vmSize string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "%[3]s"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  os_profile_windows_config {
	enable_automatic_upgrades = false
	provision_vm_agent = true

	additional_unattend_config {
	  pass = "oobeSystem"
	  component = "Microsoft-Windows-Shell-Setup"
	  setting_name = "AutoLogon"
	  content = "<AutoLogon><Username>myadmin</Username><Password><Value>Passwword1234</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount></AutoLogon>"
	}

	additional_unattend_config {
      pass = "oobeSystem"
      component = "Microsoft-Windows-Shell-Setup"
      setting_name = "FirstLogonCommands"
      content = "<FirstLogonCommands><SynchronousCommand><CommandLine>shutdown /r /t 0 /c \"initial reboot\"</CommandLine><Description>reboot</Description><Order>1</Order></SynchronousCommand></FirstLogonCommands>"
    }
  }

  network_profile {
    name = "TestNetworkProfile-%[1]d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name 		  = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, rInt, location, vmSize)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
    name = "TestNetworkProfile-%[1]d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetApplicationGatewayTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                         = "TestIPConfiguration"
      subnet_id                                    = "${azurerm_subnet.test.id}"
      application_gateway_backend_address_pool_ids = ["${azurerm_application_gateway.test.backend_address_pool.0.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

# application gateway
resource "azurerm_subnet" "gwtest" {
  name                 = "gw-subnet-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.3.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%[1]d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.gwtest.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.gwtest.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  name                           = "ssh"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "default"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
      load_balancer_inbound_nat_rules_ids    = ["${azurerm_lb_nat_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetPriorityTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false
  priority            = "Low"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetMSITemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  identity {
    type     = "systemAssigned"
  }

  extension {
    name                       = "MSILinuxExtension"
    publisher                  = "Microsoft.ManagedIdentity"
    type                       = "ManagedIdentityExtensionForLinux"
    type_handler_version       = "1.0"
    settings                   = "{\"port\": 50342}"
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME"
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME",
			"timestamp": 12345679955
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME"
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }

  extension {
    name                       = "Docker"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "DockerExtension"
    type_handler_version       = "1.0"
    auto_upgrade_minor_version = true
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
    vhd_containers    = ["should_cause_conflict"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_data_disk {
    lun               = 0
    caching           = "ReadWrite"
    create_option     = "Empty"
    disk_size_gb      = 10
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}

`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSetNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_planManagedDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  plan {
    name = "os"
    product = "rancheros"
    publisher = "rancher"
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "rancher"
    offer     = "rancheros"
    sku       = "os"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_customImage(rInt int, location string, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctpip-%[1]d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "%[3]s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.testsource.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%[4]s"
    admin_password = "%[5]s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "test" {
  name                = "accteste-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

	storage_profile_image_reference {
      id = "${azurerm_image.test.id}"
	}

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
}
`, rInt, location, hostName, userName, password)
}

func testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
    name = "accsa%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "primary-%[1]d"
      primary = true
      ip_configuration {
        name = "primary"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  network_profile {
      name = "secondary-%[1]d"
      primary = false
      ip_configuration {
        name = "secondary"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}
