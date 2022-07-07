package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var hostPoolResourceType = "azurerm_virtual_desktop_host_pool"

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
			_, err := hostpool.ParseHostPoolID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HostPoolV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
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
					string(hostpool.HostPoolTypePersonal),
					string(hostpool.HostPoolTypePooled),
				}, false),
			},

			"load_balancer_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(hostpool.LoadBalancerTypeBreadthFirst),
					string(hostpool.LoadBalancerTypeDepthFirst),
					string(hostpool.LoadBalancerTypePersistent),
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
					string(hostpool.PersonalDesktopAssignmentTypeAutomatic),
					string(hostpool.PersonalDesktopAssignmentTypeDirect),
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
					string(hostpool.PreferredAppGroupTypeDesktop),
					string(hostpool.PreferredAppGroupTypeNone),
					string(hostpool.PreferredAppGroupTypeRailApplications),
				}, false),
				Default: string(hostpool.PreferredAppGroupTypeDesktop),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualDesktopHostPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hostpool.NewHostPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_host_pool", id.ID())
	}

	personalDesktopAssignmentType := hostpool.PersonalDesktopAssignmentType(d.Get("personal_desktop_assignment_type").(string))
	payload := hostpool.HostPool{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: hostpool.HostPoolProperties{
			HostPoolType:                  hostpool.HostPoolType(d.Get("type").(string)),
			FriendlyName:                  utils.String(d.Get("friendly_name").(string)),
			Description:                   utils.String(d.Get("description").(string)),
			ValidationEnvironment:         utils.Bool(d.Get("validate_environment").(bool)),
			CustomRdpProperty:             utils.String(d.Get("custom_rdp_properties").(string)),
			MaxSessionLimit:               utils.Int64(int64(d.Get("maximum_sessions_allowed").(int))),
			StartVMOnConnect:              utils.Bool(d.Get("start_vm_on_connect").(bool)),
			LoadBalancerType:              hostpool.LoadBalancerType(d.Get("load_balancer_type").(string)),
			PersonalDesktopAssignmentType: &personalDesktopAssignmentType,
			PreferredAppGroupType:         hostpool.PreferredAppGroupType(d.Get("preferred_app_group_type").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func resourceVirtualDesktopHostPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hostpool.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPoolName, hostPoolResourceType)

	payload := hostpool.HostPoolPatch{}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChanges("custom_rdp_properties", "description", "friendly_name", "maximum_sessions_allowed", "preferred_app_group_type", "start_vm_on_connect", "validate_environment") {
		payload.Properties = &hostpool.HostPoolPatchProperties{}

		if d.HasChange("custom_rdp_properties") {
			payload.Properties.CustomRdpProperty = utils.String(d.Get("custom_rdp_properties").(string))
		}

		if d.HasChange("description") {
			payload.Properties.Description = utils.String(d.Get("description").(string))
		}

		if d.HasChange("friendly_name") {
			payload.Properties.FriendlyName = utils.String(d.Get("friendly_name").(string))
		}

		if d.HasChange("maximum_sessions_allowed") {
			payload.Properties.MaxSessionLimit = utils.Int64(int64(d.Get("maximum_sessions_allowed").(int)))
		}

		if d.HasChange("preferred_app_group_type") {
			preferredAppGroupType := hostpool.PreferredAppGroupType(d.Get("preferred_app_group_type").(string))
			payload.Properties.PreferredAppGroupType = &preferredAppGroupType
		}

		if d.HasChange("start_vm_on_connect") {
			payload.Properties.StartVMOnConnect = utils.Bool(d.Get("start_vm_on_connect").(bool))
		}

		if d.HasChange("validate_environment") {
			payload.Properties.ValidationEnvironment = utils.Bool(d.Get("validate_environment").(bool))
		}
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func resourceVirtualDesktopHostPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hostpool.ParseHostPoolID(d.Id())
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

	d.Set("name", id.HostPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		props := model.Properties
		maxSessionLimit := 0
		if props.MaxSessionLimit != nil {
			maxSessionLimit = int(*props.MaxSessionLimit)
		}

		d.Set("custom_rdp_properties", props.CustomRdpProperty)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("maximum_sessions_allowed", maxSessionLimit)
		d.Set("load_balancer_type", string(props.LoadBalancerType))
		personalDesktopAssignmentType := ""
		if props.PersonalDesktopAssignmentType != nil {
			personalDesktopAssignmentType = string(*props.PersonalDesktopAssignmentType)
		}
		d.Set("personal_desktop_assignment_type", personalDesktopAssignmentType)
		d.Set("preferred_app_group_type", string(props.PreferredAppGroupType))
		d.Set("start_vm_on_connect", props.StartVMOnConnect)
		d.Set("type", string(props.HostPoolType))
		d.Set("validate_environment", props.ValidationEnvironment)
	}

	return nil
}

func resourceVirtualDesktopHostPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hostpool.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPoolName, hostPoolResourceType)

	options := hostpool.DeleteOperationOptions{
		Force: utils.Bool(true),
	}
	if _, err = client.Delete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
