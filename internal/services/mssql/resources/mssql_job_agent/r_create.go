// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobagents"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlJobAgentCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseId := d.Get("database_id").(string)
	dbId, err := commonids.ParseSqlDatabaseID(databaseId)
	if err != nil {
		return err
	}
	id := jobagents.NewJobAgentID(dbId.SubscriptionId, dbId.ResourceGroupName, dbId.ServerName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_job_agent", id.ID())
		}
	}

	params := jobagents.JobAgent{
		Name:     &id.JobAgentName,
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &jobagents.JobAgentProperties{
			DatabaseId: databaseId,
		},
		Sku: &jobagents.Sku{
			Name: d.Get("sku").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	expandedIdentity, err := expandJobAgentIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	params.Identity = expandedIdentity

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, params, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMsSqlJobAgentRead(d, meta)
}

func expandJobAgentIdentity(input []any) (*jobagents.JobAgentIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}
	if expanded == nil || expanded.Type == identity.TypeNone {
		return nil, nil
	}

	result := &jobagents.JobAgentIdentity{
		Type: jobagents.JobAgentIdentityType(string(expanded.Type)),
	}

	if len(expanded.IdentityIds) > 0 {
		uai := make(map[string]jobagents.JobAgentUserAssignedIdentity, len(expanded.IdentityIds))
		for id := range expanded.IdentityIds {
			uai[id] = jobagents.JobAgentUserAssignedIdentity{}
		}
		result.UserAssignedIdentities = &uai
	}

	return result, nil
}
