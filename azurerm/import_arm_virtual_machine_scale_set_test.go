package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMVirtualMachineScaleSet_importBasic(t *testing.T) {
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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importBasic_managedDisk(t *testing.T) {
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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importBasic_managedDisk_withZones(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"

	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk_withZones(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importLinux(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"

	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_linux(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"os_profile.0.admin_password",
					"os_profile.0.custom_data",
				},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importLoadBalancer(t *testing.T) {
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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importOverProvision(t *testing.T) {
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
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importExtension(t *testing.T) {
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
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importMultipleExtensions(t *testing.T) {
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
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_profile.0.admin_password"},
			},
		},
	})
}
