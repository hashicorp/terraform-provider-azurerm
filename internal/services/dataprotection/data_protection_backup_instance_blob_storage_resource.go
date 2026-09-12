// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-03-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_blob_storage -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

func resourceDataProtectionBackupInstanceBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Read:   resourceDataProtectionBackupInstanceBlobStorageRead,
		Update: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Delete: resourceDataProtectionBackupInstanceBlobStorageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&backupinstanceresources.BackupInstanceId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&backupinstanceresources.BackupInstanceId{}),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: backupvaultresources.ValidateBackupVaultID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: basebackuppolicyresources.ValidateBackupPolicyID,
			},

			"storage_account_container_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				ConflictsWith: []string{"excluded_container_name_prefixes"},
			},

			"auto_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"excluded_container_name_prefixes": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				RequiredWith:  []string{"auto_protection_enabled"},
				ConflictsWith: []string{"storage_account_container_names"},
			},

			"protection_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(resourceDataProtectionBackupInstanceBlobStorageCustomizeDiff),
		),
	}
}

func resourceDataProtectionBackupInstanceBlobStorageCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	if !d.NewValueKnown("auto_protection_enabled") || !d.NewValueKnown("storage_account_container_names") || !d.NewValueKnown("excluded_container_name_prefixes") {
		return nil
	}

	oldAutoProtection, newAutoProtection := d.GetChange("auto_protection_enabled")
	autoProtectionEnabled := newAutoProtection.(bool)
	oldContainerNames, newContainerNames := d.GetChange("storage_account_container_names")

	if autoProtectionEnabled && len(newContainerNames.([]interface{})) > 0 {
		return fmt.Errorf("`storage_account_container_names` cannot be set when `auto_protection_enabled` is `true`: auto protection covers all present and future containers of the Storage Account")
	}

	if !autoProtectionEnabled && len(d.Get("excluded_container_name_prefixes").([]interface{})) > 0 {
		return fmt.Errorf("`excluded_container_name_prefixes` can only be set when `auto_protection_enabled` is `true`")
	}

	// Azure does not allow switching a Backup Instance back from auto protection to a container list (or to no
	// container selection at all) once auto protection has been enabled - the only way out is to re-create the
	// Backup Instance, which also removes its vaulted recovery points. This is deliberately an error rather than
	// `ForceNew` so that the re-creation has to be requested explicitly.
	if d.Id() != "" && oldAutoProtection.(bool) && !autoProtectionEnabled {
		return fmt.Errorf("`auto_protection_enabled` cannot be changed from `true` to `false`: enabling auto protection for a Backup Instance is irreversible in Azure. To switch back to a container list the Backup Instance has to be re-created (e.g. via `terraform apply -replace=<resource address>`), which deletes its vaulted recovery points")
	}

	// The `storage_account_container_names` can not be removed once specified - unless the Backup Instance is being
	// migrated to auto protection, which Azure supports as an in-place update.
	if len(oldContainerNames.([]interface{})) > 0 && len(newContainerNames.([]interface{})) == 0 && !autoProtectionEnabled {
		if err := d.ForceNew("storage_account_container_names"); err != nil {
			return err
		}
	}

	return nil
}

func resourceDataProtectionBackupInstanceBlobStorageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient_v2026_03_01
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := backupvaultresources.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backupinstanceresources.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.BackupInstancesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_blob_storage", id.ID())
			}
		}
	}

	storageAccountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}
	location := location.Normalize(d.Get("location").(string))
	policyId, err := basebackuppolicyresources.ParseBackupPolicyID(d.Get("backup_policy_id").(string))
	if err != nil {
		return err
	}

	parameters := backupinstanceresources.BackupInstanceResource{
		Properties: &backupinstanceresources.BackupInstance{
			DataSourceInfo: backupinstanceresources.Datasource{
				DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts/blobServices"),
				ObjectType:       pointer.To("Datasource"),
				ResourceID:       storageAccountId.ID(),
				ResourceLocation: pointer.To(location),
				ResourceName:     pointer.To(storageAccountId.StorageAccountName),
				ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
				ResourceUri:      pointer.To(storageAccountId.ID()),
			},
			FriendlyName: pointer.To(id.BackupInstanceName),
			PolicyInfo: backupinstanceresources.PolicyInfo{
				PolicyId: policyId.ID(),
			},
		},
	}

	if d.Get("auto_protection_enabled").(bool) {
		parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{
			BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
				expandBlobBackupAutoProtection(d.Get("excluded_container_name_prefixes").([]interface{})),
			},
		}
	} else if v, ok := d.GetOk("storage_account_container_names"); ok {
		parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{
			BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
				backupinstanceresources.BlobBackupDatasourceParameters{
					ContainersList: pointer.From(helpers.ExpandStringSlice(v.([]interface{}))),
				},
			},
		}
	}

	if d.IsNewResource() {
		if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating DataProtection BackupInstance (%q): %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating DataProtection BackupInstance (%q): %+v", id, err)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(backupinstanceresources.StatusConfiguringProtection), "UpdatingProtection"},
		Target:     []string{string(backupinstanceresources.StatusProtectionConfigured)},
		Refresh:    blobStorageBackupInstanceProtectionStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for BackupInstance(%q) policy protection to be completed: %+v", id, err)
	}

	return resourceDataProtectionBackupInstanceBlobStorageRead(d, meta)
}

func resourceDataProtectionBackupInstanceBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient_v2026_03_01
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupInstanceName)
	d.Set("vault_id", vaultId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("storage_account_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)
			d.Set("protection_state", pointer.FromEnum(props.CurrentProtectionState))

			containerNames := make([]string, 0)
			autoProtectionEnabled := false
			excludedContainerNamePrefixes := make([]string, 0)
			if policyParas := props.PolicyInfo.PolicyParameters; policyParas != nil {
				if dataStoreParas := policyParas.BackupDatasourceParametersList; dataStoreParas != nil {
					if dsp := pointer.From(dataStoreParas); len(dsp) > 0 {
						switch parameter := dsp[0].(type) {
						case backupinstanceresources.BlobBackupDatasourceParameters:
							containerNames = parameter.ContainersList
						case backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection:
							autoProtectionEnabled, excludedContainerNamePrefixes = flattenBlobBackupAutoProtection(parameter)
						}
					}
				}
			}
			if err := d.Set("storage_account_container_names", containerNames); err != nil {
				return fmt.Errorf("setting `storage_account_container_names`: %+v", err)
			}
			if err := d.Set("auto_protection_enabled", autoProtectionEnabled); err != nil {
				return fmt.Errorf("setting `auto_protection_enabled`: %+v", err)
			}
			if err := d.Set("excluded_container_name_prefixes", excludedContainerNamePrefixes); err != nil {
				return fmt.Errorf("setting `excluded_container_name_prefixes`: %+v", err)
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func expandBlobBackupAutoProtection(excludedContainerNamePrefixes []interface{}) backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection {
	settings := backupinstanceresources.BlobBackupRuleBasedAutoProtectionSettings{
		Enabled: true,
	}

	// rules are evaluated in the order provided; without any rules every present and future container is eligible
	if len(excludedContainerNamePrefixes) > 0 {
		rules := make([]backupinstanceresources.BlobBackupAutoProtectionRule, 0, len(excludedContainerNamePrefixes))
		for _, prefix := range excludedContainerNamePrefixes {
			rules = append(rules, backupinstanceresources.BlobBackupAutoProtectionRule{
				// `objectType` is a required plain field on the rule rather than a discriminator, so the SDK does not set it
				ObjectType: "BlobBackupAutoProtectionRule",
				Mode:       backupinstanceresources.BlobBackupRuleModeExclude,
				Type:       backupinstanceresources.BlobBackupPatternTypePrefix,
				Pattern:    prefix.(string),
			})
		}
		settings.Rules = &rules
	}

	return backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection{
		AutoProtectionSettings: settings,
	}
}

func flattenBlobBackupAutoProtection(input backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection) (enabled bool, excludedContainerNamePrefixes []string) {
	excludedContainerNamePrefixes = make([]string, 0)
	settings := input.AutoProtectionSettings

	for _, rule := range pointer.From(settings.Rules) {
		// `Exclude` + `Prefix` is the only combination the API supports today, anything else is unknown to this resource
		if rule.Mode == backupinstanceresources.BlobBackupRuleModeExclude && rule.Type == backupinstanceresources.BlobBackupPatternTypePrefix {
			excludedContainerNamePrefixes = append(excludedContainerNamePrefixes, rule.Pattern)
		}
	}

	return settings.Enabled, excludedContainerNamePrefixes
}

func resourceDataProtectionBackupInstanceBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient_v2026_03_01
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func blobStorageBackupInstanceProtectionStateRefreshFunc(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.BackupInstancesGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ProtectionStatus == nil || res.Model.Properties.ProtectionStatus.Status == nil {
			return nil, "", fmt.Errorf("reading DataProtection BackupInstance (%q) protection status: %+v", id, err)
		}

		return res, string(*res.Model.Properties.ProtectionStatus.Status), nil
	}
}
