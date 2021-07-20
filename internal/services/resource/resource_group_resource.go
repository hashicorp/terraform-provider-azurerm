package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceResourceGroupCreateUpdate,
		Read:   resourceResourceGroupRead,
		Update: resourceResourceGroupCreateUpdate,
		Delete: resourceResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceResourceGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing resource group: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := resources.Group{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("creating Resource Group %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Resource Group %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading resource group: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("location", location.NormalizeNilable(resp.Location))
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceResourceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	// In case deleting a non-empty resource group is not allowed, we should ensure the current resource group is empty.
	deleteItemsInResourceGroup := meta.(*clients.Client).Features.ResourceGroup.DeleteNestedItemsDuringDeletion
	if !deleteItemsInResourceGroup {
		resourceClient := meta.(*clients.Client).Resource.ResourcesClient
		it, err := resourceClient.ListByResourceGroupComplete(ctx, id.ResourceGroup, "", "", nil)
		if err != nil {
			return fmt.Errorf("listing resources in %s: %v", id, err)
		}
		if it.NotDone() {
			msg := ""
			if it.Value().ID != nil && it.Value().Type != nil {
				msg = fmt.Sprintf(", at least there is a resource with ID %q (type: %s)", *it.Value().ID, *it.Value().Type)
			}
			return fmt.Errorf(`%s is not empty%s. Please ensure the resource group is empty, or you can enable "delete_nested_items_during_deletion" in "resource_group" feature block to allow deleting a non-empty resource group`, id, msg)
		}
	}

	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, "")
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("deleting Resource Group %q: %+v", id.ResourceGroup, err)
	}

	err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("deleting Resource Group %q: %+v", id.ResourceGroup, err)
	}

	return nil
}
