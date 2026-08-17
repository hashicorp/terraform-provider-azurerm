// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginGroupCreate,
		Read:   resourceCdnFrontDoorOriginGroupRead,
		Update: resourceCdnFrontDoorOriginGroupUpdate,
		Delete: resourceCdnFrontDoorOriginGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := afdorigingroups.ParseOriginGroupID(id)
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
				ValidateFunc: afdorigingroups.ValidateProfileID,
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
								string(afdorigingroups.ProbeProtocolHTTP),
								string(afdorigingroups.ProbeProtocolHTTPS),
							}, false),
						},

						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(afdorigingroups.HealthProbeRequestTypeHEAD),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdorigingroups.HealthProbeRequestTypeGET),
								string(afdorigingroups.HealthProbeRequestTypeHEAD),
							}, false),
						},

						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 255),
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

	profile, err := afdorigingroups.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := afdorigingroups.NewOriginGroupID(profile.SubscriptionId, profile.ResourceGroupName, profile.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin_group", id.ID())
		}
	}

	sessionAffinity := afdorigingroups.EnabledStateDisabled
	if d.Get("session_affinity_enabled").(bool) {
		sessionAffinity = afdorigingroups.EnabledStateEnabled
	}

	props := afdorigingroups.AFDOriginGroup{
		Properties: &afdorigingroups.AFDOriginGroupProperties{
			HealthProbeSettings:   expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{})),
			LoadBalancingSettings: expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{})),
			SessionAffinityState:  pointer.To(sessionAffinity),
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: pointer.To(int64(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int))),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}

func resourceCdnFrontDoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OriginGroupName)
	d.Set("cdn_frontdoor_profile_id", afdorigingroups.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("health_probe", flattenCdnFrontDoorOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
				return fmt.Errorf("setting 'health_probe': %+v", err)
			}

			if err := d.Set("load_balancing", flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
				return fmt.Errorf("setting 'load_balancing': %+v", err)
			}

			d.Set("session_affinity_enabled", pointer.From(props.SessionAffinityState) == afdorigingroups.EnabledStateEnabled)
			d.Set("restore_traffic_time_to_healed_or_new_endpoint_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
		}
	}

	return nil
}

func resourceCdnFrontDoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	// The API requires that an explicit null be passed as the 'health_probe' value to disable the health probe
	// e.g. {"properties":{"healthProbeSettings":null}}
	if d.HasChange("health_probe") {
		props.HealthProbeSettings = expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{}))
	}

	if d.HasChange("load_balancing") {
		props.LoadBalancingSettings = expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{}))
	}

	if d.HasChange("restore_traffic_time_to_healed_or_new_endpoint_in_minutes") {
		props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes = pointer.To(int64(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int)))
	}

	if d.HasChange("session_affinity_enabled") {
		props.SessionAffinityState = pointer.To(afdorigingroups.EnabledStateDisabled)
		if d.Get("session_affinity_enabled").(bool) {
			props.SessionAffinityState = pointer.To(afdorigingroups.EnabledStateEnabled)
		}
	}

	if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}

func resourceCdnFrontDoorOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCdnFrontDoorOriginGroupHealthProbeParameters(input []interface{}) *afdorigingroups.HealthProbeParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &afdorigingroups.HealthProbeParameters{
		ProbeIntervalInSeconds: pointer.To(int64(v["interval_in_seconds"].(int))),
		ProbePath:              pointer.To(v["path"].(string)),
		ProbeProtocol:          pointer.ToEnum[afdorigingroups.ProbeProtocol](v["protocol"].(string)),
		ProbeRequestType:       pointer.ToEnum[afdorigingroups.HealthProbeRequestType](v["request_type"].(string)),
	}
}

func expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input []interface{}) *afdorigingroups.LoadBalancingSettingsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &afdorigingroups.LoadBalancingSettingsParameters{
		AdditionalLatencyInMilliseconds: pointer.To(int64(v["additional_latency_in_milliseconds"].(int))),
		SampleSize:                      pointer.To(int64(v["sample_size"].(int))),
		SuccessfulSamplesRequired:       pointer.To(int64(v["successful_samples_required"].(int))),
	}
}

func flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input *afdorigingroups.LoadBalancingSettingsParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"additional_latency_in_milliseconds": pointer.From(input.AdditionalLatencyInMilliseconds),
			"sample_size":                        pointer.From(input.SampleSize),
			"successful_samples_required":        pointer.From(input.SuccessfulSamplesRequired),
		},
	}
}

func flattenCdnFrontDoorOriginGroupHealthProbeParameters(input *afdorigingroups.HealthProbeParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"interval_in_seconds": pointer.From(input.ProbeIntervalInSeconds),
			"path":                pointer.From(input.ProbePath),
			"protocol":            pointer.FromEnum(input.ProbeProtocol),
			"request_type":        pointer.FromEnum(input.ProbeRequestType),
		},
	}
}
