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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),

			//resource fields
			"offer_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.Standard),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

			"enable_automatic_failover": {
				Type:     schema.TypeBool,
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
								string(documentdb.ConsistentPrefix),
								string(documentdb.Eventual),
								string(documentdb.Session),
								string(documentdb.Strong),
							}, true),
						},

						"max_interval_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(5, 86400),
						},

						"max_staleness_prefix": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(10, 1000000),
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountConsistencyPolicyHash,
			},

			// this actually maps to the Location field in the API/SDK on create/update so has been renamed
			// failover_policy is just the name of the field we get back from the API on a read
			"failover_policy": {
				Type:          schema.TypeSet,
				Optional:      true,
				Deprecated:    "This field has been renamed to 'location' to better match the SDK/API",
				ConflictsWith: []string{"geo_location"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"location": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        azureRMNormalizeLocation,
							DiffSuppressFunc: azureRMSuppressLocationDiff,
						},

						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountFailoverPolicyHash,
			},

			"geo_location": {
				Type: schema.TypeSet,
				//Required:     true, //todo needs to be required when failover_policy is removed
				Optional:      true,
				ConflictsWith: []string{"failover_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": {
							Type:     schema.TypeString,
							Optional: true, //TODO try and set this?
							Computed: true,
							ForceNew: true,
						},

						"location": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        azureRMNormalizeLocation,
							DiffSuppressFunc: azureRMSuppressLocationDiff,
						},

						"failover_priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"read_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"write_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountGeoLocationHash,
			},

			//computed

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func resourceArmCosmosDBAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosDBClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	ipRangeFilter := d.Get("ip_range_filter").(string)
	enableAutomaticFailover := d.Get("enable_automatic_failover").(bool)

	//hacky, todo fix up once deprecated field 'failover_policy' is removed
	var geoLocations []documentdb.Location
	var err error
	if _, ok := d.GetOk("geo_location"); ok {
		geoLocations, err = expandAzureRmCosmosDBAccountGeoLocations(name, d)
		if err != nil {
			return fmt.Errorf("Error expanding CosmosDB Account %q (Resource Group %q) geo locations: %+v", name, resourceGroup, err)
		}
	} else if _, ok := d.GetOk("failover_policy"); ok {
		geoLocations, err = expandAzureRmCosmosDBAccountFailoverPolicy(name, d)
		if err != nil {
			return fmt.Errorf("Error expanding CosmosDB Account %q (Resource Group %q) failover_policy: %+v", name, resourceGroup, err)
		}
	} else {
		//could be a CustomizeDiff?, but this is temporary
		return fmt.Errorf("Neither `geo_location` or `failover_policy` is set for CosmosDB Account '%s'", name)
	}

	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			ConsistencyPolicy:        expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                &geoLocations,
			DatabaseAccountOfferType: utils.String(offerType),
			IPRangeFilter:            utils.String(ipRangeFilter),
			EnableAutomaticFailover:  utils.Bool(enableAutomaticFailover),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("Error creating/updating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the CosmosDB Account %q (Resource Group %q) to finish creating: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading Scheduler Job Collection %q after create/update (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//todo is this still required?
	if read.ID == nil {
		return fmt.Errorf("Cannot read CosmosDB Account '%s' (resource group %s) ID", name, resourceGroup)
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

	name := id.Path["databaseAccounts"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account '%s': %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resourceGroup)
	flattenAndSetTags(d, resp.Tags)

	d.Set("kind", string(resp.Kind))
	d.Set("offer_type", string(resp.DatabaseAccountOfferType))
	d.Set("ip_range_filter", resp.IPRangeFilter)
	d.Set("endpoint", resp.DocumentEndpoint)

	if v := resp.EnableAutomaticFailover; v != nil {
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)
	}

	flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d, resp.ConsistencyPolicy)
	if _, ok := d.GetOk("geo_location"); ok {
		if err := flattenAndSetAzureRmCosmosDBAccountGeoLocations(d, resp); err != nil {
			return fmt.Errorf("Error flattening geo-locations for CosmosDB Account '%s'", name)
		}
	} else if _, ok := d.GetOk("failover_policy"); ok {
		flattenAndSetAzureRmCosmosDBAccountFailoverPolicy(d, resp.FailoverPolicies)
	} else {
		//could be a CustomizeDiff, but this is temporary
		return fmt.Errorf("Neither `geo_location` or `failover_policy` is set for CosmosDB Account '%s'", name)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_master_key", keys.PrimaryMasterKey)
		d.Set("secondary_master_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_readonly_master_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_master_key", readonlyKeys.SecondaryReadonlyMasterKey)
	}

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

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData) *documentdb.ConsistencyPolicy {
	input := d.Get("consistency_policy").(*schema.Set).List()[0].(map[string]interface{})

	consistencyLevel := input["consistency_level"].(string)
	policy := documentdb.ConsistencyPolicy{
		DefaultConsistencyLevel: documentdb.DefaultConsistencyLevel(consistencyLevel),
	}

	if stalenessPrefix, ok := input["max_staleness_prefix"].(int); ok {
		policy.MaxStalenessPrefix = utils.Int64(int64(stalenessPrefix))
	}
	if maxInterval, ok := input["max_interval_in_seconds"].(int); ok {
		policy.MaxIntervalInSeconds = utils.Int32(int32(maxInterval))
	}

	return &policy
}

