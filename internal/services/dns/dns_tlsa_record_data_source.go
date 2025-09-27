// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DnsTLSARecordDataResource{}

type DnsTLSARecordDataResource struct{}

func (DnsTLSARecordDataResource) ModelObject() interface{} {
	return &DnsTLSARecordDataSourceModel{}
}

func (d DnsTLSARecordDataResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRecordTypeID(recordsets.RecordTypeTLSA)
}

func (DnsTLSARecordDataResource) ResourceType() string {
	return "azurerm_dns_tlsa_record"
}

type DnsTLSARecordDataSourceModel struct {
	Name   string                        `tfschema:"name"`
	ZoneId string                        `tfschema:"dns_zone_id"`
	Ttl    int64                         `tfschema:"ttl"`
	Record []DnsTLSARecordResourceRecord `tfschema:"record"`
	Tags   map[string]string             `tfschema:"tags"`
	Fqdn   string                        `tfschema:"fqdn"`
}

func (DnsTLSARecordDataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"dns_zone_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: zones.ValidateDnsZoneID,
		},
	}
}

func (DnsTLSARecordDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"record": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"matching_type": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"selector": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"usage": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"certificate_association_data": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DnsTLSARecordDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DnsTLSARecordDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			zoneId, err := zones.ParseDnsZoneID(state.ZoneId)
			if err != nil {
				return fmt.Errorf("parsing dns_zone_id: %+v", err)
			}
			id := recordsets.NewRecordTypeID(subscriptionId, zoneId.ResourceGroupName, zoneId.DnsZoneName, recordsets.RecordTypeTLSA, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Ttl = pointer.From(props.TTL)
					state.Fqdn = pointer.From(props.Fqdn)

					state.Record = flattenDnsTLSARecords(props.TLSARecords)

					state.Tags = pointer.From(props.Metadata)
				}
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
