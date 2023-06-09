package resource

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceTags() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceTagsCreateOrUpdate,
		Update: resourceTagsCreateOrUpdate,
		Read:   resourceTagsRead,
		Delete: resourceTagsDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceTagsID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceTagsCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TagsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceTagsID(d.Get("resource_id").(string))
	if err != nil {
		return fmt.Errorf("could not parse resource ID for tagsId \"%s\": %v", id, err)
	}
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.GetAtScope(ctx, id.ParentResourceID())
		if err != nil {
			if !response.WasNotFound(existing.Request.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ParentResourceID(), err)
			}
		}

		for tag := range d.Get("tags").(map[string]interface{}) {
			if existing.Properties.Tags[tag] != nil {
				return tf.ImportAsExistsError("azurerm_resource_tags", tag)
			}
		}
	}

	payload := resources.TagsResource{
		Properties: &resources.Tags{
			Tags: tags.Expand(t),
		},
	}

	if _, err := client.CreateOrUpdateAtScope(ctx, id.ParentResourceID(), payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id.ParentResourceID(), err)
	}

	d.SetId(id.ID())
	return resourceTagsRead(d, meta)
}

func resourceTagsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TagsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceTagsID(d.Id())
	if err != nil {
		return fmt.Errorf(d.Id(), err)
	}

	resp, err := client.GetAtScope(ctx, id.ParentResourceID())
	if err != nil {
		if response.WasNotFound(resp.Request.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}
	d.Set("resource_id", id.ParentResourceID())
	return tags.FlattenAndSet(d, resp.Properties.Tags)
}

func resourceTagsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TagsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := d.Get("resource_id").(string)

	id, err := parse.ResourceTagsID(resourceId)
	if err != nil {
		return err
	}

	if resp, err := client.DeleteAtScope(ctx, resourceId); err != nil {
		if response.WasNotFound(resp.Request.Response) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
