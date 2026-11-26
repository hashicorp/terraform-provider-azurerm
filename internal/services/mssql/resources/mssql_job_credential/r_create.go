// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_credential

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlJobCredentialCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	jaId, err := jobcredentials.ParseJobAgentID(d.Get("job_agent_id").(string))
	if err != nil {
		return err
	}
	jobCredentialId := jobcredentials.NewCredentialID(jaId.SubscriptionId, jaId.ResourceGroupName, jaId.ServerName, jaId.JobAgentName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, jobCredentialId)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", jobCredentialId, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_job_credential", jobCredentialId.ID())
		}
	}

	woPassword, err := pluginsdk.GetWriteOnly(d, "password_wo", cty.String)
	if err != nil {
		return err
	}

	password := d.Get("password").(string)
	if !woPassword.IsNull() {
		password = woPassword.AsString()
	}

	jobCredential := jobcredentials.JobCredential{
		Name: pointer.To(jobCredentialId.CredentialName),
		Properties: &jobcredentials.JobCredentialProperties{
			Username: d.Get("username").(string),
			Password: password,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, jobCredentialId, jobCredential); err != nil {
		return fmt.Errorf("creating %s: %+v", jobCredentialId, err)
	}

	d.SetId(jobCredentialId.ID())

	return resourceMsSqlJobCredentialRead(d, meta)
}