func expandAzureRmCosmosDBAccountGeoLocations(databaseName string, d *schema.ResourceData) ([]documentdb.Location, error) {

	locations := make([]documentdb.Location, 0)
	for _, l := range d.Get("geo_location").(*schema.Set).List() {
		data := l.(map[string]interface{})

		location := documentdb.Location{
			LocationName:     utils.String(azureRMNormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int32(int32(data["failover_priority"].(int))),
		}

		if v, ok := data["id"].(string); ok {
			location.ID = utils.String(v)
		} else {
			location.ID = utils.String(fmt.Sprintf("%s-%s", databaseName, *location.DocumentEndpoint))
		}

		locations = append(locations, location)
	}

	//TODO maybe this should be in a CustomizeDiff

	// all priorities must be unique
	priorities := make(map[int]struct{}, len(locations))
	for _, location := range locations {
		priority := int(*location.FailoverPriority)
		if _, ok := priorities[priority]; ok {
			return nil, fmt.Errorf("Each `geo_location` needs to have a unique failover_prioroty")
		}

		priorities[priority] = struct{}{}
	}

	//and we have one of 0 priority
	if _, ok := priorities[0]; !ok {
		return nil, fmt.Errorf("There needs to be a `geo_location` with a failover_priority of 0")
	}

	return locations, nil
}

//todo remove when deprecated field `failover_policy` is
func expandAzureRmCosmosDBAccountFailoverPolicy(databaseName string, d *schema.ResourceData) ([]documentdb.Location, error) {

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
			return nil, fmt.Errorf("Each Failover Policy needs to be unique")
		}

		locationIds[priority] = struct{}{}
	}

	if !containsWriteLocation {
		return nil, fmt.Errorf("Failover Policy should contain a Write Location (Location '0')")
	}

	return locations, nil
}

func flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData, policy *documentdb.ConsistencyPolicy) {
	results := schema.Set{
		F: resourceAzureRMCosmosDBAccountConsistencyPolicyHash,
	}

	result := map[string]interface{}{}
	result["consistency_level"] = string(policy.DefaultConsistencyLevel)
	result["max_interval_in_seconds"] = int(*policy.MaxIntervalInSeconds) //TODO need a nil check?
	result["max_staleness_prefix"] = int(*policy.MaxStalenessPrefix)
	results.Add(result)

	d.Set("consistency_policy", &results)
}

//todo remove when failover_policy field is removed
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

func flattenAndSetAzureRmCosmosDBAccountGeoLocations(d *schema.ResourceData, account documentdb.DatabaseAccount) error {
	locationSet := schema.Set{
		F: resourceAzureRMCosmosDBAccountGeoLocationHash,
	}

	locationMap := map[string]map[string]interface{}{} //map so we can easily set the read and write locations
	for _, l := range *account.FailoverPolicies {
		lb := map[string]interface{}{
			"id":                *l.ID,
			"location":          azureRMNormalizeLocation(*l.LocationName),
			"failover_priority": int(*l.FailoverPriority),
		}

		locationMap[*l.ID] = lb
		locationSet.Add(lb)
	}

	for _, l := range *account.ReadLocations {
		if lb, ok := locationMap[*l.ID]; ok {
			lb["read_endpoint"] = *l.DocumentEndpoint
		} else {
			return fmt.Errorf("Unable to find matching location for read endpoint '%s'", *l.DocumentEndpoint)
		}
	}
	for _, l := range *account.WriteLocations {
		if lb, ok := locationMap[*l.ID]; ok {
			lb["write_endpoint"] = *l.DocumentEndpoint
		} else {
			return fmt.Errorf("Unable to find matching location for write endpoint '%s'", *l.DocumentEndpoint)
		}
	}

	return d.Set("geo_location", &locationSet)
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

//todo remove once deprecated field `failover_policy` is removed
func resourceAzureRMCosmosDBAccountFailoverPolicyHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	location := azureRMNormalizeLocation(m["location"].(string))
	priority := int32(m["priority"].(int))

	buf.WriteString(fmt.Sprintf("%s-%d", location, priority))

	return hashcode.String(buf.String())
}

func resourceAzureRMCosmosDBAccountGeoLocationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if id, ok := m["id"].(string); ok {
		buf.WriteString(id)
	}
	location := azureRMNormalizeLocation(m["location"].(string))
	priority := int32(m["failover_priority"].(int))

	buf.WriteString(fmt.Sprintf("-%s-%d", location, priority))

	return hashcode.String(buf.String())
}
