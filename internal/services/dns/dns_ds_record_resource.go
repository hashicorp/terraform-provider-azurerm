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
	_ sdk.Resource           = DnsDSRecordResource{}
	_ sdk.ResourceWithUpdate = DnsDSRecordResource{}
)

type DnsDSRecordResource struct{}

func (DnsDSRecordResource) ModelObject() interface{} {
	return &DnsDSRecordResourceModel{}
}

func (DnsDSRecordResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRecordTypeID(recordsets.RecordTypeDS)
}

func (DnsDSRecordResource) ResourceType() string {
	return "azurerm_dns_ds_record"
}

type DnsDSRecordResourceModel struct {
	Name   string                      `tfschema:"name"`
	ZoneId string                      `tfschema:"dns_zone_id"`
	Ttl    int64                       `tfschema:"ttl"`
	Record []DnsDSRecordResourceRecord `tfschema:"record"`
	Tags   map[string]string           `tfschema:"tags"`
	Fqdn   string                      `tfschema:"fqdn"`
}

type DnsDSRecordResourceRecord struct {
	Algorithm   int64  `tfschema:"algorithm"`
	KeyTag      int64  `tfschema:"key_tag"`
	DigestType  int64  `tfschema:"digest_type"`
	DigestValue string `tfschema:"digest_value"`
}

func (DnsDSRecordResource) Arguments() map[string]*pluginsdk.Schema {
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
					"algorithm": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"key_tag": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"digest_type": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"digest_value": {
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

func (DnsDSRecordResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DnsDSRecordResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DnsDSRecordResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			zoneId, err := zones.ParseDnsZoneID(model.ZoneId)
			if err != nil {
				return fmt.Errorf("parsing dns_zone_id: %+v", err)
			}
			id := recordsets.NewRecordTypeID(subscriptionId, zoneId.ResourceGroupName, zoneId.DnsZoneName, recordsets.RecordTypeDS, model.Name)

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
					Metadata:  pointer.To(model.Tags),
					TTL:       pointer.To(model.Ttl),
					DSRecords: expandDnsDSRecords(model.Record),
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

func (DnsDSRecordResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			state := DnsDSRecordResourceModel{}

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

					state.Record = flattenDnsDSRecords(props.DSRecords)

					state.Tags = pointer.From(props.Metadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (DnsDSRecordResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			var model DnsDSRecordResourceModel

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
				payload.Properties.DSRecords = expandDnsDSRecords(model.Record)
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

func (DnsDSRecordResource) Delete() sdk.ResourceFunc {
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

func flattenDnsDSRecords(records *[]recordsets.DsRecord) []DnsDSRecordResourceRecord {
	results := make([]DnsDSRecordResourceRecord, 0)

	if records != nil {
		for _, record := range *records {
			result := DnsDSRecordResourceRecord{
				Algorithm: pointer.From(record.Algorithm),
				KeyTag:    pointer.From(record.KeyTag),
			}

			if record.Digest != nil {
				result.DigestType = pointer.From(record.Digest.AlgorithmType)
				result.DigestValue = pointer.From(record.Digest.Value)
			}
			results = append(results, result)
		}
	}

	return results
}

func expandDnsDSRecords(d []DnsDSRecordResourceRecord) *[]recordsets.DsRecord {
	records := make([]recordsets.DsRecord, 0)

	for _, v := range d {
		records = append(records, recordsets.DsRecord{
			Algorithm: pointer.To(v.Algorithm),
			KeyTag:    pointer.To(v.KeyTag),
			Digest: &recordsets.Digest{
				AlgorithmType: pointer.To(v.DigestType),
				Value:         pointer.To(v.DigestValue),
			},
		})
	}

	return &records
}
