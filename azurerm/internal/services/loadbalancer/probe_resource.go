package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerProbe() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerProbeCreateUpdate,
		Read:   resourceArmLoadBalancerProbeRead,
		Update: resourceArmLoadBalancerProbeCreateUpdate,
		Delete: resourceArmLoadBalancerProbeDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerProbeID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
			return &lbId, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadBalancerValidate.LoadBalancerID,
			},

			"protocol": {
				Type:             pluginsdk.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ProbeProtocolHTTP),
					string(network.ProbeProtocolHTTPS),
					string(network.ProbeProtocolTCP),
				}, true),
			},

			"port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"request_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"interval_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(5),
			},

			"number_of_probes": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  2,
			},

			"load_balancer_rules": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: pluginsdk.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerProbeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	loadBalancerIDRaw := loadBalancerId.ID()
	id := parse.NewLoadBalancerProbeID(subscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, d.Get("name").(string))
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

	loadBalancer, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing Proe %q from state", id.LoadBalancerName, id.ProbeName)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	newProbe := expandAzureRmLoadBalancerProbe(d)
	probes := append(*loadBalancer.LoadBalancerPropertiesFormat.Probes, *newProbe)

	existingProbe, existingProbeIndex, exists := FindLoadBalancerProbeByName(&loadBalancer, id.ProbeName)
	if exists {
		if id.ProbeName == *existingProbe.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_probe", *existingProbe.ID)
			}

			// this probe is being updated/reapplied remove old copy from the slice
			probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &probes

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, loadBalancer)
	if err != nil {
		return fmt.Errorf("updating Load Balancer %q (Resource Group %q) for Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q) for Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerProbeRead(d, meta)
}

func resourceArmLoadBalancerProbeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerProbeID(d.Id())
	if err != nil {
		return err
	}

	loadBalancer, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	config, _, exists := FindLoadBalancerProbeByName(&loadBalancer, id.ProbeName)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Probe %q not found. Removing from state", id.ProbeName)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := config.ProbePropertiesFormat; props != nil {
		intervalInSeconds := 0
		if props.IntervalInSeconds != nil {
			intervalInSeconds = int(*props.IntervalInSeconds)
		}
		d.Set("interval_in_seconds", intervalInSeconds)

		numberOfProbes := 0
		if props.NumberOfProbes != nil {
			numberOfProbes = int(*props.NumberOfProbes)
		}
		d.Set("number_of_probes", numberOfProbes)

		port := 0
		if props.Port != nil {
			port = int(*props.Port)
		}
		d.Set("port", port)
		d.Set("protocol", string(props.Protocol))
		d.Set("request_path", props.RequestPath)

		// TODO: parse/make these consistent
		var loadBalancerRules []string
		if rules := props.LoadBalancingRules; rules != nil {
			for _, ruleConfig := range *rules {
				if id := ruleConfig.ID; id != nil {
					loadBalancerRules = append(loadBalancerRules, *id)
				}
			}
		}
		if err := d.Set("load_balancer_rules", loadBalancerRules); err != nil {
			return fmt.Errorf("setting `load_balancer_rules` (Load Balancer %q Probe %q): %+v", id.LoadBalancerName, id.ProbeName, err)
		}
	}

	return nil
}

func resourceArmLoadBalancerProbeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerProbeID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Probe %q: %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, id.ProbeName, err)
	}

	_, index, exists := FindLoadBalancerProbeByName(&loadBalancer, id.ProbeName)
	if !exists {
		return nil
	}

	oldProbes := *loadBalancer.LoadBalancerPropertiesFormat.Probes
	newProbes := append(oldProbes[:index], oldProbes[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.Probes = &newProbes

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, loadBalancer)
	if err != nil {
		return fmt.Errorf("updating Load Balancer %q (Resource Group %q) for deletion of Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q) for deletion of Probe %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.ProbeName, err)
	}

	return nil
}

func expandAzureRmLoadBalancerProbe(d *pluginsdk.ResourceData) *network.Probe {
	properties := network.ProbePropertiesFormat{
		NumberOfProbes:    utils.Int32(int32(d.Get("number_of_probes").(int))),
		IntervalInSeconds: utils.Int32(int32(d.Get("interval_in_seconds").(int))),
		Port:              utils.Int32(int32(d.Get("port").(int))),
	}

	if v, ok := d.GetOk("protocol"); ok {
		properties.Protocol = network.ProbeProtocol(v.(string))
	}

	if v, ok := d.GetOk("request_path"); ok {
		properties.RequestPath = utils.String(v.(string))
	}

	return &network.Probe{
		Name:                  utils.String(d.Get("name").(string)),
		ProbePropertiesFormat: &properties,
	}
}
