package hsm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/hardwaresecuritymodules/mgmt/2018-10-31-preview/hardwaresecuritymodules"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hsm/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hsm/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDedicatedHardwareSecurityModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDedicatedHardwareSecurityModuleCreate,
		Read:   resourceArmDedicatedHardwareSecurityModuleRead,
		Update: resourceArmDedicatedHardwareSecurityModuleUpdate,
		Delete: resourceArmDedicatedHardwareSecurityModuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DedicatedHardwareSecurityModuleID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DedicatedHardwareSecurityModuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(hardwaresecuritymodules.SafeNetLunaNetworkHSMA790),
				}, false),
			},

			"network_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_private_ip_addresses": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azValidate.IPv4Address,
							},
						},

						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.SubnetID,
						},
					},
				},
			},

			"stamp_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"stamp1",
					"stamp2",
				}, false),
			},

			"zones": azure.SchemaZones(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDedicatedHardwareSecurityModuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_dedicated_hardware_security_module", *existing.ID)
	}

	parameters := hardwaresecuritymodules.DedicatedHsm{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		DedicatedHsmProperties: &hardwaresecuritymodules.DedicatedHsmProperties{
			NetworkProfile: expandArmDedicatedHsmNetworkProfile(d.Get("network_profile").([]interface{})),
		},
		Sku: &hardwaresecuritymodules.Sku{
			Name: hardwaresecuritymodules.Name(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("stamp_id"); ok {
		parameters.DedicatedHsmProperties.StampID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		parameters.Zones = azure.ExpandZones(v.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Dedicated Hardware Security Module %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceArmDedicatedHardwareSecurityModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DedicatedHSMName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Hardware Security Module %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	d.Set("name", id.DedicatedHSMName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DedicatedHsmProperties; props != nil {
		if err := d.Set("network_profile", flattenArmDedicatedHsmNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting network_profile: %+v", err)
		}
		d.Set("stamp_id", props.StampID)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("zones", resp.Zones); err != nil {
		return fmt.Errorf("setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDedicatedHardwareSecurityModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	parameters := hardwaresecuritymodules.DedicatedHsmPatchParameters{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.DedicatedHSMName, parameters)
	if err != nil {
		return fmt.Errorf("updating Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating future for Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	return resourceArmDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceArmDedicatedHardwareSecurityModuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DedicatedHSMName)
	if err != nil {
		return fmt.Errorf("deleting Dedicated Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Dedicated Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	return nil
}

func expandArmDedicatedHsmNetworkProfile(input []interface{}) *hardwaresecuritymodules.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := hardwaresecuritymodules.NetworkProfile{
		Subnet: &hardwaresecuritymodules.APIEntityReference{
			ID: utils.String(v["subnet_id"].(string)),
		},
		NetworkInterfaces: expandArmDedicatedHsmNetworkInterfacePrivateIPAddresses(v["network_interface_private_ip_addresses"].(*schema.Set).List()),
	}

	return &result
}

func expandArmDedicatedHsmNetworkInterfacePrivateIPAddresses(input []interface{}) *[]hardwaresecuritymodules.NetworkInterface {
	results := make([]hardwaresecuritymodules.NetworkInterface, 0)

	for _, item := range input {
		if item != nil {
			result := hardwaresecuritymodules.NetworkInterface{
				PrivateIPAddress: utils.String(item.(string)),
			}

			results = append(results, result)
		}
	}

	return &results
}

func flattenArmDedicatedHsmNetworkProfile(input *hardwaresecuritymodules.NetworkProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var subnetId string
	if input.Subnet != nil && input.Subnet.ID != nil {
		subnetId = *input.Subnet.ID
	}

	return []interface{}{
		map[string]interface{}{
			"network_interface_private_ip_addresses": flattenArmDedicatedHsmNetworkInterfacePrivateIPAddresses(input.NetworkInterfaces),
			"subnet_id":                              subnetId,
		},
	}
}

func flattenArmDedicatedHsmNetworkInterfacePrivateIPAddresses(input *[]hardwaresecuritymodules.NetworkInterface) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.PrivateIPAddress != nil {
			results = append(results, *item.PrivateIPAddress)
		}
	}

	return results
}
