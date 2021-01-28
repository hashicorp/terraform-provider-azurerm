package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachineScaleSet_otherBootDiagnostics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Enabled
			Config: r.otherBootDiagnostics(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// Removed
			Config: r.otherBootDiagnosticsDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// Enabled
			Config: r.otherBootDiagnostics(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherBootDiagnosticsManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Enabled
			Config: r.otherBootDiagnosticsManaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// Removed
			Config: r.otherBootDiagnosticsDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// Enabled
			Config: r.otherBootDiagnosticsManaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherComputerNamePrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherComputerNamePrefix(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherComputerNamePrefixInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.otherComputerNamePrefixInvalid(data),
			ExpectError: regexp.MustCompile("unable to assume default computer name prefix"),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherCustomData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherCustomData(data, "/bin/bash"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
			"custom_data",
		),
		{
			Config: r.otherCustomData(data, "/bin/zsh"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
			"custom_data",
		),
		{
			// removed
			Config: r.authPassword(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
			"custom_data",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherDoNotRunExtensionsOnOverProvisionedMachines(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherDoNotRunExtensionsOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherDoNotRunExtensionsOnOverProvisionedMachinesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherDoNotRunExtensionsOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.otherDoNotRunExtensionsOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.otherDoNotRunExtensionsOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherPrioritySpotDeallocate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherPrioritySpot(data, "Deallocate"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherPrioritySpotDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherPrioritySpot(data, "Delete"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherPrioritySpotMaxBidPrice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// expensive, but guarantees this test will pass
			Config: r.otherPrioritySpotMaxBidPrice(data, "0.5000"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.otherPrioritySpotMaxBidPrice(data, "-1"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherPriorityRegular(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherPriorityRegular(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.authPassword(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.otherRequiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_linux_virtual_machine_scale_set"),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherSecret(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// update
			Config: r.otherSecretUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),

		{
			// removed
			Config: r.otherSecretRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// add one
			Config: r.otherTagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// remove all
			Config: r.authPassword(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVMAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVMAgent(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVMAgentDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVMAgent(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_updateHealthProbe(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.updateLoadBalancerHealthProbeSKUBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateLoadBalancerHealthProbeSKUStandard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherScaleInPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherScaleInPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scale_in_policy").HasValue("Default"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherTerminateNotification(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// turn terminate notification on
		{
			Config: r.otherTerminateNotification(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("terminate_notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("terminate_notification.0.enabled").HasValue("true"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		// turn terminate notification off
		{
			Config: r.otherTerminateNotification(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("terminate_notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("terminate_notification.0.enabled").HasValue("false"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		// turn terminate notification on again
		{
			Config: r.otherTerminateNotification(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("terminate_notification.#").HasValue("1"),
				check.That(data.ResourceName).Key("terminate_notification.0.enabled").HasValue("true"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherAutomaticRepairsPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// turn automatic repair on
		{
			Config: r.otherAutomaticRepairsPolicy(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		// turn automatic repair off
		{
			Config: r.otherAutomaticRepairsPolicy(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		// turn automatic repair on again
		{
			Config: r.otherAutomaticRepairsPolicy(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherUpgradeMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherUpgradeMode(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVmExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensions(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVmExtensionsOnlySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensionsOnlySettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVmExtensionsForceUpdateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensionsForceUpdateTag(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherVmExtensionsForceUpdateTag(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVmExtensionsMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensionsMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherVmExtensionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensions(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherVmExtensionsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherVmExtensions(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherVmExtensionsWithExtensionsTimeBudget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensionsWithExtensionsTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherVmExtensionsWithExtensionsTimeBudgetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherVmExtensionsWithExtensionsTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherVmExtensionsWithExtensionsTimeBudget(data, "PT1H"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherVmExtensionsWithExtensionsTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherExtensionsTimeBudgetWithoutExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherExtensionsTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_otherExtensionsTimeBudgetWithoutExtensionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherExtensionsTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherExtensionsTimeBudgetWithoutExtensions(data, "PT1H"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherExtensionsTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherEncryptionAtHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherEncryptionAtHost(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherEncryptionAtHostUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherEncryptionAtHost(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherEncryptionAtHost(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
		{
			Config: r.otherEncryptionAtHost(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherEncryptionAtHostWithCMK(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherEncryptionAtHostWithCMK(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_otherPlatformFaultDomainCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.otherPlatformFaultDomainCount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// TODO - extension should be changed to extension.0.protected_settings when either binary testing is available or this feature is promoted from beta
		data.ImportStep("admin_password", "extension"),
	})
}

func (r LinuxVirtualMachineScaleSetResource) otherBootDiagnostics(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.test.primary_blob_endpoint
  }
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherBootDiagnosticsManaged(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  boot_diagnostics {}
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherBootDiagnosticsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherComputerNamePrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_F2"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@ssword1234!"
  computer_name_prefix = "my-linux-computer-name-prefix"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherComputerNamePrefixInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d-too-long-to-be-a-computer-name-but-not-vmss-name"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherCustomData(data acceptance.TestData, customData string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  custom_data         = base64encode(%q)

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, customData)
}

func (r LinuxVirtualMachineScaleSetResource) otherDoNotRunExtensionsOnOverProvisionedMachines(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = true

  disable_password_authentication                   = false
  do_not_run_extensions_on_overprovisioned_machines = %t

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherPrioritySpot(data acceptance.TestData, evictionPolicy string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  eviction_policy     = %q
  priority            = "Spot"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, evictionPolicy)
}

func (r LinuxVirtualMachineScaleSetResource) otherPrioritySpotMaxBidPrice(data acceptance.TestData, maxBid string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  eviction_policy     = "Delete"
  priority            = "Spot"
  max_bid_price       = %q

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, maxBid)
}

func (r LinuxVirtualMachineScaleSetResource) otherPriorityRegular(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  priority            = "Regular"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherRequiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "import" {
  name                = azurerm_linux_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_linux_virtual_machine_scale_set.test.resource_group_name
  location            = azurerm_linux_virtual_machine_scale_set.test.location
  sku                 = azurerm_linux_virtual_machine_scale_set.test.sku
  instances           = azurerm_linux_virtual_machine_scale_set.test.instances
  admin_username      = azurerm_linux_virtual_machine_scale_set.test.admin_username
  admin_password      = azurerm_linux_virtual_machine_scale_set.test.admin_password

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.authPassword(data))
}

func (r LinuxVirtualMachineScaleSetResource) otherSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }
  }
}
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherSecretRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherSecretUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }

    certificate {
      url = azurerm_key_vault_certificate.second.secret_id
    }
  }
}
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherSecretTemplate(data acceptance.TestData) string {
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
      "purge",
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
`, r.template(data), data.RandomString)
}

func (r LinuxVirtualMachineScaleSetResource) otherTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  tags = {
    artist = "Billy"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherTagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  tags = {
    artist = "Billy"
    when   = "we all fall asleep"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherVMAgent(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false
  provision_vm_agent              = %t

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) updateLoadBalancerHealthProbeSKUBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 4
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }

  depends_on = [azurerm_lb_rule.test]
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) updateLoadBalancerHealthProbeSKUStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 4
  sku                     = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }

  depends_on = [azurerm_lb_rule.test]
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherScaleInPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    disk_size_gb         = 30
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  scale_in_policy = "Default"
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherTerminateNotification(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    disk_size_gb         = 30
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  terminate_notification {
    enabled = %t
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherAutomaticRepairsPolicy(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                    = "acctestpip-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 4
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }

  automatic_instance_repair {
    enabled = %[3]t
  }

  depends_on = [azurerm_lb_rule.test]
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherUpgradeMode(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  upgrade_mode                    = "Automatic"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    disk_size_gb         = 30
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  terminate_notification {
    enabled = %t
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensionsOnlySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensionsForceUpdateTag(data acceptance.TestData, updateTag string) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true
    force_update_tag           = %q

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger, updateTag)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensionsMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    provision_after_extensions = ["VMAccessForLinux"]

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })

  }

  extension {
    name                       = "VMAccessForLinux"
    publisher                  = "Microsoft.OSTCExtensions"
    type                       = "VMAccessForLinux"
    type_handler_version       = "1.5"
    auto_upgrade_minor_version = true

    protected_settings = jsonencode({
      "reset_ssh" = "True"
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $(date)"
    })
  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) otherVmExtensionsWithExtensionsTimeBudget(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  extensions_time_budget = "%s"
}
`, template, data.RandomInteger, duration)
}

func (r LinuxVirtualMachineScaleSetResource) otherExtensionsTimeBudgetWithoutExtensions(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  extensions_time_budget = "%s"
}
`, template, data.RandomInteger, duration)
}

func (r LinuxVirtualMachineScaleSetResource) otherEncryptionAtHost(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_DS3_V2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  encryption_at_host_enabled = %t
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherEncryptionAtHostWithCMK(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_DS3_V2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type   = "Standard_LRS"
    caching                = "ReadWrite"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  encryption_at_host_enabled = %t

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.disksOSDisk_diskEncryptionSetResource(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) otherPlatformFaultDomainCount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  platform_fault_domain_count = 3
}
`, r.template(data), data.RandomInteger)
}
