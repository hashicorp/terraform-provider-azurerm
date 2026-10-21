// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatednszonegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePrivateEndpointFlatten(ctx context.Context, metaClient *clients.Client, d *pluginsdk.ResourceData, id *privateendpoints.PrivateEndpointId, model *privateendpoints.PrivateEndpoint, fetchCompleteData bool) error {
	nicsClient := metaClient.Network.NetworkInterfaces
	dnsClient := metaClient.Network.PrivateDnsZoneGroups

	d.Set("name", id.PrivateEndpointName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", helpers.FlattenEdgeZoneModel(model.ExtendedLocation))

		if props := model.Properties; props != nil {
			if err := d.Set("custom_dns_configs", flattenCustomDnsConfigs(props.CustomDnsConfigs)); err != nil {
				return fmt.Errorf("setting `custom_dns_configs`: %+v", err)
			}

			networkInterfaceId := ""
			privateIpAddress := ""
			if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
				nic := (*nics)[0]
				if nic.Id != nil && *nic.Id != "" {
					networkInterfaceId = *nic.Id
					privateIpAddress = helpers.GetPrivateIpAddress(ctx, nicsClient, networkInterfaceId)
				}
			}

			if err := d.Set("network_interface", helpers.FlattenNetworkInterface(networkInterfaceId)); err != nil {
				return fmt.Errorf("setting `network_interface`: %+v", err)
			}

			if err := d.Set("private_service_connection", flattenPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)); err != nil {
				return fmt.Errorf("setting `private_service_connection`: %+v", err)
			}

			if err := d.Set("ip_configuration", flattenPrivateEndpointIPConfigurations(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}

			subnetId := ""
			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId = *props.Subnet.Id
			}
			d.Set("subnet_id", subnetId)
			d.Set("custom_network_interface_name", pointer.From(props.CustomNetworkInterfaceName))

			if fetchCompleteData {
				privateDnsZoneIds, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, dnsClient, *id)
				if err != nil {
					return err
				}

				privateDnsZoneConfigs := make([]any, 0)
				privateDnsZoneGroups := make([]any, 0)
				if privateDnsZoneIds != nil {
					for _, dnsZoneId := range *privateDnsZoneIds {
						flattened, err := retrieveAndFlattenPrivateDnsZone(ctx, dnsClient, dnsZoneId)
						if err != nil {
							return fmt.Errorf("reading %s for %s: %+v", dnsZoneId, id, err)
						}

						// an exceptional case but no harm in handling
						if flattened == nil {
							continue
						}

						privateDnsZoneConfigs = append(privateDnsZoneConfigs, flattened.DnsZoneConfig...)
						privateDnsZoneGroups = append(privateDnsZoneGroups, flattened.DnsZoneGroup)
					}
				}
				if err = d.Set("private_dns_zone_configs", privateDnsZoneConfigs); err != nil {
					return fmt.Errorf("setting `private_dns_zone_configs`: %+v", err)
				}
				if err = d.Set("private_dns_zone_group", privateDnsZoneGroups); err != nil {
					return fmt.Errorf("setting `private_dns_zone_group`: %+v", err)
				}
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenPrivateEndpointIPConfigurations(ipConfigurations *[]privateendpoints.PrivateEndpointIPConfiguration) []any {
	results := make([]any, 0)
	if ipConfigurations == nil {
		return results
	}

	for _, item := range *ipConfigurations {
		if props := item.Properties; props != nil {
			results = append(results, map[string]any{
				"name":               item.Name,
				"private_ip_address": props.PrivateIPAddress,
				"subresource_name":   props.GroupId,
				"member_name":        props.MemberName,
			})
		}
	}

	return results
}

func flattenCustomDnsConfigs(customDnsConfigs *[]privateendpoints.CustomDnsConfigPropertiesFormat) []any {
	results := make([]any, 0)
	if customDnsConfigs == nil {
		return results
	}

	for _, item := range *customDnsConfigs {
		results = append(results, map[string]any{
			"fqdn":         item.Fqdn,
			"ip_addresses": pluginsdk.FlattenSlice(item.IPAddresses),
		})
	}

	return results
}

