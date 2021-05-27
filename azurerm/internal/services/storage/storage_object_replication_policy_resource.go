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

func resourceStorageObjectReplicationPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageObjectReplicationPolicyCreate,
		Read:   resourceStorageObjectReplicationPolicyRead,
		Update: resourceStorageObjectReplicationPolicyUpdate,
		Delete: resourceStorageObjectReplicationPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ObjectReplicationPolicyID(id)
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
						"copy_over_from_time": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "OnlyNewObjects",
							ValidateFunc: validate.ObjectReplicationPolicyCopyOverFromTime,
						},

						"filter_prefix_matches": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"rule_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			// there will be two ids, the destination and source storage account object replication policy ids, we keep the destination one as the resource id
			// and we keep the source one here
			"source_object_replication_policy_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageObjectReplicationPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationPolicyClient
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
			return fmt.Errorf("checking for present of existing Storage Object Replication Policy for destination %q): %+v", dstAccount, err)
		}
	}
	for _, existing := range *existingList.Value {
		if existing.ID != nil && *existing.ID != "" && existing.SourceAccount != nil && *existing.SourceAccount == srcAccount.Name && existing.DestinationAccount != nil && *existing.DestinationAccount == dstAccount.Name {
			return tf.ImportAsExistsError("azurerm_storage_object_replication_policy", *existing.ID)
		}
	}

	props := storage.ObjectReplicationPolicy{
		ObjectReplicationPolicyProperties: &storage.ObjectReplicationPolicyProperties{
			SourceAccount:      utils.String(srcAccount.Name),
			DestinationAccount: utils.String(dstAccount.Name),
			Rules:              expandArmObjectReplicationPolicyRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// create in dest storage account
	resp, err := client.CreateOrUpdate(ctx, dstAccount.ResourceGroup, dstAccount.Name, "default", props)
	if err != nil {
		return fmt.Errorf("creating Storage Object Replication Policy for destination storage account name %q: %+v", dstAccount.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Storage Object Replication Policy for destination storage account name %q ID", dstAccount.Name)
	}

	// here the id is computed from the service
	// we keep the destination ORP ID as the resource id
	d.SetId(*resp.ID)

	if resp.Name == nil || *resp.Name == "" {
		return fmt.Errorf("empty or nil Name returned for Storage Object Replication Policy for destination storage account name %q ID", dstAccount.Name)
	}

	// create in source storage account, update policy Id and ruleId
	props.Rules = resp.Rules
	if _, err := client.CreateOrUpdate(ctx, srcAccount.ResourceGroup, srcAccount.Name, *resp.Name, props); err != nil {
		return fmt.Errorf("creating Storage Object Replication Policy %q for source storage account name %q: %+v", *resp.Name, srcAccount.Name, err)
	}

	return resourceStorageObjectReplicationPolicyRead(d, meta)
}

func resourceStorageObjectReplicationPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationPolicyClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	srcId, err := parse.ObjectReplicationPolicyID(d.Get("source_object_replication_policy_id").(string))
	if err != nil {
		return err
	}

	props := storage.ObjectReplicationPolicy{
		ObjectReplicationPolicyProperties: &storage.ObjectReplicationPolicyProperties{
			SourceAccount:      utils.String(srcId.StorageAccountName),
			DestinationAccount: utils.String(id.StorageAccountName),
			Rules:              expandArmObjectReplicationPolicyRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// update in dest storage account
	resp, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, id.Name, props)
	if err != nil {
		return fmt.Errorf("updating %q for destination storage account name %q: %+v", id, id.StorageAccountName, err)
	}

	// update in source storage account, update policy Id and ruleId
	props.Rules = resp.Rules
	if _, err := client.CreateOrUpdate(ctx, srcId.ResourceGroup, srcId.StorageAccountName, srcId.Name, props); err != nil {
		return fmt.Errorf("updating %q for source storage account name %q: %+v", srcId, srcId.StorageAccountName, err)
	}

	return resourceStorageObjectReplicationPolicyRead(d, meta)
}

func resourceStorageObjectReplicationPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Storage.ObjectReplicationPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] storage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if props := resp.ObjectReplicationPolicyProperties; props != nil {
		// in import an existing resource, the source storage account resource group is not returned in response
		// so we need to find the source storage account id by its name
		srcAccount, err := storageClient.FindAccount(ctx, *props.SourceAccount)
		if err != nil {
			return fmt.Errorf("retrieving Source Account %q for %q: %s", *props.SourceAccount, id, err)
		}
		if srcAccount == nil {
			return fmt.Errorf("unable to locate Storage Account %q", *props.SourceAccount)
		}

		d.Set("source_storage_account_id", srcAccount.ID)
		d.Set("destination_storage_account_id", parse.NewStorageAccountID(subscriptionId, id.ResourceGroup, id.StorageAccountName).ID())
		if err := d.Set("rules", flattenObjectReplicationPolicyRules(props.Rules)); err != nil {
			return fmt.Errorf("setting `rules`: %+v", err)
		}
		d.Set("source_object_replication_policy_id", parse.NewObjectReplicationPolicyID(subscriptionId, srcAccount.ResourceGroup, *props.SourceAccount, *resp.Name).ID())
	}
	return nil
}

func resourceStorageObjectReplicationPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ObjectReplicationPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.StorageAccountName, id.Name); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	srcId, err := parse.ObjectReplicationPolicyID(d.Get("source_object_replication_policy_id").(string))
	if err != nil {
		return err
	}
	if _, err := client.Delete(ctx, srcId.ResourceGroup, srcId.StorageAccountName, srcId.Name); err != nil {
		return fmt.Errorf("deleting %q : %+v", srcId, err)
	}
	return nil
}

func expandArmObjectReplicationPolicyRuleArray(input []interface{}) *[]storage.ObjectReplicationPolicyRule {
	results := make([]storage.ObjectReplicationPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := storage.ObjectReplicationPolicyRule{
			SourceContainer:      utils.String(v["source_container_name"].(string)),
			DestinationContainer: utils.String(v["destination_container_name"].(string)),
			Filters: &storage.ObjectReplicationPolicyFilter{
				MinCreationTime: utils.String(expandArmObjectReplicationPolicyMinCreationTime(v["copy_over_from_time"].(string))),
			},
		}

		if r, ok := v["rule_name"].(string); ok && r != "" {
			result.RuleID = utils.String(r)
		}

		if f, ok := v["filter_prefix_matches"]; ok {
			result.Filters.PrefixMatch = utils.ExpandStringSlice(f.(*pluginsdk.Set).List())
		}

		results = append(results, result)
	}
	return &results
}

func expandArmObjectReplicationPolicyMinCreationTime(input string) string {
	switch input {
	case "Everything":
		return "1601-01-01T00:00:00Z"
	case "OnlyNewObjects":
		return ""
	default:
		return input
	}
}

func flattenObjectReplicationPolicyRules(input *[]storage.ObjectReplicationPolicyRule) []interface{} {
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
			"destination_container_name": destinationContainer,
			"source_container_name":      sourceContainer,
			"copy_over_from_time":        flattenArmObjectReplicationPolicyMinCreationTime(minCreationTime),
			"filter_prefix_matches":      pluginsdk.NewSet(pluginsdk.HashString, prefix),
			"rule_name":                  ruleId,
		}
		results = append(results, v)
	}
	return results
}

func flattenArmObjectReplicationPolicyMinCreationTime(input string) string {
	switch input {
	case "1601-01-01T00:00:00Z":
		return "Everything"
	case "":
		return "OnlyNewObjects"
	default:
		return input
	}
}
