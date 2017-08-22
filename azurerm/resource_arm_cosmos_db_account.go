package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/arm/cosmos-db"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				ValidateFunc: validateAzureRmCosmosDBAccountName,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"offer_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cosmosdb.Standard),
				}, true),
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
								string(cosmosdb.BoundedStaleness),
								string(cosmosdb.Eventual),
								string(cosmosdb.Session),
								string(cosmosdb.Strong),
							}, true),
						},

						"max_interval_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 100),
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
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	offerType := d.Get("offer_type").(string)
	ipRangeFilter := d.Get("ip_range_filter").(string)

	consistencyPolicy := expandAzureRmCosmosDBAccountConsistencyPolicy(d)
	failoverPolicies, err := expandAzureRmCosmosDBAccountFailoverPolicies(name, d)
	if err != nil {
		return err
	}
	tags := d.Get("tags").(map[string]interface{})

	parameters := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: &location,
		DatabaseAccountCreateUpdateProperties: &cosmosdb.DatabaseAccountCreateUpdateProperties{
			ConsistencyPolicy:        &consistencyPolicy,
			DatabaseAccountOfferType: &offerType,
			Locations:                &failoverPolicies,
			IPRangeFilter:            &ipRangeFilter,
		},
		Tags: expandTags(tags),
	}

	_, error := client.CreateOrUpdate(resGroup, name, parameters, make(chan struct{}))
	err = <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
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
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["databaseAccounts"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account '%s': %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resGroup)
	d.Set("offer_type", string(resp.DatabaseAccountOfferType))
	d.Set("ip_range_filter", resp.IPRangeFilter)
	flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d, resp.ConsistencyPolicy)
	flattenAndSetAzureRmCosmosDBAccountFailoverPolicy(d, resp.FailoverPolicies)

	keys, err := client.ListKeys(resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_master_key", keys.PrimaryMasterKey)
		d.Set("secondary_master_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(resGroup, name)
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

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["databaseAccounts"]

	deleteResp, error := client.Delete(resGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-error

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for CosmosDB Account '%s': %+v", name, err)
	}

	return nil
}

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData) cosmosdb.ConsistencyPolicy {
	inputs := d.Get("consistency_policy").(*schema.Set).List()
	input := inputs[0].(map[string]interface{})

	consistencyLevel := input["consistency_level"].(string)

	policy := cosmosdb.ConsistencyPolicy{
		DefaultConsistencyLevel: cosmosdb.DefaultConsistencyLevel(consistencyLevel),
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

func expandAzureRmCosmosDBAccountFailoverPolicies(databaseName string, d *schema.ResourceData) ([]cosmosdb.Location, error) {
	input := d.Get("failover_policy").(*schema.Set).List()
	locations := make([]cosmosdb.Location, 0, len(input))

	for _, configRaw := range input {
		data := configRaw.(map[string]interface{})

		locationName := azureRMNormalizeLocation(data["location"].(string))
		id := fmt.Sprintf("%s-%s", databaseName, locationName)
		failoverPriority := int32(data["priority"].(int))

		location := cosmosdb.Location{
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

func flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData, policy *cosmosdb.ConsistencyPolicy) {
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

func flattenAndSetAzureRmCosmosDBAccountFailoverPolicy(d *schema.ResourceData, list *[]cosmosdb.FailoverPolicy) {
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

func validateAzureRmCosmosDBAccountName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	r, _ := regexp.Compile("[a-z0-9-]")
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("CosmosDB Account Name can only contain lower-case characters, numbers and the `-` character."))
	}

	length := len(value)
	if length > 50 || 3 > length {
		errors = append(errors, fmt.Errorf("CosmosDB Account Name can only be between 3 and 50 seconds."))
	}

	return
}
