package datamigration

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/datamigration/mgmt/2018-04-19/datamigration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataMigrationService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataMigrationServiceCreate,
		Read:   resourceArmDataMigrationServiceRead,
		Update: resourceArmDataMigrationServiceUpdate,
		Delete: resourceArmDataMigrationServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"virtual_subnet_id": {
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

			"kind": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Cloud",
				ValidateFunc: validation.StringInSlice([]string{"Cloud"}, false),
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDataMigrationServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataMigration.ServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_migration_service", *existing.ID)
		}
	}

	skuName := d.Get("sku_name").(string)
	virtualSubnetID := d.Get("virtual_subnet_id").(string)
	location := d.Get("location").(string)

	parameters := datamigration.Service{
		Location: utils.String(location),
		ServiceProperties: &datamigration.ServiceProperties{
			VirtualSubnetID: utils.String(virtualSubnetID),
		},
		Sku: &datamigration.ServiceSku{Name: utils.String(skuName)},
	}
	if kind, ok := d.GetOk("kind"); ok {
		parameters.Kind = utils.String(kind.(string))
	}
	if t, ok := d.GetOk("tags"); ok {
		parameters.Tags = tags.Expand(t.(map[string]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error creating Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Migration Service (Service Name %q / Group Name %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmDataMigrationServiceRead(d, meta)
}

func resourceArmDataMigrationServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataMigration.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["services"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Data Migration Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	d.Set("id", resp.ID)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("resource_group_name", resourceGroup)
	d.Set("name", name)
	if serviceProperties := resp.ServiceProperties; serviceProperties != nil {
		d.Set("virtual_subnet_id", serviceProperties.VirtualSubnetID)
	}
	if resp.Sku != nil && resp.Sku.Name != nil {
		d.Set("sku_name", resp.Sku.Name)
	}
	d.Set("type", resp.Type)
	d.Set("kind", resp.Kind)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDataMigrationServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataMigration.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["services"]

	parameters := datamigration.Service{}

	if t, ok := d.GetOk("tags"); ok {
		parameters.Tags = tags.Expand(t.(map[string]interface{}))
	}

	future, err := client.Update(ctx, parameters, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error updating Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	return resourceArmDataMigrationServiceRead(d, meta)
}

func resourceArmDataMigrationServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataMigration.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["services"]

	// For dms tasks created via terraform, the deletion order has already handled the dependency between dms task and dms service. In which case, dms service will not start to delete until
	// dms tasks have been deleted.
	// For dms tasks created out of terraform, it is user's responsibility to ensure the deletion order. And for dms service, it makes sense to just force delete outstanding tasks.
	toDeleteRunningTasks := true
	future, err := client.Delete(ctx, resourceGroup, name, &toDeleteRunningTasks)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Data Migration Service (Service Name %q / Group Name %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
