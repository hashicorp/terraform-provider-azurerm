package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/eventgrid"
	"github.com/hashicorp/terraform/helper/schema"
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

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

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
	client := meta.(*ArmClient).eventGridTopicsClient

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties := eventgrid.Topic{
		Location:        &location,
		TopicProperties: &eventgrid.TopicProperties{},
		Tags:            expandTags(tags),
	}

	log.Printf("[INFO] preparing arguments for AzureRM EventGrid Topic creation with Properties: %+v.", properties)

	_, createErr := client.CreateOrUpdate(resourceGroup, name, properties, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resourceGroup, name)
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
	client := meta.(*ArmClient).eventGridTopicsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["topics"]

	resp, err := client.Get(resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] EventGrid Topic '%s' was not found (resource group '%s')", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Topic '%s': %+v", name, err)
	}

	keys, err := client.ListSharedAccessKeys(resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Access Keys for EventGrid Topic '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.TopicProperties; props != nil {
		d.Set("endpoint", props.Endpoint)
	}

	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmEventGridTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridTopicsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["topics"]

	deleteResp, deleteErr := client.Delete(resGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-deleteErr

	if utils.ResponseWasNotFound(resp) {
		return nil
	}

	return err
}
