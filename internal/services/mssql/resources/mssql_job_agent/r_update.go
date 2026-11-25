// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlJobAgentUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseId := d.Get("database_id").(string)
	dbId, err := commonids.ParseSqlDatabaseID(databaseId)
	if err != nil {
		return err
	}
	id := jobagents.NewJobAgentID(dbId.SubscriptionId, dbId.ResourceGroupName, dbId.ServerName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", id)
	}
	params := existing.Model

	if d.HasChanges("identity") {
		expandedIdentity, err := expandJobAgentIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		params.Identity = expandedIdentity
	}

	if d.HasChanges("sku") {
		params.Sku = &jobagents.Sku{
			Name: d.Get("sku").(string),
		}
	}

	if d.HasChanges("tags") {
		params.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, *params); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMsSqlJobAgentRead(d, meta)
}
