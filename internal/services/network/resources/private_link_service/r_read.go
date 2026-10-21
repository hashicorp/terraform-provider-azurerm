// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_link_service

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateLinkServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, privatelinkservices.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrivateLinkServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("alias", props.Alias)

			d.Set("proxy_protocol_enabled", props.EnableProxyProtocol)

			d.Set("destination_ip_address", pointer.From(props.DestinationIPAddress))

			var autoApprovalSub []any
			if autoApproval := props.AutoApproval; autoApproval != nil {
				autoApprovalSub = pluginsdk.FlattenSlice(autoApproval.Subscriptions)
			}
			if err := d.Set("auto_approval_subscription_ids", autoApprovalSub); err != nil {
				return fmt.Errorf("setting `auto_approval_subscription_ids`: %+v", err)
			}

			var subscriptions []any
			if visibility := props.Visibility; visibility != nil {
				subscriptions = pluginsdk.FlattenSlice(visibility.Subscriptions)
			}
			if err := d.Set("visibility_subscription_ids", subscriptions); err != nil {
				return fmt.Errorf("setting `visibility_subscription_ids`: %+v", err)
			}

			if err := d.Set("fqdns", pluginsdk.FlattenSlice(props.Fqdns)); err != nil {
				return fmt.Errorf("setting `fqdns`: %+v", err)
			}

			if err := d.Set("nat_ip_configuration", flattenPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `nat_ip_configuration`: %+v", err)
			}

			if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
				return fmt.Errorf("setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenPrivateLinkServiceIPConfiguration(input *[]privatelinkservices.PrivateLinkServiceIPConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		privateIpAddress := ""
		privateIpVersion := ""
		subnetId := ""
		primary := false

		if props := item.Properties; props != nil {
			if props.PrivateIPAddress != nil {
				privateIpAddress = *props.PrivateIPAddress
			}

			privateIpVersion = pointer.FromEnum(props.PrivateIPAddressVersion)

			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId = *props.Subnet.Id
			}

			if props.Primary != nil {
				primary = *props.Primary
			}
		}

		results = append(results, map[string]any{
			"name":                       pointer.From(item.Name),
			"primary":                    primary,
			"private_ip_address":         privateIpAddress,
			"private_ip_address_version": privateIpVersion,
			"subnet_id":                  subnetId,
		})
	}

	return results
}

func flattenPrivateLinkServiceFrontendIPConfiguration(input *[]privatelinkservices.FrontendIPConfiguration) *pluginsdk.Set {
	results := &pluginsdk.Set{F: pluginsdk.HashString}
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.Id; id != nil {
			results.Add(*id)
		}
	}

	return results
}

func privateLinkServiceWaitForReadyRefreshFunc(ctx context.Context, client *privatelinkservices.PrivateLinkServicesClient, id privatelinkservices.PrivateLinkServiceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id, privatelinkservices.DefaultGetOperationOptions())
		if err != nil {
			// the API is eventually consistent during recreates..
			if response.WasNotFound(res.HttpResponse) {
				return res, "Pending", nil
			}

			return nil, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if props := model.Properties; props != nil && props.ProvisioningState != nil {
				if state := *props.ProvisioningState; state != "" {
					return res, string(state), nil
				}
			}
		}

		return res, "Pending", nil
	}
}

func validatePrivateLinkNatIpConfiguration(d *pluginsdk.ResourceDiff) error {
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ipConfigurations := d.Get("nat_ip_configuration").([]any)

	for i, item := range ipConfigurations {
		v := item.(map[string]any)
		p := fmt.Sprintf("nat_ip_configuration.%d.private_ip_address", i)
		s := fmt.Sprintf("nat_ip_configuration.%d.subnet_id", i)
		isPrimary := v["primary"].(bool)
		in := v["name"].(string)

		if d.HasChange(p) {
			o, n := d.GetChange(p)
			if o != "" && n == "" {
				return fmt.Errorf("for Private Link Service %q (Resource Group %q) nat_ip_configuration %q `private_ip_address` once assigned cannot be removed", name, resourceGroup, in)
			}
		}

		if isPrimary && d.HasChange(s) {
			o, _ := d.GetChange(s)
			if o != "" {
				return fmt.Errorf("for Private Link Service %q (Resource Group %q) nat_ip_configuration %q `subnet_id` once assigned cannot be changed", name, resourceGroup, in)
			}
		}
	}

	return nil
}
