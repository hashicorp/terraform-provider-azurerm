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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = DnsARecordResource{}
	_ sdk.ResourceWithUpdate = DnsARecordResource{}
)

type DnsARecordResource struct{}

func (DnsARecordResource) ModelObject() interface{} {
	return &DnsARecordResourceModel{}
}

func (DnsARecordResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRecordTypeID(recordsets.RecordTypeA)
}

func (DnsARecordResource) ResourceType() string {
	return "azurerm_dns_a_record"
}

type DnsARecordResourceModel struct {
	Name              string            `tfschema:"name"`
	ZoneName          string            `tfschema:"zone_name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Ttl               int64             `tfschema:"ttl"`
	Records           []string          `tfschema:"records"`
	Tags              map[string]string `tfschema:"tags"`
	Fqdn              string            `tfschema:"fqdn"`
	TargetResourceId  string            `tfschema:"target_resource_id"`
}

func (DnsARecordResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"records": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Elem:         &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"records", "target_resource_id"},
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ExactlyOneOf: []string{"records", "target_resource_id"},
		},

		"tags": commonschema.Tags(),
	}
}

func (DnsARecordResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DnsARecordResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DnsARecordResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			zoneId := zones.NewDnsZoneID(subscriptionId, model.ResourceGroupName, model.ZoneName)

			id := recordsets.NewRecordTypeID(subscriptionId, zoneId.ResourceGroupName, zoneId.DnsZoneName, recordsets.RecordTypeA, model.Name)
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
					Metadata:       pointer.To(model.Tags),
					TTL:            pointer.To(model.Ttl),
					ARecords:       expandAzureRmDnsARecords(model.Records),
					TargetResource: &recordsets.SubResource{},
				},
			}

			if model.TargetResourceId != "" {
				parameters.Properties.TargetResource.Id = pointer.To(model.TargetResourceId)
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (DnsARecordResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			state := DnsARecordResourceModel{}

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
			state.Name = id.RelativeRecordSetName
			state.ZoneName = id.DnsZoneName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Ttl = pointer.From(props.TTL)
					state.Fqdn = pointer.From(props.Fqdn)

					state.Records = flattenAzureRmDnsARecords(props.ARecords)
					state.TargetResourceId = pointer.From(props.TargetResource.Id)

					state.Tags = pointer.From(props.Metadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (DnsARecordResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.RecordSets

			var model DnsARecordResourceModel

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

			if metadata.ResourceData.HasChange("records") {
				payload.Properties.ARecords = expandAzureRmDnsARecords(model.Records)
			}

			if metadata.ResourceData.HasChange("ttl") {
				payload.Properties.TTL = pointer.To(model.Ttl)
			}

			if metadata.ResourceData.HasChange("target_resource_id") {
				payload.Properties.TargetResource = &recordsets.SubResource{}
				if targetId := model.TargetResourceId; targetId != "" {
					payload.Properties.TargetResource.Id = pointer.To(targetId)
				}
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

func (DnsARecordResource) Delete() sdk.ResourceFunc {
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

func expandAzureRmDnsARecords(input []string) *[]recordsets.ARecord {
	records := make([]recordsets.ARecord, len(input))

	for i, v := range input {
		records[i] = recordsets.ARecord{
			IPv4Address: &v,
		}
	}

	return &records
}

func flattenAzureRmDnsARecords(records *[]recordsets.ARecord) []string {
	if records == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, record := range *records {
		if record.IPv4Address == nil {
			continue
		}

		results = append(results, *record.IPv4Address)
	}

	return results
}
