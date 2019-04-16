package azurerm

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"log"
	"regexp"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerRegistryCreate,
		Read:   resourceArmContainerRegistryRead,
		Update: resourceArmContainerRegistryUpdate,
		Delete: resourceArmContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceAzureRMContainerRegistryMigrateState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMContainerRegistryName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"sku": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(containerregistry.Classic),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerregistry.Classic),
					string(containerregistry.Basic),
					string(containerregistry.Standard),
					string(containerregistry.Premium),
				}, true),
			},

			"admin_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"georeplication_locations": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Set: azureRMHashLocation,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"storage_account": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "`storage_account` has been replaced by `storage_account_id`.",
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"login_server": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_username": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"network_access_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerregistry.DefaultActionAllow),
								string(containerregistry.DefaultActionDeny),
							}, false),
						},

						"subnet_rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerregistry.Allow),
										}, false),
									},
									"subnet_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
							},
						},

						"ip_rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerregistry.Allow),
										}, false),
									},
									"ip_range": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.CIDR,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			sku := d.Get("sku").(string)
			geoReplicationLocations := d.Get("georeplication_locations").(*schema.Set)
			// if locations have been specified for geo-replication then, the SKU has to be Premium
			if geoReplicationLocations != nil && geoReplicationLocations.Len() > 0 && !strings.EqualFold(sku, string(containerregistry.Premium)) {
				return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
			}

			return nil
		},
	}
}

func resourceArmContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Container Registry creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Container Registry %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})
	geoReplicationLocations := d.Get("georeplication_locations").(*schema.Set)

	networkRuleSet := expandNetworkRuleSet(d)
	//NetworkRuleSet is only supported by Premium SDK
	if &networkRuleSet != nil && !strings.EqualFold(sku, string(containerregistry.Premium)) {
		return fmt.Errorf("`network_rule_set` can only be specified for a Premium Sku.")
	}

	parameters := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		RegistryProperties: &containerregistry.RegistryProperties{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
			NetworkRuleSet:   networkRuleSet,
		},

		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		if !strings.EqualFold(sku, string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` can only be specified for a Classic (unmanaged) Sku.")
		}

		parameters.StorageAccount = &containerregistry.StorageAccountProperties{
			ID: utils.String(v.(string)),
		}
	} else {
		if strings.EqualFold(sku, string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` must be specified for a Classic (unmanaged) Sku.")
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// locations have been specified for geo-replication
	if geoReplicationLocations != nil && geoReplicationLocations.Len() > 0 {
		// the ACR is being created so no previous geo-replication locations
		oldGeoReplicationLocations := []interface{}{}
		err = applyGeoReplicationLocations(meta, resourceGroup, name, oldGeoReplicationLocations, geoReplicationLocations.List())
		if err != nil {
			return fmt.Errorf("Error applying geo replications for Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Registry %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func resourceArmContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Container Registry update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})

	old, new := d.GetChange("georeplication_locations")
	hasGeoReplicationChanges := d.HasChange("georeplication_locations")
	oldGeoReplicationLocations := old.(*schema.Set)
	newGeoReplicationLocations := new.(*schema.Set)

	networkRuleSet := expandNetworkRuleSet(d)
	//NetworkRuleSet is only supported by Premium SDK
	if &networkRuleSet != nil && !strings.EqualFold(sku, string(containerregistry.Premium)) {
		return fmt.Errorf("`network_rule_set` can only be specified for a Premium Sku.")
	}

	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
			NetworkRuleSet:   networkRuleSet,
		},
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		if !strings.EqualFold(sku, string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` can only be specified for a Classic (unmanaged) Sku.")
		}

		parameters.StorageAccount = &containerregistry.StorageAccountProperties{
			ID: utils.String(v.(string)),
		}
	} else {
		if strings.EqualFold(sku, string(containerregistry.Classic)) {
			return fmt.Errorf("`storage_account_id` must be specified for a Classic (unmanaged) Sku.")
		}
	}

	// geo replication is only supported by Premium Sku
	if hasGeoReplicationChanges && newGeoReplicationLocations.Len() > 0 && !strings.EqualFold(sku, string(containerregistry.Premium)) {
		return fmt.Errorf("ACR geo-replication can only be applied when using the Premium Sku.")
	}

	// if the registry had replications and is updated to another Sku than premium - remove old locations
	if !strings.EqualFold(sku, string(containerregistry.Premium)) && oldGeoReplicationLocations != nil && oldGeoReplicationLocations.Len() > 0 {
		err := applyGeoReplicationLocations(meta, resourceGroup, name, oldGeoReplicationLocations.List(), newGeoReplicationLocations.List())
		if err != nil {
			return fmt.Errorf("Error applying geo replications for Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if strings.EqualFold(sku, string(containerregistry.Premium)) && hasGeoReplicationChanges {
		err = applyGeoReplicationLocations(meta, resourceGroup, name, oldGeoReplicationLocations.List(), newGeoReplicationLocations.List())
		if err != nil {
			return fmt.Errorf("Error applying geo replications for Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Registry %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func applyGeoReplicationLocations(meta interface{}, resourceGroup string, name string, oldGeoReplicationLocations []interface{}, newGeoReplicationLocations []interface{}) error {
	replicationClient := meta.(*ArmClient).containerRegistryReplicationsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing to apply geo-replications for AzureRM Container Registry.")

	createLocations := make(map[string]bool)

	// loop on the new location values
	for _, nl := range newGeoReplicationLocations {
		newLocation := azureRMNormalizeLocation(nl)
		createLocations[newLocation] = true // the location needs to be created
	}

	// loop on the old location values
	for _, ol := range oldGeoReplicationLocations {
		// oldLocation was created from a previous deployment
		oldLocation := azureRMNormalizeLocation(ol)

		// if the list of locations to create already contains the location
		if _, ok := createLocations[oldLocation]; ok {
			createLocations[oldLocation] = false // the location do not need to be created, it already exists
		}
	}

	// create new geo-replication locations
	for locationToCreate := range createLocations {
		// if false, the location does not need to be created, continue
		if !createLocations[locationToCreate] {
			continue
		}

		// create the new replication location
		replication := containerregistry.Replication{
			Location: &locationToCreate,
			Name:     &locationToCreate,
		}

		future, err := replicationClient.Create(ctx, resourceGroup, name, locationToCreate, replication)
		if err != nil {
			return fmt.Errorf("Error creating Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToCreate, err)
		}

		if err = future.WaitForCompletionRef(ctx, replicationClient.Client); err != nil {
			return fmt.Errorf("Error waiting for creation of Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToCreate, err)
		}
	}

	// loop on the list of previously deployed locations
	for _, ol := range oldGeoReplicationLocations {
		oldLocation := azureRMNormalizeLocation(ol)
		// if the old location is still in the list of locations, then continue
		if _, ok := createLocations[oldLocation]; ok {
			continue
		}

		// the old location is not in the list of locations, delete it
		future, err := replicationClient.Delete(ctx, resourceGroup, name, oldLocation)
		if err != nil {
			return fmt.Errorf("Error deleting Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, oldLocation, err)
		}

		if err = future.WaitForCompletionRef(ctx, replicationClient.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, oldLocation, err)
		}
	}

	return nil
}

func resourceArmContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	replicationClient := meta.(*ArmClient).containerRegistryReplicationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["registries"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry %q was not found in Resource Group %q", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	location := resp.Location
	if location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("admin_enabled", resp.AdminUserEnabled)
	d.Set("login_server", resp.LoginServer)

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Tier))
	}

	if account := resp.StorageAccount; account != nil {
		d.Set("storage_account_id", account.ID)
	}

	if *resp.AdminUserEnabled {
		credsResp, errList := client.ListCredentials(ctx, resourceGroup, name)
		if errList != nil {
			return fmt.Errorf("Error making Read request on Azure Container Registry %s for Credentials: %s", name, errList)
		}

		d.Set("admin_username", credsResp.Username)
		for _, v := range *credsResp.Passwords {
			d.Set("admin_password", v.Value)
			break
		}
	} else {
		d.Set("admin_username", "")
		d.Set("admin_password", "")
	}

	flattenAndSetTags(d, resp.Tags)

	replications, err := replicationClient.List(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Container Registry %s for replications: %s", name, err)
	}

	replicationValues := replications.Values()

	// if there is more than one location (the main one and the replicas)
	if replicationValues != nil || len(replicationValues) > 1 {
		georeplication_locations := &schema.Set{F: schema.HashString}

		for _, value := range replicationValues {
			if value.Location != nil {
				valueLocation := azureRMNormalizeLocation(*value.Location)
				if location != nil && valueLocation != azureRMNormalizeLocation(*location) {
					georeplication_locations.Add(valueLocation)
				}
			}
		}

		d.Set("georeplication_locations", georeplication_locations)
	}

	return nil
}

func resourceArmContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerRegistryClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["registries"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry '%s': %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRMContainerRegistryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters only are allowed in %q: %q", k, value))
	}

	if 5 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 5 characters: %q", k, value))
	}

	if len(value) >= 50 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 50 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func expandNetworkRuleSet(d *schema.ResourceData) *containerregistry.NetworkRuleSet {
	configs := d.Get("network_access_profile").([]interface{})
	config := configs[0].(map[string]interface{})

	virtualNetworkRuleConfigs := config["subnet_rule"].([]interface{})
	virtualNetworkRules := make([]containerregistry.VirtualNetworkRule, 0)

	for _, virtualNetworkRuleInterface := range virtualNetworkRuleConfigs {
		config := virtualNetworkRuleInterface.(map[string]interface{})
		virtualNetworkRules =
			append(virtualNetworkRules, containerregistry.VirtualNetworkRule{
				Action:                   containerregistry.Action(config["action"].(string)),
				VirtualNetworkResourceID: utils.String(config["subnet_id"].(string)),
			})
	}

	ipRuleConfigs := config["ip_rule"].([]interface{})
	ipRules := make([]containerregistry.IPRule, 0)

	for _, ipRuleInterface := range ipRuleConfigs {
		config := ipRuleInterface.(map[string]interface{})
		ipRules =
			append(ipRules, containerregistry.IPRule{
				Action:           containerregistry.Action(config["action"].(string)),
				IPAddressOrRange: utils.String(config["ip_range"].(string)),
			})
	}

	networkRuleSet := containerregistry.NetworkRuleSet{
		DefaultAction:       containerregistry.DefaultAction(config["default_action"].(string)),
		VirtualNetworkRules: &virtualNetworkRules,
		IPRules:             &ipRules,
	}
	return &networkRuleSet
}
