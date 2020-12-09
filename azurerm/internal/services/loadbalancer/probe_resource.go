package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerProbe() *schema.Resource {
	return &schema.Resource{
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadBalancerValidate.LoadBalancerID,
			},

			"protocol": {
				Type:             schema.TypeString,
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
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"request_path": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"interval_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(5),
			},

			"number_of_probes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"load_balancer_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerProbeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	loadBalancerIDRaw := loadBalancerId.ID("")
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	newProbe := expandAzureRmLoadBalancerProbe(d)
	probes := append(*loadBalancer.LoadBalancerPropertiesFormat.Probes, *newProbe)

	existingProbe, existingProbeIndex, exists := FindLoadBalancerProbeByName(loadBalancer, name)
	if exists {
		if name == *existingProbe.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_probe", *existingProbe.ID)
			}

			// this probe is being updated/reapplied remove old copy from the slice
			probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &probes

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", loadBalancerId.Name, loadBalancerId.ResourceGroup)
	}

	var createdProbeId string
	for _, Probe := range *read.LoadBalancerPropertiesFormat.Probes {
		if *Probe.Name == name {
			createdProbeId = *Probe.ID
		}
	}

	if createdProbeId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Probe ID %q", createdProbeId)
	}

	d.SetId(createdProbeId)

	return resourceArmLoadBalancerProbeRead(d, meta)
}

func resourceArmLoadBalancerProbeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerProbeID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
		return nil
	}

	config, _, exists := FindLoadBalancerProbeByName(loadBalancer, id.ProbeName)
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
			return fmt.Errorf("Error setting `load_balancer_rules` (Load Balancer Probe %q): %+v", id.ProbeName, err)
		}
	}

	return nil
}

func resourceArmLoadBalancerProbeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerProbeID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID("")
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, loadBalancerId)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := FindLoadBalancerProbeByName(loadBalancer, id.ProbeName)
	if !exists {
		return nil
	}

	oldProbes := *loadBalancer.LoadBalancerPropertiesFormat.Probes
	newProbes := append(oldProbes[:index], oldProbes[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.Probes = &newProbes

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %s (resource group %s) ID", id.LoadBalancerName, id.ResourceGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerProbe(d *schema.ResourceData) *network.Probe {
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
