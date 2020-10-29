package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2019-12-10-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualDesktopApplicationGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualDesktopApplicationGroupCreateUpdate,
		Read:   resourceArmVirtualDesktopApplicationGroupRead,
		Update: resourceArmVirtualDesktopApplicationGroupCreateUpdate,
		Delete: resourceArmVirtualDesktopApplicationGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.VirtualDesktopApplicationGroupID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{1,260}$"),
						"Virtual desktop application group name must be 1 - 260 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(desktopvirtualization.ApplicationGroupTypeDesktop),
					string(desktopvirtualization.ApplicationGroupTypeRemoteApp),
				}, false),
			},

			"host_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},

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

func resourceArmVirtualDesktopApplicationGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Host Pool creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Checking for presence of existing Virtual Desktop Host Pool %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_application_group", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.ApplicationGroup{
		Location: &location,
		Tags:     tags.Expand(t),
		ApplicationGroupProperties: &desktopvirtualization.ApplicationGroupProperties{
			ApplicationGroupType: desktopvirtualization.ApplicationGroupType(d.Get("type").(string)),
			FriendlyName:         utils.String(d.Get("friendly_name").(string)),
			Description:          utils.String(d.Get("description").(string)),
			HostPoolArmPath:      utils.String(d.Get("host_pool_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, context); err != nil {
		return fmt.Errorf("Creating Virtual Desktop Host Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	result, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Retrieving Virtual Desktop Host Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot read Virtual Desktop Host Pool %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*result.ID)

	return resourceArmVirtualDesktopApplicationGroupRead(d, meta)
}

func resourceArmVirtualDesktopApplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualDesktopApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ApplicationGroupProperties; props != nil {
		d.Set("type", string(props.ApplicationGroupType))

		if fn := props.FriendlyName; fn != nil {
			d.Set("friendly_name", fn)
		}

		if desc := props.Description; desc != nil {
			d.Set("description", desc)
		}

		if hstpl := props.HostPoolArmPath; hstpl != nil {
			d.Set("host_pool_id", hstpl)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualDesktopApplicationGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualDesktopApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Deleting Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
