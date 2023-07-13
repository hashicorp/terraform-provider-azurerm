// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			_, err := backupinstances.ParseBackupInstanceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: backupinstances.ValidateBackupVaultID,
			},

			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: databases.ValidateDatabaseID,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: backuppolicies.ValidateBackupPolicyID,
			},

			"database_credential_key_vault_secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
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
	vaultId, _ := backupinstances.ParseBackupVaultID(d.Get("vault_id").(string))

	id := backupinstances.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_postgresql", id.ID())
		}
	}

	databaseId, _ := databases.ParseDatabaseID(d.Get("database_id").(string))
	location := location.Normalize(d.Get("location").(string))
	serverId := servers.NewServerID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.ServerName)
	policyId, _ := backuppolicies.ParseBackupPolicyID(d.Get("backup_policy_id").(string))

	parameters := backupinstances.BackupInstanceResource{
		Properties: &backupinstances.BackupInstance{
			DataSourceInfo: backupinstances.Datasource{
				DatasourceType:   utils.String("Microsoft.DBforPostgreSQL/servers/databases"),
				ObjectType:       utils.String("Datasource"),
				ResourceID:       databaseId.ID(),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(databaseId.DatabaseName),
				ResourceType:     utils.String("Microsoft.DBforPostgreSQL/servers/databases"),
				ResourceUri:      utils.String(""),
			},
			DataSourceSetInfo: &backupinstances.DatasourceSet{
				DatasourceType:   utils.String("Microsoft.DBForPostgreSQL/servers"),
				ObjectType:       utils.String("DatasourceSet"),
				ResourceID:       serverId.ID(),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(serverId.ServerName),
				ResourceType:     utils.String("Microsoft.DBForPostgreSQL/servers"),
				ResourceUri:      utils.String(""),
			},
			FriendlyName: utils.String(id.BackupInstanceName),
			PolicyInfo: backupinstances.PolicyInfo{
				PolicyId: policyId.ID(),
			},
		},
	}

	if v, ok := d.GetOk("database_credential_key_vault_secret_id"); ok {
		parameters.Properties.DatasourceAuthCredentials = backupinstances.SecretStoreBasedAuthCredentials{
			SecretStoreResource: &backupinstances.SecretStoreResource{
				Uri:             utils.String(v.(string)),
				SecretStoreType: backupinstances.SecretStoreTypeAzureKeyVault,
			},
		}
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(backupinstances.StatusConfiguringProtection), "UpdatingProtection"},
		Target:     []string{string(backupinstances.StatusProtectionConfigured)},
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

	id, err := backupinstances.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := backupinstances.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupInstanceName)
	d.Set("vault_id", vaultId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("database_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)

			if props.DatasourceAuthCredentials != nil {
				credential := props.DatasourceAuthCredentials.(backupinstances.SecretStoreBasedAuthCredentials)
				if credential.SecretStoreResource != nil {
					d.Set("database_credential_key_vault_secret_id", credential.SecretStoreResource.Uri)
				}
			} else {
				log.Printf("[DEBUG] Skipping setting database_credential_key_vault_secret_id since this DatasourceAuthCredentials is not supported")
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupInstancePostgreSQLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstances.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting DataProtection BackupInstance (%q): %+v", id, err)
	}

	return nil
}

func policyProtectionStateRefreshFunc(ctx context.Context, client *backupinstances.BackupInstancesClient, id backupinstances.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ProtectionStatus == nil || res.Model.Properties.ProtectionStatus.Status == nil {
			return nil, "", fmt.Errorf("reading DataProtection BackupInstance (%q) protection status: %+v", id, err)
		}

		return res, string(*res.Model.Properties.ProtectionStatus.Status), nil
	}
}
