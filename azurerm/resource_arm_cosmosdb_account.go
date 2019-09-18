package azurerm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbAccountCreate,
		Read:   resourceArmCosmosDbAccountRead,
		Update: resourceArmCosmosDbAccountUpdate,
		Delete: resourceArmCosmosDbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,50}$"),
					"Cosmos DB Account name must be 3 - 50 characters long, contain only lowercase letters, numbers and hyphens.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": tags.Schema(),

			//resource fields
			"offer_type": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.Standard),
				}, true),
			},

			"kind": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(documentdb.GlobalDocumentDB),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.GlobalDocumentDB),
					string(documentdb.MongoDB),
				}, true),
			},

			"ip_range_filter": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\/([1-2][0-9]|3[0-2]))?\b[,]?)*$`),
					"Cosmos DB ip_range_filter must be a set of CIDR IP addresses separated by commas with no spaces: '10.0.0.1,10.0.0.2,10.20.0.0/16'",
				),
			},

			"enable_automatic_failover": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"consistency_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consistency_level": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
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
							ValidateFunc: validation.IntBetween(5, 86400), // single region values
						},

						"max_staleness_prefix": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(10, 1000000), // single region values
						},
					},
				},
			},

			// this actually maps to the Location field in the API/SDK on create/update so has been renamed
			// failover_policy is just the name of the field we get back from the API on a read
			"failover_policy": {
				Type:          schema.TypeSet,
				Optional:      true,
				Deprecated:    "This field has been renamed to 'geo_location' to match Azure's usage",
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
							StateFunc:        azure.NormalizeLocation,
							DiffSuppressFunc: azure.SuppressLocationDiff,
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
				Computed:      true,
				ConflictsWith: []string{"failover_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-a-z0-9]{3,50}$"),
								"Cosmos DB location prefix (ID) must be 3 - 50 characters long, contain only lowercase letters, numbers and hyphens.",
							),
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": azure.SchemaLocation(),

						"failover_priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountGeoLocationHash,
			},

			"capabilities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								"EnableTable",
								"EnableGremlin",
								"EnableCassandra",
								"EnableAggregationPipeline",
								"MongoDBv3.4",
								"mongoEnableDocLevelTTL",
							}, true),
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountCapabilitiesHash,
			},

			"is_virtual_network_filter_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"virtual_network_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
				Set: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
			},

			"enable_multiple_write_locations": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			//computed
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"read_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"write_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"primary_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"connection_strings": {
				Type:      schema.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmCosmosDbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing CosmosDB Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cosmosdb_account", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	ipRangeFilter := d.Get("ip_range_filter").(string)
	isVirtualNetworkFilterEnabled := d.Get("is_virtual_network_filter_enabled").(bool)
	enableAutomaticFailover := d.Get("enable_automatic_failover").(bool)
	enableMultipleWriteLocations := d.Get("enable_multiple_write_locations").(bool)

	r, err := client.CheckNameExists(ctx, name)
	if err != nil {
		return fmt.Errorf("Error checking if CosmosDB Account %q already exists (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if !utils.ResponseWasNotFound(r) {
		return fmt.Errorf("CosmosDB Account %s already exists, please import the resource via terraform import", name)
	}

	//hacky, todo fix up once deprecated field 'failover_policy' is removed
	var geoLocations []documentdb.Location
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
		return fmt.Errorf("Neither `geo_location` or `failover_policy` is set for CosmosDB Account %s", name)
	}

	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:      utils.String(offerType),
			IPRangeFilter:                 utils.String(ipRangeFilter),
			IsVirtualNetworkFilterEnabled: utils.Bool(isVirtualNetworkFilterEnabled),
			EnableAutomaticFailover:       utils.Bool(enableAutomaticFailover),
			ConsistencyPolicy:             expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                     &geoLocations,
			Capabilities:                  expandAzureRmCosmosDBAccountCapabilities(d),
			VirtualNetworkRules:           expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:  utils.Bool(enableMultipleWriteLocations),
		},
		Tags: tags.Expand(t),
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not

	cp := account.DatabaseAccountCreateUpdateProperties.ConsistencyPolicy
	if len(geoLocations) > 1 && cp != nil && cp.DefaultConsistencyLevel == documentdb.BoundedStaleness {
		if msp := cp.MaxStalenessPrefix; msp != nil && *msp < 100000 {
			return fmt.Errorf("Error max_staleness_prefix (%d) must be greater then 100000 when more then one geo_location is used", *msp)
		}
		if mis := cp.MaxIntervalInSeconds; mis != nil && *mis < 300 {
			return fmt.Errorf("Error max_interval_in_seconds (%d) must be greater then 300 (5min) when more then one geo_location is used", *mis)
		}
	}

	resp, err := resourceArmCosmosDbAccountApiUpsert(client, ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("Error creating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//for some reason capabilities doesn't always work on create, so lets patch it
	//tracked: https://github.com/Azure/azure-sdk-for-go/issues/2864
	future, err := client.Patch(ctx, resourceGroup, name, documentdb.DatabaseAccountPatchParameters{
		DatabaseAccountPatchProperties: &documentdb.DatabaseAccountPatchProperties{
			Capabilities: account.Capabilities,
		},
	})
	if err != nil {
		return fmt.Errorf("Error Patching CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(context.Background(), client.Client); err != nil {
		return fmt.Errorf("Error waiting on patch future CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	id := resp.ID
	if id == nil {
		return fmt.Errorf("Cannot read CosmosDB Account '%s' (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*id)

	return resourceArmCosmosDbAccountRead(d, meta)
}

func resourceArmCosmosDbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Cosmos DB Account update.")

	//move to function
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	kind := d.Get("kind").(string)
	offerType := d.Get("offer_type").(string)
	ipRangeFilter := d.Get("ip_range_filter").(string)
	isVirtualNetworkFilterEnabled := d.Get("is_virtual_network_filter_enabled").(bool)
	enableAutomaticFailover := d.Get("enable_automatic_failover").(bool)
	enableMultipleWriteLocations := d.Get("enable_multiple_write_locations").(bool)

	//hacky, todo fix up once deprecated field 'failover_policy' is removed
	var newLocations []documentdb.Location
	var err error
	if _, ok := d.GetOk("geo_location"); ok {
		newLocations, err = expandAzureRmCosmosDBAccountGeoLocations(name, d)
		if err != nil {
			return fmt.Errorf("Error expanding CosmosDB Account %q (Resource Group %q) geo locations: %+v", name, resourceGroup, err)
		}
	} else if _, ok := d.GetOk("failover_policy"); ok {
		newLocations, err = expandAzureRmCosmosDBAccountFailoverPolicy(name, d)
		if err != nil {
			return fmt.Errorf("Error expanding CosmosDB Account %q (Resource Group %q) failover_policy: %+v", name, resourceGroup, err)
		}
	} else {
		//could be a CustomizeDiff?, but this is temporary
		return fmt.Errorf("Neither `geo_location` or `failover_policy` is set for CosmosDB Account '%s'", name)
	}

	//get existing locations (if exists)
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account '%s': %s", name, err)
	}

	oldLocations := make([]documentdb.Location, 0)
	oldLocationsMap := map[string]documentdb.Location{}
	for _, l := range *resp.FailoverPolicies {
		location := documentdb.Location{
			ID:               l.ID,
			LocationName:     l.LocationName,
			FailoverPriority: l.FailoverPriority,
		}

		oldLocations = append(oldLocations, location)
		oldLocationsMap[azure.NormalizeLocation(*location.LocationName)] = location
	}

	//cannot update properties and add/remove replication locations at the same time
	//so first just update any changed properties
	account := documentdb.DatabaseAccountCreateUpdateParameters{
		Location: utils.String(location),
		Kind:     documentdb.DatabaseAccountKind(kind),
		DatabaseAccountCreateUpdateProperties: &documentdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:      utils.String(offerType),
			IPRangeFilter:                 utils.String(ipRangeFilter),
			IsVirtualNetworkFilterEnabled: utils.Bool(isVirtualNetworkFilterEnabled),
			EnableAutomaticFailover:       utils.Bool(enableAutomaticFailover),
			Capabilities:                  expandAzureRmCosmosDBAccountCapabilities(d),
			ConsistencyPolicy:             expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                     &oldLocations,
			VirtualNetworkRules:           expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:  utils.Bool(enableMultipleWriteLocations),
		},
		Tags: tags.Expand(t),
	}

	if _, err = resourceArmCosmosDbAccountApiUpsert(client, ctx, resourceGroup, name, account); err != nil {
		return fmt.Errorf("Error updating CosmosDB Account %q properties (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//determine if any locations have been renamed/priority reordered and remove them
	removedOne := false
	for _, l := range newLocations {
		if ol, ok := oldLocationsMap[*l.LocationName]; ok {
			if *l.FailoverPriority != *ol.FailoverPriority {
				if *l.FailoverPriority == 0 {
					return fmt.Errorf("Cannot change the failover priority of primary Cosmos DB account %q location %s to %d (Resource Group %q)", name, *l.LocationName, *l.FailoverPriority, resourceGroup)
				}
				delete(oldLocationsMap, *l.LocationName)
				removedOne = true
				continue
			}
			if *l.ID == "" && *ol.ID == resourceArmCosmosDbAccountGenerateDefaultId(name, *l.LocationName) {
				continue
			}
			if *l.ID != *ol.ID {
				if *l.FailoverPriority == 0 {
					return fmt.Errorf("Cannot change the prefix/ID of the primary Cosmos DB account %q location %s (Resource Group %q)", name, *l.LocationName, resourceGroup)
				}
				delete(oldLocationsMap, *l.LocationName)
				removedOne = true
			}
		}

	}

	if removedOne {
		locationsUnchanged := make([]documentdb.Location, 0, len(oldLocationsMap))
		for _, value := range oldLocationsMap {
			locationsUnchanged = append(locationsUnchanged, value)
		}

		account.DatabaseAccountCreateUpdateProperties.Locations = &locationsUnchanged
		if _, err = resourceArmCosmosDbAccountApiUpsert(client, ctx, resourceGroup, name, account); err != nil {
			return fmt.Errorf("Error removing CosmosDB Account %q renamed locations (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	//add any new/renamed locations
	account.DatabaseAccountCreateUpdateProperties.Locations = &newLocations
	upsertResponse, err := resourceArmCosmosDbAccountApiUpsert(client, ctx, resourceGroup, name, account)
	if err != nil {
		return fmt.Errorf("Error updating CosmosDB Account %q locations (Resource Group %q): %+v", name, resourceGroup, err)
	}

	id := (*upsertResponse).ID
	if id == nil {
		return fmt.Errorf("Cannot read CosmosDB Account '%s' (resource group %s) ID", name, resourceGroup)
	}

	//for some reason capabilities doesn't always work on create, so lets patch it
	//tracked: https://github.com/Azure/azure-sdk-for-go/issues/2864
	future, err := client.Patch(ctx, resourceGroup, name, documentdb.DatabaseAccountPatchParameters{
		DatabaseAccountPatchProperties: &documentdb.DatabaseAccountPatchProperties{
			Capabilities: account.Capabilities,
		},
	})
	if err != nil {
		return fmt.Errorf("Error Patching CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(context.Background(), client.Client); err != nil {
		return fmt.Errorf("Error waiting on patch future CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*id)

	return resourceArmCosmosDbAccountRead(d, meta)
}

func resourceArmCosmosDbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("resource_group_name", resourceGroup)

	d.Set("kind", string(resp.Kind))
	d.Set("offer_type", string(resp.DatabaseAccountOfferType))
	d.Set("ip_range_filter", resp.IPRangeFilter)
	d.Set("endpoint", resp.DocumentEndpoint)

	if v := resp.IsVirtualNetworkFilterEnabled; v != nil {
		d.Set("is_virtual_network_filter_enabled", resp.IsVirtualNetworkFilterEnabled)
	}

	if v := resp.EnableAutomaticFailover; v != nil {
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)
	}

	if v := resp.EnableMultipleWriteLocations; v != nil {
		d.Set("enable_multiple_write_locations", resp.EnableMultipleWriteLocations)
	}

	if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(resp.ConsistencyPolicy)); err != nil {
		return fmt.Errorf("Error setting CosmosDB Account %q `consistency_policy` (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if _, ok := d.GetOk("failover_policy"); ok {
		if err = d.Set("failover_policy", flattenAzureRmCosmosDBAccountFailoverPolicy(resp.FailoverPolicies)); err != nil {
			return fmt.Errorf("Error setting `failover_policy`: %+v", err)
		}
	} else {
		//if failover policy isn't default to using geo_location
		if err = d.Set("geo_location", flattenAzureRmCosmosDBAccountGeoLocations(d, resp)); err != nil {
			return fmt.Errorf("Error setting `geo_location`: %+v", err)
		}
	}

	if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilities(resp.Capabilities)); err != nil {
		return fmt.Errorf("Error setting `capabilities`: %+v", err)
	}

	if err = d.Set("virtual_network_rule", flattenAzureRmCosmosDBAccountVirtualNetworkRules(resp.VirtualNetworkRules)); err != nil {
		return fmt.Errorf("Error setting `virtual_network_rule`: %+v", err)
	}

	readEndpoints := make([]string, 0)
	if p := resp.ReadLocations; p != nil {
		for _, l := range *p {
			if l.DocumentEndpoint == nil {
				continue
			}

			readEndpoints = append(readEndpoints, *l.DocumentEndpoint)
		}
	}
	if err := d.Set("read_endpoints", readEndpoints); err != nil {
		return fmt.Errorf("Error setting `read_endpoints`: %s", err)
	}

	writeEndpoints := make([]string, 0)
	if p := resp.WriteLocations; p != nil {
		for _, l := range *p {
			if l.DocumentEndpoint == nil {
				continue
			}

			writeEndpoints = append(writeEndpoints, *l.DocumentEndpoint)
		}
	}
	if err := d.Set("write_endpoints", writeEndpoints); err != nil {
		return fmt.Errorf("Error setting `write_endpoints`: %s", err)
	}

	// ListKeys returns a data structure containing a DatabaseAccountListReadOnlyKeysResult pointer
	// implying that it also returns the read only keys, however this appears to not be the case
	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	}
	d.Set("primary_master_key", keys.PrimaryMasterKey)
	d.Set("secondary_master_key", keys.SecondaryMasterKey)

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Read Only Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", name, err)
	}
	d.Set("primary_readonly_master_key", readonlyKeys.PrimaryReadonlyMasterKey)
	d.Set("secondary_readonly_master_key", readonlyKeys.SecondaryReadonlyMasterKey)

	connStringResp, err := client.ListConnectionStrings(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(keys.Response) {
			log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Unable to List connection strings for CosmosDB Account %s: %s", name, err)
	}

	var connStrings []string
	if connStringResp.ConnectionStrings != nil {
		connStrings = make([]string, len(*connStringResp.ConnectionStrings))
		for i, v := range *connStringResp.ConnectionStrings {
			connStrings[i] = *v.ConnectionString
		}

	}
	d.Set("connection_strings", connStrings)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCosmosDbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmos.DatabaseClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["databaseAccounts"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for CosmosDB Account '%s': %+v", name, err)
	}

	//the SDK now will return a `WasNotFound` response even when still deleting
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"NotFound"},
		Timeout:    180 * time.Minute,
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err2 := client.Get(ctx, resourceGroup, name)
			if err2 != nil {

				if utils.ResponseWasNotFound(resp.Response) {
					return resp, "NotFound", nil
				}
				return nil, "", fmt.Errorf("Error reading CosmosDB Account %q after delete (Resource Group %q): %+v", name, resourceGroup, err2)
			}

			return resp, "Deleting", nil
		},
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Waiting forCosmosDB Account %q to delete (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func resourceArmCosmosDbAccountApiUpsert(client *documentdb.DatabaseAccountsClient, ctx context.Context, resourceGroup string, name string, account documentdb.DatabaseAccountCreateUpdateParameters) (*documentdb.DatabaseAccount, error) {
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, account)
	if err != nil {
		return nil, fmt.Errorf("Error creating/updating CosmosDB Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return nil, fmt.Errorf("Error waiting for the CosmosDB Account %q (Resource Group %q) to finish creating/updating: %+v", name, resourceGroup, err)
	}

	//if a replication location is added or removed it can take some time to provision
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting"},
		Target:     []string{"Succeeded"},
		Timeout:    180 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      30 * time.Second, // required because it takes some time before the 'creating' location shows up
		Refresh: func() (interface{}, string, error) {

			resp, err2 := client.Get(ctx, resourceGroup, name)
			if err2 != nil {
				return nil, "", fmt.Errorf("Error reading CosmosDB Account %q after create/update (Resource Group %q): %+v", name, resourceGroup, err2)
			}

			status := "Succeeded"
			for _, l := range append(*resp.ReadLocations, *resp.WriteLocations...) {
				if status = *l.ProvisioningState; status == "Creating" || status == "Updating" || status == "Deleting" {
					break //return the first non successful status.
				}
			}

			return resp, status, nil
		},
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for the CosmosDB Account %q (Resource Group %q) to provision: %+v", name, resourceGroup, err)
	}

	r := resp.(documentdb.DatabaseAccount)
	return &r, nil
}

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *schema.ResourceData) *documentdb.ConsistencyPolicy {
	i := d.Get("consistency_policy").([]interface{})
	if len(i) <= 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

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

func resourceArmCosmosDbAccountGenerateDefaultId(databaseName string, location string) string {
	return fmt.Sprintf("%s-%s", databaseName, location)
}

func expandAzureRmCosmosDBAccountGeoLocations(databaseName string, d *schema.ResourceData) ([]documentdb.Location, error) {

	locations := make([]documentdb.Location, 0)
	for _, l := range d.Get("geo_location").(*schema.Set).List() {
		data := l.(map[string]interface{})

		location := documentdb.Location{
			LocationName:     utils.String(azure.NormalizeLocation(data["location"].(string))),
			FailoverPriority: utils.Int32(int32(data["failover_priority"].(int))),
		}

		if v, ok := data["prefix"].(string); ok {
			data["id"] = v
		} else {
			data["id"] = utils.String(resourceArmCosmosDbAccountGenerateDefaultId(databaseName, *location.LocationName))
		}
		location.ID = utils.String(data["id"].(string))

		locations = append(locations, location)
	}

	//TODO maybe this should be in a CustomizeDiff
	// all priorities & locations must be unique
	byPriorities := make(map[int]interface{}, len(locations))
	byName := make(map[string]interface{}, len(locations))
	for _, location := range locations {

		priority := int(*location.FailoverPriority)
		name := *location.LocationName

		if _, ok := byPriorities[priority]; ok {
			return nil, fmt.Errorf("Each `geo_location` needs to have a unique failover_prioroty. Multiple instances of '%d' found", priority)
		}

		if _, ok := byName[name]; ok {
			return nil, fmt.Errorf("Each `geo_location` needs to be in unique location. Multiple instances of '%s' found", name)
		}

		byPriorities[priority] = location
		byName[name] = location
	}

	//and must have one of 0 priority
	if _, ok := byPriorities[0]; !ok {
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

		locationName := azure.NormalizeLocation(data["location"].(string))
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

func expandAzureRmCosmosDBAccountCapabilities(d *schema.ResourceData) *[]documentdb.Capability {

	capabilities := d.Get("capabilities").(*schema.Set).List()
	s := make([]documentdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]interface{})
		s = append(s, documentdb.Capability{Name: utils.String(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *schema.ResourceData) *[]documentdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*schema.Set).List()

	s := make([]documentdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]interface{})
		s[i] = documentdb.VirtualNetworkRule{ID: utils.String(m["id"].(string))}
	}
	return &s
}

func flattenAzureRmCosmosDBAccountConsistencyPolicy(policy *documentdb.ConsistencyPolicy) []interface{} {

	result := map[string]interface{}{}
	result["consistency_level"] = string(policy.DefaultConsistencyLevel)
	if policy.MaxIntervalInSeconds != nil {
		result["max_interval_in_seconds"] = int(*policy.MaxIntervalInSeconds)
	}
	if policy.MaxStalenessPrefix != nil {
		result["max_staleness_prefix"] = int(*policy.MaxStalenessPrefix)
	}

	return []interface{}{result}
}

//todo remove when failover_policy field is removed
func flattenAzureRmCosmosDBAccountFailoverPolicy(list *[]documentdb.FailoverPolicy) *schema.Set {
	results := schema.Set{
		F: resourceAzureRMCosmosDBAccountFailoverPolicyHash,
	}

	for _, i := range *list {
		result := map[string]interface{}{
			"id":       *i.ID,
			"location": azure.NormalizeLocation(*i.LocationName),
			"priority": int(*i.FailoverPriority),
		}

		results.Add(result)
	}

	return &results
}

func flattenAzureRmCosmosDBAccountGeoLocations(d *schema.ResourceData, account documentdb.DatabaseAccount) *schema.Set {
	locationSet := schema.Set{
		F: resourceAzureRMCosmosDBAccountGeoLocationHash,
	}

	//we need to propagate the `prefix` field so fetch existing
	prefixMap := map[string]string{}
	if locations, ok := d.GetOk("geo_location"); ok {
		for _, lRaw := range locations.(*schema.Set).List() {
			lb := lRaw.(map[string]interface{})
			prefixMap[lb["location"].(string)] = lb["prefix"].(string)
		}
	}

	for _, l := range *account.FailoverPolicies {
		id := *l.ID
		lb := map[string]interface{}{
			"id":                id,
			"location":          azure.NormalizeLocation(*l.LocationName),
			"failover_priority": int(*l.FailoverPriority),
		}

		//if id is not the default then it must be set via prefix
		if id != resourceArmCosmosDbAccountGenerateDefaultId(d.Get("name").(string), lb["location"].(string)) {
			lb["prefix"] = id
		}

		locationSet.Add(lb)
	}

	return &locationSet
}

func flattenAzureRmCosmosDBAccountCapabilities(capabilities *[]documentdb.Capability) *schema.Set {
	s := schema.Set{
		F: resourceAzureRMCosmosDBAccountCapabilitiesHash,
	}

	for _, c := range *capabilities {
		if v := c.Name; v != nil {
			e := map[string]interface{}{
				"name": *v,
			}
			s.Add(e)
		}
	}

	return &s
}

func flattenAzureRmCosmosDBAccountVirtualNetworkRules(rules *[]documentdb.VirtualNetworkRule) *schema.Set {
	results := schema.Set{
		F: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
	}

	if rules != nil {
		for _, r := range *rules {
			rule := map[string]interface{}{
				"id": *r.ID,
			}
			results.Add(rule)
		}
	}

	return &results
}

//todo remove once deprecated field `failover_policy` is removed
func resourceAzureRMCosmosDBAccountFailoverPolicyHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		location := azure.NormalizeLocation(m["location"].(string))
		priority := int32(m["priority"].(int))

		buf.WriteString(fmt.Sprintf("%s-%d", location, priority))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMCosmosDBAccountGeoLocationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		prefix := ""
		if v, ok := m["prefix"].(string); ok {
			prefix = v
		}
		location := azure.NormalizeLocation(m["location"].(string))
		priority := int32(m["failover_priority"].(int))

		buf.WriteString(fmt.Sprintf("%s-%s-%d", prefix, location, priority))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMCosmosDBAccountCapabilitiesHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(strings.ToLower(m["id"].(string)))
	}

	return hashcode.String(buf.String())
}
