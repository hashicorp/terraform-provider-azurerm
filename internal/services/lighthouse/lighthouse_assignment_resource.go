// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLighthouseAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLighthouseAssignmentCreate,
		Read:   resourceLighthouseAssignmentRead,
		Delete: resourceLighthouseAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := registrationassignments.ParseScopedRegistrationAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"lighthouse_definition_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: registrationdefinitions.ValidateScopedRegistrationDefinitionID,
			},

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.Any(commonids.ValidateSubscriptionID, commonids.ValidateResourceGroupID),
			},
		},
	}
}

func resourceLighthouseAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseAssignmentName := d.Get("name").(string)
	if lighthouseAssignmentName == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Lighthouse Assignment: %+v", err)
		}

		lighthouseAssignmentName = uuid
	}

	id := registrationassignments.NewScopedRegistrationAssignmentID(d.Get("scope").(string), lighthouseAssignmentName)
	options := registrationassignments.GetOperationOptions{
		ExpandRegistrationDefinition: utils.Bool(false),
	}
	existing, err := client.Get(ctx, id, options)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_lighthouse_assignment", id.ID())
	}

	parameters := registrationassignments.RegistrationAssignment{
		Properties: &registrationassignments.RegistrationAssignmentProperties{
			RegistrationDefinitionId: d.Get("lighthouse_definition_id").(string),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLighthouseAssignmentRead(d, meta)
}

func resourceLighthouseAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registrationassignments.ParseScopedRegistrationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := registrationassignments.GetOperationOptions{
		ExpandRegistrationDefinition: utils.Bool(false),
	}
	resp, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RegistrationAssignmentId)
	d.Set("scope", id.Scope)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("lighthouse_definition_id", props.RegistrationDefinitionId)
		}
	}

	return nil
}

func resourceLighthouseAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registrationassignments.ParseScopedRegistrationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    lighthouseAssignmentDeleteRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}

func lighthouseAssignmentDeleteRefreshFunc(ctx context.Context, client *registrationassignments.RegistrationAssignmentsClient, id registrationassignments.ScopedRegistrationAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		options := registrationassignments.GetOperationOptions{
			ExpandRegistrationDefinition: utils.Bool(true),
		}
		resp, err := client.Get(ctx, id, options)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling to check the deletion of %s: %+v", id, err)
		}

		return resp, "Deleting", nil
	}
}
