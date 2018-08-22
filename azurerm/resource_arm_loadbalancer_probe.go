package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerProbe() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerProbeCreateUpdate,
		Read:   resourceArmLoadBalancerProbeRead,
		Update: resourceArmLoadBalancerProbeCreateUpdate,
		Delete: resourceArmLoadBalancerProbeDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"protocol": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				StateFunc:        ignoreCaseStateFunc,
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
					ValidateFunc: validation.NoZeroValues,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerProbeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)

	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error updating Load Balancer Probes: props was nil")
	}

	probes := props.Probes
	if probes == nil {
		return fmt.Errorf("Error updating Load Balancer Probes: props.Probes was nil")
	}
	existingProbes := *probes

	newProbe := expandAzureRmLoadBalancerProbe(d)
	if d.IsNewResource() {
		// firstly check if it exists
		for _, probe := range existingProbes {
			if probe.Name != nil && *probe.Name == name {
				return tf.ImportAsExistsError("azurerm_lb_probe", *probe.ID)
			}
		}

		// then append it if not
		existingProbes = append(existingProbes, *newProbe)
	} else {
		// swap it in-place
		index := -1
		for i, probe := range existingProbes {
			if probe.Name != nil && *probe.Name == name {
				index = i
				break
			}
		}

		if index == -1 {
			// should have been caught by the Read
			return fmt.Errorf("Error: Probe %q was not found on Load Balancer %q", name, loadBalancerID)
		}

		existingProbes[index] = *newProbe
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &existingProbes
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", loadBalancerName, resGroup)
	}

	var probeId string
	if props := read.LoadBalancerPropertiesFormat; props != nil {
		if probes := props.Probes; probes != nil {
			for _, probe := range *probes {
				if *probe.Name == name {
					probeId = *probe.ID
				}
			}
		}
	}

	if probeId == "" {
		return fmt.Errorf("Cannot find created Load Balancer Probe ID %q", probeId)
	}

	d.SetId(probeId)

	return resourceArmLoadBalancerProbeRead(d, meta)
}

func resourceArmLoadBalancerProbeRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["probes"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	var probe *network.Probe
	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if probes := props.Probes; probes != nil {
			for _, p := range *probes {
				if p.Name != nil && *p.Name == name {
					probe = &p
					break
				}
			}
		}
	}

	if probe == nil {
		log.Printf("[INFO] Load Balancer Probe %q not found. Removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := probe.ProbePropertiesFormat; properties != nil {
		d.Set("protocol", properties.Protocol)
		d.Set("interval_in_seconds", properties.IntervalInSeconds)
		d.Set("number_of_probes", properties.NumberOfProbes)
		d.Set("port", properties.Port)
		d.Set("request_path", properties.RequestPath)

		loadBalancerRules := make([]string, 0)
		if rules := properties.LoadBalancingRules; rules != nil {
			for _, ruleConfig := range *rules {
				if id := ruleConfig.ID; id != nil {
					loadBalancerRules = append(loadBalancerRules, *id)
				}
			}
		}
		if err := d.Set("load_balancer_rules", loadBalancerRules); err != nil {
			return fmt.Errorf("Error setting `load_balancer_rules`: %+v", err)
		}
	}

	return nil
}

func resourceArmLoadBalancerProbeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)

	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	props := loadBalancer.LoadBalancerPropertiesFormat
	if props == nil {
		log.Printf("[INFO] props was nil for Load Balancer %q - removing from state", loadBalancerID)
		d.SetId("")
		return nil
	}

	probes := props.Probes
	if probes == nil {
		log.Printf("[INFO] props.Probes was nil for Load Balancer %q - removing from state", loadBalancerID)
		d.SetId("")
		return nil
	}

	newProbes := make([]network.Probe, 0)
	for _, probe := range *probes {
		if probe.Name != nil && *probe.Name != name {
			newProbes = append(newProbes, probe)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &newProbes

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer: %+v", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %s (resource group %s) ID", loadBalancerName, resGroup)
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
		Name: utils.String(d.Get("name").(string)),
		ProbePropertiesFormat: &properties,
	}
}
