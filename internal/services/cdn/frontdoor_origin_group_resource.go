package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorOriginGroupCreate,
		Read:   resourceFrontdoorOriginGroupRead,
		Update: resourceFrontdoorOriginGroupUpdate,
		Delete: resourceFrontdoorOriginGroupDelete,

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

			"health_probe_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"probe_interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"probe_path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"probe_protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"probe_request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"load_balancing_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"additional_latency_in_milliseconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"successful_samples_required": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"response_based_afd_origin_error_detection_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"http_error_ranges": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"begin": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
									},

									"end": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"response_based_detected_error_types": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"response_based_failover_threshold_percentage": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"session_affinity_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"traffic_restoration_time_to_healed_or_new_endpoints_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceFrontdoorOriginGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	sdkId := afdorigingroups.NewOriginGroupID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))
	id := parse.NewFrontdoorOriginGroupID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_origin_group", id.ID())
		}
	}

	sessionAffinityStateValue := afdorigingroups.EnabledState(d.Get("session_affinity_state").(string))
	props := afdorigingroups.AFDOriginGroup{
		Properties: &afdorigingroups.AFDOriginGroupProperties{
			HealthProbeSettings:                                   expandOriginGroupHealthProbeParameters(d.Get("health_probe_settings").([]interface{})),
			LoadBalancingSettings:                                 expandOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing_settings").([]interface{})),
			ResponseBasedAfdOriginErrorDetectionSettings:          expandOriginGroupResponseBasedOriginErrorDetectionParameters(d.Get("response_based_afd_origin_error_detection_settings").([]interface{})),
			SessionAffinityState:                                  &sessionAffinityStateValue,
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: utils.Int64(int64(d.Get("traffic_restoration_time_to_healed_or_new_endpoints_in_minutes").(int))),
		},
	}
	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorOriginGroupRead(d, meta)
}

func resourceFrontdoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginGroupID(d.Id())
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

	d.Set("name", id.OriginGroupName)

	d.Set("frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("deployment_status", props.DeploymentStatus)

			if err := d.Set("health_probe_settings", flattenOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
				return fmt.Errorf("setting `health_probe_settings`: %+v", err)
			}

			if err := d.Set("load_balancing_settings", flattenOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
				return fmt.Errorf("setting `load_balancing_settings`: %+v", err)
			}
			d.Set("profile_name", props.ProfileName)
			d.Set("provisioning_state", props.ProvisioningState)

			if err := d.Set("response_based_afd_origin_error_detection_settings", flattenOriginGroupResponseBasedOriginErrorDetectionParameters(props.ResponseBasedAfdOriginErrorDetectionSettings)); err != nil {
				return fmt.Errorf("setting `response_based_afd_origin_error_detection_settings`: %+v", err)
			}
			d.Set("session_affinity_state", props.SessionAffinityState)
			d.Set("traffic_restoration_time_to_healed_or_new_endpoints_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
		}
	}
	return nil
}

func resourceFrontdoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	sessionAffinityStateValue := afdorigingroups.EnabledState(d.Get("session_affinity_state").(string))
	props := afdorigingroups.AFDOriginGroupUpdateParameters{
		Properties: &afdorigingroups.AFDOriginGroupUpdatePropertiesParameters{
			HealthProbeSettings:                                   expandOriginGroupHealthProbeParameters(d.Get("health_probe_settings").([]interface{})),
			LoadBalancingSettings:                                 expandOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing_settings").([]interface{})),
			ResponseBasedAfdOriginErrorDetectionSettings:          expandOriginGroupResponseBasedOriginErrorDetectionParameters(d.Get("response_based_afd_origin_error_detection_settings").([]interface{})),
			SessionAffinityState:                                  &sessionAffinityStateValue,
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: utils.Int64(int64(d.Get("traffic_restoration_time_to_healed_or_new_endpoints_in_minutes").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *sdkId, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorOriginGroupRead(d, meta)
}

func resourceFrontdoorOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *sdkId); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandOriginGroupHealthProbeParameters(input []interface{}) *afdorigingroups.HealthProbeParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	probeProtocolValue := afdorigingroups.ProbeProtocol(v["probe_protocol"].(string))
	probeRequestTypeValue := afdorigingroups.HealthProbeRequestType(v["probe_request_type"].(string))
	return &afdorigingroups.HealthProbeParameters{
		ProbeIntervalInSeconds: utils.Int64(int64(v["probe_interval_in_seconds"].(int))),
		ProbePath:              utils.String(v["probe_path"].(string)),
		ProbeProtocol:          &probeProtocolValue,
		ProbeRequestType:       &probeRequestTypeValue,
	}
}

func expandOriginGroupLoadBalancingSettingsParameters(input []interface{}) *afdorigingroups.LoadBalancingSettingsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &afdorigingroups.LoadBalancingSettingsParameters{
		AdditionalLatencyInMilliseconds: utils.Int64(int64(v["additional_latency_in_milliseconds"].(int))),
		SampleSize:                      utils.Int64(int64(v["sample_size"].(int))),
		SuccessfulSamplesRequired:       utils.Int64(int64(v["successful_samples_required"].(int))),
	}
}

