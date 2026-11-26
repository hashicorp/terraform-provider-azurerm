// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_definition

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLighthouseDefinitionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registrationdefinitions.ParseScopedRegistrationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("lighthouse_definition_id", id.RegistrationDefinitionId)
	d.Set("scope", id.Scope)

	if model := resp.Model; model != nil {
		if err := d.Set("plan", flattenLighthouseDefinitionPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("authorization", flattenLighthouseDefinitionAuthorization(props.Authorizations)); err != nil {
				return fmt.Errorf("setting `authorization`: %+v", err)
			}
			if err := d.Set("eligible_authorization", flattenLighthouseDefinitionEligibleAuthorization(props.EligibleAuthorizations)); err != nil {
				return fmt.Errorf("setting `eligible_authorization`: %+v", err)
			}
			d.Set("description", props.Description)
			d.Set("name", props.RegistrationDefinitionName)
			d.Set("managing_tenant_id", props.ManagedByTenantId)
		}
	}

	return nil
}
