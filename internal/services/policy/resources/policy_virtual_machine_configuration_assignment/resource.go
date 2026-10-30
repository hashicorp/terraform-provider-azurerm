// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePolicyVirtualMachineConfigurationAssignmentCreateUpdate,
		Read:   resourcePolicyVirtualMachineConfigurationAssignmentRead,
		Update: resourcePolicyVirtualMachineConfigurationAssignmentCreateUpdate,
		Delete: resourcePolicyVirtualMachineConfigurationAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := guestconfigurationassignments.ParseVirtualMachineProviders2GuestConfigurationAssignmentID(id)
			return err
		}),

		Schema: resourcePolicyVirtualMachineConfigurationAssignmentSchema(),
	}
}
