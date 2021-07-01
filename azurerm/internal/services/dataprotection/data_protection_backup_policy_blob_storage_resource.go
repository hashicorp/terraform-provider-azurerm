package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	helperValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataProtectionBackupPolicyBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupPolicyBlobStorageCreate,
		Read:   resourceDataProtectionBackupPolicyBlobStorageRead,
		Delete: resourceDataProtectionBackupPolicyBlobStorageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BackupPolicyID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
					"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BackupVaultID,
			},

			"retention_duration": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helperValidate.ISO8601Duration,
			},
		},
	}
}
func resourceDataProtectionBackupPolicyBlobStorageCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := parse.BackupVaultID(d.Get("vault_id").(string))
	id := parse.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroup, vaultId.Name, name)

	existing, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_blob_storage", id.ID())
	}

	parameters := dataprotection.BaseBackupPolicyResource{
		Properties: &dataprotection.BackupPolicy{
			PolicyRules: &[]dataprotection.BasicBasePolicyRule{
				dataprotection.AzureRetentionRule{
					Name:       utils.String("Default"),
					ObjectType: dataprotection.ObjectTypeAzureRetentionRule,
					IsDefault:  utils.Bool(true),
					Lifecycles: &[]dataprotection.SourceLifeCycle{
						{
							DeleteAfter: dataprotection.AbsoluteDeleteOption{
								Duration:   utils.String(d.Get("retention_duration").(string)),
								ObjectType: dataprotection.ObjectTypeAbsoluteDeleteOption,
							},
							SourceDataStore: &dataprotection.DataStoreInfoBase{
								DataStoreType: "OperationalStore",
								ObjectType:    utils.String("DataStoreInfoBase"),
							},
							TargetDataStoreCopySettings: &[]dataprotection.TargetCopySetting{},
						},
					},
				},
			},
			DatasourceTypes: &[]string{"Microsoft.Storage/storageAccounts/blobServices"},
			ObjectType:      dataprotection.ObjectTypeBackupPolicy,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.BackupVaultName, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyBlobStorageRead(d, meta)
}

func resourceDataProtectionBackupPolicyBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
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
		return fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	vaultId := parse.NewBackupVaultID(id.SubscriptionId, id.ResourceGroup, id.BackupVaultName)
	d.Set("name", id.Name)
	d.Set("vault_id", vaultId.ID())
	if resp.Properties != nil {
		if props, ok := resp.Properties.AsBackupPolicy(); ok {
			if err := d.Set("retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `default_retention_duration`: %+v", err)
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.BackupVaultName, id.ResourceGroup, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return nil
}

func flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(input *[]dataprotection.BasicBasePolicyRule) interface{} {
	if input == nil {
		return nil
	}

	for _, item := range *input {
		if retentionRule, ok := item.AsAzureRetentionRule(); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if retentionRule.Lifecycles != nil && len(*retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (*retentionRule.Lifecycles)[0].DeleteAfter.AsAbsoluteDeleteOption(); ok {
					return *deleteOption.Duration
				}
			}
		}
	}
	return nil
}
