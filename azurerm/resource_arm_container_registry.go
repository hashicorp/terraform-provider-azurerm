package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
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
					Type:             schema.TypeString,
					StateFunc:        azureRMNormalizeLocation,
					DiffSuppressFunc: azureRMSuppressLocationDiff,
				},
				Set: schema.HashString,
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

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			sku := d.Get("sku").(string)
			if _, ok := d.GetOk("storage_account_id"); ok {
				if strings.ToLower(sku) != strings.ToLower(string(containerregistry.Classic)) {
					return fmt.Errorf("`storage_account_id` can only be specified for a Classic (unmanaged) Sku.")
				}
			} else {
				if strings.ToLower(sku) == strings.ToLower(string(containerregistry.Classic)) {
					return fmt.Errorf("`storage_account_id` must be specified for a Classic (unmanaged) Sku.")
				}
			}

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
	location := azureRMNormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	adminUserEnabled := d.Get("admin_enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})
	geoReplicationLocations := d.Get("georeplication_locations").(*schema.Set)

	parameters := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.SkuName(sku),
			Tier: containerregistry.SkuTier(sku),
		},
		RegistryProperties: &containerregistry.RegistryProperties{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		parameters.StorageAccount = &containerregistry.StorageAccountProperties{
			ID: utils.String(v.(string)),
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// locations have been specified for geo-replication
	if geoReplicationLocations != nil && geoReplicationLocations.Len() > 0 {
		// the ACR is beeing created so no previous geo-replication locations
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
	oldGeoReplicationLocations := old.(*schema.Set)
	newGeoReplicationLocations := new.(*schema.Set)

	parameters := containerregistry.RegistryUpdateParameters{
		RegistryPropertiesUpdateParameters: &containerregistry.RegistryPropertiesUpdateParameters{
			AdminUserEnabled: utils.Bool(adminUserEnabled),
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
	if newGeoReplicationLocations != nil && newGeoReplicationLocations.Len() > 0 && !strings.EqualFold(sku, string(containerregistry.Premium)) {
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

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Container Registry %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if strings.EqualFold(sku, string(containerregistry.Premium)) {
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

	replicationLocationsToCreate := getReplicationLocationsToCreate(oldGeoReplicationLocations, newGeoReplicationLocations)
	replicationLocationsToDelete := getReplicationLocationsToDelete(oldGeoReplicationLocations, newGeoReplicationLocations)

	// create new geo-replication locations
	for _, locationToCreate := range replicationLocationsToCreate {
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

	// delete not needed geo-replication locations
	for _, locationToDelete := range replicationLocationsToDelete {
		future, err := replicationClient.Delete(ctx, resourceGroup, name, locationToDelete)
		if err != nil {
			return fmt.Errorf("Error deleting Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToDelete, err)
		}

		if err = future.WaitForCompletionRef(ctx, replicationClient.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Container Registry Replication %q (Resource Group %q, Location %q): %+v", name, resourceGroup, locationToDelete, err)
		}
	}

	return nil
}

func getReplicationLocationsToDelete(oldGeoreplicationLocations []interface{}, newGeoreplicationLocations []interface{}) []string {
	replicationLocationsToDelete := []string{}

	for _, oldLocation := range oldGeoreplicationLocations {
		oldLocation := azureRMNormalizeLocation(oldLocation)
		doDelete := true
		for _, newLocation := range newGeoreplicationLocations {
			newLocation := azureRMNormalizeLocation(newLocation)
			if newLocation == oldLocation {
				doDelete = false
			}
		}

		if doDelete {
			replicationLocationsToDelete = append(replicationLocationsToDelete, oldLocation)
		}
	}

	return replicationLocationsToDelete
}

func getReplicationLocationsToCreate(oldGeoreplicationLocations []interface{}, newGeoreplicationLocations []interface{}) []string {
	replicationLocationsToCreate := []string{}

	for _, newLocation := range newGeoreplicationLocations {
		newLocation := azureRMNormalizeLocation(newLocation)
		doCreate := true
		for _, oldLocation := range oldGeoreplicationLocations {
			oldLocation := azureRMNormalizeLocation(oldLocation)
			if oldLocation == newLocation {
				doCreate = false
			}
		}

		if doCreate {
			replicationLocationsToCreate = append(replicationLocationsToCreate, newLocation)
		}
	}

	return replicationLocationsToCreate
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
		credsResp, err := client.ListCredentials(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error making Read request on Azure Container Registry %s for Credentials: %s", name, err)
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

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRMContainerRegistryName(v interface{}, k string) (ws []string, errors []error) {
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

	return
}
