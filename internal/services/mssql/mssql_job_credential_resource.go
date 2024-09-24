// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobcredentials"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlJobCredential() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlJobCredentialCreateUpdate,
		Read:   resourceMsSqlJobCredentialRead,
		Update: resourceMsSqlJobCredentialCreateUpdate,
		Delete: resourceMsSqlJobCredentialDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.JobCredentialID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"job_agent_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.JobAgentID,
			},

			"username": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"password": {
				Type:      pluginsdk.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceMsSqlJobCredentialCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Job Credential creation.")

	jaId, err := jobcredentials.ParseJobAgentID(d.Get("job_agent_id").(string))
	if err != nil {
		return err
	}
	jobCredentialId := jobcredentials.NewCredentialID(jaId.SubscriptionId, jaId.ResourceGroupName, jaId.ServerName, jaId.JobAgentName, d.Get("name").(string))

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, jobCredentialId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing MsSql %s: %+v", jobCredentialId, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_job_credential", jobCredentialId.ID())
		}
	}

	jobCredential := jobcredentials.JobCredential{
		Name: utils.String(jobCredentialId.CredentialName),
		Properties: &jobcredentials.JobCredentialProperties{
			Username: username,
			Password: password,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, jobCredentialId, jobCredential); err != nil {
		return fmt.Errorf("creating MsSql %s: %+v", jobCredentialId, err)
	}

	d.SetId(jobCredentialId.ID())

	return resourceMsSqlJobCredentialRead(d, meta)
}

func resourceMsSqlJobCredentialRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	return nil
}

func resourceMsSqlJobCredentialDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobcredentials.ParseCredentialID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
