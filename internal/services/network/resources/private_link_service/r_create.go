// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_link_service

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateLinkServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatelinkservices.NewPrivateLinkServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, privatelinkservices.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_link_service", id.ID())
		}
	}

	parameters := privatelinkservices.PrivateLinkService{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &privatelinkservices.PrivateLinkServiceProperties{
			AutoApproval: &privatelinkservices.ResourceSet{
				Subscriptions: pluginsdk.ExpandStringSlice(d.Get("auto_approval_subscription_ids").(*pluginsdk.Set).List()),
			},
			EnableProxyProtocol: pointer.To(d.Get("proxy_protocol_enabled").(bool)),
			Visibility: &privatelinkservices.ResourceSet{
				Subscriptions: pluginsdk.ExpandStringSlice(d.Get("visibility_subscription_ids").(*pluginsdk.Set).List()),
			},
			IPConfigurations:                     expandPrivateLinkServiceIPConfiguration(d.Get("nat_ip_configuration").([]any)),
			LoadBalancerFrontendIPConfigurations: expandPrivateLinkServiceFrontendIPConfiguration(d.Get("load_balancer_frontend_ip_configuration_ids").(*pluginsdk.Set).List()),
			Fqdns:                                pluginsdk.ExpandStringSlice(d.Get("fqdns").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("destination_ip_address"); ok {
		parameters.Properties.DestinationIPAddress = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// TODO: confirm whether this can be removed
	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for %s to finish applying", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	return resourcePrivateLinkServiceRead(d, meta)
}

func expandPrivateLinkServiceIPConfiguration(input []any) *[]privatelinkservices.PrivateLinkServiceIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]privatelinkservices.PrivateLinkServiceIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]any)
		privateIpAddress := v["private_ip_address"].(string)
		subnetId := v["subnet_id"].(string)
		privateIpAddressVersion := v["private_ip_address_version"].(string)
		name := v["name"].(string)
		primary := v["primary"].(bool)

		result := privatelinkservices.PrivateLinkServiceIPConfiguration{
			Name: pointer.To(name),
			Properties: &privatelinkservices.PrivateLinkServiceIPConfigurationProperties{
				PrivateIPAddress:        pointer.To(privateIpAddress),
				PrivateIPAddressVersion: pointer.ToEnum[privatelinkservices.IPVersion](privateIpAddressVersion),
				Subnet: &privatelinkservices.Subnet{
					Id: pointer.To(subnetId),
				},
				Primary: pointer.To(primary),
			},
		}

		if privateIpAddress != "" {
			result.Properties.PrivateIPAllocationMethod = pointer.To(privatelinkservices.IPAllocationMethodStatic)
		} else {
			result.Properties.PrivateIPAllocationMethod = pointer.To(privatelinkservices.IPAllocationMethodDynamic)
		}

		results = append(results, result)
	}

	return &results
}

func expandPrivateLinkServiceFrontendIPConfiguration(input []any) *[]privatelinkservices.FrontendIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]privatelinkservices.FrontendIPConfiguration, 0)

	for _, item := range input {
		result := privatelinkservices.FrontendIPConfiguration{
			Id: pointer.To(item.(string)),
		}

		results = append(results, result)
	}

	return &results
}
