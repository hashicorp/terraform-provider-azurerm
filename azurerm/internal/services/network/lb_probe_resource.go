package network

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
				ValidateFunc: azure.ValidateResourceID,
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
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
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
			if features.ShouldResourcesBeImported() && d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_probe", *existingProbe.ID)
			}

			// this probe is being updated/reapplied remove old copy from the slice
			probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &probes
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(loadBalancerID)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", loadBalancerName, resGroup)
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
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["probes"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d, d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := FindLoadBalancerProbeByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Probe %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := config.ProbePropertiesFormat; properties != nil {
		d.Set("protocol", properties.Protocol)
		d.Set("interval_in_seconds", properties.IntervalInSeconds)
		d.Set("number_of_probes", properties.NumberOfProbes)
		d.Set("port", properties.Port)
		d.Set("request_path", properties.RequestPath)

		var loadBalancerRules []string
		if rules := properties.LoadBalancingRules; rules != nil {
			for _, ruleConfig := range *rules {
				if id := ruleConfig.ID; id != nil {
					loadBalancerRules = append(loadBalancerRules, *id)
				}
			}
		}
		if err := d.Set("load_balancer_rules", loadBalancerRules); err != nil {
			return fmt.Errorf("Error setting `load_balancer_rules` (Load Balancer Probe %q): %+v", name, err)
		}
	}

	return nil
}

func resourceArmLoadBalancerProbeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerID := d.Get("loadbalancer_id").(string)
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(d, loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := FindLoadBalancerProbeByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldProbes := *loadBalancer.LoadBalancerPropertiesFormat.Probes
	newProbes := append(oldProbes[:index], oldProbes[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.Probes = &newProbes

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
		Name:                  utils.String(d.Get("name").(string)),
		ProbePropertiesFormat: &properties,
	}
}
