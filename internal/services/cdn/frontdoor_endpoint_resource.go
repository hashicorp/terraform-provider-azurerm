package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontdoorEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorEndpointCreate,
		Read:   resourceFrontdoorEndpointRead,
		Update: resourceFrontdoorEndpointUpdate,
		Delete: resourceFrontdoorEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := afdendpoints.ParseAfdEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: profiles.ValidateProfileID,
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"origin_response_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceFrontdoorEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	sdkId := afdendpoints.NewAfdEndpointID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))
	id := parse.NewFrontdoorEndpointID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_endpoint", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location"))
	enabledStateValue := afdendpoints.EnabledState(d.Get("enabled_state").(string))
	props := afdendpoints.AFDEndpoint{
		Location: location,
		Properties: &afdendpoints.AFDEndpointProperties{
			EnabledState: &enabledStateValue,
			// OriginResponseTimeoutSeconds: utils.Int64(int64(d.Get("origin_response_timeout_seconds").(int))),
		},
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorEndpointRead(d, meta)
}

func resourceFrontdoorEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *sdkId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AfdEndpointName)

	d.Set("frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("enabled_state", props.EnabledState)
			d.Set("host_name", props.HostName)
			// d.Set("origin_response_timeout_seconds", props.OriginResponseTimeoutSeconds)
			d.Set("profile_name", props.ProfileName)
			d.Set("provisioning_state", props.ProvisioningState)
		}

		// TODO: Fix Tag Type
		// if err := tags.FlattenAndSet(d, tagsHelper.Flatten(model.Tags)); err != nil {
		// 	return err
		// }
	}
	return nil
}

func resourceFrontdoorEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	enabledStateValue := afdendpoints.EnabledState(d.Get("enabled_state").(string))
	props := afdendpoints.AFDEndpointUpdateParameters{
		Properties: &afdendpoints.AFDEndpointPropertiesUpdateParameters{
			EnabledState: &enabledStateValue,
			// OriginResponseTimeoutSeconds: utils.Int64(int64(d.Get("origin_response_timeout_seconds").(int))),
		},
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.UpdateThenPoll(ctx, *sdkId, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorEndpointRead(d, meta)
}

func resourceFrontdoorEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *sdkId); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
