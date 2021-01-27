package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceJitNetworkAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceJitNetworkAccessPolicyCreateOrUpdate,
		Read:   resourceJitNetworkAccessPolicyRead,
		Update: resourceJitNetworkAccessPolicyCreateOrUpdate,
		Delete: resourceJitNetworkAccessPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.JitNetworkAccessPolicyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"virtual_machines": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: computeValidate.VirtualMachineID,
						},

						"ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_source_address_prefixes": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.Any(
												validate.CIDR,
												validate.IPv4Address,
												validation.StringInSlice([]string{"*"}, false),
											),
										},
									},

									"max_request_access_duration": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ISO8601Duration,
									},

									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validate.PortNumber,
									},

									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(security.All),
											string(security.TCP),
											string(security.UDP),
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceJitNetworkAccessPolicyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))

	client := meta.(*clients.Client).SecurityCenter.NewJitNetworkAccessPoliciesClient(location)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := parse.NewJitNetworkAccessPolicyID(subscriptionId, resourceGroup, location, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_jit_network_access_policy", id.ID())
		}
	}

	jitNetworkAccessPolicy := security.JitNetworkAccessPolicy{
		// for now, kind could only be "Basic"
		Kind: utils.String("Basic"),
		JitNetworkAccessPolicyProperties: &security.JitNetworkAccessPolicyProperties{
			VirtualMachines: expandJitNetworkAccessPolicyVirtualMachine(d.Get("virtual_machines").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, jitNetworkAccessPolicy); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.String(), err)
	}

	// the API returns 'done' but it's not actually finished provisioning yet
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Updating"},
		Target:  []string{"Succeeded"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, resourceGroup, name)
			if err != nil {
				return resp, "Error", fmt.Errorf("retrieving %s: %+v", id.String(), err)
			}

			if properties := resp.JitNetworkAccessPolicyProperties; properties != nil {
				return resp, *properties.ProvisioningState, nil
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
		return fmt.Errorf("waiting for %s to finish provisioning: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceJitNetworkAccessPolicyRead(d, meta)
}

func resourceJitNetworkAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JitNetworkAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).SecurityCenter.NewJitNetworkAccessPoliciesClient(id.LocationName)
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] JIT Network Access Policy %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.Normalize(id.LocationName))
	if props := resp.JitNetworkAccessPolicyProperties; props != nil {
		if err := d.Set("virtual_machines", flattenJitNetworkAccessPolicyVirtualMachine(props.VirtualMachines)); err != nil {
			return fmt.Errorf("settings 'virtual_machines': %+v", err)
		}
	}

	return nil
}

func resourceJitNetworkAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JitNetworkAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).SecurityCenter.NewJitNetworkAccessPoliciesClient(id.LocationName)
	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandJitNetworkAccessPolicyVirtualMachine(input []interface{}) *[]security.JitNetworkAccessPolicyVirtualMachine {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]security.JitNetworkAccessPolicyVirtualMachine, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, security.JitNetworkAccessPolicyVirtualMachine{
			ID:    utils.String(v["id"].(string)),
			Ports: expandJitNetworkAccessPolicyVirtualMachinePort(v["ports"].(*schema.Set).List()),
		})
	}
	return &results
}

func expandJitNetworkAccessPolicyVirtualMachinePort(input []interface{}) *[]security.JitNetworkAccessPortRule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]security.JitNetworkAccessPortRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		portRule := security.JitNetworkAccessPortRule{
			Number:                   utils.Int32(int32(v["port"].(int))),
			Protocol:                 security.Protocol(v["protocol"].(string)),
			MaxRequestAccessDuration: utils.String(v["max_request_access_duration"].(string)),
		}

		allowedSourceAddressPrefixes := v["allowed_source_address_prefixes"].(*schema.Set).List()
		if len(allowedSourceAddressPrefixes) == 1 {
			portRule.AllowedSourceAddressPrefix = utils.String(allowedSourceAddressPrefixes[0].(string))
		} else {
			portRule.AllowedSourceAddressPrefixes = utils.ExpandStringSlice(allowedSourceAddressPrefixes)
		}

		results = append(results, portRule)
	}
	return &results
}

func flattenJitNetworkAccessPolicyVirtualMachine(input *[]security.JitNetworkAccessPolicyVirtualMachine) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, virtualMachine := range *input {
		var id string
		if virtualMachine.ID != nil {
			id = *virtualMachine.ID
		}
		result = append(result, map[string]interface{}{
			"id":    id,
			"ports": flattenJitNetworkAccessPolicyVirtualMachinePort(virtualMachine.Ports),
		})
	}

	return result
}

func flattenJitNetworkAccessPolicyVirtualMachinePort(input *[]security.JitNetworkAccessPortRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		var portNumber int32
		var maxRequestAccessDuration string
		var allowedSourceAddressPrefixes []string

		if item.Number != nil {
			portNumber = *item.Number
		}

		if item.AllowedSourceAddressPrefix != nil {
			allowedSourceAddressPrefixes = append(allowedSourceAddressPrefixes, *item.AllowedSourceAddressPrefix)
		} else if item.AllowedSourceAddressPrefixes != nil {
			allowedSourceAddressPrefixes = *item.AllowedSourceAddressPrefixes
		}

		if item.MaxRequestAccessDuration != nil {
			maxRequestAccessDuration = *item.MaxRequestAccessDuration
		}

		result = append(result, map[string]interface{}{
			"port":                            portNumber,
			"protocol":                        string(item.Protocol),
			"allowed_source_address_prefixes": utils.FlattenStringSlice(&allowedSourceAddressPrefixes),
			"max_request_access_duration":     maxRequestAccessDuration,
		})
	}
	return result
}
