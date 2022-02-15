package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var applicationGroupType = "azurerm_virtual_desktop_application_group"

func resourceVirtualDesktopApplicationGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Read:   resourceVirtualDesktopApplicationGroupRead,
		Update: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Delete: resourceVirtualDesktopApplicationGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationGroupID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationGroupV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(desktopvirtualization.ApplicationGroupTypeDesktop),
					string(desktopvirtualization.ApplicationGroupTypeRemoteApp),
				}, false),
			},

			"host_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.HostPoolID,
			},

			"friendly_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"default_desktop_display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualDesktopApplicationGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[INFO] preparing arguments for Virtual Desktop Application Group creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, applicationGroupType)
	defer locks.UnlockByName(name, applicationGroupType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewApplicationGroupID(subscriptionId, resourceGroup, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Virtual Desktop Application Group %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ApplicationGroupProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_application_group", resourceId)
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
		return fmt.Errorf("creating Virtual Desktop Application Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if desktopvirtualization.ApplicationGroupType(d.Get("type").(string)) == desktopvirtualization.ApplicationGroupTypeDesktop {
		if desktopFriendlyName := utils.String(d.Get("default_desktop_display_name").(string)); desktopFriendlyName != nil {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktop, err := desktopClient.Get(ctx, resourceGroup, name, "sessionDesktop")
			if err != nil {
				if !utils.ResponseWasNotFound(desktop.Response) {
					return fmt.Errorf("checking for presence of default desktop in Application Group %q (Resource Group %q): %s", name, resourceGroup, err)
				}
			}

			desktopPatch := desktopvirtualization.DesktopPatch{
				DesktopPatchProperties: &desktopvirtualization.DesktopPatchProperties{
					FriendlyName: desktopFriendlyName,
				},
			}

			if _, err := desktopClient.Update(ctx, resourceGroup, name, "sessionDesktop", &desktopPatch); err != nil {
				return fmt.Errorf("setting default desktop friendly name for Application Group %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	d.SetId(resourceId)
	return resourceVirtualDesktopApplicationGroupRead(d, meta)
}

func resourceVirtualDesktopApplicationGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Application Group %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Desktop Application Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ApplicationGroupProperties; props != nil {
		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("type", string(props.ApplicationGroupType))
		if props.ApplicationGroupType == desktopvirtualization.ApplicationGroupTypeDesktop {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktop, err := desktopClient.Get(ctx, id.ResourceGroup, id.Name, "sessionDesktop")
			// if the default desktop was found then set the display name attribute
			if err == nil {
				if desktopProps := desktop.DesktopProperties; desktopProps != nil {
					d.Set("default_desktop_display_name", desktopProps.FriendlyName)
				}
			}
		}

		hostPoolIdStr := ""
		if props.HostPoolArmPath != nil {
			// TODO: raise an API bug
			hostPoolId, err := parse.HostPoolIDInsensitively(*props.HostPoolArmPath)
			if err != nil {
				return fmt.Errorf("parsing Host Pool ID %q: %+v", *props.HostPoolArmPath, err)
			}

			hostPoolIdStr = hostPoolId.ID()
		}
		d.Set("host_pool_id", hostPoolIdStr)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualDesktopApplicationGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient

	id, err := parse.ApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, applicationGroupType)
	defer locks.UnlockByName(id.Name, applicationGroupType)

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Virtual Desktop Application Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