func flattenPrivateLinkEndpointServiceConnection(serviceConnections *[]privateendpoints.PrivateLinkServiceConnection, manualServiceConnections *[]privateendpoints.PrivateLinkServiceConnection, privateIPAddress string) []any {
	results := make([]any, 0)
	if serviceConnections == nil && manualServiceConnections == nil {
		return results
	}

	if serviceConnections != nil {
		for _, item := range *serviceConnections {
			privateConnectionId := ""
			subResourceNames := make([]any, 0)

			if props := item.Properties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = pluginsdk.FlattenSlice(v)
				}
				if props.PrivateLinkServiceId != nil {
					privateConnectionId = *props.PrivateLinkServiceId
				}
			}
			attrs := map[string]any{
				"name":                 pointer.From(item.Name),
				"is_manual_connection": false,
				"private_ip_address":   privateIPAddress,
				"subresource_names":    subResourceNames,
			}
			if strings.HasSuffix(privateConnectionId, ".azure.privatelinkservice") {
				attrs["private_connection_resource_alias"] = privateConnectionId
			} else {
				privateConnectionId = normalizePrivateConnectionId(privateConnectionId)
				attrs["private_connection_resource_id"] = privateConnectionId
			}

			results = append(results, attrs)
		}
	}

	if manualServiceConnections != nil {
		for _, item := range *manualServiceConnections {
			privateConnectionId := ""
			requestMessage := ""
			subResourceNames := make([]any, 0)

			if props := item.Properties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = pluginsdk.FlattenSlice(v)
				}
				if props.PrivateLinkServiceId != nil {
					privateConnectionId = *props.PrivateLinkServiceId
				}
				if props.RequestMessage != nil {
					requestMessage = *props.RequestMessage
				}
			}

			attrs := map[string]any{
				"name":                 pointer.From(item.Name),
				"is_manual_connection": true,
				"private_ip_address":   privateIPAddress,
				"request_message":      requestMessage,
				"subresource_names":    subResourceNames,
			}
			if strings.HasSuffix(privateConnectionId, ".azure.privatelinkservice") {
				attrs["private_connection_resource_alias"] = privateConnectionId
			} else {
				privateConnectionId = normalizePrivateConnectionId(privateConnectionId)
				attrs["private_connection_resource_id"] = privateConnectionId
			}

			results = append(results, attrs)
		}
	}

	return results
}

func flattenPrivateDnsZoneGroupRecordSets(input *[]privatednszonegroups.RecordSet) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		ttl := 0
		if v.Ttl != nil {
			ttl = int(*v.Ttl)
		}

		ipAddresses := make([]string, 0)
		if v.IPAddresses != nil {
			ipAddresses = *v.IPAddresses
		}

		output = append(output, map[string]any{
			"fqdn":         pointer.From(v.Fqdn),
			"ip_addresses": ipAddresses,
			"name":         pointer.From(v.RecordSetName),
			"ttl":          ttl,
			"type":         pointer.From(v.RecordType),
		})
	}

	return output
}

type flattenedPrivateDnsZoneGroup struct {
	DnsZoneConfig []any
	DnsZoneGroup  map[string]any
}

func retrieveAndFlattenPrivateDnsZone(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privatednszonegroups.PrivateDnsZoneGroupId) (*flattenedPrivateDnsZoneGroup, error) {
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	privateDnsZoneIds := make([]string, 0)
	dnsZoneConfigs := make([]any, 0)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.PrivateDnsZoneConfigs != nil {
			for _, config := range *props.PrivateDnsZoneConfigs {
				if config.Name == nil {
					// necessary to build up the ID
					continue
				}
				if config.Properties == nil || config.Properties.PrivateDnsZoneId == nil {
					// necessary for a bunch of other things
					continue
				}
				configProps := *config.Properties
				name := *config.Name
				privateDnsZoneId := *configProps.PrivateDnsZoneId

				privateDnsZoneIds = append(privateDnsZoneIds, privateDnsZoneId)

				recordSets := flattenPrivateDnsZoneGroupRecordSets(configProps.RecordSets)
				dnsZoneConfigs = append(dnsZoneConfigs, map[string]any{
					"id":                  parse.NewPrivateDnsZoneConfigID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, id.PrivateDnsZoneGroupName, name).ID(),
					"name":                name,
					"private_dns_zone_id": privateDnsZoneId,
					"record_sets":         recordSets,
				})
			}
		}
	}

	return &flattenedPrivateDnsZoneGroup{
		DnsZoneConfig: dnsZoneConfigs,
		DnsZoneGroup: map[string]any{
			"id":                   id.ID(),
			"name":                 id.PrivateDnsZoneGroupName,
			"private_dns_zone_ids": privateDnsZoneIds,
		},
	}, nil
}
