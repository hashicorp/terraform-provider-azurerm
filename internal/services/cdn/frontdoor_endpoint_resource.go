package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"origin_response_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  120,
			},

			"frontdoor_profile_name": {
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

	id := afdendpoints.NewAfdEndpointID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_endpoint", id.ID())
		}
	}

	location := azure.NormalizeLocation("global")

	props := afdendpoints.AFDEndpoint{
		Name:     utils.String(d.Get("name").(string)),
		Location: location,
		Properties: &afdendpoints.AFDEndpointProperties{
			EnabledState: ConvertBoolToEndpointsEnabledState(d.Get("enabled").(bool)),
			// Bug in API, the OriginResponseTimeoutSeconds is not currently exposed
			// OriginResponseTimeoutSeconds: utils.Int64(int64(d.Get("origin_response_timeout_seconds").(int))),

		},
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceFrontdoorEndpointRead(d, meta)
}

func resourceFrontdoorEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	profileId := profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.EndpointName)
	d.Set("frontdoor_profile_id", profileId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("enabled", ConvertEndpointsEnabledStateToBool(props.EnabledState))
			d.Set("host_name", props.HostName)

			// BUG: Profile Name is not being returned by that API pull it from the ID
			d.Set("frontdoor_profile_name", id.ProfileName)

			// BUG API does not currently expose this field so temporarily hardcoding to default value
			// d.Set("origin_response_timeout_seconds", props.OriginResponseTimeoutSeconds)
			d.Set("origin_response_timeout_seconds", 60)
		}

		if err := tags.FlattenAndSet(d, ConvertFrontdoorTags(model.Tags)); err != nil {
			return err
		}
	}

	d.SetId(id.ID())

	return nil
}

func resourceFrontdoorEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	props := afdendpoints.AFDEndpointUpdateParameters{
		Properties: &afdendpoints.AFDEndpointPropertiesUpdateParameters{
			EnabledState: ConvertBoolToEndpointsEnabledState(d.Get("enabled").(bool)),
			ProfileName:  utils.String(profileId.ProfileName),
			// OriginResponseTimeoutSeconds: utils.Int64(int64(d.Get("origin_response_timeout_seconds").(int))),
		},
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
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
