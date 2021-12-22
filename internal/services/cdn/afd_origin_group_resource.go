package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdOriginGroups() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdOriginGroupsCreate,
		Read:   resourceAfdOriginGroupsRead,
		Update: resourceAfdOriginGroupsUpdate,
		Delete: resourceAfdOriginGroupsDelete,

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
			_, err := parse.AfdOriginGroupsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ProfileID,
			},

			"session_affinity_state": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"load_balancing": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  4,
						},
						"successful_samples_required": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"additional_latency_in_ms": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  50,
						},
					},
				},
			},
			"health_probe": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "/",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.HealthProbeRequestTypeNotSet),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.HealthProbeRequestTypeGET),
								string(cdn.HealthProbeRequestTypeHEAD),
								string(cdn.HealthProbeRequestTypeNotSet),
							}, true),
						},
						"protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.ProbeProtocolHTTP),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ProbeProtocolHTTP),
								string(cdn.ProbeProtocolHTTPS),
							}, true),
						},
						"interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  100,
						},
					},
				},
			},
		},
	}
}

func resourceAfdOriginGroupsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string) // OriginGroupName
	sessionAffinityState := d.Get("session_affinity_state").(bool)

	// parse profile_id
	profileId := d.Get("profile_id").(string)
	profile, err := parse.ProfileID(profileId)
	if err != nil {
		return err
	}

	id := parse.NewAfdOriginGroupsID(profile.SubscriptionId, profile.ResourceGroup, profile.Name, name)

	loadbalancing := d.Get("load_balancing").([]interface{})
	healthprobes := d.Get("health_probe").([]interface{})

	originGroup := cdn.AFDOriginGroup{
		Name: &name,
		AFDOriginGroupProperties: &cdn.AFDOriginGroupProperties{
			LoadBalancingSettings: expandLoadBalancingSettings(loadbalancing),
			HealthProbeSettings:   expandHealthProbeSettings(healthprobes),
		},
	}

	if sessionAffinityState {
		originGroup.SessionAffinityState = cdn.EnabledStateEnabled
	} else {
		originGroup.SessionAffinityState = cdn.EnabledStateDisabled
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, originGroup)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id.ProfileName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id.OriginGroupName, err)
	}

	d.SetId(id.ID())
	return resourceAfdOriginGroupsRead(d, meta)
}

func expandLoadBalancingSettings(input []interface{}) *cdn.LoadBalancingSettingsParameters {
	if len(input) == 0 {
		return nil
	}

	loadBalancingSettings := cdn.LoadBalancingSettingsParameters{}

	config := input[0].(map[string]interface{})

	samplesize := int32(config["sample_size"].(int))
	samplesrequired := int32(config["successful_samples_required"].(int))
	latencyinms := int32(config["additional_latency_in_ms"].(int))

	loadBalancingSettings.SampleSize = &samplesize
	loadBalancingSettings.SuccessfulSamplesRequired = &samplesrequired
	loadBalancingSettings.AdditionalLatencyInMilliseconds = &latencyinms

	return &loadBalancingSettings
}

func expandHealthProbeSettings(input []interface{}) *cdn.HealthProbeParameters {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	probeinterval := int32(config["interval_in_seconds"].(int))
	probepath := config["path"].(string)
	proberequesttype := config["request_type"].(string)
	probeprotocol := config["protocol"].(string)

	healthProbeParameters := cdn.HealthProbeParameters{}
	healthProbeParameters.ProbeIntervalInSeconds = &probeinterval
	healthProbeParameters.ProbePath = &probepath

	// ProbeRequestType
	switch proberequesttype {
	case "GET":
		healthProbeParameters.ProbeRequestType = cdn.HealthProbeRequestTypeGET
	case "HEAD":
		healthProbeParameters.ProbeRequestType = cdn.HealthProbeRequestTypeHEAD
	case "NotSet":
		healthProbeParameters.ProbeRequestType = cdn.HealthProbeRequestTypeNotSet
	default:
		healthProbeParameters.ProbeRequestType = cdn.HealthProbeRequestTypeNotSet
	}

	// ProbeProtocol
	switch probeprotocol {
	case "Http":
		healthProbeParameters.ProbeProtocol = cdn.ProbeProtocolHTTP
	case "Https":
		healthProbeParameters.ProbeProtocol = cdn.ProbeProtocolHTTPS
	case "NotSet":
		healthProbeParameters.ProbeProtocol = cdn.ProbeProtocolNotSet
	default:
		healthProbeParameters.ProbeProtocol = cdn.ProbeProtocolNotSet
	}

	return &healthProbeParameters
}

func resourceAfdOriginGroupsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.AfdOriginGroupsID(d.Id())
	if err != nil {
		return err
	}

	loadbalancing := d.Get("load_balancing").([]interface{})
	healthprobes := d.Get("health_probe").([]interface{})
	sessionAffinityState := d.Get("session_affinity_state").(bool)

	properties := cdn.AFDOriginGroupUpdateParameters{
		AFDOriginGroupUpdatePropertiesParameters: &cdn.AFDOriginGroupUpdatePropertiesParameters{
			LoadBalancingSettings: expandLoadBalancingSettings(loadbalancing),
			HealthProbeSettings:   expandHealthProbeSettings(healthprobes),
		},
	}

	if d.HasChange("session_affinity_state") {
		if sessionAffinityState {
			properties.SessionAffinityState = cdn.EnabledStateEnabled
		} else {
			properties.SessionAffinityState = cdn.EnabledStateDisabled
		}
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, properties)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id.OriginGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id.OriginGroupName, err)
	}

	d.SetId(id.ID())

	return resourceAfdOriginGroupsRead(d, meta)
}

func resourceAfdOriginGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginGroupsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure CDN OriginGroups %q (Resource Group %q): %+v", id.OriginGroupName, id.ResourceGroup, err)
	}

	if healthSet := resp.HealthProbeSettings; healthSet != nil {
		d.Set("health_probe", flattenHealthProbeSettings(resp.HealthProbeSettings))
	}
	if loadSet := resp.LoadBalancingSettings; loadSet != nil {
		d.Set("load_balancing", flattenLoadBalancingSettings(resp.LoadBalancingSettings))
	}

	d.Set("name", id.OriginGroupName)
	d.Set("profile_id", parse.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if resp.SessionAffinityState == cdn.EnabledStateEnabled {
		d.Set("session_affinity_state", true)
	} else {
		d.Set("session_affinity_state", false)
	}

	return nil
}

func flattenLoadBalancingSettings(input *cdn.LoadBalancingSettingsParameters) []interface{} {
	results := make([]interface{}, 0)

	var sampleSize, additionalLatencyInMilliseconds, successfulSamplesRequired int

	if i := input; i != nil {

		if i.SampleSize != nil {
			sampleSize = int(*i.SampleSize)
		}

		if i.AdditionalLatencyInMilliseconds != nil {
			additionalLatencyInMilliseconds = int(*i.AdditionalLatencyInMilliseconds)
		}

		if i.SuccessfulSamplesRequired != nil {
			successfulSamplesRequired = int(*i.SuccessfulSamplesRequired)
		}

	}

	results = append(results, map[string]interface{}{
		"sample_size":                 sampleSize,
		"successful_samples_required": successfulSamplesRequired,
		"additional_latency_in_ms":    additionalLatencyInMilliseconds,
	})

	return results
}

func flattenHealthProbeSettings(input *cdn.HealthProbeParameters) []interface{} {
	results := make([]interface{}, 0)

	var probeIntervalInSeconds int32
	var probePath, probeProtocol, probeRequestType string

	if i := input; i != nil {

		if i.ProbeIntervalInSeconds != nil {
			probeIntervalInSeconds = *i.ProbeIntervalInSeconds
		}

		if i.ProbePath != nil {
			probePath = *i.ProbePath
		}

		switch i.ProbeProtocol {
		case cdn.ProbeProtocolHTTP:
			probeProtocol = "Http"
		case cdn.ProbeProtocolHTTPS:
			probeProtocol = "Https"
		}

		switch i.ProbeRequestType {
		case cdn.HealthProbeRequestTypeGET:
			probeRequestType = "GET"
		case cdn.HealthProbeRequestTypeHEAD:
			probeRequestType = "HEAD"
		case cdn.HealthProbeRequestTypeNotSet:
			probeRequestType = "NotSet"
		}
	}

	results = append(results, map[string]interface{}{
		"path":                probePath,
		"request_type":        probeRequestType,
		"protocol":            probeProtocol,
		"interval_in_seconds": probeIntervalInSeconds,
	})

	return results
}

func resourceAfdOriginGroupsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginGroupsID(d.Id())
	if err != nil {
		return err
	}

	originsClient := meta.(*clients.Client).Cdn.AFDOriginsClient
	origins, err := originsClient.ListByOriginGroup(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		return err
	}

	for _, o := range origins.Values() {
		future, err := originsClient.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, *o.Name)
		if err != nil {
			return fmt.Errorf("deleting %s: %+v", *o.Name, err)
		}

		if err = future.WaitForCompletionRef(ctx, originsClient.Client); err != nil {
			return fmt.Errorf("waiting for the deletion of %s: %+v", *o.Name, err)
		}
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
