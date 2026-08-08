// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name private_dns_cname_record -properties "name" -compare-values "subscription_id:private_dns_zone_id,resource_group_name:private_dns_zone_id,private_dns_zone_name:private_dns_zone_id,record_type:id"

const azurePrivateDnsCNameRecordResourceName = "azurerm_private_dns_cname_record"

func resourcePrivateDnsCNameRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsCNameRecordCreateUpdate,
		Read:   resourcePrivateDnsCNameRecordRead,
		Update: resourcePrivateDnsCNameRecordCreateUpdate,
		Delete: resourcePrivateDnsCNameRecordDelete,

		Importer: pluginsdk.ImporterValidatingIdentityThen(&privatedns.RecordTypeId{}, resourcePrivateDnsCNameRecordImporter),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&privatedns.RecordTypeId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"private_dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
			},

			"record": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"ttl": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, math.MaxInt32),
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsCNameRecordImporter(_ context.Context, d *pluginsdk.ResourceData, _ interface{}) ([]*pluginsdk.ResourceData, error) {
	resourceId, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}
	if resourceId.RecordType != privatedns.RecordTypeCNAME {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s wrong type received: expected %s received %s", resourceId, privatedns.RecordTypeCNAME, resourceId.RecordType)
	}
	return []*pluginsdk.ResourceData{d}, nil
}

func resourcePrivateDnsCNameRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	privateDNSZoneID, err := privatezones.ParsePrivateDnsZoneID(d.Get("private_dns_zone_id").(string))
	if err != nil {
		return err
	}

	id := privatedns.NewRecordTypeID(privateDNSZoneID.SubscriptionId, privateDNSZoneID.ResourceGroupName, privateDNSZoneID.PrivateDnsZoneName, privatedns.RecordTypeCNAME, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.RecordSetsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(azurePrivateDnsCNameRecordResourceName, id.ID())
			}
		}
	}

	parameters := privatedns.RecordSet{
		Name: pointer.To(id.RelativeRecordSetName),
		Properties: &privatedns.RecordSetProperties{
			Metadata: tags.Expand(d.Get("tags").(map[string]interface{})),
			Ttl:      pointer.To(int64(d.Get("ttl").(int))),
			CnameRecord: &privatedns.CnameRecord{
				Cname: pointer.To(d.Get("record").(string)),
			},
		},
	}

	options := privatedns.RecordSetsCreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}
	if _, err := client.RecordSetsCreateOrUpdate(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	}

	return resourcePrivateDnsCNameRecordRead(d, meta)
}

func resourcePrivateDnsCNameRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.RecordSetsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return resourcePrivateDnsCNameRecordFlatten(d, id, resp.Model)
}

func resourcePrivateDnsCNameRecordFlatten(d *pluginsdk.ResourceData, id *privatedns.RecordTypeId, model *privatedns.RecordSet) error {
	d.Set("name", id.RelativeRecordSetName)
	d.Set("private_dns_zone_id", privatezones.NewPrivateDnsZoneID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if record := props.CnameRecord; record != nil {
				d.Set("record", record.Cname)
			}

			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourcePrivateDnsCNameRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	options := privatedns.RecordSetsDeleteOperationOptions{IfMatch: pointer.To("")}

	if _, err = dnsClient.RecordSetsDelete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
