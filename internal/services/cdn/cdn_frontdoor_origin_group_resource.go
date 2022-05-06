package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorOriginGroupCreate,
		Read:   resourceCdnFrontdoorOriginGroupRead,
		Update: resourceCdnFrontdoorOriginGroupUpdate,
		Delete: resourceCdnFrontdoorOriginGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorOriginGroupID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorProfileID,
			},

			"health_probe": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(5, 31536000),
						},

						"path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "/",
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.ProbeProtocolHTTP),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ProbeProtocolHTTP),
								string(cdn.ProbeProtocolHTTPS),
								string(cdn.ProbeProtocolNotSet), // TODO: what does this do? feels like we should remove this, the default & make it required?
							}, false),
						},

						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.HealthProbeRequestTypeHEAD),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.HealthProbeRequestTypeGET),
								string(cdn.HealthProbeRequestTypeHEAD),
								string(cdn.HealthProbeRequestTypeNotSet), // TODO: what does this do? feels like we should remove this, the default & make it required?
							}, false),
						},
					},
				},
			},

			"load_balancing": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"additional_latency_in_milliseconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      50,
							ValidateFunc: validation.IntBetween(0, 1000),
						},

						"sample_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      4,
							ValidateFunc: validation.IntBetween(0, 255),
						},

						"successful_samples_required": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntBetween(0, 255),
						},
					},
				},
			},

			"cdn_frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"session_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"restore_traffic_time_to_healed_or_new_endpoint_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(0, 50),
			},
		},
	}
}

func resourceCdnFrontdoorOriginGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontdoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorOriginGroupID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin_group", id.ID())
		}
	}

	props := cdn.AFDOriginGroup{
		AFDOriginGroupProperties: &cdn.AFDOriginGroupProperties{
			HealthProbeSettings:   expandCdnFrontdoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{})),
			LoadBalancingSettings: expandCdnFrontdoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{})),
			SessionAffinityState:  convertBoolToEnabledState(d.Get("session_affinity_enabled").(bool)),
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: utils.Int32(int32(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int))),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorOriginGroupRead(d, meta)
}

func resourceCdnFrontdoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OriginGroupName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDOriginGroupProperties; props != nil {
		if err := d.Set("health_probe", flattenCdnFrontdoorOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
			return fmt.Errorf("setting `health_probe`: %+v", err)
		}

		if err := d.Set("load_balancing", flattenCdnFrontdoorOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("setting `load_balancing`: %+v", err)
		}

		// TODO: BLOCKER - BUG: API does not return the profile name, pull it from the ID
		d.Set("cdn_frontdoor_profile_name", id.ProfileName)

		d.Set("session_affinity_enabled", convertEnabledStateToBool(&props.SessionAffinityState))
		d.Set("restore_traffic_time_to_healed_or_new_endpoint_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
	}

	return nil
}

func resourceCdnFrontdoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	params := &cdn.AFDOriginGroupUpdatePropertiesParameters{}

	if d.HasChange("health_probe") {
		params.HealthProbeSettings = expandCdnFrontdoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{}))
	}

	if d.HasChange("load_balancing") {
		params.LoadBalancingSettings = expandCdnFrontdoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{}))
	}

	if d.HasChange("restore_traffic_time_to_healed_or_new_endpoint_in_minutes") {
		params.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes = utils.Int32(int32(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int)))
	}

	if d.HasChange("session_affinity_enabled") {
		params.SessionAffinityState = convertBoolToEnabledState(d.Get("session_affinity_enabled").(bool))
	}

	payload := cdn.AFDOriginGroupUpdateParameters{
		AFDOriginGroupUpdatePropertiesParameters: params,
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontdoorOriginGroupRead(d, meta)
}

func resourceCdnFrontdoorOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}
	return nil
}

func expandCdnFrontdoorOriginGroupHealthProbeParameters(input []interface{}) *cdn.HealthProbeParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	probeProtocolValue := cdn.ProbeProtocol(v["protocol"].(string))
	probeRequestTypeValue := cdn.HealthProbeRequestType(v["request_type"].(string))
	return &cdn.HealthProbeParameters{
		ProbeIntervalInSeconds: utils.Int32(int32(v["interval_in_seconds"].(int))),
		ProbePath:              utils.String(v["path"].(string)),
		ProbeProtocol:          probeProtocolValue,
		ProbeRequestType:       probeRequestTypeValue,
	}
}

func expandCdnFrontdoorOriginGroupLoadBalancingSettingsParameters(input []interface{}) *cdn.LoadBalancingSettingsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &cdn.LoadBalancingSettingsParameters{
		AdditionalLatencyInMilliseconds: utils.Int32(int32(v["additional_latency_in_milliseconds"].(int))),
		SampleSize:                      utils.Int32(int32(v["sample_count"].(int))),
		SuccessfulSamplesRequired:       utils.Int32(int32(v["successful_samples_required"].(int))),
	}
}

func flattenCdnFrontdoorOriginGroupLoadBalancingSettingsParameters(input *cdn.LoadBalancingSettingsParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	additionalLatencyInMilliseconds := 0
	if input.AdditionalLatencyInMilliseconds != nil {
		additionalLatencyInMilliseconds = int(*input.AdditionalLatencyInMilliseconds)
	}

	sampleSize := 0
	if input.SampleSize != nil {
		sampleSize = int(*input.SampleSize)
	}

	successfulSamplesRequired := 0
	if input.SuccessfulSamplesRequired != nil {
		successfulSamplesRequired = int(*input.SuccessfulSamplesRequired)
	}
	return []interface{}{
		map[string]interface{}{
			"additional_latency_in_milliseconds": additionalLatencyInMilliseconds,
			"sample_count":                       sampleSize, // TODO: why isn't this called "sample_size"?
			"successful_samples_required":        successfulSamplesRequired,
		},
	}
}

func flattenCdnFrontdoorOriginGroupHealthProbeParameters(input *cdn.HealthProbeParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	intervalInSeconds := 0
	if input.ProbeIntervalInSeconds != nil {
		intervalInSeconds = int(*input.ProbeIntervalInSeconds)
	}

	path := ""
	if input.ProbePath != nil {
		path = *input.ProbePath
	}

	protocol := ""
	if input.ProbeProtocol != "" {
		protocol = string(input.ProbeProtocol)
	}

	requestType := ""
	if input.ProbeRequestType != "" {
		requestType = string(input.ProbeRequestType)
	}

	return []interface{}{
		map[string]interface{}{
			"interval_in_seconds": intervalInSeconds,
			"path":                path,
			"protocol":            protocol,
			"request_type":        requestType,
		},
	}
}
