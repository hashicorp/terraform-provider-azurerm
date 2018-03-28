package azurerm

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDBAccountCreateUpdate,
		Read:   resourceArmCosmosDBAccountRead,
		Update: resourceArmCosmosDBAccountCreateUpdate,
		Delete: resourceArmCosmosDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDBAccountName,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"offer_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.Standard),
				}, true),
			},

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(documentdb.GlobalDocumentDB),
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.GlobalDocumentDB),
					string(documentdb.MongoDB),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"ip_range_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"consistency_policy": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consistency_level": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(documentdb.BoundedStaleness),
								string(documentdb.Eventual),
								string(documentdb.Session),
								string(documentdb.Strong),
							}, true),
						},

						"max_interval_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 86400),
						},

						"max_staleness_prefix": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(1, 2147483647),
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountConsistencyPolicyHash,
			},

			"failover_policy": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"location": {
							Type:      schema.TypeString,
							Required:  true,
							StateFunc: azureRMNormalizeLocation,
						},

						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountFailoverPolicyHash,
			},

			"primary_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_readonly_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_readonly_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmCosmosDBAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosDBClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	ipRangeFilter := d.Get("ip_range_filter").(string)

	consistencyPolicy := expandAzureRmCosmosDBAccountConsistencyPolicy(d)
	failoverPolicies, err := expandAzureRmCosmosDBAccountFailoverPolicies(name, d)
	if err != nil {
		return err
	}
	tags := d.Get("tags").(map[string]interface{})

	parameters := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			ConsistencyPolicy:        &consistencyPolicy,
			Locations:                &failoverPolicies,
			DatabaseAccountOfferType: utils.String(offerType),
			IPRangeFilter:            utils.String(ipRangeFilter),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read CosmosDB Account '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmCosmosDBAccountRead(d, meta)
}

func resourceArmCosmosDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosDBClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["databaseAccounts"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account '%s': %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resGroup)
	d.Set("kind", string(resp.Kind))
	d.Set("offer_type", string(resp.DatabaseAccountOfferType))
	d.Set("ip_range_filter", resp.IPRangeFilter)
	flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d, resp.ConsistencyPolicy)
	flattenAndSetAzureRmCosmosDBAccountFailoverPolicy(d, resp.FailoverPolicies)

	keys, err := client.ListKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_master_key", keys.PrimaryMasterKey)
		d.Set("secondary_master_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_readonly_master_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_master_key", readonlyKeys.SecondaryReadonlyMasterKey)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmCosmosDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosDBClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["databaseAccounts"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for CosmosDB Account '%s': %+v", name, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for CosmosDB Account '%s': %+v", name, err)
	}
	return nil
}

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData) documentdb.ConsistencyPolicy {
	inputs := d.Get("consistency_policy").(*schema.Set).List()
	input := inputs[0].(map[string]interface{})

	consistencyLevel := input["consistency_level"].(string)

	policy := documentdb.ConsistencyPolicy{
		DefaultConsistencyLevel: documentdb.DefaultConsistencyLevel(consistencyLevel),
	}

	if stalenessPrefix := input["max_staleness_prefix"].(int); stalenessPrefix > 0 {
		maxStalenessPrefix := int64(stalenessPrefix)
		policy.MaxStalenessPrefix = &maxStalenessPrefix
	}

	if maxInterval := input["max_interval_in_seconds"].(int); maxInterval > 0 {
		maxIntervalInSeconds := int32(maxInterval)
		policy.MaxIntervalInSeconds = &maxIntervalInSeconds
	}

	return policy
}

func expandAzureRmCosmosDBAccountFailoverPolicies(databaseName string, d *schema.ResourceData) ([]documentdb.Location, error) {
	input := d.Get("failover_policy").(*schema.Set).List()
	locations := make([]documentdb.Location, 0, len(input))

	for _, configRaw := range input {
		data := configRaw.(map[string]interface{})

		locationName := azureRMNormalizeLocation(data["location"].(string))
		id := fmt.Sprintf("%s-%s", databaseName, locationName)
		failoverPriority := int32(data["priority"].(int))

		location := documentdb.Location{
			ID:               &id,
			LocationName:     &locationName,
			FailoverPriority: &failoverPriority,
		}

		locations = append(locations, location)
	}

	containsWriteLocation := false
	writeFailoverPriority := int32(0)
	for _, location := range locations {
		if *location.FailoverPriority == writeFailoverPriority {
			containsWriteLocation = true
			break
		}
	}

	// all priorities must be unique
	locationIds := make(map[int]struct{}, len(locations))
	for _, location := range locations {
		priority := int(*location.FailoverPriority)
		if _, ok := locationIds[priority]; ok {
			err := fmt.Errorf("Each Failover Policy needs to be unique")
			return nil, err
		}

		locationIds[priority] = struct{}{}
	}

	if !containsWriteLocation {
		err := fmt.Errorf("Failover Policy should contain a Write Location (Location '0')")
		return nil, err
	}

	return locations, nil
}

func flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData, policy *documentdb.ConsistencyPolicy) {
	results := schema.Set{
		F: resourceAzureRMCosmosDBAccountConsistencyPolicyHash,
	}

	result := map[string]interface{}{}
	result["consistency_level"] = string(policy.DefaultConsistencyLevel)
	result["max_interval_in_seconds"] = int(*policy.MaxIntervalInSeconds)
	result["max_staleness_prefix"] = int(*policy.MaxStalenessPrefix)
	results.Add(result)

	d.Set("consistency_policy", &results)
}

func flattenAndSetAzureRmCosmosDBAccountFailoverPolicy(d *schema.ResourceData, list *[]documentdb.FailoverPolicy) {
	results := schema.Set{
		F: resourceAzureRMCosmosDBAccountFailoverPolicyHash,
	}

	for _, i := range *list {
		result := map[string]interface{}{
			"id":       *i.ID,
			"location": azureRMNormalizeLocation(*i.LocationName),
			"priority": int(*i.FailoverPriority),
		}

		results.Add(result)
	}

	d.Set("failover_policy", &results)
}

func resourceAzureRMCosmosDBAccountConsistencyPolicyHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	consistencyLevel := m["consistency_level"].(string)
	maxInterval := m["max_interval_in_seconds"].(int)
	maxStalenessPrefix := m["max_staleness_prefix"].(int)

	buf.WriteString(fmt.Sprintf("%s-%d-%d", consistencyLevel, maxInterval, maxStalenessPrefix))

	return hashcode.String(buf.String())
}

func resourceAzureRMCosmosDBAccountFailoverPolicyHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	locationName := m["location"].(string)
	location := azureRMNormalizeLocation(locationName)
	priority := int32(m["priority"].(int))

	buf.WriteString(fmt.Sprintf("%s-%d", location, priority))

	return hashcode.String(buf.String())
}
