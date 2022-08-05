package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			_, err := application.ParseApplicationID(id)
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
				ValidateFunc: applicationgroup.ValidateApplicationGroupID,
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
					string(application.CommandLineSettingAllow),
					string(application.CommandLineSettingDoNotAllow),
					string(application.CommandLineSettingRequire),
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

	applicationGroup, _ := applicationgroup.ParseApplicationGroupID(d.Get("application_group_id").(string))
	id := application.NewApplicationID(subscriptionId, applicationGroup.ResourceGroupName, applicationGroup.ApplicationGroupName, d.Get("name").(string))

	locks.ByName(id.ApplicationName, applicationType)
	defer locks.UnlockByName(id.ApplicationName, applicationType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_application", id.ID())
		}
	}

	payload := application.Application{
		Properties: application.ApplicationProperties{
			FriendlyName:         utils.String(d.Get("friendly_name").(string)),
			Description:          utils.String(d.Get("description").(string)),
			FilePath:             utils.String(d.Get("path").(string)),
			CommandLineSetting:   application.CommandLineSetting(d.Get("command_line_argument_policy").(string)),
			CommandLineArguments: utils.String(d.Get("command_line_arguments").(string)),
			ShowInPortal:         utils.Bool(d.Get("show_in_portal").(bool)),
			IconPath:             utils.String(d.Get("icon_path").(string)),
			IconIndex:            utils.Int64(int64(d.Get("icon_index").(int))),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualDesktopApplicationRead(d, meta)
}

func resourceVirtualDesktopApplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ApplicationName)
	d.Set("application_group_id", applicationgroup.NewApplicationGroupID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName).ID())

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("path", props.FilePath)
		d.Set("command_line_argument_policy", string(props.CommandLineSetting))
		d.Set("command_line_arguments", props.CommandLineArguments)
		d.Set("show_in_portal", props.ShowInPortal)
		d.Set("icon_path", props.IconPath)
		d.Set("icon_index", props.IconIndex)
	}

	return nil
}

func resourceVirtualDesktopApplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationsClient

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ApplicationName, applicationType)
	defer locks.UnlockByName(id.ApplicationName, applicationType)

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
