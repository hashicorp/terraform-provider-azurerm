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

var (
	_ sdk.Resource           = DnsTLSARecordResource{}
	_ sdk.ResourceWithUpdate = DnsTLSARecordResource{}
)

type DnsTLSARecordResource struct{}

func (DnsTLSARecordResource) ModelObject() interface{} {
	return &DnsTLSARecordResourceModel{}
}

func (DnsTLSARecordResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRecordTypeID(recordsets.RecordTypeTLSA)
}

func (DnsTLSARecordResource) ResourceType() string {
	return "azurerm_dns_tlsa_record"
}

type DnsTLSARecordResourceModel struct {
	Name   string                        `tfschema:"name"`
	ZoneId string                        `tfschema:"dns_zone_id"`
	Ttl    int64                         `tfschema:"ttl"`
	Record []DnsTLSARecordResourceRecord `tfschema:"record"`
	Tags   map[string]string             `tfschema:"tags"`
	Fqdn   string                        `tfschema:"fqdn"`
}

type DnsTLSARecordResourceRecord struct {
	MatchingType        int64  `tfschema:"matching_type"`
	Selector            int64  `tfschema:"selector"`
	Usage               int64  `tfschema:"usage"`
	CertAssociationData string `tfschema:"certificate_association_data"`
}

func (DnsTLSARecordResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"dns_zone_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: zones.ValidateDnsZoneID,
		},

		"record": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"matching_type": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"selector": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"usage": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"certificate_association_data": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (DnsTLSARecordResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DnsTLSARecordResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DnsTLSARecordResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			zoneId, err := zones.ParseDnsZoneID(model.ZoneId)
			if err != nil {
				return fmt.Errorf("parsing dns_zone_id: %+v", err)
			}

			id := recordsets.NewRecordTypeID(subscriptionId, zoneId.ResourceGroupName, zoneId.DnsZoneName, recordsets.RecordTypeTLSA, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := recordsets.RecordSet{
				Name: pointer.To(model.Name),
				Properties: &recordsets.RecordSetProperties{
					Metadata:    pointer.To(model.Tags),
					TTL:         pointer.To(model.Ttl),
					TLSARecords: expandDnsTLSARecords(model.Record),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (DnsTLSARecordResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			state := DnsTLSARecordResourceModel{}

			id, err := recordsets.ParseRecordTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			zoneId := zones.NewDnsZoneID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName)
			state.Name = id.RelativeRecordSetName
			state.ZoneId = zoneId.ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Ttl = pointer.From(props.TTL)
					state.Fqdn = pointer.From(props.Fqdn)

					state.Record = flattenDnsTLSARecords(props.TLSARecords)

					state.Tags = pointer.From(props.Metadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (DnsTLSARecordResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			var model DnsTLSARecordResourceModel

			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := recordsets.ParseRecordTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("record") {
				payload.Properties.TLSARecords = expandDnsTLSARecords(model.Record)
			}

			if metadata.ResourceData.HasChange("ttl") {
				payload.Properties.TTL = pointer.To(model.Ttl)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Properties.Metadata = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (DnsTLSARecordResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			id, err := recordsets.ParseRecordTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id, recordsets.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func flattenDnsTLSARecords(records *[]recordsets.TlsaRecord) []DnsTLSARecordResourceRecord {
	results := make([]DnsTLSARecordResourceRecord, 0)

	if records != nil {
		for _, record := range *records {
			results = append(results, DnsTLSARecordResourceRecord{
				MatchingType:        pointer.From(record.MatchingType),
				Selector:            pointer.From(record.Selector),
				Usage:               pointer.From(record.Usage),
				CertAssociationData: pointer.From(record.CertAssociationData),
			})
		}
	}

	return results
}

func expandDnsTLSARecords(d []DnsTLSARecordResourceRecord) *[]recordsets.TlsaRecord {
	records := make([]recordsets.TlsaRecord, 0)

	for _, v := range d {
		records = append(records, recordsets.TlsaRecord{
			MatchingType:        pointer.To(v.MatchingType),
			Selector:            pointer.To(v.Selector),
			Usage:               pointer.To(v.Usage),
			CertAssociationData: pointer.To(v.CertAssociationData),
		})
	}

	return &records
}
