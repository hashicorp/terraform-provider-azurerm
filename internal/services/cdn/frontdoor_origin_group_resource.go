package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			_, err := afdorigingroups.ParseOriginGroupID(id)
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

			"health_probe": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      240,
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
							Default:  string(afdorigingroups.ProbeProtocolHttps),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdorigingroups.ProbeProtocolHttp),
								string(afdorigingroups.ProbeProtocolHttps),
								string(afdorigingroups.ProbeProtocolNotSet),
							}, false),
						},

						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(afdorigingroups.HealthProbeRequestTypeGET),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdorigingroups.HealthProbeRequestTypeGET),
								string(afdorigingroups.HealthProbeRequestTypeHEAD),
								string(afdorigingroups.HealthProbeRequestTypeNotSet),
							}, false),
						},
					},
				},
			},

			"load_balancing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"additional_latency_in_milliseconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1000),
						},

						"sample_size": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      16,
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

			"frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"response_based_origin_error_detection": {
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
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      300,
										ValidateFunc: validation.IntBetween(100, 999),
									},

									"end": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      599,
										ValidateFunc: validation.IntBetween(100, 999),
									},
								},
							},
						},

						"detected_error_types": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(afdorigingroups.ResponseBasedDetectedErrorTypesTcpAndHttpErrors),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdorigingroups.ResponseBasedDetectedErrorTypesNone),
								string(afdorigingroups.ResponseBasedDetectedErrorTypesTcpAndHttpErrors),
								string(afdorigingroups.ResponseBasedDetectedErrorTypesTcpErrorsOnly),
							}, false),
						},

						"failover_threshold_percentage": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},

			"session_affinity": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"restore_traffic_or_new_endpoints_time": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(0, 50),
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

	id := afdorigingroups.NewOriginGroupID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_origin_group", id.ID())
		}
	}

	props := afdorigingroups.AFDOriginGroup{
		Properties: &afdorigingroups.AFDOriginGroupProperties{
			HealthProbeSettings:                                   expandOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{})),
			LoadBalancingSettings:                                 expandOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{})),
			ResponseBasedAfdOriginErrorDetectionSettings:          expandOriginGroupResponseBasedOriginErrorDetectionParameters(d.Get("response_based_origin_error_detection").([]interface{})),
			SessionAffinityState:                                  ConvertBoolToOriginGroupsEnabledState(d.Get("session_affinity").(bool)),
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: utils.Int64(int64(d.Get("restore_traffic_or_new_endpoints_time").(int))),
		},
	}
	if err := client.CreateThenPoll(ctx, id, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorOriginGroupRead(d, meta)
}

func resourceFrontdoorOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("deployment_status", props.DeploymentStatus)

			if err := d.Set("health_probe", flattenOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
				return fmt.Errorf("setting `health_probe`: %+v", err)
			}

			if err := d.Set("load_balancing", flattenOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
				return fmt.Errorf("setting `load_balancing`: %+v", err)
			}

			if err := d.Set("response_based_origin_error_detection", flattenOriginGroupResponseBasedOriginErrorDetectionParameters(props.ResponseBasedAfdOriginErrorDetectionSettings)); err != nil {
				return fmt.Errorf("setting `response_based_origin_error_detection`: %+v", err)
			}

			d.Set("frontdoor_profile_name", props.ProfileName)
			d.Set("session_affinity", ConvertOriginGroupsEnabledStateToBool(props.SessionAffinityState))
			d.Set("restore_traffic_or_new_endpoints_time", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
		}
	}
	return nil
}

func resourceFrontdoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	props := afdorigingroups.AFDOriginGroupUpdateParameters{
		Properties: &afdorigingroups.AFDOriginGroupUpdatePropertiesParameters{
			HealthProbeSettings:                                   expandOriginGroupHealthProbeParameters(d.Get("health_probe").([]interface{})),
			LoadBalancingSettings:                                 expandOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]interface{})),
			ResponseBasedAfdOriginErrorDetectionSettings:          expandOriginGroupResponseBasedOriginErrorDetectionParameters(d.Get("response_based_origin_error_detection").([]interface{})),
			SessionAffinityState:                                  ConvertBoolToOriginGroupsEnabledState(d.Get("session_affinity").(bool)),
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: utils.Int64(int64(d.Get("restore_traffic_or_new_endpoints_time").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorOriginGroupRead(d, meta)
}

func resourceFrontdoorOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandOriginGroupHealthProbeParameters(input []interface{}) *afdorigingroups.HealthProbeParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	probeProtocolValue := afdorigingroups.ProbeProtocol(v["protocol"].(string))
	probeRequestTypeValue := afdorigingroups.HealthProbeRequestType(v["request_type"].(string))
	return &afdorigingroups.HealthProbeParameters{
		ProbeIntervalInSeconds: utils.Int64(int64(v["interval_in_seconds"].(int))),
		ProbePath:              utils.String(v["path"].(string)),
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

	responseBasedDetectedErrorTypesValue := afdorigingroups.ResponseBasedDetectedErrorTypes(v["detected_error_types"].(string))
	return &afdorigingroups.ResponseBasedOriginErrorDetectionParameters{
		HttpErrorRanges:                          expandOriginGroupHttpErrorRangeParametersArray(v["http_error_ranges"].([]interface{})),
		ResponseBasedDetectedErrorTypes:          &responseBasedDetectedErrorTypesValue,
		ResponseBasedFailoverThresholdPercentage: utils.Int64(int64(v["failover_threshold_percentage"].(int))),
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
		result["detected_error_types"] = *input.ResponseBasedDetectedErrorTypes
	}

	if input.ResponseBasedFailoverThresholdPercentage != nil {
		result["failover_threshold_percentage"] = *input.ResponseBasedFailoverThresholdPercentage
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
		result["interval_in_seconds"] = *input.ProbeIntervalInSeconds
	}

	if input.ProbePath != nil {
		result["path"] = *input.ProbePath
	}

	if input.ProbeProtocol != nil {
		result["protocol"] = *input.ProbeProtocol
	}

	if input.ProbeRequestType != nil {
		result["request_type"] = *input.ProbeRequestType
	}

	return append(results, result)
}
