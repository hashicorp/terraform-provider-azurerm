// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyVirtualMachineConfigurationAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := commonids.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID(subscriptionId, vmId.ResourceGroupName, vmId.VirtualMachineName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for present of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_policy_virtual_machine_configuration_assignment", id.ID())
			}
		}
	}
	guestConfiguration := expandGuestConfigurationAssignment(d.Get("configuration").([]any), id.GuestConfigurationAssignmentName)
	assignment := guestconfigurationassignments.GuestConfigurationAssignment{
		Name:     *pointer.To(id.GuestConfigurationAssignmentName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &guestconfigurationassignments.GuestConfigurationAssignmentProperties{
			GuestConfiguration: guestConfiguration,
		},
	}

	// I need to determine if the passed in guest config is a built-in config or not
	// since the attribute is computed and optional I need to check the value of the
	// contentURI to see if it is on a service team owned storage account or not
	// all built-in guest configuration will always be on a service team owned
	// storage account
	if guestConfiguration.ContentUri != nil || *guestConfiguration.ContentUri != "" {
		if strings.Contains(strings.ToLower(*guestConfiguration.ContentUri), "oaasguestconfig") {
			assignment.Properties.GuestConfiguration.ContentHash = nil
			assignment.Properties.GuestConfiguration.ContentUri = nil
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, assignment); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourcePolicyVirtualMachineConfigurationAssignmentRead(d, meta)
}

func expandGuestConfigurationAssignment(input []any, name string) *guestconfigurationassignments.GuestConfigurationNavigation {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	result := guestconfigurationassignments.GuestConfigurationNavigation{
		Name:                   pointer.To(name),
		Version:                pointer.To(v["version"].(string)),
		ConfigurationParameter: expandGuestConfigurationAssignmentConfigurationParameters(v["parameter"].(*pluginsdk.Set).List()),
	}

	if v, ok := v["assignment_type"]; ok {
		result.AssignmentType = pointer.ToEnum[guestconfigurationassignments.AssignmentType](v.(string))
	}

	if v, ok := v["content_hash"]; ok {
		result.ContentHash = pointer.To(v.(string))
	}

	if v, ok := v["content_uri"]; ok {
		result.ContentUri = pointer.To(v.(string))
	}

	return &result
}

func expandGuestConfigurationAssignmentConfigurationParameters(input []any) *[]guestconfigurationassignments.ConfigurationParameter {
	results := make([]guestconfigurationassignments.ConfigurationParameter, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, guestconfigurationassignments.ConfigurationParameter{
			Name:  pointer.To(v["name"].(string)),
			Value: pointer.To(v["value"].(string)),
		})
	}
	return &results
}
