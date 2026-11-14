// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyVirtualMachineConfigurationAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := guestconfigurationassignments.ParseVirtualMachineProviders2GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	vmId := commonids.NewVirtualMachineID(subscriptionId, id.ResourceGroupName, id.VirtualMachineName)
	d.Set("name", id.GuestConfigurationAssignmentName)
	d.Set("virtual_machine_id", vmId.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("configuration", flattenGuestConfigurationAssignment(props.GuestConfiguration)); err != nil {
				return fmt.Errorf("setting `configuration`: %+v", err)
			}
		}
	}
	return nil
}

func flattenGuestConfigurationAssignment(input *guestconfigurationassignments.GuestConfigurationNavigation) []any {
	if input == nil {
		return make([]any, 0)
	}

	version := pointer.From(input.Version)
	assignmentType := pointer.From(input.AssignmentType)
	contentHash := pointer.From(input.ContentHash)
	contentUri := pointer.From(input.ContentUri)
	return []any{
		map[string]any{
			"assignment_type": string(assignmentType),
			"content_hash":    contentHash,
			"content_uri":     contentUri,
			"parameter":       flattenGuestConfigurationAssignmentConfigurationParameters(input.ConfigurationParameter),
			"version":         version,
		},
	}
}

func flattenGuestConfigurationAssignmentConfigurationParameters(input *[]guestconfigurationassignments.ConfigurationParameter) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		name := pointer.From(item.Name)
		value := pointer.From(item.Value)
		results = append(results, map[string]any{
			"name":  name,
			"value": value,
		})
	}
	return results
}
