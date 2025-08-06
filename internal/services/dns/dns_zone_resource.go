// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate         = DnsZoneResource{}
	_ sdk.ResourceWithStateMigration = DnsZoneResource{}
)

type DnsZoneResource struct{}

func (DnsZoneResource) ModelObject() interface{} {
	return &DnsZoneResourceModel{}
}

func (DnsZoneResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return zones.ValidateDnsZoneID
}

func (DnsZoneResource) ResourceType() string {
	return "azurerm_dns_zone"
}

func (r DnsZoneResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 2,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.DnsZoneV0ToV1{},
			1: migration.DnsZoneV1ToV2{},
		},
	}
}

func (DnsZoneResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"soa_record": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DnsZoneSOARecordEmail,
					},

					"expire_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      2419200,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"minimum_ttl": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"refresh_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"retry_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"serial_number": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"ttl": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},

					"tags": commonschema.Tags(),

					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Computed: true, // (@jackofallops) - This should not be set or updatable to meet API design, see https://learn.microsoft.com/en-us/azure/dns/dns-zones-records#soa-records
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (DnsZoneResource) Attributes() map[string]*pluginsdk.Schema {
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
	}
}

type DnsZoneResourceModel struct {
	Name                  string                           `tfschema:"name"`
	ResourceGroupName     string                           `tfschema:"resource_group_name"`
	NumberOfRecordSets    int64                            `tfschema:"number_of_record_sets"`
	MaxNumberOfRecordSets int64                            `tfschema:"max_number_of_record_sets"`
	NameServers           []string                         `tfschema:"name_servers"`
	SoaRecord             []DnsZoneSoaRecordResourceRecord `tfschema:"soa_record"`
	Tags                  map[string]string                `tfschema:"tags"`
}

type DnsZoneSoaRecordResourceRecord struct {
	Email        string            `tfschema:"email"`
	ExpireTime   int64             `tfschema:"expire_time"`
	MinimumTtl   int64             `tfschema:"minimum_ttl"`
	RefreshTime  int64             `tfschema:"refresh_time"`
	RetryTime    int64             `tfschema:"retry_time"`
	SerialNumber int64             `tfschema:"serial_number"`
	Ttl          int64             `tfschema:"ttl"`
	Fqdn         string            `tfschema:"fqdn"`
	HostName     string            `tfschema:"host_name"`
	Tags         map[string]string `tfschema:"tags"`
}

func (r DnsZoneResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DnsZoneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := zones.NewDnsZoneID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := zones.Zone{
				Location: location.Normalize("global"),
				Tags:     pointer.To(model.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters, zones.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if len(model.SoaRecord) == 1 {
				soaRecordID := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				soaRecordResp, err := recordSetsClient.Get(ctx, soaRecordID)
				if err != nil {
					return fmt.Errorf("retrieving %s to update SOA: %+v", id, err)
				}

				props := soaRecordResp.Model.Properties
				if props == nil || props.SOARecord == nil {
					return fmt.Errorf("could not read SOA properties for %s", id)
				}

				inputSOARecord := expandDNSZoneSOARecord(model.SoaRecord[0])

				inputSOARecord.Host = props.SOARecord.Host

				rsParameters := recordsets.RecordSet{
					Properties: &recordsets.RecordSetProperties{
						TTL:       pointer.To(model.SoaRecord[0].Ttl),
						Metadata:  pointer.To(model.SoaRecord[0].Tags),
						SOARecord: inputSOARecord,
					},
				}

				if len(id.DnsZoneName+strings.TrimSuffix(*rsParameters.Properties.SOARecord.Email, ".")) > 253 {
					return fmt.Errorf("`email` which is concatenated with DNS Zone `name` cannot exceed 253 characters excluding a trailing period")
				}

				soaRecordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				if _, err := recordSetsClient.CreateOrUpdate(ctx, soaRecordId, rsParameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
					return fmt.Errorf("creating %s: %+v", soaRecordId, err)
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (DnsZoneResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			zonesClient := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets

			state := DnsZoneResourceModel{}

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := zonesClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			soaRecord := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
			soaRecordResp, err := recordSetsClient.Get(ctx, soaRecord)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.SoaRecord = flattenDNSZoneSOARecord(soaRecordResp.Model)

			state.Name = id.DnsZoneName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.NumberOfRecordSets = pointer.From(props.NumberOfRecordSets)
					state.MaxNumberOfRecordSets = pointer.From(props.MaxNumberOfRecordSets)
					state.NameServers = pointer.From(props.NameServers)
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DnsZoneResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets

			var model DnsZoneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
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

			payload := zones.Zone{
				Location: existing.Model.Location,
				Tags:     existing.Model.Tags,
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload, zones.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("soa_record") && len(model.SoaRecord) == 1 {
				soaRecordID := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				soaRecordResp, err := recordSetsClient.Get(ctx, soaRecordID)
				if err != nil {
					return fmt.Errorf("retrieving %s to update SOA: %+v", id, err)
				}

				props := soaRecordResp.Model.Properties
				if props == nil || props.SOARecord == nil {
					return fmt.Errorf("could not read SOA properties for %s", id)
				}

				inputSOARecord := expandDNSZoneSOARecord(model.SoaRecord[0])

				inputSOARecord.Host = props.SOARecord.Host

				rsParameters := recordsets.RecordSet{
					Properties: &recordsets.RecordSetProperties{
						TTL:       pointer.To(model.SoaRecord[0].Ttl),
						Metadata:  pointer.To(model.SoaRecord[0].Tags),
						SOARecord: inputSOARecord,
					},
				}

				if len(id.DnsZoneName+strings.TrimSuffix(*rsParameters.Properties.SOARecord.Email, ".")) > 253 {
					return fmt.Errorf("`email` which is concatenated with DNS Zone `name` cannot exceed 253 characters excluding a trailing period")
				}

				soaRecordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				if _, err := recordSetsClient.CreateOrUpdate(ctx, soaRecordId, rsParameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
					return fmt.Errorf("updating %s: %+v", soaRecordId, err)
				}
			}

			return nil
		},
	}
}

func (r DnsZoneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, zones.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDNSZoneSOARecord(input DnsZoneSoaRecordResourceRecord) *recordsets.SoaRecord {
	result := &recordsets.SoaRecord{
		Email:        pointer.To(input.Email),
		ExpireTime:   pointer.To(input.ExpireTime),
		MinimumTTL:   pointer.To(input.MinimumTtl),
		RefreshTime:  pointer.To(input.RefreshTime),
		RetryTime:    pointer.To(input.RetryTime),
		SerialNumber: pointer.To(input.SerialNumber),
	}

	return result
}

func flattenDNSZoneSOARecord(input *recordsets.RecordSet) []DnsZoneSoaRecordResourceRecord {
	output := make([]DnsZoneSoaRecordResourceRecord, 0)
	if input != nil {
		if props := input.Properties; props != nil {
			result := DnsZoneSoaRecordResourceRecord{
				Ttl:  pointer.From(props.TTL),
				Tags: pointer.From(props.Metadata),
				Fqdn: pointer.From(props.Fqdn),
			}

			if record := props.SOARecord; record != nil {
				result.Email = pointer.From(record.Email)
				result.HostName = pointer.From(record.Host)
				result.ExpireTime = pointer.From(record.ExpireTime)
				result.MinimumTtl = pointer.From(record.MinimumTTL)
				result.RefreshTime = pointer.From(record.RefreshTime)
				result.RetryTime = pointer.From(record.RetryTime)
				result.SerialNumber = pointer.From(record.SerialNumber)
			}

			output = append(output, result)
		}
	}

	return output
}
