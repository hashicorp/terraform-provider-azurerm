// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2020-06-25/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePolicyVirtualMachineConfigurationAssignment() *pluginsdk.Resource {
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
			_, err := guestconfigurationassignments.ParseProviders2GuestConfigurationAssignmentID(id)
			return err
		}),

		Schema: resourcePolicyVirtualMachineConfigurationAssignmentSchema(),
	}
}

func resourcePolicyVirtualMachineConfigurationAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: computeValidate.VirtualMachineID,
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
							string(guestconfigurationassignments.AssignmentTypeAudit),
							string(guestconfigurationassignments.AssignmentTypeDeployAndAutoCorrect),
							string(guestconfigurationassignments.AssignmentTypeApplyAndAutoCorrect),
							string(guestconfigurationassignments.AssignmentTypeApplyAndMonitor),
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

func resourcePolicyVirtualMachineConfigurationAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := computeParse.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	id := guestconfigurationassignments.NewProviders2GuestConfigurationAssignmentID(subscriptionId, vmId.ResourceGroup, vmId.Name, d.Get("name").(string))

	if d.IsNewResource() {
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
	guestConfiguration := expandGuestConfigurationAssignment(d.Get("configuration").([]interface{}), id.GuestConfigurationAssignmentName)
	assignment := guestconfigurationassignments.GuestConfigurationAssignment{
		Name:     utils.String(id.GuestConfigurationAssignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
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

	d.SetId(id.ID())

	return resourcePolicyVirtualMachineConfigurationAssignmentRead(d, meta)
}

func resourcePolicyVirtualMachineConfigurationAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := guestconfigurationassignments.ParseProviders2GuestConfigurationAssignmentID(d.Id())
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

	vmId := computeParse.NewVirtualMachineID(subscriptionId, id.ResourceGroupName, id.VirtualMachineName)
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

func resourcePolicyVirtualMachineConfigurationAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := guestconfigurationassignments.ParseProviders2GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandGuestConfigurationAssignment(input []interface{}, name string) *guestconfigurationassignments.GuestConfigurationNavigation {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := guestconfigurationassignments.GuestConfigurationNavigation{
		Name:                   utils.String(name),
		Version:                utils.String(v["version"].(string)),
		ConfigurationParameter: expandGuestConfigurationAssignmentConfigurationParameters(v["parameter"].(*pluginsdk.Set).List()),
	}

	if v, ok := v["assignment_type"]; ok {
		result.AssignmentType = pointer.To(guestconfigurationassignments.AssignmentType(v.(string)))
	}

	if v, ok := v["content_hash"]; ok {
		result.ContentHash = utils.String(v.(string))
	}

	if v, ok := v["content_uri"]; ok {
		result.ContentUri = utils.String(v.(string))
	}

	return &result
}

func expandGuestConfigurationAssignmentConfigurationParameters(input []interface{}) *[]guestconfigurationassignments.ConfigurationParameter {
	results := make([]guestconfigurationassignments.ConfigurationParameter, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, guestconfigurationassignments.ConfigurationParameter{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		})
	}
	return &results
}

func flattenGuestConfigurationAssignment(input *guestconfigurationassignments.GuestConfigurationNavigation) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var version string
	if input.Version != nil {
		version = *input.Version
	}
	var assignmentType guestconfigurationassignments.AssignmentType
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
			"parameter":       flattenGuestConfigurationAssignmentConfigurationParameters(input.ConfigurationParameter),
			"version":         version,
		},
	}
}

func flattenGuestConfigurationAssignmentConfigurationParameters(input *[]guestconfigurationassignments.ConfigurationParameter) []interface{} {
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
