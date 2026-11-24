// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLighthouseDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseDefinitionID := d.Get("lighthouse_definition_id").(string)
	if lighthouseDefinitionID == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Lighthouse Definition: %+v", err)
		}

		lighthouseDefinitionID = uuid
	}

	id := registrationdefinitions.NewScopedRegistrationDefinitionID(d.Get("scope").(string), lighthouseDefinitionID)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_lighthouse_definition", id.ID())
			}
		}
	}
	authorizations := expandLighthouseDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List())
	parameters := registrationdefinitions.RegistrationDefinition{
		Plan: expandLighthouseDefinitionPlan(d.Get("plan").([]any)),
		Properties: &registrationdefinitions.RegistrationDefinitionProperties{
			Description:                pointer.To(d.Get("description").(string)),
			Authorizations:             authorizations,
			RegistrationDefinitionName: pointer.To(d.Get("name").(string)),
			ManagedByTenantId:          d.Get("managing_tenant_id").(string),
		},
	}

	if v, ok := d.GetOk("eligible_authorization"); ok {
		parameters.Properties.EligibleAuthorizations = expandLighthouseDefinitionEligibleAuthorization(v.(*pluginsdk.Set).List())
	}

	// NOTE: this API call uses DefinitionId then Scope - check in the future
	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceLighthouseDefinitionRead(d, meta)
}
