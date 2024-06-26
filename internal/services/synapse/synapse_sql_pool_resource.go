// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	mssqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	DefaultCreateMode            = "Default"
	RecoveryCreateMode           = "Recovery"
	PointInTimeRestoreCreateMode = "PointInTimeRestore"
)

func resourceSynapseSqlPool() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceSynapseSqlPoolCreate,
		Read:   resourceSynapseSqlPoolRead,
		Update: resourceSynapseSqlPoolUpdate,
		Delete: resourceSynapseSqlPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.SqlPoolID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			d.Set("create_mode", DefaultCreateMode)
			if v, ok := d.GetOk("create_mode"); ok && v.(string) != "" {
				d.Set("create_mode", v)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlPoolName,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DW100c",
					"DW200c",
					"DW300c",
					"DW400c",
					"DW500c",
					"DW1000c",
					"DW1500c",
					"DW2000c",
					"DW2500c",
					"DW3000c",
					"DW5000c",
					"DW6000c",
					"DW7500c",
					"DW10000c",
					"DW15000c",
					"DW30000c",
				}, false),
			},

			"create_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  DefaultCreateMode,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					DefaultCreateMode,
					RecoveryCreateMode,
					PointInTimeRestoreCreateMode,
				}, false),
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C The default of this is configurable by the user, so this should remain
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: mssqlValidate.DatabaseCollation(),
			},

			"recovery_database_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"restore"},
				ValidateFunc: validation.Any(
					validate.SqlPoolID,
					validate.SqlPoolRecoverableDatabaseID,
					commonids.ValidateSqlDatabaseID,
					mssqlValidate.RecoverableDatabaseID,
				),
			},

			"restore": {
				Type:          pluginsdk.TypeList,
				ForceNew:      true,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"recovery_database_id"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"point_in_time": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},

						"source_database_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.Any(
								validate.SqlPoolID,
								commonids.ValidateSqlDatabaseID,
							),
						},
					},
				},
			},

			"data_encrypted": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"geo_backup_policy_enabled": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(synapseSqlPoolCustomizeDiff),
	}

	if !features.FourPointOhBeta() {
		// NOTE: In v3.0 providers this will be an Optional field with a 'Default'
		// of 'GRS' to match existing v3.0 behavior, the 'ForceNew' logic will be
		// applied in the CustomizeDiff function...
		resource.Schema["storage_account_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Default:  string(synapse.StorageAccountTypeGRS),
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(synapse.StorageAccountTypeLRS),
				string(synapse.StorageAccountTypeGRS),
			}, false),
		}
	} else {
		resource.Schema["storage_account_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(synapse.StorageAccountTypeLRS),
				string(synapse.StorageAccountTypeGRS),
			}, false),
		}
	}

	return resource
}

func synapseSqlPoolCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
	_, value := d.GetChange("geo_backup_policy_enabled")
	geoBackupEnabled := value.(bool)

	_, value = d.GetChange("storage_account_type")
	storageAccountType := synapse.StorageAccountType(value.(string))

	if storageAccountType == synapse.StorageAccountTypeLRS && geoBackupEnabled {
		return fmt.Errorf("`geo_backup_policy_enabled` cannot be `true` if the `storage_account_type` is `LRS`")
	}

	return nil
}

func resourceSynapseSqlPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	sqlPTDEClient := meta.(*clients.Client).Synapse.SqlPoolTransparentDataEncryptionClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	geoBackUpClient := meta.(*clients.Client).Synapse.SqlPoolGeoBackupPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	existing, err := sqlClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_synapse_sql_pool", id.ID())
	}

	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", workspaceId, err)
	}

	geoBackupEnabled := d.Get("geo_backup_policy_enabled").(bool)
	storageAccountType := synapse.StorageAccountType(d.Get("storage_account_type").(string))

	mode := d.Get("create_mode").(string)
	sqlPoolInfo := synapse.SQLPool{
		Location: workspace.Location,
		SQLPoolResourceProperties: &synapse.SQLPoolResourceProperties{
			CreateMode:         synapse.CreateMode(*utils.String(mode)),
			StorageAccountType: storageAccountType,
		},
		Sku: &synapse.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	switch mode {
	case DefaultCreateMode:
		sqlPoolInfo.SQLPoolResourceProperties.Collation = utils.String(d.Get("collation").(string))
	case RecoveryCreateMode:
		recoveryDatabaseId := constructSourceDatabaseId(d.Get("recovery_database_id").(string))

		if recoveryDatabaseId == "" {
			return fmt.Errorf("`recovery_database_id` must be set when `create_mode` is %q", RecoveryCreateMode)
		}

		sqlPoolInfo.SQLPoolResourceProperties.RecoverableDatabaseID = utils.String(recoveryDatabaseId)
	case PointInTimeRestoreCreateMode:
		restore := d.Get("restore").([]interface{})
		if len(restore) == 0 || restore[0] == nil {
			return fmt.Errorf("`restore` block must be set when `create_mode` is %q", PointInTimeRestoreCreateMode)
		}

		v := restore[0].(map[string]interface{})
		sourceDatabaseId := constructSourceDatabaseId(v["source_database_id"].(string))
		vTime, parseErr := date.ParseTime(time.RFC3339, v["point_in_time"].(string))

		if parseErr != nil {
			return fmt.Errorf("parsing time format: %+v", parseErr)
		}

		sqlPoolInfo.SQLPoolResourceProperties.RestorePointInTime = &date.Time{Time: vTime}
		sqlPoolInfo.SQLPoolResourceProperties.SourceDatabaseID = utils.String(sourceDatabaseId)
	}

	future, err := sqlClient.Create(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, sqlPoolInfo)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, sqlClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	if d.Get("data_encrypted").(bool) {
		parameter := synapse.TransparentDataEncryption{
			TransparentDataEncryptionProperties: &synapse.TransparentDataEncryptionProperties{
				Status: synapse.TransparentDataEncryptionStatusEnabled,
			},
		}

		if _, err := sqlPTDEClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameter); err != nil {
			return fmt.Errorf("setting `data_encrypted`: %+v", err)
		}
	}

	// Only update the Geo Backup Policy if it has been disabled since it is
	// already enabled by default...
	if !geoBackupEnabled {
		geoBackupParams := synapse.GeoBackupPolicy{
			GeoBackupPolicyProperties: &synapse.GeoBackupPolicyProperties{
				State: synapse.GeoBackupPolicyStateDisabled,
			},
		}

		if _, err := geoBackUpClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, geoBackupParams); err != nil {
			return fmt.Errorf("setting `geo_backup_policy_enabled`: %+v", err)
		}
	}

	d.SetId(id.ID())
	return resourceSynapseSqlPoolRead(d, meta)
}

func resourceSynapseSqlPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	geoBackUpClient := meta.(*clients.Client).Synapse.SqlPoolGeoBackupPoliciesClient
	sqlPTDEClient := meta.(*clients.Client).Synapse.SqlPoolTransparentDataEncryptionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("data_encrypted") {
		status := synapse.TransparentDataEncryptionStatusDisabled
		if d.Get("data_encrypted").(bool) {
			status = synapse.TransparentDataEncryptionStatusEnabled
		}

		parameter := synapse.TransparentDataEncryption{
			TransparentDataEncryptionProperties: &synapse.TransparentDataEncryptionProperties{
				Status: status,
			},
		}

		if _, err := sqlPTDEClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameter); err != nil {
			return fmt.Errorf("updating `data_encrypted`: %+v", err)
		}
	}

	if d.HasChange("geo_backup_policy_enabled") {
		state := synapse.GeoBackupPolicyStateEnabled
		if !d.Get("geo_backup_policy_enabled").(bool) {
			state = synapse.GeoBackupPolicyStateDisabled
		}

		geoBackupParams := synapse.GeoBackupPolicy{
			GeoBackupPolicyProperties: &synapse.GeoBackupPolicyProperties{
				State: state,
			},
		}

		if _, err := geoBackUpClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, geoBackupParams); err != nil {
			return fmt.Errorf("updating `geo_backup_policy_enabled`: %+v", err)
		}
	}

	if d.HasChanges("sku_name", "tags") {
		sqlPoolInfo := synapse.SQLPoolPatchInfo{
			Sku: &synapse.Sku{
				Name: utils.String(d.Get("sku_name").(string)),
			},
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		if _, err := sqlClient.Update(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, sqlPoolInfo); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}

		// wait for sku scale completion
		if d.HasChange("sku_name") {
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					"Scaling",
				},
				Target: []string{
					"Online",
				},
				Refresh:                   synapseSqlPoolScaleStateRefreshFunc(ctx, sqlClient, id.ResourceGroup, id.WorkspaceName, id.Name),
				MinTimeout:                5 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for scaling of %s: %+v", *id, err)
			}
		}
	}

	return resourceSynapseSqlPoolRead(d, meta)
}

func resourceSynapseSqlPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	sqlPTDEClient := meta.(*clients.Client).Synapse.SqlPoolTransparentDataEncryptionClient
	geoBackUpClient := meta.(*clients.Client).Synapse.SqlPoolGeoBackupPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := sqlClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	transparentDataEncryption, err := sqlPTDEClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Transparent Data Encryption settings of %s: %+v", *id, err)
	}

	geoBackupPolicy, err := geoBackUpClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Geo Backup Policy of %s: %+v", *id, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID()
	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", workspaceId)
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	if props := resp.SQLPoolResourceProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("storage_account_type", props.StorageAccountType)
	}

	geoBackupEnabled := true
	if geoBackupProps := geoBackupPolicy.GeoBackupPolicyProperties; geoBackupProps != nil {
		geoBackupEnabled = geoBackupProps.State == synapse.GeoBackupPolicyStateEnabled
	}
	d.Set("geo_backup_policy_enabled", geoBackupEnabled)

	if tdeProps := transparentDataEncryption.TransparentDataEncryptionProperties; tdeProps != nil {
		d.Set("data_encrypted", tdeProps.Status == synapse.TransparentDataEncryptionStatusEnabled)
	}

	// whole "restore" block is not returned. to avoid conflict, so set it from the old state
	d.Set("restore", d.Get("restore").([]interface{}))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSynapseSqlPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := sqlClient.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, sqlClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func synapseSqlPoolScaleStateRefreshFunc(ctx context.Context, client *synapse.SQLPoolsClient, resourceGroup, workspaceName, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
		if err != nil {
			return resp, "failed", err
		}

		if resp.SQLPoolResourceProperties == nil || resp.SQLPoolResourceProperties.Status == nil {
			return resp, "failed", nil
		}

		return resp, *resp.SQLPoolResourceProperties.Status, nil
	}
}

// sqlPool backend service is a proxy to sql database
// backend service restore and backup only accept id format of sql database
// so if the id is sqlPool, we need to construct the corresponding sql database id
func constructSourceDatabaseId(id string) string {
	sqlPoolId, err := parse.SqlPoolID(id)
	if err != nil {
		return id
	}

	return commonids.NewSqlDatabaseID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name).ID()
}
