// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

const azureWindowsVirtualMachineResourceName = "azurerm_windows_virtual_machine"

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name authPassword

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWindowsVirtualMachineCreate,
		Read:   resourceWindowsVirtualMachineRead,
		Update: resourceWindowsVirtualMachineUpdate,
		Delete: resourceWindowsVirtualMachineDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualmachines.VirtualMachineId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentityThen(&virtualmachines.VirtualMachineId{}, helpers.ImportVirtualMachine(virtualmachines.OperatingSystemTypesWindows, azureWindowsVirtualMachineResourceName)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(resourceWindowsVirtualMachineCustomizeDiff),
		),

		Schema: windowsVirtualMachineSchema(),
	}
}

func resourceWindowsVirtualMachineCustomizeDiff(ctx context.Context, diff *pluginsdk.ResourceDiff, meta any) error {
	rawConfig := diff.GetRawConfig().AsValueMap()
	adminUsername, ok := rawConfig["admin_username"]
	if !ok || adminUsername.IsNull() {
		return nil
	}

	if osManagedDiskId, ok := rawConfig["os_managed_disk_id"]; ok && (!osManagedDiskId.IsNull() || !osManagedDiskId.IsKnown()) {
		return nil
	}

	adminPasswordSet := false
	if adminPassword, ok := rawConfig["admin_password"]; ok && !adminPassword.IsNull() {
		adminPasswordSet = true
	}

	adminPasswordWOVersionSet := false
	if adminPasswordWOVersion, ok := rawConfig["admin_password_wo_version"]; ok && !adminPasswordWOVersion.IsNull() {
		adminPasswordWOVersionSet = true
	}

	if !adminPasswordSet && !adminPasswordWOVersionSet {
		return errors.New("one of `admin_password` or `admin_password_wo_version` must be specified when `admin_username` is specified")
	}

	return nil
}
