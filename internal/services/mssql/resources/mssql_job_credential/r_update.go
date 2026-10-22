// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_credential

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlJobCredentialUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	jaId, err := jobcredentials.ParseJobAgentID(d.Get("job_agent_id").(string))
	if err != nil {
		return err
	}
	jobCredentialId := jobcredentials.NewCredentialID(jaId.SubscriptionId, jaId.ResourceGroupName, jaId.ServerName, jaId.JobAgentName, d.Get("name").(string))

	existing, err := client.Get(ctx, jobCredentialId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", jobCredentialId, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", jobCredentialId)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", jobCredentialId)
	}
	payload := existing.Model

	if d.HasChange("username") {
		payload.Properties.Username = d.Get("username").(string)
	}

	if d.HasChange("password") {
		payload.Properties.Password = d.Get("password").(string)
	}

	if d.HasChange("password_wo_version") {
		woPassword, err := pluginsdk.GetWriteOnly(d, "password_wo", cty.String)
		if err != nil {
			return err
		}

		if !woPassword.IsNull() {
			payload.Properties.Password = woPassword.AsString()
		}
	}

	if _, err := client.CreateOrUpdate(ctx, jobCredentialId, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", jobCredentialId, err)
	}

	return resourceMsSqlJobCredentialRead(d, meta)
}
