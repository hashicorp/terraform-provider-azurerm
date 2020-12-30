package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAzureRMVirtualMachineScaleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					// testing default scaleset values
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(data.ResourceName, true),
					resource.TestCheckResourceAttr(data.ResourceName, "eviction_policy", ""),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualMachineScaleSet_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_scale_set"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_evictionPolicyDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_evictionPolicyDelete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "eviction_policy", "Delete"),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_standardSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_standardSSD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_withPPG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_withPPG(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "proximity_placement_group_id"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicPublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicPublicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetIsPrimary(data.ResourceName, true),
					testCheckAzureRMVirtualMachineScaleSetPublicIPName(data.ResourceName, "TestPublicIPConfiguration"),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicPublicIP_simpleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetIsPrimary(data.ResourceName, true),
					testCheckAzureRMVirtualMachineScaleSetPublicIPName(data.ResourceName, "TestPublicIPConfiguration"),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP_updated_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetPublicIPName(data.ResourceName, "TestPublicIPConfiguration"),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_updateNetworkProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetIPForwarding(data.ResourceName, false),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyNetworkProfile_true_ipforwarding(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetIPForwarding(data.ResourceName, true),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_updateNetworkProfile_ipconfiguration_dns_name_label(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	expectedDNS := fmt.Sprintf("test-domain-label-%[1]d", data.RandomInteger)
	updatedDNS := fmt.Sprintf("test-updated-domain-label-%[1]d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetDomainNameLabel(data.ResourceName, expectedDNS),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP_updatedDNS_label(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetDomainNameLabel(data.ResourceName, updatedDNS),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_verify_key_data_changed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	initialKeyData := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
	expectedKeyData := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetSshKey(data.ResourceName, initialKeyData),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linuxKeyDataUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetSshKey(data.ResourceName, expectedKeyData),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicApplicationSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicApplicationSecurity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetAcceleratedNetworking(data.ResourceName, true),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicIPForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicIPForwarding(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicDNSSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicDNSSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_bootDiagnostic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "boot_diagnostics.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),

					// single placement group should default to true
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(data.ResourceName, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(data.ResourceName, false),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_linuxUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linuxUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customDataUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(data, "Standard_D1_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk_resize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(data, "Standard_D1_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(data, "Standard_D2_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_planManagedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_planManagedDisk(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")
	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := fmt.Sprintf("testadmin%d", data.RandomInteger)
	password := fmt.Sprintf("Password1234!%d", data.RandomInteger)
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMImage_standaloneImage_setup(data, userName, password, hostName, "LRS"),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_customImage(data, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMImageExists("azurerm_image", true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_applicationGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetApplicationGatewayTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetHasApplicationGateway(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetHasLoadbalancer(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancerManagedDataDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetHasDataDisks(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_overprovision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetOverprovision(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_priority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetPriorityTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "Low"),
					resource.TestCheckResourceAttr(data.ResourceName, "eviction_policy", "Deallocate"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_SystemAssignedMSI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetSystemAssignedMSI(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "0"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_UserAssignedMSI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetUserAssignedMSI(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.principal_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleAssignedMSI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetMultipleAssignedMSI(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned, UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtensionTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extensionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtensionTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleExtensions_provision_after_extension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate_provision_after_extension(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(data),
				ExpectError: regexp.MustCompile("Conflict between `vhd_containers`"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetNonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMVirtualMachineScaleSetNonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_linux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"os_profile.0.admin_password",
				"os_profile.0.custom_data",
			),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_AutoUpdates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_rollingAutoUpdates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data, "Manual"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_policy_mode", "Manual"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "rolling_upgrade_policy.#"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data, "Automatic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_policy_mode", "Automatic"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "rolling_upgrade_policy.#"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data, "Rolling"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_policy_mode", "Rolling"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_batch_instance_percent", "21"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_instance_percent", "22"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_upgraded_instance_percent", "23"),
				),
			},
			{
				PreConfig: func() { time.Sleep(1 * time.Minute) }, // VM Scale Set updates are not allowed while there is a Rolling Upgrade in progress.
				Config:    testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data, "Automatic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_policy_mode", "Automatic"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_batch_instance_percent", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_instance_percent", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_upgraded_instance_percent", "20"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data, "Manual"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "upgrade_policy_mode", "Manual"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_batch_instance_percent", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_instance_percent", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "rolling_upgrade_policy.0.max_unhealthy_upgraded_instance_percent", "20"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_importBasic_managedDisk_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk_withZones(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep("os_profile.0.admin_password"),
		},
	})
}

func testGetAzureRMVirtualMachineScaleSet(s *terraform.State, resourceName string) (result *compute.VirtualMachineScaleSet, err error) {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMVirtualMachineScaleSetIPForwarding(name string, boolean bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ipForwarding := *(*n)[0].EnableIPForwarding
		if ipForwarding != boolean {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, boolean, ipForwarding)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetDomainNameLabel(name string, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ipconfig := *(*n)[0].IPConfigurations
		dnsLabel := *ipconfig[0].PublicIPAddressConfiguration.DNSSettings.DomainNameLabel
		if dnsLabel != expected {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, expected, dnsLabel)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetSshKey(name string, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.OsProfile.LinuxConfiguration
		if n == nil {
			return fmt.Errorf("Bad: Could not get linux configuration for scale set %v", name)
		}

		publicKey := *n.SSH.PublicKeys
		keyData := *(publicKey[0]).KeyData
		if keyData != expected {
			return fmt.Errorf("Bad: ssh_keys.key_data set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, expected, keyData)
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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMVirtualMachineScaleSet_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  zones               = []

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set" "import" {
  name                = azurerm_virtual_machine_scale_set.test.name
  location            = azurerm_virtual_machine_scale_set.test.location
  resource_group_name = azurerm_virtual_machine_scale_set.test.resource_group_name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSet_evictionPolicyDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  priority            = "Low"
  eviction_policy     = "Delete"

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_standardSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                   = "acctvmss-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  upgrade_policy_mode    = "Manual"
  single_placement_group = false

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_withPPG(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "accPPG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                   = "acctvmss-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  upgrade_policy_mode    = "Manual"
  single_placement_group = false

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicPublicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  tags = {
    state = "create"
  }

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 0
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
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP_updated_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  tags = {
    state = "update"
  }

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 0
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
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicEmptyNetworkProfile_true_ipforwarding(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  tags = {
    state = "update"
  }

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 0
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name          = "TestNetworkProfile-%[1]d"
    primary       = true
    ip_forwarding = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicEmptyPublicIP_updatedDNS_label(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  tags = {
    state = "create"
  }

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 0
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
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-updated-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicApplicationSecurity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_application_security_group" "test" {
  location            = azurerm_resource_group.test.location
  name                = "TestApplicationSecurityGroup"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name                           = "TestIPConfiguration"
      primary                        = true
      subnet_id                      = azurerm_subnet.test.id
      application_security_group_ids = [azurerm_application_security_group.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D4_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name                   = "TestNetworkProfile-%[1]d"
    primary                = true
    accelerated_networking = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicIPForwarding(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D4_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name          = "TestNetworkProfile-%[1]d"
    primary       = true
    ip_forwarding = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicDNSSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D4_v2"
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

    dns_settings {
      dns_servers = ["8.8.8.8", "8.8.4.4"]
    }

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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

  boot_diagnostics {
    storage_uri = azurerm_storage_account.test.primary_blob_endpoint
  }

  network_profile {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
    name                      = "TestNetworkProfile-%[1]d"
    primary                   = true
    network_security_group_id = azurerm_network_security_group.test.id

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = azurerm_subnet.test.id
      primary   = true

      public_ip_address_configuration {
        name              = "TestPublicIPConfiguration"
        domain_name_label = "test-domain-label-%[1]d"
        idle_timeout      = 4
      }
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_windows_config {
    enable_automatic_upgrades = false
    provision_vm_agent        = true

    winrm {
      protocol = "http"
    }
  }

  network_profile {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                   = "acctvmss-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  upgrade_policy_mode    = "Manual"
  single_placement_group = false

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_linux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
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
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_linuxUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
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
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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

  tags = {
    ThisIs = "a test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
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
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_linuxKeyDataUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  upgrade_policy_mode = "Automatic"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = "1"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    custom_data          = "updated custom data!"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk_withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  zones               = ["1", "2"]

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(data acceptance.TestData, vmSize string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "%[3]s"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_windows_config {
    enable_automatic_upgrades = false
    provision_vm_agent        = true

    additional_unattend_config {
      pass         = "oobeSystem"
      component    = "Microsoft-Windows-Shell-Setup"
      setting_name = "AutoLogon"
      content      = "<AutoLogon><Username>myadmin</Username><Password><Value>Passwword1234</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount></AutoLogon>"
    }

    additional_unattend_config {
      pass         = "oobeSystem"
      component    = "Microsoft-Windows-Shell-Setup"
      setting_name = "FirstLogonCommands"
      content      = "<FirstLogonCommands><SynchronousCommand><CommandLine>shutdown /r /t 0 /c \"initial reboot\"</CommandLine><Description>reboot</Description><Order>1</Order></SynchronousCommand></FirstLogonCommands>"
    }
  }

  network_profile {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, vmSize)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetApplicationGatewayTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary                                      = true
      subnet_id                                    = azurerm_subnet.test.id
      application_gateway_backend_address_pool_ids = [azurerm_application_gateway.test.backend_address_pool[0].id]
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
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.3.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = azurerm_subnet.gwtest.id
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = azurerm_public_ip.test.id
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = azurerm_subnet.gwtest.id

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

  tags = {
    environment = "tf01"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  name                           = "ssh"
  loadbalancer_id                = azurerm_lb.test.id
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "default"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetPriorityTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  overprovision       = false
  priority            = "Low"
  eviction_policy     = "Deallocate"

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetSystemAssignedMSI(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }

  extension {
    name                 = "MSILinuxExtension"
    publisher            = "Microsoft.ManagedIdentity"
    type                 = "ManagedIdentityExtensionForLinux"
    type_handler_version = "1.0"
    settings             = "{\"port\": 50342}"
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetUserAssignedMSI(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[3]s"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  extension {
    name                 = "MSILinuxExtension"
    publisher            = "Microsoft.ManagedIdentity"
    type                 = "ManagedIdentityExtensionForLinux"
    type_handler_version = "1.0"
    settings             = "{\"port\": 50342}"
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate_provision_after_extension(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
    provision_after_extensions = ["CustomScript"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
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
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSetNonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_planManagedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  plan {
    name      = "os"
    product   = "rancheros"
    publisher = "rancher"
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "rancher"
    offer     = "rancheros"
    sku       = "os"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_customImage(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%[3]s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
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

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "test" {
  name                = "accteste-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  storage_profile_image_reference {
    id = azurerm_image.test.id
  }

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
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
}
`, data.RandomInteger, data.Locations.Primary, hostName, userName, password)
}

func testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
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
    name    = "primary-%[1]d"
    primary = true

    ip_configuration {
      name      = "primary"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  network_profile {
    name    = "secondary-%[1]d"
    primary = false

    ip_configuration {
      name      = "secondary"
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_rollingAutoUpdates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/16"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 4
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_rule" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "AccTestLBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "PublicIPAddress"
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbapool"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  upgrade_policy_mode  = "Rolling"
  automatic_os_upgrade = true
  health_probe_id      = azurerm_lb_probe.test.id
  depends_on           = [azurerm_lb_rule.test]

  rolling_upgrade_policy {
    max_batch_instance_percent              = 21
    max_unhealthy_instance_percent          = 22
    max_unhealthy_upgraded_instance_percent = 23
    pause_time_between_batches              = "PT30S"
  }

  sku {
    name     = "Standard_F2"
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
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      primary                                = true
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
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualMachineScaleSet_upgradeModeUpdate(data acceptance.TestData, mode string) string {
	policy := ""
	if mode == "Rolling" {
		policy = `
  rolling_upgrade_policy {
    max_batch_instance_percent              = 21
    max_unhealthy_instance_percent          = 22
    max_unhealthy_upgraded_instance_percent = 23
  }`
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/16"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 4
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_rule" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "AccTestLBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "PublicIPAddress"
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbapool"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  upgrade_policy_mode = "%[3]s"
  health_probe_id     = azurerm_lb_probe.test.id
  depends_on          = [azurerm_lb_rule.test]

  %[4]s

  sku {
    name     = "Standard_F2"
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
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      primary                                = true
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
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, mode, policy)
}

func testAccAzureRMVirtualMachineScaleSetMultipleAssignedMSI(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
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
      primary   = true
      subnet_id = azurerm_subnet.test.id
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
