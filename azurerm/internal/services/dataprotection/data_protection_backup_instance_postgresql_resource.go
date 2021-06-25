package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/validate"
	postgresParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	postgresValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataProtectionBackupInstancePostgreSQL() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionBackupInstancePostgreSQLCreateUpdate,
		Read:   resourceDataProtectionBackupInstancePostgreSQLRead,
		Update: resourceDataProtectionBackupInstancePostgreSQLCreateUpdate,
		Delete: resourceDataProtectionBackupInstancePostgreSQLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BackupInstanceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": location.Schema(),

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BackupVaultID,
			},

			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: postgresValidate.DatabaseID,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.BackupPolicyID,
			},
		},
	}
}
func resourceDataProtectionBackupInstancePostgreSQLCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := parse.BackupVaultID(d.Get("vault_id").(string))

	id := parse.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroup, vaultId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_postgresql", id.ID())
		}
	}

	databaseId, _ := postgresParse.DatabaseID(d.Get("database_id").(string))
	location := location.Normalize(d.Get("location").(string))
	serverId := postgresParse.NewServerID(databaseId.SubscriptionId, databaseId.ResourceGroup, databaseId.ServerName)
	policyId, _ := parse.BackupPolicyID(d.Get("backup_policy_id").(string))

	parameters := dataprotection.BackupInstanceResource{
		Properties: &dataprotection.BackupInstance{
			DataSourceInfo: &dataprotection.Datasource{
				DatasourceType:   utils.String("Microsoft.DBforPostgreSQL/servers/databases"),
				ObjectType:       utils.String("Datasource"),
				ResourceID:       utils.String(databaseId.ID()),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(databaseId.Name),
				ResourceType:     utils.String("Microsoft.DBforPostgreSQL/servers/databases"),
				ResourceURI:      utils.String(""),
			},
			DataSourceSetInfo: &dataprotection.DatasourceSet{
				DatasourceType:   utils.String("Microsoft.DBForPostgreSQL/servers"),
				ObjectType:       utils.String("DatasourceSet"),
				ResourceID:       utils.String(serverId.ID()),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(serverId.Name),
				ResourceType:     utils.String("Microsoft.DBForPostgreSQL/servers"),
				ResourceURI:      utils.String(""),
			},
			FriendlyName: utils.String(id.Name),
			PolicyInfo: &dataprotection.PolicyInfo{
				PolicyID: utils.String(policyId.ID()),
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.BackupVaultName, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupInstance (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the DataProtection BackupInstance (%q): %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(dataprotection.ConfiguringProtection), string(dataprotection.UpdatingProtection)},
		Target:     []string{string(dataprotection.ProtectionConfigured)},
		Refresh:    policyProtectionStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for BackupInstance(%q) policy protection to be completed: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupInstancePostgreSQLRead(d, meta)
}

func resourceDataProtectionBackupInstancePostgreSQLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := parse.NewBackupVaultID(id.SubscriptionId, id.ResourceGroup, id.BackupVaultName)
	d.Set("name", id.Name)
	d.Set("vault_id", vaultId.ID())
	if props := resp.Properties; props != nil {
		if props.DataSourceInfo != nil {
			d.Set("database_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
		}
		if props.PolicyInfo != nil {
			d.Set("backup_policy_id", props.PolicyInfo.PolicyID)
		}
	}
	return nil
}

func resourceDataProtectionBackupInstancePostgreSQLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting DataProtection BackupInstance (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the DataProtection BackupInstance (%q): %+v", id.Name, err)
	}
	return nil
}

func policyProtectionStateRefreshFunc(ctx context.Context, client *dataprotection.BackupInstancesClient, id parse.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
		}
		if res.Properties == nil || res.Properties.ProtectionStatus == nil {
			return nil, "", fmt.Errorf("error reading DataProtection BackupInstance (%q) protection status: %+v", id, err)
		}

		return res, string(res.Properties.ProtectionStatus.Status), nil
	}
}
