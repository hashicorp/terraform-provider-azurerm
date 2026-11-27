// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	policyClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyVirtualMachineConfigurationAssignmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := policyClient.From(meta).GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := guestconfigurationassignments.ParseVirtualMachineProviders2GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
