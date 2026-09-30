// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name virtual_machine_extension -service-package-name compute -properties "name" -compare-values "virtual_machine_name:virtual_machine_id,resource_group_name:virtual_machine_id,subscription_id:virtual_machine_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceVirtualMachineExtensionsCreateUpdate,
		Read:     resourceVirtualMachineExtensionsRead,
		Update:   resourceVirtualMachineExtensionsCreateUpdate,
		Delete:   resourceVirtualMachineExtensionsDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&virtualmachineextensions.ExtensionId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualmachineextensions.ExtensionId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: virtualMachineExtensionSchema(),
	}
}
