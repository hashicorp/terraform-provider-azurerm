// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name authSSH

const azureLinuxVirtualMachineResourceName = "azurerm_linux_virtual_machine"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLinuxVirtualMachineCreate,
		Read:   resourceLinuxVirtualMachineRead,
		Update: resourceLinuxVirtualMachineUpdate,
		Delete: resourceLinuxVirtualMachineDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualmachines.VirtualMachineId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentityThen(&virtualmachines.VirtualMachineId{}, helpers.ImportVirtualMachine(virtualmachines.OperatingSystemTypesLinux, azureLinuxVirtualMachineResourceName)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: linuxVirtualMachineSchema(),
	}
}
