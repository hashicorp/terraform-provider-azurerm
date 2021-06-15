package databasemigration

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datamigration/mgmt/2018-04-19/datamigration"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDatabaseMigrationService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabaseMigrationServiceCreate,
		Read:   resourceDatabaseMigrationServiceRead,
		Update: resourceDatabaseMigrationServiceUpdate,
		Delete: resourceDatabaseMigrationServiceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServiceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
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

func resourceDatabaseMigrationServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_database_migration_service", id.ID())
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

	future, err := client.CreateOrUpdate(ctx, parameters, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if serviceProperties := resp.ServiceProperties; serviceProperties != nil {
		d.Set("subnet_id", serviceProperties.VirtualSubnetID)
	}
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDatabaseMigrationServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return err
	}

	// TODO: fix this behaviour in 3.0 - Terraform should remove even if there's running tasks
	// this last param is `delete the resource even if it contains running tasks`
	toDeleteRunningTasks := false
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, &toDeleteRunningTasks)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
		}
	}

	return nil
}
