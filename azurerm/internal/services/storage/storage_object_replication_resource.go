package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStorageObjectReplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageObjectReplicationCreate,
		Read:   resourceStorageObjectReplicationRead,
		Update: resourceStorageObjectReplicationUpdate,
		Delete: resourceStorageObjectReplicationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ObjectReplicationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"source_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"destination_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"rules": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"source_container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.StorageContainerName,
						},

						"destination_container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.StorageContainerName,
						},

						// Possible values are "" means "OnlyNewObjects", "1601-01-01T00:00:00Z" means "Everything" and timeStamp "2020-10-21T16:00:00Z"
						"copy_blobs_created_after": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "OnlyNewObjects",
							ValidateFunc: validate.ObjectReplicationCopyBlobsCreatedAfter,
						},

						"filter_out_blobs_with_prefix": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"source_object_replication_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"destination_object_replication_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageObjectReplicationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	srcAccount, err := parse.StorageAccountID(d.Get("source_storage_account_id").(string))
	if err != nil {
		return err
	}
	dstAccount, err := parse.StorageAccountID(d.Get("destination_storage_account_id").(string))
	if err != nil {
		return err
	}

	existingList, err := client.List(ctx, dstAccount.ResourceGroup, dstAccount.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existingList.Response) {
			return fmt.Errorf("checking for present of existing Storage Object Replication for destination %q): %+v", dstAccount, err)
		}
	}
	for _, existing := range *existingList.Value {
		if existing.ID != nil && *existing.ID != "" && existing.SourceAccount != nil && *existing.SourceAccount == srcAccount.Name && existing.DestinationAccount != nil && *existing.DestinationAccount == dstAccount.Name {
			return tf.ImportAsExistsError("azurerm_storage_object_replication", *existing.ID)
		}
	}

	props := storage.ObjectReplicationPolicy{
		ObjectReplicationPolicyProperties: &storage.ObjectReplicationPolicyProperties{
			SourceAccount:      utils.String(srcAccount.Name),
			DestinationAccount: utils.String(dstAccount.Name),
			Rules:              expandArmObjectReplicationRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// create in dest storage account
	dstResp, err := client.CreateOrUpdate(ctx, dstAccount.ResourceGroup, dstAccount.Name, "default", props)
	if err != nil {
		return fmt.Errorf("creating Storage Object Replication for destination storage account name %q: %+v", dstAccount.Name, err)
	}

	if dstResp.ID == nil || *dstResp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Storage Object Replication for destination storage account name %q ID", dstAccount.Name)
	}

	if dstResp.Name == nil || *dstResp.Name == "" {
		return fmt.Errorf("empty or nil Name returned for Storage Object Replication for destination storage account name %q ID", dstAccount.Name)
	}

	// create in source storage account, update policy Id and ruleId which are computed from destination ORP
	props.Rules = dstResp.Rules
	srcResp, err := client.CreateOrUpdate(ctx, srcAccount.ResourceGroup, srcAccount.Name, *dstResp.Name, props)
	if err != nil {
		return fmt.Errorf("creating Storage Object Replication %q for source storage account name %q: %+v", *dstResp.Name, srcAccount.Name, err)
	}

	if srcResp.ID == nil || *srcResp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Storage Object Replication for destination storage account name %q ID", dstAccount.Name)
	}

	// here we concat the `srcId;dstId` as the resource id
	d.SetId(fmt.Sprintf("%s;%s", *srcResp.ID, *dstResp.ID))

	return resourceStorageObjectReplicationRead(d, meta)
}

func resourceStorageObjectReplicationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	props := storage.ObjectReplicationPolicy{
		ObjectReplicationPolicyProperties: &storage.ObjectReplicationPolicyProperties{
			SourceAccount:      utils.String(id.SrcStorageAccountName),
			DestinationAccount: utils.String(id.DstStorageAccountName),
			Rules:              expandArmObjectReplicationRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// update in dest storage account
	resp, err := client.CreateOrUpdate(ctx, id.DstResourceGroup, id.DstStorageAccountName, id.DstName, props)
	if err != nil {
		return fmt.Errorf("updating %q for destination storage account name %q: %+v", id, id.DstStorageAccountName, err)
	}

	// update in source storage account, update policy Id and ruleId
	props.Rules = resp.Rules
	if _, err := client.CreateOrUpdate(ctx, id.SrcResourceGroup, id.SrcStorageAccountName, id.SrcName, props); err != nil {
		return fmt.Errorf("updating %q for source storage account name %q: %+v", id, id.SrcStorageAccountName, err)
	}

	return resourceStorageObjectReplicationRead(d, meta)
}

func resourceStorageObjectReplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	dstResp, err := client.Get(ctx, id.DstResourceGroup, id.DstStorageAccountName, id.DstName)
	if err != nil {
		if utils.ResponseWasNotFound(dstResp.Response) {
			log.Printf("[INFO] storage object replication %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	srcResp, err := client.Get(ctx, id.SrcResourceGroup, id.SrcStorageAccountName, id.SrcName)
	if err != nil {
		if utils.ResponseWasNotFound(srcResp.Response) {
			log.Printf("[INFO] storage object replication %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if props := dstResp.ObjectReplicationPolicyProperties; props != nil {
		d.Set("source_storage_account_id", parse.NewStorageAccountID(id.SrcSubscriptionId, id.SrcResourceGroup, id.SrcStorageAccountName).ID())
		d.Set("destination_storage_account_id", parse.NewStorageAccountID(id.DstSubscriptionId, id.DstResourceGroup, id.DstStorageAccountName).ID())
		if err := d.Set("rules", flattenObjectReplicationRules(props.Rules)); err != nil {
			return fmt.Errorf("setting `rules`: %+v", err)
		}
		d.Set("source_object_replication_id", id.SourceObjectReplicationID())
		d.Set("destination_object_replication_id", id.DestinationObjectReplicationID())
	}
	return nil
}

func resourceStorageObjectReplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.DstResourceGroup, id.DstStorageAccountName, id.DstName); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if _, err := client.Delete(ctx, id.SrcResourceGroup, id.SrcStorageAccountName, id.SrcName); err != nil {
		return fmt.Errorf("deleting %q : %+v", id, err)
	}
	return nil
}

func expandArmObjectReplicationRuleArray(input []interface{}) *[]storage.ObjectReplicationPolicyRule {
	results := make([]storage.ObjectReplicationPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := storage.ObjectReplicationPolicyRule{
			SourceContainer:      utils.String(v["source_container_name"].(string)),
			DestinationContainer: utils.String(v["destination_container_name"].(string)),
			Filters: &storage.ObjectReplicationPolicyFilter{
				MinCreationTime: utils.String(expandArmObjectReplicationMinCreationTime(v["copy_blobs_created_after"].(string))),
			},
		}

		if r, ok := v["name"].(string); ok && r != "" {
			result.RuleID = utils.String(r)
		}

		if f, ok := v["filter_out_blobs_with_prefix"]; ok {
			result.Filters.PrefixMatch = utils.ExpandStringSlice(f.(*pluginsdk.Set).List())
		}

		results = append(results, result)
	}
	return &results
}

func expandArmObjectReplicationMinCreationTime(input string) string {
	switch input {
	case "Everything":
		return "1601-01-01T00:00:00Z"
	case "OnlyNewObjects":
		return ""
	default:
		return input
	}
}

func flattenObjectReplicationRules(input *[]storage.ObjectReplicationPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var destinationContainer string
		if item.DestinationContainer != nil {
			destinationContainer = *item.DestinationContainer
		}

		var sourceContainer string
		if item.SourceContainer != nil {
			sourceContainer = *item.SourceContainer
		}

		var ruleId string
		if item.RuleID != nil {
			ruleId = *item.RuleID
		}

		var minCreationTime string
		if item.Filters != nil && item.Filters.MinCreationTime != nil {
			minCreationTime = *item.Filters.MinCreationTime
		}

		var prefix []interface{}
		if item.Filters != nil && item.Filters.PrefixMatch != nil {
			prefix = utils.FlattenStringSlice(item.Filters.PrefixMatch)
		}

		v := map[string]interface{}{
			"destination_container_name":   destinationContainer,
			"source_container_name":        sourceContainer,
			"copy_blobs_created_after":     flattenArmObjectReplicationMinCreationTime(minCreationTime),
			"filter_out_blobs_with_prefix": pluginsdk.NewSet(pluginsdk.HashString, prefix),
			"name":                         ruleId,
		}
		results = append(results, v)
	}
	return results
}

func flattenArmObjectReplicationMinCreationTime(input string) string {
	switch input {
	case "1601-01-01T00:00:00Z":
		return "Everything"
	case "":
		return "OnlyNewObjects"
	default:
		return input
	}
}
