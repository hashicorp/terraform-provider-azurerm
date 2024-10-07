// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hsm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDedicatedHardwareSecurityModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDedicatedHardwareSecurityModuleCreate,
		Read:   resourceDedicatedHardwareSecurityModuleRead,
		Update: resourceDedicatedHardwareSecurityModuleUpdate,
		Delete: resourceDedicatedHardwareSecurityModuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dedicatedhsms.ParseDedicatedHSMID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DedicatedHardwareSecurityModuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dedicatedhsms.SkuNameSafeNetLunaNetworkHSMASevenNineZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKOneCPSSixZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKOneCPSTwoFiveZeroZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKTwoCPSSixZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZero),
					string(dedicatedhsms.SkuNamePayShieldOneZeroKLMKTwoCPSTwoFiveZeroZero),
				}, false),
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"network_interface_private_ip_addresses": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azValidate.IPv4Address,
							},
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},
					},
				},
			},

			"management_network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"network_interface_private_ip_addresses": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azValidate.IPv4Address,
							},
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},
					},
				},
			},

			"stamp_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"stamp1",
					"stamp2",
				}, false),
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDedicatedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dedicatedhsms.NewDedicatedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.DedicatedHsmGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_dedicated_hardware_security_module", id.ID())
	}

	skuName := dedicatedhsms.SkuName(d.Get("sku_name").(string))
	if _, ok := d.GetOk("management_network_profile"); ok {
		if skuName == dedicatedhsms.SkuNameSafeNetLunaNetworkHSMASevenNineZero {
			return fmt.Errorf("management_network_profile should not be specified when sku_name is %s", skuName)
		}
	}

	parameters := dedicatedhsms.DedicatedHsm{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: dedicatedhsms.DedicatedHsmProperties{
			NetworkProfile:           expandDedicatedHsmNetworkProfile(d.Get("network_profile").([]interface{})),
			ManagementNetworkProfile: expandDedicatedHsmNetworkProfile(d.Get("management_network_profile").([]interface{})),
		},
		Sku: &dedicatedhsms.Sku{
			Name: &skuName,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("stamp_id"); ok {
		parameters.Properties.StampId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := zones.ExpandUntyped(v.(*pluginsdk.Set).List())
		if len(zones) > 0 {
			parameters.Zones = &zones
		}
	}

	if err := client.DedicatedHsmCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceDedicatedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DedicatedHsmGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DedicatedHSMName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		props := model.Properties

		if err := d.Set("management_network_profile", flattenDedicatedHsmNetworkProfile(props.ManagementNetworkProfile)); err != nil {
			return fmt.Errorf("setting management_network_profile: %+v", err)
		}

		if err := d.Set("network_profile", flattenDedicatedHsmNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting network_profile: %+v", err)
		}
		d.Set("stamp_id", props.StampId)

		skuName := ""
		if model.Sku.Name != nil {
			skuName = string(*model.Sku.Name)
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceDedicatedHardwareSecurityModuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	parameters := dedicatedhsms.DedicatedHsmPatchParameters{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.DedicatedHsmUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceDedicatedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DedicatedHsmDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDedicatedHsmNetworkProfile(input []interface{}) *dedicatedhsms.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := dedicatedhsms.NetworkProfile{
		Subnet: &dedicatedhsms.ApiEntityReference{
			Id: utils.String(v["subnet_id"].(string)),
		},
		NetworkInterfaces: expandDedicatedHsmNetworkInterfacePrivateIPAddresses(v["network_interface_private_ip_addresses"].(*pluginsdk.Set).List()),
	}

	return &result
}

func expandDedicatedHsmNetworkInterfacePrivateIPAddresses(input []interface{}) *[]dedicatedhsms.NetworkInterface {
	results := make([]dedicatedhsms.NetworkInterface, 0)

	for _, item := range input {
		if item != nil {
			results = append(results, dedicatedhsms.NetworkInterface{
				PrivateIPAddress: utils.String(item.(string)),
			})
		}
	}

	return &results
}

func flattenDedicatedHsmNetworkProfile(input *dedicatedhsms.NetworkProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var subnetId string
	if input.Subnet != nil && input.Subnet.Id != nil {
		subnetId = *input.Subnet.Id
	}

	return []interface{}{
		map[string]interface{}{
			"network_interface_private_ip_addresses": flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input.NetworkInterfaces),
			"subnet_id":                              subnetId,
		},
	}
}

func flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input *[]dedicatedhsms.NetworkInterface) []interface{} {
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
