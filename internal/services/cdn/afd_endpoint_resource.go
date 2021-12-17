package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdEndpoints() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdEndpointsCreate,
		Read:   resourceAfdEndpointsRead,
		Update: resourceAfdEndpointsUpdate,
		Delete: resourceAfdEndpointsDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdEndpointsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ProfileID,
			},

			"origin_response_timeout_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(16, 240),
			},
		},
	}
}

func resourceAfdEndpointsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	afdEndpointsClient := meta.(*clients.Client).Cdn.AFDEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := "global" // location is always global (required)
	enabledState := d.Get("enabled").(bool)
	originResponseTimeoutInSeconds := int32(d.Get("origin_response_timeout_in_seconds").(int))

	// parse profile_id
	profileId := d.Get("profile_id").(string)
	profile, err := parse.ProfileID(profileId)
	if err != nil {
		return err
	}

	id := parse.NewAfdEndpointsID(profile.SubscriptionId, profile.ResourceGroup, profile.Name, d.Get("name").(string))
	existing, err := afdEndpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_endpoint", id.ID())
	}

	endpoint := cdn.AFDEndpoint{
		Location: &location,
		AFDEndpointProperties: &cdn.AFDEndpointProperties{
			OriginResponseTimeoutSeconds: &originResponseTimeoutInSeconds,
		},
	}

	if enabledState {
		endpoint.EnabledState = cdn.EnabledStateEnabled
	} else {
		endpoint.EnabledState = cdn.EnabledStateDisabled
	}

	future, err := afdEndpointsClient.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, afdEndpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAfdEndpointsRead(d, meta)
}

func resourceAfdEndpointsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure CDN Endpoint %q (Resource Group %q): %+v", id.AfdEndpointName, id.ResourceGroup, err)
	}

	d.Set("origin_response_timeout_in_seconds", resp.OriginResponseTimeoutSeconds)

	if resp.EnabledState == cdn.EnabledStateEnabled {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}
	d.Set("name", id.AfdEndpointName)

	return nil
}

func resourceAfdEndpointsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointsID(d.Id())
	if err != nil {
		return err
	}

	endpointUpdate := cdn.AFDEndpointUpdateParameters{}
	endpointUpdateParameters := cdn.AFDEndpointPropertiesUpdateParameters{}

	if d.HasChange("enabled") {
		enabledState := d.Get("enabled").(bool)
		if enabledState {
			endpointUpdateParameters.EnabledState = cdn.EnabledStateEnabled
		} else {
			endpointUpdateParameters.EnabledState = cdn.EnabledStateDisabled
		}
	}

	if d.HasChange("origin_response_timeout_in_seconds") {
		originResponseTimeoutInSeconds := int32(d.Get("origin_response_timeout_in_seconds").(int))
		endpointUpdateParameters.OriginResponseTimeoutSeconds = utils.Int32(originResponseTimeoutInSeconds)
	}

	endpointUpdate.AFDEndpointPropertiesUpdateParameters = &endpointUpdateParameters

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, endpointUpdate)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceAfdEndpointsRead(d, meta)
}

func resourceAfdEndpointsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointsID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
