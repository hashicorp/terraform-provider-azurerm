package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventGridTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventGridTopicCreateUpdate,
		Read:   resourceArmEventGridTopicRead,
		Update: resourceArmEventGridTopicCreateUpdate,
		Delete: resourceArmEventGridTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": tags.Schema(),

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmEventGridTopicCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGrid.TopicsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventGrid Topic %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventgrid_topic", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties := eventgrid.Topic{
		Location:        &location,
		TopicProperties: &eventgrid.TopicProperties{},
		Tags:            tags.Expand(t),
	}

	log.Printf("[INFO] preparing arguments for AzureRM EventGrid Topic creation with Properties: %+v.", properties)

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read EventGrid Topic %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventGridTopicRead(d, meta)
}

func resourceArmEventGridTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGrid.TopicsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["topics"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] EventGrid Topic '%s' was not found (resource group '%s')", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Topic '%s': %+v", name, err)
	}

	keys, err := client.ListSharedAccessKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Access Keys for EventGrid Topic '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.TopicProperties; props != nil {
		d.Set("endpoint", props.Endpoint)
	}

	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmEventGridTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGrid.TopicsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["topics"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Topic %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Topic %q: %+v", name, err)
	}

	return nil
}
