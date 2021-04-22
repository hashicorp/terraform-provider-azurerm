package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlJobCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlJobCredentialCreateUpdate,
		Read:   resourceMsSqlJobCredentialRead,
		Update: resourceMsSqlJobCredentialCreateUpdate,
		Delete: resourceMsSqlJobCredentialDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.JobCredentialID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"job_agent_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.JobAgentID,
			},

			"username": {
				Type:     schema.TypeString,
				Required: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceMsSqlJobCredentialCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Job Credential creation.")

	jaId, err := parse.JobAgentID(d.Get("job_agent_id").(string))
	if err != nil {
		return err
	}
	jobCredentialId := parse.NewJobCredentialID(jaId.SubscriptionId, jaId.ResourceGroup, jaId.ServerName, jaId.Name, d.Get("name").(string))

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, jobCredentialId.ResourceGroup, jobCredentialId.ServerName, jobCredentialId.JobAgentName, jobCredentialId.CredentialName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing MsSql %s: %+v", jobCredentialId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mssql_job_credential", jobCredentialId.ID())
		}
	}

	jobCredential := sql.JobCredential{
		Name: utils.String(jobCredentialId.CredentialName),
		JobCredentialProperties: &sql.JobCredentialProperties{
			Username: utils.String(username),
			Password: utils.String(password),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, jobCredentialId.ResourceGroup, jobCredentialId.ServerName, jobCredentialId.JobAgentName, jobCredentialId.CredentialName, jobCredential); err != nil {
		return fmt.Errorf("creating MsSql %s: %+v", jobCredentialId, err)
	}

	d.SetId(jobCredentialId.ID())

	return resourceMsSqlJobCredentialRead(d, meta)
}

func resourceMsSqlJobCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobCredentialID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.JobAgentName, id.CredentialName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Job Credential %s (Job Agent %q / MsSql Server %q / Resource Group %q): %s", id.CredentialName, id.JobAgentName, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("username", resp.Username)

	jobAgentId := parse.NewJobAgentID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.JobAgentName)
	d.Set("job_agent_id", jobAgentId.ID())

	return nil
}

func resourceMsSqlJobCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobCredentialsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobCredentialID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.ServerName, id.JobAgentName, id.CredentialName)
	if err != nil {
		return fmt.Errorf("deleting Job Credential %s: %+v", id.CredentialName, err)
	}

	return nil
}
