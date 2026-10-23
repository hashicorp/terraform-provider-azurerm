// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	mariadbServers "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatednszonegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	postgresqlServers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/redis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	return resourcePrivateEndpointFlatten(ctx, meta.(*clients.Client), d, id, resp.Model, true)
}

func validatePrivateLinkServiceId(endpoints []privateendpoints.PrivateLinkServiceConnection) error {
	for _, connection := range endpoints {
		if connection.Properties == nil || connection.Properties.PrivateLinkServiceId == nil {
			return fmt.Errorf("properties/id was nil for %+v", connection)
		}
		_, errors := azure.ValidateResourceID(*connection.Properties.PrivateLinkServiceId, "PrivateLinkServiceID")
		if len(errors) == 0 {
			continue
		}
		_, errors = validate.PrivateConnectionResourceAlias(*connection.Properties.PrivateLinkServiceId, "PrivateLinkServiceID")
		if len(errors) != 0 {
			return fmt.Errorf("PrivateLinkServiceId Invalid: %q", *connection.Properties.PrivateLinkServiceId)
		}
	}
	return nil
}

func getCosmosDbResIdInPrivateServiceConnections(p *privateendpoints.PrivateEndpointProperties) []string {
	var ids []string
	exists := make(map[string]struct{})

	for _, l := range *p.PrivateLinkServiceConnections {
		if l.Properties.PrivateLinkServiceId == nil {
			continue
		}
		id := *l.Properties.PrivateLinkServiceId
		if _, err := cosmosdb.ParseDatabaseAccountID(id); err == nil {
			_, ok := exists[id]
			if !ok {
				ids = append(ids, id)
				exists[id] = struct{}{}
			}
		}
	}
	for _, l := range *p.ManualPrivateLinkServiceConnections {
		if l.Properties.PrivateLinkServiceId == nil {
			continue
		}
		id := *l.Properties.PrivateLinkServiceId
		if _, err := cosmosdb.ParseDatabaseAccountID(id); err == nil {
			_, ok := exists[id]
			if !ok {
				ids = append(ids, id)
				exists[id] = struct{}{}
			}
		}
	}
	// Sort ids, force adding lock in consistent order to avoid potential deadlock
	sort.Strings(ids)
	return ids
}

func createPrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId, inputRaw []any) error {
	if len(inputRaw) != 1 {
		return fmt.Errorf("expected a single Private DNS Zone Groups but got %d", len(inputRaw))
	}
	item := inputRaw[0].(map[string]any)

	dnsZoneGroupId := privatednszonegroups.NewPrivateDnsZoneGroupID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, item["name"].(string))
	privateDnsZoneIdsRaw := item["private_dns_zone_ids"].([]any)
	privateDnsZoneConfigs := make([]privatednszonegroups.PrivateDnsZoneConfig, 0)
	for _, item := range privateDnsZoneIdsRaw {
		v := item.(string)

		privateDnsZone, err := privatezones.ParsePrivateDnsZoneID(v)
		if err != nil {
			return err
		}

		privateDnsZoneConfigs = append(privateDnsZoneConfigs, privatednszonegroups.PrivateDnsZoneConfig{
			Name: pointer.To(privateDnsZone.PrivateDnsZoneName),
			Properties: &privatednszonegroups.PrivateDnsZonePropertiesFormat{
				PrivateDnsZoneId: pointer.To(privateDnsZone.ID()),
			},
		})
	}

	parameters := privatednszonegroups.PrivateDnsZoneGroup{
		Name: pointer.To(id.PrivateEndpointName),
		Properties: &privatednszonegroups.PrivateDnsZoneGroupPropertiesFormat{
			PrivateDnsZoneConfigs: &privateDnsZoneConfigs,
		},
	}
	if err := client.CreateOrUpdateThenPoll(ctx, dnsZoneGroupId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return nil
}

func deletePrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId) error {
	// lookup and delete the (should be, Single) Private DNS Zone Group associated with this Private Endpoint
	privateDnsZoneIds, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, client, id)
	if err != nil {
		return err
	}

	for _, privateDnsZoneId := range *privateDnsZoneIds {
		if err := client.DeleteThenPoll(ctx, privateDnsZoneId); err != nil {
			return fmt.Errorf("deleting %s: %+v", privateDnsZoneId, err)
		}
	}

	return nil
}

func retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId) (*[]privatednszonegroups.PrivateDnsZoneGroupId, error) {
	output := make([]privatednszonegroups.PrivateDnsZoneGroupId, 0)

	privateEndpointId := privatednszonegroups.NewPrivateEndpointID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName)
	dnsZones, err := client.ListComplete(ctx, privateEndpointId)
	if err != nil {
		if response.WasNotFound(dnsZones.LatestHttpResponse) {
			return &output, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	for _, zone := range dnsZones.Items {
		if zone.Id != nil {
			groupId, err := privatednszonegroups.ParsePrivateDnsZoneGroupID(*zone.Id)
			if err != nil {
				return nil, err
			}

			output = append(output, *groupId)
		}
	}

	return &output, nil
}

// normalize the PrivateConnectionId due to the casing change at service side
func normalizePrivateConnectionId(privateConnectionId string) string {
	// intentionally including the extra segment to handle Redis vs Redis Enterprise (which is within the same RP)
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.cache/redis/") {
		if cacheId, err := redis.ParseRediIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = cacheId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbforpostgresql") {
		if serverId, err := postgresqlServers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbformysql") {
		if serverId, err := servers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbformariadb") {
		if serverId, err := mariadbServers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.kusto") {
		if clusterId, err := commonids.ParseKustoClusterIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = clusterId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.signalrservice") {
		if serviceId, err := signalr.ParseSignalRIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serviceId.ID()
		}
	}
	return privateConnectionId
}
