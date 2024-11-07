// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdn "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/migration"
	frontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDnsCNameRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDnsCNameRecordCreate,
		Read:   resourceDnsCNameRecordRead,
		Update: resourceDnsCNameRecordUpdate,
		Delete: resourceDnsCNameRecordDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			parsed, err := recordsets.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if parsed.RecordType != recordsets.RecordTypeCNAME {
				return fmt.Errorf("this resource only supports 'CNAME' records")
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CNAMERecordV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
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

			"record": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"record", "target_resource_id"},
			},

			"ttl": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"record", "target_resource_id"},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDnsCNameRecordCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	id := recordsets.NewRecordTypeID(subscriptionId, resGroup, zoneName, recordsets.RecordTypeCNAME, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_dns_cname_record", id.ID())
	}

	ttl := int64(d.Get("ttl").(int))
	record := d.Get("record").(string)
	t := d.Get("tags").(map[string]interface{})
	targetResourceId := d.Get("target_resource_id").(string)

	parameters := recordsets.RecordSet{
		Name: &name,
		Properties: &recordsets.RecordSetProperties{
			Metadata:       tags.Expand(t),
			TTL:            &ttl,
			CNAMERecord:    &recordsets.CnameRecord{},
			TargetResource: &recordsets.SubResource{},
		},
	}

	if record != "" {
		parameters.Properties.CNAMERecord.Cname = pointer.To(record)
	}

	if targetResourceId != "" {
		parameters.Properties.TargetResource.Id = pointer.To(targetResourceId)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDnsCNameRecordRead(d, meta)
}

func resourceDnsCNameRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("fqdn", props.Fqdn)
			d.Set("ttl", props.TTL)

			cname := ""
			if props.CNAMERecord != nil && props.CNAMERecord.Cname != nil {
				cname = *props.CNAMERecord.Cname
			}
			d.Set("record", cname)

			targetResourceId := ""
			if props.TargetResource != nil && props.TargetResource.Id != nil {
				// TODO update this once https://github.com/hashicorp/go-azure-helpers/issues/189 is resolved
				targetResourceId = *props.TargetResource.Id
				if recordTypeID, err := recordsets.ParseRecordTypeIDInsensitively(targetResourceId); err == nil {
					targetResourceId = recordTypeID.ID()
				} else if trafficManagerID, err := endpoints.ParseEndpointTypeIDInsensitively(targetResourceId); err == nil {
					targetResourceId = trafficManagerID.ID()
				} else if cdnID, err := cdn.EndpointIDInsensitively(targetResourceId); err == nil {
					targetResourceId = cdnID.ID()
				} else if frontDoorID, err := frontdoor.FrontendEndpointIDInsensitively(targetResourceId); err == nil {
					targetResourceId = frontDoorID.ID()
				}
			}
			d.Set("target_resource_id", targetResourceId)

			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceDnsCNameRecordUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("record") {
		payload.Properties.CNAMERecord = &recordsets.CnameRecord{}
		if record := d.Get("record").(string); record != "" {
			payload.Properties.CNAMERecord.Cname = pointer.To(record)
		}
	}

	if d.HasChange("target_resource_id") {
		payload.Properties.TargetResource = &recordsets.SubResource{}
		if targetId := d.Get("target_resource_id").(string); targetId != "" {
			payload.Properties.TargetResource.Id = pointer.To(targetId)
		}
	}

	if d.HasChange("ttl") {
		payload.Properties.TTL = pointer.To(int64(d.Get("ttl").(int)))
	}

	if d.HasChange("tags") {
		payload.Properties.Metadata = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDnsCNameRecordRead(d, meta)
}

func resourceDnsCNameRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id, recordsets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
