package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2020-11-02-preview/desktopvirtualization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var applicationType = "azurerm_virtual_desktop_application"

func resourceVirtualDesktopApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopApplicationCreateUpdate,
		Read:   resourceVirtualDesktopApplicationRead,
		Update: resourceVirtualDesktopApplicationCreateUpdate,
		Delete: resourceVirtualDesktopApplicationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationID(id)
			return err
		}),

		SchemaVersion: 0,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{1,260}$"),
						"Virtual desktop application name must be 1 - 260 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"application_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationGroupID,
			},

			"friendly_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Computed:     true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"path": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"command_line_argument_policy": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(desktopvirtualization.Allow),
					string(desktopvirtualization.DoNotAllow),
					string(desktopvirtualization.Require),
				}, false),
			},

			"command_line_arguments": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"show_in_portal": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"icon_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"icon_index": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceVirtualDesktopApplicationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[INFO] preparing arguments for Virtual Desktop Application creation")

	name := d.Get("name").(string)
	applicationGroup, _ := parse.ApplicationGroupID(d.Get("application_group_id").(string))

	locks.ByName(name, applicationType)
	defer locks.UnlockByName(name, applicationType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewApplicationID(subscriptionId, applicationGroup.ResourceGroup, applicationGroup.Name, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, applicationGroup.ResourceGroup, applicationGroup.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Virtual Desktop Application %q (Application Group %q) (Resource Group %q): %s", name, applicationGroup.Name, applicationGroup.ResourceGroup, err)
			}
		}

		if existing.ApplicationProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_application", resourceId)
		}
	}

	context := desktopvirtualization.Application{
		ApplicationProperties: &desktopvirtualization.ApplicationProperties{
			FriendlyName:         utils.String(d.Get("friendly_name").(string)),
			Description:          utils.String(d.Get("description").(string)),
			FilePath:             utils.String(d.Get("path").(string)),
			CommandLineSetting:   desktopvirtualization.CommandLineSetting(d.Get("command_line_argument_policy").(string)),
			CommandLineArguments: utils.String(d.Get("command_line_arguments").(string)),
			ShowInPortal:         utils.Bool(d.Get("show_in_portal").(bool)),
			IconPath:             utils.String(d.Get("icon_path").(string)),
			IconIndex:            utils.Int32(int32(d.Get("icon_index").(int))),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, applicationGroup.ResourceGroup, applicationGroup.Name, name, context); err != nil {
		return fmt.Errorf("creating Virtual Desktop Application %q (Application Group %q) (Resource Group %q): %+v", name, applicationGroup.Name, applicationGroup.ResourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceVirtualDesktopApplicationRead(d, meta)
}

func resourceVirtualDesktopApplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ApplicationGroupName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Application %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Desktop Application %q (Application Group %q) (Resource Group %q): %+v", id.Name, id.ApplicationGroupName, id.ResourceGroup, err)
	}

	applicationGroup := parse.ApplicationGroupId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		Name:           id.ApplicationGroupName,
	}

	d.Set("name", id.Name)
	d.Set("application_group_id", applicationGroup.ID())

	if props := resp.ApplicationProperties; props != nil {
		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("path", props.FilePath)
		d.Set("command_line_argument_policy", props.CommandLineSetting)
		d.Set("command_line_arguments", props.CommandLineArguments)
		d.Set("show_in_portal", props.ShowInPortal)
		d.Set("icon_path", props.IconPath)
		d.Set("icon_index", props.IconIndex)
	}

	return nil
}

func resourceVirtualDesktopApplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, applicationType)
	defer locks.UnlockByName(id.Name, applicationType)

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if _, err = client.Delete(ctx, id.ResourceGroup, id.ApplicationGroupName, id.Name); err != nil {
		return fmt.Errorf("deleting Virtual Desktop Application %q (Application Group %q) (Resource Group %q): %+v", id.Name, id.ApplicationGroupName, id.ResourceGroup, err)
	}

	return nil
}
