package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var hostpoolResourceType = "azurerm_virtual_desktop_host_pool"

func resourceVirtualDesktopHostPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopHostPoolCreate,
		Read:   resourceVirtualDesktopHostPoolRead,
		Update: resourceVirtualDesktopHostPoolUpdate,
		Delete: resourceVirtualDesktopHostPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HostPoolID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HostPoolV0ToV1{},
		}),

		Schema: func() map[string]*pluginsdk.Schema {
			s := map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"location": azure.SchemaLocation(),

				"resource_group_name": azure.SchemaResourceGroupName(),

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(desktopvirtualization.HostPoolTypePersonal),
						string(desktopvirtualization.HostPoolTypePooled),
					}, false),
				},

				"load_balancer_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
						string(desktopvirtualization.LoadBalancerTypeDepthFirst),
						string(desktopvirtualization.LoadBalancerTypePersistent),
					}, false),
				},

				"friendly_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(1, 64),
				},

				"description": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(1, 512),
				},

				"validate_environment": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"custom_rdp_properties": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"personal_desktop_assignment_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(desktopvirtualization.PersonalDesktopAssignmentTypeAutomatic),
						string(desktopvirtualization.PersonalDesktopAssignmentTypeDirect),
					}, false),
				},

				"maximum_sessions_allowed": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      999999,
					ValidateFunc: validation.IntBetween(0, 999999),
				},

				"start_vm_on_connect": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"preferred_app_group_type": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Preferred App Group type to display",
					ValidateFunc: validation.StringInSlice([]string{
						string(desktopvirtualization.PreferredAppGroupTypeDesktop),
						string(desktopvirtualization.PreferredAppGroupTypeNone),
						string(desktopvirtualization.PreferredAppGroupTypeRailApplications),
					}, false),
					Default: string(desktopvirtualization.PreferredAppGroupTypeDesktop),
				},

				"tags": tags.Schema(),
			}

			if !features.ThreePointOhBeta() {
				s["registration_info"] = &schema.Schema{
					Type:        pluginsdk.TypeList,
					Optional:    true,
					Description: "This block is now non-functional and will be removed in version 3.0 of the Azure Provider - use the `azurerm_virtual_desktop_host_pool_registration_info` resource instead.",
					MaxItems:    1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"expiration_date": {
								Type:       pluginsdk.TypeString,
								Optional:   true,
								Deprecated: "This field is now non-functional and will be removed in version 3.0 of the Azure Provider - use the `azurerm_virtual_desktop_host_pool_registration_info` resource instead.",
							},

							"reset_token": {
								Type:     pluginsdk.TypeBool,
								Computed: true,
							},

							"token": {
								Type:      pluginsdk.TypeString,
								Sensitive: true,
								Computed:  true,
							},
						},
					},
				}
			}

			return s
		}(),
	}
}

func resourceVirtualDesktopHostPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewHostPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_host_pool", id.ID())
	}

	context := desktopvirtualization.HostPool{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		HostPoolProperties: &desktopvirtualization.HostPoolProperties{
			HostPoolType:                  desktopvirtualization.HostPoolType(d.Get("type").(string)),
			FriendlyName:                  utils.String(d.Get("friendly_name").(string)),
			Description:                   utils.String(d.Get("description").(string)),
			ValidationEnvironment:         utils.Bool(d.Get("validate_environment").(bool)),
			CustomRdpProperty:             utils.String(d.Get("custom_rdp_properties").(string)),
			MaxSessionLimit:               utils.Int32(int32(d.Get("maximum_sessions_allowed").(int))),
			StartVMOnConnect:              utils.Bool(d.Get("start_vm_on_connect").(bool)),
			LoadBalancerType:              desktopvirtualization.LoadBalancerType(d.Get("load_balancer_type").(string)),
			PersonalDesktopAssignmentType: desktopvirtualization.PersonalDesktopAssignmentType(d.Get("personal_desktop_assignment_type").(string)),
			PreferredAppGroupType:         desktopvirtualization.PreferredAppGroupType(d.Get("preferred_app_group_type").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, context); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func resourceVirtualDesktopHostPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostPoolID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, hostpoolResourceType)
	defer locks.UnlockByName(id.Name, hostpoolResourceType)

	update := &desktopvirtualization.HostPoolPatch{}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChanges("custom_rdp_properties", "description", "friendly_name", "maximum_sessions_allowed", "preferred_app_group_type", "start_vm_on_connect", "validate_environment") {
		update.HostPoolPatchProperties = &desktopvirtualization.HostPoolPatchProperties{}

		if d.HasChange("custom_rdp_properties") {
			update.HostPoolPatchProperties.CustomRdpProperty = utils.String(d.Get("custom_rdp_properties").(string))
		}

		if d.HasChange("description") {
			update.HostPoolPatchProperties.Description = utils.String(d.Get("description").(string))
		}

		if d.HasChange("friendly_name") {
			update.HostPoolPatchProperties.FriendlyName = utils.String(d.Get("friendly_name").(string))
		}

		if d.HasChange("maximum_sessions_allowed") {
			update.HostPoolPatchProperties.MaxSessionLimit = utils.Int32(int32(d.Get("maximum_sessions_allowed").(int)))
		}

		if d.HasChange("preferred_app_group_type") {
			update.HostPoolPatchProperties.PreferredAppGroupType = desktopvirtualization.PreferredAppGroupType(d.Get("preferred_app_group_type").(string))
		}

		if d.HasChange("start_vm_on_connect") {
			update.HostPoolPatchProperties.StartVMOnConnect = utils.Bool(d.Get("start_vm_on_connect").(bool))
		}

		if d.HasChange("validate_environment") {
			update.HostPoolPatchProperties.ValidationEnvironment = utils.Bool(d.Get("validate_environment").(bool))
		}
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func resourceVirtualDesktopHostPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.HostPoolProperties; props != nil {
		maxSessionLimit := 0
		if props.MaxSessionLimit != nil {
			maxSessionLimit = int(*props.MaxSessionLimit)
		}

		d.Set("custom_rdp_properties", props.CustomRdpProperty)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("maximum_sessions_allowed", maxSessionLimit)
		d.Set("load_balancer_type", string(props.LoadBalancerType))
		d.Set("personal_desktop_assignment_type", string(props.PersonalDesktopAssignmentType))
		d.Set("preferred_app_group_type", string(props.PreferredAppGroupType))
		d.Set("start_vm_on_connect", props.StartVMOnConnect)
		d.Set("type", string(props.HostPoolType))
		d.Set("validate_environment", props.ValidationEnvironment)

		if !features.ThreePointOhBeta() {
			d.Set("registration_info", []interface{}{})
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualDesktopHostPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HostPoolID(d.Id())

	locks.ByName(id.Name, hostpoolResourceType)
	defer locks.UnlockByName(id.Name, hostpoolResourceType)

	if err != nil {
		return err
	}

	forceDeleteSessionHost := utils.Bool(true)
	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name, forceDeleteSessionHost); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
