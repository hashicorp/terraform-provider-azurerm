// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationassignments"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLighthouseAssignmentCreate(d *pluginsdk.ResourceData, meta any) error {
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
		ExpandRegistrationDefinition: pointer.To(false),
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, options)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_lighthouse_assignment", id.ID())
		}
	}

	parameters := registrationassignments.RegistrationAssignment{
		Properties: &registrationassignments.RegistrationAssignmentProperties{
			RegistrationDefinitionId: d.Get("lighthouse_definition_id").(string),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLighthouseAssignmentRead(d, meta)
}
