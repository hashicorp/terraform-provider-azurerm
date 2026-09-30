// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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

		Schema: windowsVirtualMachineSchema(),
	}
}
