package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccWindowsVirtualMachine_otherPatchModeManual(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherPatchModeManual(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherPatchModeAutomaticByOS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherPatchModeAutomaticByOS(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherPatchModeAutomaticByPlatform(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherPatchModeAutomaticByPlatform(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherPatchModeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherPatchModeAutomaticByOS(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherPatchModeAutomaticByPlatform(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherPatchModeManual(data), // this update requires force replacement actually
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherAdditionalUnattendContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherAdditionalUnattendContent(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"additional_unattend_content.0.content",
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherAllowExtensionOperationsDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherAllowExtensionOperationsDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherAllowExtensionOperationsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherAllowExtensionOperationsUpdatedWithoutVmAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsDisabledWithoutVmAgent(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherAllowExtensionOperationsEnabledWithoutVmAgent(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_extension_operations", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherExtensionsTimeBudget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherExtensionsTimeBudget(data, "PT30M"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extensions_time_budget", "PT30M"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherExtensionsTimeBudgetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherExtensionsTimeBudget(data, "PT30M"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extensions_time_budget", "PT30M"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherExtensionsTimeBudget(data, "PT50M"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extensions_time_budget", "PT50M"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherExtensionsTimeBudget(data, "PT30M"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extensions_time_budget", "PT30M"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherBootDiagnostics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testWindowsVirtualMachine_otherBootDiagnostics(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Disabled
				Config: testWindowsVirtualMachine_otherBootDiagnosticsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Enabled
				Config: testWindowsVirtualMachine_otherBootDiagnostics(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherBootDiagnosticsManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testWindowsVirtualMachine_otherBootDiagnosticsManaged(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Disabled
				Config: testWindowsVirtualMachine_otherBootDiagnosticsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Enabled
				Config: testWindowsVirtualMachine_otherBootDiagnosticsManaged(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherComputerNameDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherComputerNameDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "computer_name"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherComputerNameDefaultInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testWindowsVirtualMachine_otherComputerNameDefaultInvalid(data),
				ExpectError: regexp.MustCompile("unable to assume default computer name"),
			},
		},
	})
}

func TestAccWindowsVirtualMachine_otherComputerNameCustom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherComputerNameCustom(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "computer_name", "custom123"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherCustomData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherCustomData(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password", "custom_data"),
		},
	})
}

func TestAccWindowsVirtualMachine_otherEnableAutomaticUpdatesDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherEnableAutomaticUpdatesDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_automatic_updates", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherEnableAutomaticUpdatesDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherEnableAutomaticUpdatesDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_automatic_updates", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherLicenseTypeNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherLicenseType(data, "None"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherLicenseTypeWindowsClient(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherLicenseTypeWindowsClient(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherLicenseTypeWindowsServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherLicenseType(data, "Windows_Server"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherLicenseTypeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherLicenseTypeDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherLicenseTypeWindowsClient(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "Windows_Client"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherLicenseTypeDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherPrioritySpot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherPrioritySpot(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherPrioritySpotMaxBidPrice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// expensive, but guarantees this test will pass
				Config: testWindowsVirtualMachine_otherPrioritySpotMaxBidPrice(data, "0.5000"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// no limit
				Config: testWindowsVirtualMachine_otherPrioritySpotMaxBidPrice(data, "-1"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherProvisionVMAgentDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherProvisionVMAgentDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "provision_vm_agent", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherProvisionVMAgentDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherProvisionVMAgentDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "provision_vm_agent", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testWindowsVirtualMachine_requiresImport),
		},
	})
}

func TestAccWindowsVirtualMachine_otherSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherSecret(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherSecretUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherSecretRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "secret.#", "0"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherTags(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherTimeZone(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherUltraSsdDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherUltraSsdDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_capabilities.0.ultra_ssd_enabled", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherUltraSsdEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherUltraSsdEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_capabilities.0.ultra_ssd_enabled", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherUltraSsdUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherUltraSsdDefault(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_capabilities.0.ultra_ssd_enabled", "false"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_otherUltraSsdEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_capabilities.0.ultra_ssd_enabled", "true"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherWinRMHTTP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherWinRMHTTP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherWinRMHTTPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherWinRMHTTPS(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherEncryptionAtHostEnabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherEncryptionAtHostEnabledWithCMK(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsVirtualMachine_otherEncryptionAtHostEnabledWithCMK(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testWindowsVirtualMachine_otherPatchModeManual(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  enable_automatic_updates = false
  patch_mode               = "Manual"
}
`, template)
}

func testWindowsVirtualMachine_otherPatchModeAutomaticByOS(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  patch_mode = "AutomaticByOS"
}
`, template)
}

func testWindowsVirtualMachine_otherPatchModeAutomaticByPlatform(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  patch_mode = "AutomaticByPlatform"
}
`, template)
}

func TestAccWindowsVirtualMachine_otherGracefulShutdownDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherGracefulShutdown(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_otherGracefulShutdownEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_otherGracefulShutdown(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testWindowsVirtualMachine_otherAdditionalUnattendContent(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  additional_unattend_content {
    setting = "AutoLogon"
    content = "<AutoLogon><Username>myadmin</Username><Password><Value>P@ssword1234!</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount></AutoLogon>"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherAllowExtensionOperationsDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherAllowExtensionOperationsDisabled(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                       = local.vm_name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  admin_password             = "P@$$w0rd1234!"
  allow_extension_operations = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherAllowExtensionOperationsDisabledWithoutVmAgent(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                       = local.vm_name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  admin_password             = "P@$$w0rd1234!"
  allow_extension_operations = false
  provision_vm_agent         = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherAllowExtensionOperationsEnabledWithoutVmAgent(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                       = local.vm_name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  admin_password             = "P@$$w0rd1234!"
  allow_extension_operations = true
  provision_vm_agent         = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherExtensionsTimeBudget(data acceptance.TestData, duration string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  extensions_time_budget = "%s"
}
`, template, duration)
}

func testWindowsVirtualMachine_otherBootDiagnostics(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherBootDiagnosticsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.test.primary_blob_endpoint
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherBootDiagnosticsManaged(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherBootDiagnosticsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  boot_diagnostics {}

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherBootDiagnosticsDisabled(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherBootDiagnosticsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherBootDiagnosticsTemplate(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, template, data.RandomString)
}

func testWindowsVirtualMachine_otherComputerNameDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherComputerNameDefaultInvalid(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = "${local.vm_name}-this-too-long-to-be-a-computer-name"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherComputerNameCustom(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  computer_name       = "custom123"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherCustomData(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  custom_data         = base64encode("/bin/bash")
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherEnableAutomaticUpdatesDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherEnableAutomaticUpdatesDisabled(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                     = local.vm_name
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  size                     = "Standard_F2"
  admin_username           = "adminuser"
  admin_password           = "P@$$w0rd1234!"
  enable_automatic_updates = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherLicenseTypeDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsDesktop"
    offer     = "Windows-10"
    sku       = "rs5-pro"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherLicenseType(data acceptance.TestData, licenseType string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  license_type        = %q
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, licenseType)
}

func testWindowsVirtualMachine_otherLicenseTypeWindowsClient(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  license_type        = "Windows_Client"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsDesktop"
    offer     = "Windows-10"
    sku       = "rs5-pro"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherPrioritySpot(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  eviction_policy     = "Deallocate"
  priority            = "Spot"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherPrioritySpotMaxBidPrice(data acceptance.TestData, maxBidPrice string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  eviction_policy     = "Deallocate"
  priority            = "Spot"
  max_bid_price       = %q
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, maxBidPrice)
}

func testWindowsVirtualMachine_otherProvisionVMAgentDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherProvisionVMAgentDisabled(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                       = local.vm_name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  provision_vm_agent         = false
  allow_extension_operations = false
  admin_username             = "adminuser"
  admin_password             = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_requiresImport(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_authPassword(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "import" {
  name                  = azurerm_windows_virtual_machine.test.name
  resource_group_name   = azurerm_windows_virtual_machine.test.resource_group_name
  location              = azurerm_windows_virtual_machine.test.location
  size                  = azurerm_windows_virtual_machine.test.size
  admin_username        = "adminuser"
  admin_password        = "P@$$w0rd1234!"
  network_interface_ids = azurerm_windows_virtual_machine.test.network_interface_ids

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherSecretTemplate(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}
`, template, data.RandomString)
}

func testWindowsVirtualMachine_otherSecret(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.first.secret_id
    }
  }
}
`, template)
}

func testWindowsVirtualMachine_otherSecretRemoved(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherSecretUpdated(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_otherSecretTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.first.secret_id
    }

    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.second.secret_id
    }
  }
}
`, template)
}

func testWindowsVirtualMachine_otherTags(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  tags = {
    Hello = "World"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherTagsUpdated(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  tags = {
    Hello     = "World"
    Dimension = "C-137"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherTimeZone(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  timezone            = "Hawaiian Standard Time"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherUltraSsdDefault(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherUltraSsdEnabled(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }
}
`, template)
}

func testWindowsVirtualMachine_otherWinRMHTTP(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  winrm_listener {
    protocol = "Http"
  }
}
`, template)
}

func testWindowsVirtualMachine_otherWinRMHTTPS(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "backup",
      "create",
      "decrypt",
      "delete",
      "encrypt",
      "get",
      "import",
      "list",
      "purge",
      "recover",
      "restore",
      "sign",
      "unwrapKey",
      "update",
      "verify",
      "wrapKey",
    ]

    secret_permissions = [
      "backup",
      "delete",
      "get",
      "list",
      "purge",
      "recover",
      "restore",
      "set",
    ]

    certificate_permissions = [
      "create",
      "delete",
      "deleteissuers",
      "get",
      "getissuers",
      "import",
      "list",
      "listissuers",
      "managecontacts",
      "manageissuers",
      "setissuers",
      "update",
    ]
  }

  enabled_for_deployment          = true
  enabled_for_template_deployment = true
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "example"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=${local.vm_name}"
      validity_in_months = 12
    }
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.test.secret_id
    }
  }

  winrm_listener {
    protocol = "Http"
  }

  winrm_listener {
    certificate_url = azurerm_key_vault_certificate.test.secret_id
    protocol        = "Https"
  }
}
`, template, data.RandomString)
}

func testAccWindowsVirtualMachine_otherEncryptionAtHostEnabled(data acceptance.TestData, enabled bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_DS3_V2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  encryption_at_host_enabled = %t
}
`, template, enabled)
}

func testAccWindowsVirtualMachine_otherEncryptionAtHostEnabledWithCMK(data acceptance.TestData, enabled bool) string {
	template := testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_DS3_V2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching                = "ReadWrite"
    storage_account_type   = "Standard_LRS"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  encryption_at_host_enabled = %t

  depends_on = [
    azurerm_role_assignment.disk-encryption-read-keyvault,
    azurerm_key_vault_access_policy.disk-encryption,
  ]
}
`, template, enabled)
}

func testWindowsVirtualMachine_otherGracefulShutdown(data acceptance.TestData, gracefulShutdown bool) string {
	return fmt.Sprintf(`
locals {
  vm_name = "acctestvm%s"
}

provider "azurerm" {
  features {
    virtual_machine {
      graceful_shutdown = %t
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, data.RandomString, gracefulShutdown, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
