// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMssqlJobAgentSetFlatten(d *pluginsdk.ResourceData, id *jobagents.JobAgentId, model *jobagents.JobAgent) error {
	d.Set("name", id.JobAgentName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("database_id", props.DatabaseId)
		}

		flattenedIdentity, err := flattenJobAgentIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		d.Set("identity", flattenedIdentity)

		if sku := model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenJobAgentIdentity(input *jobagents.JobAgentIdentity) ([]any, error) {
	if input == nil {
		result, err := identity.FlattenUserAssignedMap(&identity.UserAssignedMap{Type: identity.TypeNone})
		if err != nil {
			return nil, err
		}
		return *result, nil
	}

	identityIds := make(map[string]identity.UserAssignedIdentityDetails)
	if input.UserAssignedIdentities != nil {
		for id := range *input.UserAssignedIdentities {
			identityIds[id] = identity.UserAssignedIdentityDetails{}
		}
	}

	result, err := identity.FlattenUserAssignedMap(&identity.UserAssignedMap{
		Type:        identity.Type(string(input.Type)),
		IdentityIds: identityIds,
	})
	if err != nil {
		return nil, err
	}
	return *result, nil
}
