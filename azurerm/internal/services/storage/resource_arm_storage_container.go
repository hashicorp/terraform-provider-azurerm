package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

func resourceArmStorageContainer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageContainerCreate,
		Read:          resourceArmStorageContainerRead,
		Delete:        resourceArmStorageContainerDelete,
		Update:        resourceArmStorageContainerUpdate,
		MigrateState:  ResourceStorageContainerMigrateState,
		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validate.StorageContainerName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateArmStorageAccountName,
			},

			"container_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					string(containers.Blob),
					string(containers.Container),
					"private",
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"default_encryption_scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"encryption_scope_for_all_blobs": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"immutability_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"since_creation_in_days": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 146000),
						},

						"allow_protected_append_writes": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						// "etag" is required in immutability policy delete
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"legal_hold": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: storageValidate.StorageContainerLegalHoldTag(),
							},
						},
					},
				},
			},

			"metadata": MetaDataComputedSchema(),

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_immutability_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remaining_retention_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceArmStorageContainerCreate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	mgmtContainerClient := storageClient.BlobContainersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaDataPtr(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q", accountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Containers Client: %s", err)
	}

	id := client.GetResourceID(accountName, containerName)
	existing, err := mgmtContainerClient.Get(ctx, account.ResourceGroup, accountName, containerName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", containerName, accountName, account.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_storage_container", *existing.ID)
	}

	log.Printf("[INFO] Creating Container %q in Storage Account %q", containerName, accountName)
	params := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			PublicAccess: accessLevel,
			Metadata:     metaData,
		},
	}
	if v := d.Get("encryption_scope_for_all_blobs").(bool); v {
		params.ContainerProperties.DenyEncryptionScopeOverride = utils.Bool(v)
	}

	if es, ok := d.GetOk("default_encryption_scope"); ok {
		params.ContainerProperties.DefaultEncryptionScope = utils.String(es.(string))
	}

	if _, err := mgmtContainerClient.Create(ctx, account.ResourceGroup, accountName, containerName, params); err != nil {
		return fmt.Errorf("creating Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", containerName, accountName, account.ResourceGroup, err)
	}

	if immutability := expandStorageContainerImmutabilityPolicy(d.Get("immutability_policy").([]interface{})); immutability != nil {
		if _, err := mgmtContainerClient.CreateOrUpdateImmutabilityPolicy(ctx, account.ResourceGroup, accountName, containerName, immutability, ""); err != nil {
			return fmt.Errorf("setting `immutability_policy` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", containerName, accountName, account.ResourceGroup, err)
		}
	}

	if expandLegalHold := expandStorageContainerLegalHold(d.Get("legal_hold").([]interface{})); expandLegalHold != nil {

		if _, err := mgmtContainerClient.SetLegalHold(ctx, account.ResourceGroup, accountName, containerName, *expandLegalHold); err != nil {
			return fmt.Errorf("setting `legal_hold` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", containerName, accountName, account.ResourceGroup, err)
		}
	}

	d.SetId(id)

	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	mgmtContainerClient := storageClient.BlobContainersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	params := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{},
	}

	if d.HasChange("container_access_type") {
		log.Printf("[DEBUG] Updating the Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		accessLevelRaw := d.Get("container_access_type").(string)
		params.ContainerProperties.PublicAccess = expandStorageContainerAccessLevel(accessLevelRaw)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		params.ContainerProperties.Metadata = ExpandMetaDataPtr(metaDataRaw)
	}

	if _, err := mgmtContainerClient.Update(ctx, account.ResourceGroup, id.AccountName, id.ContainerName, params); err != nil {
		return fmt.Errorf("updating Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	if d.HasChange("immutability_policy") {
		if immutabilityResp, err := mgmtContainerClient.DeleteImmutabilityPolicy(ctx, account.ResourceGroup, id.AccountName, id.ContainerName, d.Get("immutability_policy.0.etag").(string)); err != nil {
			if !utils.ResponseWasNotFound(immutabilityResp.Response) {
				return fmt.Errorf("deleting `immutability_policy` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", id.ContainerName, id.AccountName, account.ResourceGroup, err)
			}
		}
		if immutabilityPolicy := expandStorageContainerImmutabilityPolicy(d.Get("immutability_policy").([]interface{})); immutabilityPolicy != nil {
			if _, err := mgmtContainerClient.CreateOrUpdateImmutabilityPolicy(ctx, account.ResourceGroup, id.AccountName, id.ContainerName, immutabilityPolicy, ""); err != nil {
				return fmt.Errorf("setting `immutability_policy` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", id.ContainerName, id.AccountName, account.ResourceGroup, err)
			}
		}
	}

	if d.HasChange("legal_hold") {
		o, n := d.GetChange("legal_hold")
		if expandClearLegalHold := expandStorageContainerLegalHold(o.([]interface{})); expandClearLegalHold != nil {
			if _, err := mgmtContainerClient.ClearLegalHold(ctx, account.ResourceGroup, id.AccountName, id.ContainerName, *expandClearLegalHold); err != nil {
				return fmt.Errorf("clearing `legal_hold` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", id.ContainerName, id.AccountName, account.ResourceGroup, err)
			}
		}

		if expandLegalHold := expandStorageContainerLegalHold(n.([]interface{})); expandLegalHold != nil {
			if _, err := mgmtContainerClient.SetLegalHold(ctx, account.ResourceGroup, id.AccountName, id.ContainerName, *expandLegalHold); err != nil {
				return fmt.Errorf("setting `legal_hold` in Storage Account Container %q (Storage Account %q/ Resource Group Name %q): %+v", id.ContainerName, id.AccountName, account.ResourceGroup, err)
			}
		}
	}
	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	mgmtContainerClient := storageClient.BlobContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Account %q for Storage Container %q - assuming removed & removing from state", id.AccountName, id.ContainerName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	resp, err := mgmtContainerClient.Get(ctx, account.ResourceGroup, id.AccountName, id.ContainerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", id.ContainerName, id.AccountName, account.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Container %q (Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("name", id.ContainerName)
	d.Set("storage_account_name", id.AccountName)

	if props := resp.ContainerProperties; props != nil {
		d.Set("container_access_type", flattenStorageContainerAccessLevel(props.PublicAccess))
		if err := d.Set("metadata", FlattenMetaDataPtr(props.Metadata)); err != nil {
			return fmt.Errorf("setting `metadata`: %+v", err)
		}
		d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
		d.Set("has_legal_hold", props.HasLegalHold)
		d.Set("default_encryption_scope", props.DefaultEncryptionScope)
		d.Set("encryption_scope_for_all_blobs", props.DenyEncryptionScopeOverride)
		d.Set("deleted", props.Deleted)
		d.Set("remaining_retention_days", props.RemainingRetentionDays)
		if err := d.Set("immutability_policy", flattenStorageContainerImmutabilityPolicy(props.ImmutabilityPolicy, props.HasImmutabilityPolicy)); err != nil {
			return fmt.Errorf("setting `immutability_policy` in Container %q", id.ContainerName)
		}

		if err := d.Set("legal_hold", flattenStorageContainerLegalHold(props.LegalHold)); err != nil {
			return fmt.Errorf("setting `legal_hold` in Container %q", id.ContainerName)
		}
	}

	resourceManagerId := client.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, id.AccountName, id.ContainerName)
	d.Set("resource_manager_id", resourceManagerId)

	return nil
}

func resourceArmStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	mgmtContainerClient := storageClient.BlobContainersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	if _, err := mgmtContainerClient.Delete(ctx, account.ResourceGroup, id.AccountName, id.ContainerName); err != nil {
		return fmt.Errorf("deleting Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func expandStorageContainerAccessLevel(input string) storage.PublicAccess {
	// for historical reasons, "private" above is an empty string in the API
	// so the enum doesn't 1:1 match. You could argue the SDK should handle this
	// but this is suitable for now
	if input == "private" {
		return storage.PublicAccessNone
	}

	return storage.PublicAccess(input)
}

func flattenStorageContainerAccessLevel(input storage.PublicAccess) string {
	// for historical reasons, "private" above is an empty string in the API
	if input == storage.PublicAccessNone {
		return "private"
	}

	return string(input)
}

func expandStorageContainerImmutabilityPolicy(input []interface{}) *storage.ImmutabilityPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	policy := input[0].(map[string]interface{})
	return &storage.ImmutabilityPolicy{
		ImmutabilityPolicyProperty: &storage.ImmutabilityPolicyProperty{
			ImmutabilityPeriodSinceCreationInDays: utils.Int32(int32(policy["since_creation_in_days"].(int))),
			AllowProtectedAppendWrites:            utils.Bool(policy["allow_protected_append_writes"].(bool)),
		},
	}
}

func flattenStorageContainerImmutabilityPolicy(input *storage.ImmutabilityPolicyProperties, hasImmutabilityPolicy *bool) []interface{} {
	if input == nil || hasImmutabilityPolicy == nil || !*hasImmutabilityPolicy {
		return []interface{}{}
	}
	var immutabilityPeriodSinceCreationInDays int32
	var allowProtectedAppendWrites bool
	if props := input.ImmutabilityPolicyProperty; props != nil {
		if v := props.ImmutabilityPeriodSinceCreationInDays; v != nil {
			immutabilityPeriodSinceCreationInDays = *v
		}
		if v := props.AllowProtectedAppendWrites; v != nil {
			allowProtectedAppendWrites = *v
		}
	}
	var etag string
	if input.Etag != nil {
		etag = *input.Etag
	}

	return []interface{}{
		map[string]interface{}{
			"since_creation_in_days":        immutabilityPeriodSinceCreationInDays,
			"allow_protected_append_writes": allowProtectedAppendWrites,
			"etag":                          etag,
		},
	}
}

func expandStorageContainerLegalHold(input []interface{}) *storage.LegalHold {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	lh := input[0].(map[string]interface{})
	return &storage.LegalHold{
		Tags: utils.ExpandStringSlice(lh["tags"].(*schema.Set).List()),
	}

}

func flattenStorageContainerLegalHold(input *storage.LegalHoldProperties) []interface{} {
	if input == nil || input.HasLegalHold == nil || !*input.HasLegalHold {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"tags": flattenStorageContainerLegalHoldTags(input.Tags),
		},
	}
}

func flattenStorageContainerLegalHoldTags(input *[]storage.TagProperty) *schema.Set {
	res := &schema.Set{F: schema.HashString}
	if input == nil || len(*input) == 0 {
		return res
	}

	for _, t := range *input {
		if t.Tag != nil {
			res.Add(*t.Tag)
		}
	}

	return res
}
