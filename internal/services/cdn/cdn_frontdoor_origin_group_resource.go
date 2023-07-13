// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginGroupCreate,
		Read:   resourceCdnFrontDoorOriginGroupRead,
		Update: resourceCdnFrontDoorOriginGroupUpdate,
		Delete: resourceCdnFrontDoorOriginGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorOriginGroupID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorOriginGroupName,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorProfileID,
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

						"sample_size": {
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

			// Optional
			"health_probe": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ProbeProtocolHTTP),
								string(cdn.ProbeProtocolHTTPS),
							}, false),
						},

						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.HealthProbeRequestTypeHEAD),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.HealthProbeRequestTypeGET),
								string(cdn.HealthProbeRequestTypeHEAD),
							}, false),
						},

						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 31536000),
						},

						"path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "/",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
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

func resourceCdnFrontDoorOriginGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profile, err := parse.FrontDoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorOriginGroupID(profile.SubscriptionId, profile.ResourceGroup, profile.ProfileName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin_group", id.ID())
	}

	props := cdn.AFDOriginGroup{
		AFDOriginGroupProperties: &cdn.AFDOriginGroupProperties{
			HealthProbeSettings:   expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{})),
			LoadBalancingSettings: expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{})),
			SessionAffinityState:  expandEnabledBool(d.Get("session_affinity_enabled").(bool)),
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
	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}

func resourceCdnFrontDoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginGroupID(d.Id())
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
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDOriginGroupProperties; props != nil {
		if err := d.Set("health_probe", flattenCdnFrontDoorOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
			return fmt.Errorf("setting 'health_probe': %+v", err)
		}

		if err := d.Set("load_balancing", flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
			return fmt.Errorf("setting 'load_balancing': %+v", err)
		}

		d.Set("session_affinity_enabled", flattenEnabledBool(props.SessionAffinityState))
		d.Set("restore_traffic_time_to_healed_or_new_endpoint_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
	}

	return nil
}

func resourceCdnFrontDoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	workaroundClient := azuresdkhacks.NewCdnFrontDoorOriginGroupsWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	params := &azuresdkhacks.AFDOriginGroupUpdatePropertiesParameters{}

	// The API requires that an explicit null be passed as the 'health_probe' value to disable the health probe
	// e.g. {"properties":{"healthProbeSettings":null}}
	if d.HasChange("health_probe") {
		params.HealthProbeSettings = expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{}))
	}

	if d.HasChange("load_balancing") {
		params.LoadBalancingSettings = expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{}))
	}

	if d.HasChange("restore_traffic_time_to_healed_or_new_endpoint_in_minutes") {
		params.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes = utils.Int32(int32(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int)))
	}

	if d.HasChange("session_affinity_enabled") {
		params.SessionAffinityState = expandEnabledBool(d.Get("session_affinity_enabled").(bool))
	}

	payload := &azuresdkhacks.AFDOriginGroupUpdateParameters{
		AFDOriginGroupUpdatePropertiesParameters: params,
	}

	future, err := workaroundClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, *payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}

func resourceCdnFrontDoorOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginGroupID(d.Id())
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

func expandCdnFrontDoorOriginGroupHealthProbeParameters(input []interface{}) *cdn.HealthProbeParameters {
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

func expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input []interface{}) *cdn.LoadBalancingSettingsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &cdn.LoadBalancingSettingsParameters{
		AdditionalLatencyInMilliseconds: utils.Int32(int32(v["additional_latency_in_milliseconds"].(int))),
		SampleSize:                      utils.Int32(int32(v["sample_size"].(int))),
		SuccessfulSamplesRequired:       utils.Int32(int32(v["successful_samples_required"].(int))),
	}
}

func flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input *cdn.LoadBalancingSettingsParameters) []interface{} {
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
			"sample_size":                        sampleSize,
			"successful_samples_required":        successfulSamplesRequired,
		},
	}
}

func flattenCdnFrontDoorOriginGroupHealthProbeParameters(input *cdn.HealthProbeParameters) []interface{} {
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
