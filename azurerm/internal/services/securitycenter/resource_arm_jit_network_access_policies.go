package securitycenter

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmJitNetworkAccessPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmJitNetworkAccessPoliciesCreateOrUpdate,
		Read:   resourceArmJitNetworkAccessPoliciesRead,
		Update: resourceArmJitNetworkAccessPoliciesCreateOrUpdate,
		Delete: resourceArmJitNetworkAccessPoliciesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"asc_location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_machines": virtualMachinesPolicySchema(),
		},
	}
}

func flattenPorts(input *[]security.JitNetworkAccessPortRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)
	for _, port := range *input {
		result = append(result, map[string]interface{}{
			"port":                          port.Number,
			"protocol":                      port.Protocol,
			"allowed_source_address_prefix": port.AllowedSourceAddressPrefix,
			// "allowed_source_address_prefixes":
			"max_request_access_duration": port.MaxRequestAccessDuration,
		})
	}
	return result
}

func extractVirtualMachineName(id *string) string {
	subStrings := strings.Split(*id, "/")
	name := subStrings[len(subStrings)-1]
	if len(name) > 0 {
		fmt.Printf("Virtual Machine Name: %s\n", name)
		return name
	}
	return ""
}

func flattenVirtualMachines(input *[]security.JitNetworkAccessPolicyVirtualMachine) []map[string]interface{} {
	if input == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)

	for _, virtualMachine := range *input {
		result = append(result, map[string]interface{}{
			"name":  extractVirtualMachineName(virtualMachine.ID),
			"ports": flattenPorts(virtualMachine.Ports),
		})
	}

	//fmt.Printf("Virtual Machines result: %s\n", spew.Sdump(result))
	return result
}

func virtualMachinesPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"ports": portsPolicySchema(),
			},
		},
	}
}

func portsPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port": {
					Type:         schema.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.PortNumber,
				},
				"protocol": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(security.TCP),
					ValidateFunc: validation.StringInSlice([]string{
						string(security.All),
						string(security.TCP),
						string(security.UDP),
					}, false),
				},
				"allowed_source_address_prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"max_request_access_duration": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validate.ISO8601Duration,
				},
			},
		},
	}
}

func expandJitNetworkAccessPortRules(portData map[string]interface{}) *[]security.JitNetworkAccessPortRule {
	ports := make([]security.JitNetworkAccessPortRule, 0)

	if v, ok := portData["ports"].([]interface{}); ok {
		for _, portConfig := range v {
			data := portConfig.(map[string]interface{})
			number := int32(data["port"].(int))
			protocol := security.Protocol(data["protocol"].(string))
			allowed_source_address_prefix := data["allowed_source_address_prefix"].(string)
			max_request_access_duration := data["max_request_access_duration"].(string)
			port := security.JitNetworkAccessPortRule{
				Number:                     &number,
				Protocol:                   protocol,
				AllowedSourceAddressPrefix: &allowed_source_address_prefix,
				MaxRequestAccessDuration:   &max_request_access_duration,
			}
			ports = append(ports, port)
		}
	}

	return &ports
}

func resourceArmJitNetworkAccessPoliciesCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachinesConfig := d.Get("virtual_machines").([]interface{})
	virtualMachines := make([]security.JitNetworkAccessPolicyVirtualMachine, 0)

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	for _, virtualMachineConfig := range virtualMachinesConfig {
		data := virtualMachineConfig.(map[string]interface{})
		var virtualMachineID string = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{virtualMachineName}"

		vmIdReplacer := strings.NewReplacer("{subscriptionId}", subscriptionId,
			"{resourceGroupName}", d.Get("resource_group_name").(string),
			"{virtualMachineName}", data["name"].(string))
		vmID := vmIdReplacer.Replace(virtualMachineID)

		virtualMachine := security.JitNetworkAccessPolicyVirtualMachine{
			ID:    &vmID,
			Ports: expandJitNetworkAccessPortRules(data),
		}
		virtualMachines = append(virtualMachines, virtualMachine)
	}

	accessRequests := make([]security.JitNetworkAccessRequest, 0)
	// TODO: build access request
	jitNetworkAccessPolicyProperties := security.JitNetworkAccessPolicyProperties{
		VirtualMachines: &virtualMachines,
		Requests:        &accessRequests,
	}

	resourceGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)
	location := d.Get("asc_location").(string)
	name := d.Get("name").(string)

	policyId := fmt.Sprintf(
		"/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Security/locations/%v/jitNetworkAccessPolicies/%v",
		subscriptionId, resourceGroup, location, name,
	)

	policyType := "Microsoft.Security/locations/jitNetworkAccessPolicies"

	jitNetworkAccessPolicy := security.JitNetworkAccessPolicy{
		ID:                               &policyId,
		Name:                             &name,
		Type:                             &policyType,
		Kind:                             &kind,
		Location:                         &location,
		JitNetworkAccessPolicyProperties: &jitNetworkAccessPolicyProperties,
	}

	ascLocation := d.Get("asc_location").(string)

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, name, ascLocation, jitNetworkAccessPolicy)
	if err != nil {
		return fmt.Errorf("Error creating/updating Security Center Subscription pricing: %+v", err)
	}

	// the API returns 'done' but it's not actually finished provisioning yet
	// Updating and Succeeded are the only states discovered so far, and may not be complete list
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string("Updating"),
		},
		Target: []string{
			string("Succeeded"),
		},
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, resourceGroup, name, ascLocation)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("Error retrieving JIT Network Access  Policy %q (Location %q / Resource Group %q): %+v", name, ascLocation, resourceGroup, err2)
			}

			if properties := resp.JitNetworkAccessPolicyProperties; properties != nil {
				return resp, string(*properties.ProvisioningState), nil
			}

			return resp, "Unknown", nil
		},
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for JIT Network Access  Policy %q (Location %q / Resource Group %q) to finish provisioning: %+v", name, ascLocation, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, ascLocation)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read JIT Network Access Policy %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmJitNetworkAccessPoliciesRead(d, meta)
}

func resourceArmJitNetworkAccessPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JitNetworkAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.Location)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] JIT Network Access Policy %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving JIT Network Access Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("asc_location", azure.NormalizeLocation(*location))
	} else {
		d.Set("asc_location", id.Location)
	}
	d.Set("kind", "Basic")

	// flatten JitNetworkAccessPolicyProperties to create the virtualmachines struct for tf state
	if resp.JitNetworkAccessPolicyProperties == nil {
		return fmt.Errorf("retrieving Jit Network Access Policy %q (Resource Group %q): `JitNetworkAccessPolicyProperties` was nil", id.Name, id.ResourceGroup)
	}

	if err := d.Set("virtual_machines", flattenVirtualMachines(resp.JitNetworkAccessPolicyProperties.VirtualMachines)); err != nil {
		return fmt.Errorf("settings 'JitNetworkAccessPolicyProperties': %+v", err)
	}

	return nil
}

func resourceArmJitNetworkAccessPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[DEBUG] Security Center Subscription JIT Network Access Deletion invocation")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("asc_location").(string)

	resp, err := client.Delete(ctx, resGroup, name, location)

	if err != nil {
		return fmt.Errorf("Error deleting JIT Network Access Policy %q (Resource Group %q): %s", name, resGroup, err)
	} else {
		fmt.Printf("Delete Succeeded %d\n", resp.StatusCode)
	}

	return nil
}
