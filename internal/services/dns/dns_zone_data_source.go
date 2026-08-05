// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DnsZoneDataResource{}

type DnsZoneDataResource struct{}

func (DnsZoneDataResource) ModelObject() interface{} {
	return &DnsZoneDataResourceModel{}
}

func (d DnsZoneDataResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return zones.ValidateDnsZoneID
}

func (DnsZoneDataResource) ResourceType() string {
	return "azurerm_dns_zone"
}

func (DnsZoneDataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			// TODO: we need a CommonSchema type for this which doesn't have ForceNew
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func (DnsZoneDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"name_servers": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

type DnsZoneDataResourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	NumberOfRecordSets    int64             `tfschema:"number_of_record_sets"`
	MaxNumberOfRecordSets int64             `tfschema:"max_number_of_record_sets"`
	NameServers           []string          `tfschema:"name_servers"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (DnsZoneDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DnsZoneDataResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := zones.NewDnsZoneID(subscriptionId, state.ResourceGroupName, state.Name)
			var zone *zones.Zone
			if id.ResourceGroupName != "" {
				resp, err := client.Get(ctx, id)
				if err != nil {
					if response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("%s was not found", id)
					}
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}

				zone = resp.Model
			} else {
				result, resourceGroupName, err := findZone(ctx, client, id.SubscriptionId, id.DnsZoneName)
				if err != nil {
					return err
				}

				if resourceGroupName == nil {
					return fmt.Errorf("unable to locate the Resource Group for DNS Zone %q in Subscription %q", id.DnsZoneName, subscriptionId)
				}

				zone = result
				id.ResourceGroupName = pointer.From(resourceGroupName)
				state.ResourceGroupName = pointer.From(resourceGroupName)
			}

			if zone == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			metadata.SetID(id)

			if props := zone.Properties; props != nil {
				state.NumberOfRecordSets = pointer.From(props.NumberOfRecordSets)
				state.MaxNumberOfRecordSets = pointer.From(props.MaxNumberOfRecordSets)
				state.NameServers = pointer.From(props.NameServers)
			}

			state.Tags = pointer.From(zone.Tags)

			return metadata.Encode(&state)
		},
	}
}

func findZone(ctx context.Context, client *zones.ZonesClient, subscriptionId, name string) (*zones.Zone, *string, error) {
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	zonesIterator, err := client.ListComplete(ctx, subscriptionResourceId, zones.DefaultListOperationOptions())
	if err != nil {
		return nil, nil, fmt.Errorf("listing DNS Zones: %+v", err)
	}

	var found zones.Zone
	for _, zone := range zonesIterator.Items {
		if zone.Name != nil && *zone.Name == name {
			if found.Id != nil {
				return nil, nil, fmt.Errorf("found multiple DNS zones with name %q, please specify the resource group", name)
			}
			found = zone
		}
	}

	if found.Id == nil {
		return nil, nil, fmt.Errorf("could not find DNS zone with name: %q", name)
	}

	id, err := zones.ParseDnsZoneIDInsensitively(*found.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %q as a DNS Zone ID: %+v", *found.Id, err)
	}
	return &found, &id.ResourceGroupName, nil
}
