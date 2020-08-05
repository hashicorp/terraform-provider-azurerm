package postgres

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"time"
)

func resourceArmPostgreSQLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerKeyCreateUpdate,
		Read:   resourceArmPostgreSQLServerKeyRead,
		Update: resourceArmPostgreSQLServerKeyCreateUpdate,
		Delete: resourceArmPostgreSQLServerKeyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServerKeyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PostgresServerServerID,
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildId,
			},
		},
	}
}

func resourceArmPostgreSQLServerKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	serverID, err := parse.PostgresServerServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverID.ResourceGroup, serverID.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presense of existing PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", name, serverID.ResourceGroup, serverID.Name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_server_key", *existing.ID)
		}
	}

	param := postgresql.ServerKey{
		ServerKeyProperties: &postgresql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           utils.String(d.Get("key_vault_key_id").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, serverID.Name, name, param, serverID.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", name, serverID.ResourceGroup, serverID.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", name, serverID.ResourceGroup, serverID.Name, err)
	}

	resp, err := client.Get(ctx, serverID.ResourceGroup, serverID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", name, serverID.ResourceGroup, serverID.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("nil or empty ID returned for PostgreSQL ServerKey %q (Resource Group %q / Server %q): %+v", name, serverID.ResourceGroup, serverID.Name, err)
	}

	d.SetId(*resp.ID)
	return resourceArmPostgreSQLServerKeyRead(d, meta)
}

func resourceArmPostgreSQLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*clients.Client).Postgres.ServerKeysClient
	//ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	//defer cancel()
	//
	//id, err := parse.ServerKeyID(d.Id())

	return nil
}

func resourceArmPostgreSQLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*clients.Client).Postgres.ServerKeysClient
	//ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	//defer cancel()

	return nil
}
