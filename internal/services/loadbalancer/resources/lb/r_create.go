// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := loadbalancers.NewLoadBalancerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	// Casting because of a bug in provider ID generation due to a data issue that needs investigation.
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, loadbalancers.ProviderLoadBalancerId(id), loadbalancers.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_lb", id.ID())
		}
	}

	if strings.EqualFold(d.Get("sku_tier").(string), string(loadbalancers.LoadBalancerSkuTierGlobal)) {
		if !strings.EqualFold(d.Get("sku").(string), string(loadbalancers.LoadBalancerSkuNameStandard)) {
			return fmt.Errorf("global load balancing is only supported for standard SKU load balancers")
		}
	}

	sku := loadbalancers.LoadBalancerSku{
		Name: pointer.ToEnum[loadbalancers.LoadBalancerSkuName](d.Get("sku").(string)),
		Tier: pointer.ToEnum[loadbalancers.LoadBalancerSkuTier](d.Get("sku_tier").(string)),
	}

	properties := loadbalancers.LoadBalancerPropertiesFormat{}

	if _, ok := d.GetOk("frontend_ip_configuration"); ok {
		properties.FrontendIPConfigurations = expandAzureRmLoadBalancerFrontendIpConfigurations(d)
	}

	loadBalancer := loadbalancers.LoadBalancer{
		Name:             pointer.To(id.LoadBalancerName),
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:             tags.Expand(d.Get("tags").(map[string]any)),
		Sku:              pointer.To(sku),
		Properties:       pointer.To(properties),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, loadbalancers.ProviderLoadBalancerId(id), loadBalancer, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRead(d, meta)
}

func expandAzureRmLoadBalancerFrontendIpConfigurations(d *pluginsdk.ResourceData) *[]loadbalancers.FrontendIPConfiguration {
	configsRaw, ok := d.GetOk("frontend_ip_configuration")
	if !ok {
		return nil
	}

	configs := configsRaw.([]any)
	frontEndConfigs := make([]loadbalancers.FrontendIPConfiguration, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)

		properties := loadbalancers.FrontendIPConfigurationPropertiesFormat{}

		if v := data["private_ip_address_allocation"].(string); v != "" {
			properties.PrivateIPAllocationMethod = pointer.ToEnum[loadbalancers.IPAllocationMethod](v)
		}

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &loadbalancers.SubResource{
				Id: pointer.To(v),
			}
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = pointer.To(v)
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &loadbalancers.PublicIPAddress{
				Id: pointer.To(v),
			}
		}

		if v := data["public_ip_prefix_id"].(string); v != "" {
			properties.PublicIPPrefix = &loadbalancers.SubResource{
				Id: pointer.To(v),
			}
		}

		if v := data["subnet_id"].(string); v != "" {
			properties.PrivateIPAddressVersion = pointer.To(loadbalancers.IPVersionIPvFour)
			if v := data["private_ip_address_version"].(string); v != "" {
				properties.PrivateIPAddressVersion = pointer.ToEnum[loadbalancers.IPVersion](v)
			}
			properties.Subnet = &loadbalancers.Subnet{
				Id: pointer.To(v),
			}
		}

		frontEndConfig := loadbalancers.FrontendIPConfiguration{
			Name:       pointer.To(data["name"].(string)),
			Properties: pointer.To(properties),
		}

		zones := zones.ExpandUntyped(data["zones"].(*pluginsdk.Set).List())
		if len(zones) > 0 {
			frontEndConfig.Zones = &zones
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}
