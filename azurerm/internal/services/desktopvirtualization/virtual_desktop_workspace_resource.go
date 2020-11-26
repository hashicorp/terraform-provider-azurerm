package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2019-12-10-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var workspaceResourceType = "azurerm_virtual_desktop_workspace"

func resourceArmDesktopVirtualizationWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDesktopVirtualizationWorkspaceCreateUpdate,
		Read:   resourceArmDesktopVirtualizationWorkspaceRead,
		Update: resourceArmDesktopVirtualizationWorkspaceCreateUpdate,
		Delete: resourceArmDesktopVirtualizationWorkspaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty, // TODO: determine more accurate requirements in time
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"friendly_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDesktopVirtualizationWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace create/update")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewWorkspaceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Virtual Desktop Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.WorkspaceProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace", id.ID(""))
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.Workspace{
		Location: &location,
		Tags:     tags.Expand(t),
		WorkspaceProperties: &desktopvirtualization.WorkspaceProperties{
			Description:  utils.String(d.Get("description").(string)),
			FriendlyName: utils.String(d.Get("friendly_name").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, context); err != nil {
		return fmt.Errorf("creating Virtual Desktop Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID(""))

	return resourceArmDesktopVirtualizationWorkspaceRead(d, meta)
}

func resourceArmDesktopVirtualizationWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Workspace %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Desktop Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDesktopVirtualizationWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, workspaceResourceType)
	defer locks.UnlockByName(id.Name, workspaceResourceType)

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Desktop Virtualization Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
