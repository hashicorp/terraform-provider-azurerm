package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVPNGatewayNatRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVPNGatewayNatRuleCreateUpdate,
		Read:   resourceVPNGatewayNatRuleRead,
		Update: resourceVPNGatewayNatRuleCreateUpdate,
		Delete: resourceVPNGatewayNatRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VpnGatewayNatRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"vpn_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VpnGatewayID,
			},

			"external_address_space_mappings": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"internal_address_space_mappings": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"ip_configuration_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Instance0",
					"Instance1",
				}, false),
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.VpnNatRuleModeEgressSnat),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VpnNatRuleModeEgressSnat),
					string(network.VpnNatRuleModeIngressSnat),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.VpnNatRuleTypeStatic),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VpnNatRuleTypeStatic),
					string(network.VpnNatRuleTypeDynamic),
				}, false),
			},
		},
	}
}

func resourceVPNGatewayNatRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.NatRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vpnGatewayId, err := parse.VpnGatewayID(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewVpnGatewayNatRuleID(subscriptionId, d.Get("resource_group_name").(string), vpnGatewayId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_vpn_gateway_nat_rule", id.ID())
		}
	}

	props := network.VpnGatewayNatRule{
		Name: utils.String(d.Get("name").(string)),
		VpnGatewayNatRuleProperties: &network.VpnGatewayNatRuleProperties{
			ExternalMappings: expandVpnGatewayNatRuleAddressSpaceMappings(d.Get("external_address_space_mappings").(*pluginsdk.Set).List()),
			InternalMappings: expandVpnGatewayNatRuleAddressSpaceMappings(d.Get("internal_address_space_mappings").(*pluginsdk.Set).List()),
			Mode:             network.VpnNatRuleMode(d.Get("mode").(string)),
			Type:             network.VpnNatRuleType(d.Get("type").(string)),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.VpnGatewayNatRuleProperties.IPConfigurationID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName, props)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNGatewayNatRuleRead(d, meta)
}

func resourceVPNGatewayNatRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NatRuleName)
	d.Set("resource_group_name", id.ResourceGroup)

	gatewayId := parse.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroup, id.VpnGatewayName)
	d.Set("vpn_gateway_id", gatewayId.ID())

	if props := resp.VpnGatewayNatRuleProperties; props != nil {
		d.Set("ip_configuration_id", props.IPConfigurationID)
		d.Set("mode", props.Mode)
		d.Set("type", props.Type)

		if err := d.Set("external_address_space_mappings", flattenVpnGatewayNatRuleAddressSpaceMappings(props.ExternalMappings)); err != nil {
			return fmt.Errorf("setting `external_address_space_mappings`: %+v", err)
		}

		if err := d.Set("internal_address_space_mappings", flattenVpnGatewayNatRuleAddressSpaceMappings(props.InternalMappings)); err != nil {
			return fmt.Errorf("setting `internal_address_space_mappings`: %+v", err)
		}
	}

	return nil
}

func resourceVPNGatewayNatRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
	}

	return nil
}

func expandVpnGatewayNatRuleAddressSpaceMappings(input []interface{}) *[]network.VpnNatRuleMapping {
	results := make([]network.VpnNatRuleMapping, 0)

	for _, v := range input {
		results = append(results, network.VpnNatRuleMapping{
			AddressSpace: utils.String(v.(string)),
		})
	}

	return &results
}

func flattenVpnGatewayNatRuleAddressSpaceMappings(input *[]network.VpnNatRuleMapping) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.AddressSpace != nil {
			results = append(results, *item.AddressSpace)
		}
	}

	return results
}
