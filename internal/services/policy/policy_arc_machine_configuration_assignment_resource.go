// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationhcrpassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2025-01-13/machines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyArcMachineConfigurationAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePolicyArcMachineConfigurationAssignmentCreateUpdate,
		Read:   resourcePolicyArcMachineConfigurationAssignmentRead,
		Update: resourcePolicyArcMachineConfigurationAssignmentCreateUpdate,
		Delete: resourcePolicyArcMachineConfigurationAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := guestconfigurationhcrpassignments.ParseProviders2GuestConfigurationAssignmentID(id)
			return err
		}),

		Schema: resourcePolicyArcMachineConfigurationAssignmentSchema(),
	}
}

func resourcePolicyArcMachineConfigurationAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"assignment_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(guestconfigurationhcrpassignments.AssignmentTypeAudit),
							string(guestconfigurationhcrpassignments.AssignmentTypeDeployAndAutoCorrect),
							string(guestconfigurationhcrpassignments.AssignmentTypeApplyAndAutoCorrect),
							string(guestconfigurationhcrpassignments.AssignmentTypeApplyAndMonitor),
						}, false),
					},

					"content_hash": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"parameter": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func resourcePolicyArcMachineConfigurationAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationHCRPAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := machines.ParseMachineID(d.Get("machine_id").(string))
	if err != nil {
		return err
	}

	id := guestconfigurationhcrpassignments.NewProviders2GuestConfigurationAssignmentID(subscriptionId, vmId.ResourceGroupName, vmId.MachineName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_policy_arc_machine_configuration_assignment", id.ID())
		}
	}
	guestConfiguration := expandGuestConfigurationHCRPAssignment(d.Get("configuration").([]interface{}), id.GuestConfigurationAssignmentName)
	assignment := guestconfigurationhcrpassignments.GuestConfigurationAssignment{
		Name:     *pointer.To(id.GuestConfigurationAssignmentName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &guestconfigurationhcrpassignments.GuestConfigurationAssignmentProperties{
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

	d.SetId(id.ID())

	return resourcePolicyArcMachineConfigurationAssignmentRead(d, meta)
}

func resourcePolicyArcMachineConfigurationAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationHCRPAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := guestconfigurationhcrpassignments.ParseProviders2GuestConfigurationAssignmentID(d.Id())
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

	vmId := guestconfigurationhcrpassignments.NewMachineID(subscriptionId, id.ResourceGroupName, id.MachineName)
	d.Set("name", id.GuestConfigurationAssignmentName)
	d.Set("machine_id", vmId.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("configuration", flattenGuestConfigurationHCRPAssignment(props.GuestConfiguration)); err != nil {
				return fmt.Errorf("setting `configuration`: %+v", err)
			}
		}
	}
	return nil
}

func resourcePolicyArcMachineConfigurationAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.GuestConfigurationHCRPAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := guestconfigurationhcrpassignments.ParseProviders2GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandGuestConfigurationHCRPAssignment(input []interface{}, name string) *guestconfigurationhcrpassignments.GuestConfigurationNavigation {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := guestconfigurationhcrpassignments.GuestConfigurationNavigation{
		Name:                   pointer.To(name),
		Version:                pointer.To(v["version"].(string)),
		ConfigurationParameter: expandGuestConfigurationHCRPAssignmentConfigurationParameters(v["parameter"].(*pluginsdk.Set).List()),
	}

	if v, ok := v["assignment_type"]; ok {
		result.AssignmentType = pointer.To(guestconfigurationhcrpassignments.AssignmentType(v.(string)))
	}

	if v, ok := v["content_hash"]; ok {
		result.ContentHash = pointer.To(v.(string))
	}

	if v, ok := v["content_uri"]; ok {
		result.ContentUri = pointer.To(v.(string))
	}

	return &result
}

func expandGuestConfigurationHCRPAssignmentConfigurationParameters(input []interface{}) *[]guestconfigurationhcrpassignments.ConfigurationParameter {
	results := make([]guestconfigurationhcrpassignments.ConfigurationParameter, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, guestconfigurationhcrpassignments.ConfigurationParameter{
			Name:  pointer.To(v["name"].(string)),
			Value: pointer.To(v["value"].(string)),
		})
	}
	return &results
}

func flattenGuestConfigurationHCRPAssignment(input *guestconfigurationhcrpassignments.GuestConfigurationNavigation) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var version string
	if input.Version != nil {
		version = *input.Version
	}
	var assignmentType guestconfigurationhcrpassignments.AssignmentType
	if input.AssignmentType != nil {
		assignmentType = *input.AssignmentType
	}
	var contentHash string
	if input.ContentHash != nil {
		contentHash = *input.ContentHash
	}
	var contentUri string
	if input.ContentUri != nil {
		contentUri = *input.ContentUri
	}
	return []interface{}{
		map[string]interface{}{
			"assignment_type": string(assignmentType),
			"content_hash":    contentHash,
			"content_uri":     contentUri,
			"parameter":       flattenGuestConfigurationHCRPAssignmentConfigurationParameters(input.ConfigurationParameter),
			"version":         version,
		},
	}
}

func flattenGuestConfigurationHCRPAssignmentConfigurationParameters(input *[]guestconfigurationhcrpassignments.ConfigurationParameter) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		var value string
		if item.Value != nil {
			value = *item.Value
		}
		results = append(results, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	return results
}
