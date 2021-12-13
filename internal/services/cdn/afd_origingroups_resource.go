package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			_, err := parse.OriginGroupsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"session_affinity_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
						},
						"successful_samples_required": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
						"additional_latency_in_ms": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
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
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "/",
						},
						"request_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "HEAD",
						},
						"protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "Https",
						},
						"interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  240,
						},
					},
				},
			},
		},
	}
}

func resourceAfdOriginGroupsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string) // OriginGroupName

	id := parse.NewOriginGroupsID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), name)

	loadbalancing := d.Get("load_balancing").([]interface{})
	healthprobes := d.Get("health_probe").([]interface{})

	originGroup := cdn.AFDOriginGroup{
		Name: &name,
		AFDOriginGroupProperties: &cdn.AFDOriginGroupProperties{
			LoadBalancingSettings: expandLoadBalancingSettings(loadbalancing),
			HealthProbeSettings:   expandHealthProbeSettings(healthprobes),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, name, originGroup)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id.ProfileName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id.ProfileName, err)
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

	healthProbeParameters := cdn.HealthProbeParameters{}

	config := input[0].(map[string]interface{})

	probeinterval := int32(config["interval_in_seconds"].(int))

	probepath := config["path"].(string)

	healthProbeParameters.ProbeIntervalInSeconds = &probeinterval
	healthProbeParameters.ProbePath = &probepath

	// ProbeRequestType
	proberequesttype := config["request_type"].(string)
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
	probeprotocol := config["protocol"].(string)
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
	return nil
}

func resourceAfdOriginGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	name := d.Get("name").(string)
	defer cancel()

	id, err := parse.OriginGroupsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure CDN OriginGroups %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("profile_name", id.ProfileName)

	return nil
}

func resourceAfdOriginGroupsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	name := d.Get("name").(string)
	defer cancel()

	id, err := parse.OriginGroupsID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
