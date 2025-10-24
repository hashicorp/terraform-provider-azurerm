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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DnsARecordDataResource{}

type DnsARecordDataResource struct{}

func (DnsARecordDataResource) ModelObject() interface{} {
	return &DnsARecordDataSourceModel{}
}

func (d DnsARecordDataResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRecordTypeID(recordsets.RecordTypeA)
}

func (DnsARecordDataResource) ResourceType() string {
	return "azurerm_dns_a_record"
}

type DnsARecordDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ZoneName          string            `tfschema:"zone_name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Ttl               int64             `tfschema:"ttl"`
	Records           []string          `tfschema:"records"`
	Tags              map[string]string `tfschema:"tags"`
	Fqdn              string            `tfschema:"fqdn"`
	TargetResourceId  string            `tfschema:"target_resource_id"`
}

func (DnsARecordDataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"zone_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (DnsARecordDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"records": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DnsARecordDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DnsARecordDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := recordsets.NewRecordTypeID(subscriptionId, state.ResourceGroupName, state.ZoneName, recordsets.RecordTypeA, state.Name)

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

					state.Records = flattenAzureRmDnsARecords(props.ARecords)
					state.TargetResourceId = pointer.From(props.TargetResource.Id)

					state.Tags = pointer.From(props.Metadata)
				}
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
