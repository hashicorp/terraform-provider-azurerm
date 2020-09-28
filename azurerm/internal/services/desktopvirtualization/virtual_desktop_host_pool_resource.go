package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2019-01-23-preview/desktopvirtualization"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualDesktopHostPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualDesktopHostPoolCreateUpdate,
		Read:   resourceArmVirtualDesktopHostPoolRead,
		Update: resourceArmVirtualDesktopHostPoolCreateUpdate,
		Delete: resourceArmVirtualDesktopHostPoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.VirtualDesktopHostPoolID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevSpaceName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Personal",
					"Shared",
				}, false),
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

			"validate_environment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"load_balancer_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BreadthFirst",
					"DepthFirst",
				}, false),
			},

			"personal_desktop_assignment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Automatic",
					"Direct",
				}, false),
			},

			"maximum_sessions_allowed": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 999999),
			},

			"preferred_app_group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Preferred App Group type to display",
				ValidateFunc: validation.StringInSlice([]string{
					"Desktop",
					"None",
					"RailApplications",
				}, false),
				Default: "Desktop",
			},

			"registration_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_date": {
							Type:         schema.TypeString,
							ValidateFunc: validation.IsRFC3339Time,
							Required:     true,
						},

						"reset_token": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualDesktopHostPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
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
			return tf.ImportAsExistsError("azurerm_virtual_host_pool", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	hostPoolType := d.Get("type").(string)
	friendlyName := d.Get("friendly_name").(string)
	description := d.Get("description").(string)
	ValidationEnvironment := d.Get("validate_environment").(bool)
	LoadBalancerType := d.Get("load_balancer_type").(string)
	PersonalDesktopAssignmentType := d.Get("personal_desktop_assignment_type").(string)
	PreferredAppGroupType := d.Get("preferred_app_group_type").(string)
	maxSessionLimit := int32(d.Get("maximum_sessions_allowed").(int))

	context := desktopvirtualization.HostPool{
		Location: &location,
		Tags:     tags.Expand(t),
		HostPoolProperties: &desktopvirtualization.HostPoolProperties{
			HostPoolType:                  desktopvirtualization.HostPoolType(hostPoolType),
			FriendlyName:                  &friendlyName,
			Description:                   &description,
			ValidationEnvironment:         &ValidationEnvironment,
			MaxSessionLimit:               &maxSessionLimit,
			LoadBalancerType:              desktopvirtualization.LoadBalancerType(LoadBalancerType),
			PersonalDesktopAssignmentType: desktopvirtualization.PersonalDesktopAssignmentType(PersonalDesktopAssignmentType),
			PreferredAppGroupType:         desktopvirtualization.PreferredAppGroupType(PreferredAppGroupType),
			RegistrationInfo:              expandVirtualDesktopHostPoolRegistrationInfo(d),
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
		return fmt.Errorf("Reading Virtual Desktop Host Pool %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*result.ID)

	return resourceArmVirtualDesktopHostPoolRead(d, meta)
}

func resourceArmVirtualDesktopHostPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualDesktopHostPoolID(d.Id())
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

	if props := resp.HostPoolProperties; props != nil {
		if hstplt := props.HostPoolType; hstplt != nil {
			d.Set("type", hstplt)
		}

		if fn := props.FriendlyName; fn != nil {
			d.Set("friendly_name", fn)
		}

		if desc := props.Description; desc != nil {
			d.Set("description", desc)
		}

		if ve := props.ValidationEnvironment; ve != nil {
			d.Set("validate_environment", ve)
		}

		if mxsl := props.MaxSessionLimit; mxsl != nil {
			d.Set("maximum_sessions_allowed", mxsl)
		}

		if lbt := props.LoadBalancerType; lbt != nil {
			d.Set("load_balancer_type", lbt)
		}

		if pdat := props.PersonalDesktopAssignmentType; pdat != nil {
			d.Set("personal_desktop_assignment_type", pdat)
		}

		if pagt := props.PreferredAppGroupType; pagt != nil {
			d.Set("preferred_app_group_type", pagt)
		}

		if err := d.Set("registration_info", flattenVirtualDesktopHostPoolRegistrationInfo(props.RegistrationInfo)); err != nil {
			return fmt.Errorf("Setting `registration_info`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualDesktopHostPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualDesktopHostPoolID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name, utils.Bool(true)); err != nil {
		return fmt.Errorf("Deleting Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandVirtualDesktopHostPoolRegistrationInfo(d *schema.ResourceData) *desktopvirtualization.RegistrationInfo {
	ri := d.Get("registration_info").([]interface{})
	if len(ri) == 0 {
		return nil
	}

	v := ri[0].(map[string]interface{})

	expdt, _ := date.ParseTime(time.RFC3339, v["expiration_date"].(string))

	configuration := desktopvirtualization.RegistrationInfo{
		ExpirationTime: &date.Time{
			Time: expdt,
		},
		ResetToken: utils.Bool(true),
	}

	return &configuration
}

func flattenVirtualDesktopHostPoolRegistrationInfo(input *desktopvirtualization.RegistrationInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.ExpirationTime == nil || input.Token == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"expiration_date": input.ExpirationTime.Format(time.RFC3339),
			"token":           *input.Token,
		},
	}
}