func expandOriginGroupResponseBasedOriginErrorDetectionParameters(input []interface{}) *afdorigingroups.ResponseBasedOriginErrorDetectionParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	responseBasedDetectedErrorTypesValue := afdorigingroups.ResponseBasedDetectedErrorTypes(v["response_based_detected_error_types"].(string))
	return &afdorigingroups.ResponseBasedOriginErrorDetectionParameters{
		HttpErrorRanges:                          expandOriginGroupHttpErrorRangeParametersArray(v["http_error_ranges"].([]interface{})),
		ResponseBasedDetectedErrorTypes:          &responseBasedDetectedErrorTypesValue,
		ResponseBasedFailoverThresholdPercentage: utils.Int64(int64(v["response_based_failover_threshold_percentage"].(int))),
	}
}

func expandOriginGroupHttpErrorRangeParametersArray(input []interface{}) *[]afdorigingroups.HttpErrorRangeParameters {
	results := make([]afdorigingroups.HttpErrorRangeParameters, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, afdorigingroups.HttpErrorRangeParameters{
			Begin: utils.Int64(int64(v["begin"].(int))),
			End:   utils.Int64(int64(v["end"].(int))),
		})
	}
	return &results
}

func flattenOriginGroupLoadBalancingSettingsParameters(input *afdorigingroups.LoadBalancingSettingsParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.AdditionalLatencyInMilliseconds != nil {
		result["additional_latency_in_milliseconds"] = *input.AdditionalLatencyInMilliseconds
	}

	if input.SampleSize != nil {
		result["sample_size"] = *input.SampleSize
	}

	if input.SuccessfulSamplesRequired != nil {
		result["successful_samples_required"] = *input.SuccessfulSamplesRequired
	}
	return append(results, result)
}

func flattenOriginGroupResponseBasedOriginErrorDetectionParameters(input *afdorigingroups.ResponseBasedOriginErrorDetectionParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["http_error_ranges"] = flattenOriginGroupHttpErrorRangeParametersArray(input.HttpErrorRanges)

	if input.ResponseBasedDetectedErrorTypes != nil {
		result["response_based_detected_error_types"] = *input.ResponseBasedDetectedErrorTypes
	}

	if input.ResponseBasedFailoverThresholdPercentage != nil {
		result["response_based_failover_threshold_percentage"] = *input.ResponseBasedFailoverThresholdPercentage
	}
	return append(results, result)
}

func flattenOriginGroupHttpErrorRangeParametersArray(inputs *[]afdorigingroups.HttpErrorRangeParameters) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.Begin != nil {
			result["begin"] = *input.Begin
		}

		if input.End != nil {
			result["end"] = *input.End
		}
		results = append(results, result)
	}

	return results
}

func flattenOriginGroupHealthProbeParameters(input *afdorigingroups.HealthProbeParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.ProbeIntervalInSeconds != nil {
		result["probe_interval_in_seconds"] = *input.ProbeIntervalInSeconds
	}

	if input.ProbePath != nil {
		result["probe_path"] = *input.ProbePath
	}

	if input.ProbeProtocol != nil {
		result["probe_protocol"] = *input.ProbeProtocol
	}

	if input.ProbeRequestType != nil {
		result["probe_request_type"] = *input.ProbeRequestType
	}

	return append(results, result)
}
