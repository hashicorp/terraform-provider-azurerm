package databasemigration

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/datamigration/mgmt/2018-04-19/datamigration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDatabaseMigrationService() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseMigrationServiceCreate,
		Read:   resourceDatabaseMigrationServiceRead,
		Update: resourceDatabaseMigrationServiceUpdate,
		Delete: resourceDatabaseMigrationServiceDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceID(id)
			return err
		}),

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
				ValidateFunc: validateDatabaseMigrationServiceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// No const defined in go sdk, the literal listed below is derived from the response of listskus endpoint.
					// See: https://docs.microsoft.com/en-us/rest/api/datamigration/resourceskus/listskus
					"Premium_4vCores",
					"Standard_1vCores",
					"Standard_2vCores",
					"Standard_4vCores",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDatabaseMigrationServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Database Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_database_migration_service", *existing.ID)
		}
	}

	skuName := d.Get("sku_name").(string)
	subnetID := d.Get("subnet_id").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	parameters := datamigration.Service{
		Location: utils.String(location),
		ServiceProperties: &datamigration.ServiceProperties{
			VirtualSubnetID: utils.String(subnetID),
		},
		Sku: &datamigration.ServiceSku{
			Name: utils.String(skuName),
		},
		Kind: utils.String("Cloud"), // currently only "Cloud" is supported, hence hardcode here
	}
	if t, ok := d.GetOk("tags"); ok {
		parameters.Tags = tags.Expand(t.(map[string]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error creating Database Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Database Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Database Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Database Migration Service (Service Name %q / Group Name %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Database Migration Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Database Migration Service (Service Name %q / Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("unexpected empty ID retrieved for Database Migration Service (Service Name %q / Group Name %q)", id.Name, id.ResourceGroup)
	}
	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("name", id.Name)
	if serviceProperties := resp.ServiceProperties; serviceProperties != nil {
		d.Set("subnet_id", serviceProperties.VirtualSubnetID)
	}
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDatabaseMigrationServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return err
	}

	parameters := datamigration.Service{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, parameters, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error updating Database Migration Service (Service Name %q / Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Database Migration Service (Service Name %q / Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return err
	}

	// Always leave outstanding migration tasks untouched when deleting DMS. This is to avoid unexpectedly delete any tasks managed out of terraform.
	toDeleteRunningTasks := false
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, &toDeleteRunningTasks)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Database Migration Service (Service Name %q / Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Database Migration Service (Service Name %q / Group Name %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
