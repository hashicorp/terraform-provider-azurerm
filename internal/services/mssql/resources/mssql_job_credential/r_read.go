// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_credential

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlJobCredentialRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobcredentials.ParseCredentialID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %s", *id, err)
	}

	d.Set("name", id.CredentialName)
	jobAgentId := jobcredentials.NewJobAgentID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName)
	d.Set("job_agent_id", jobAgentId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("username", props.Username)
		}
	}

	d.Set("password_wo_version", d.Get("password_wo_version").(int))

	return nil
}
