// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_cname_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	cdn "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	frontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceDnsCNameRecordFlatten(d *pluginsdk.ResourceData, id *recordsets.RecordTypeId, model *recordsets.RecordSet) error {
	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model != nil {
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
				} else if trafficManagerID, err := trafficmanagers.ParseEndpointTypeIDInsensitively(targetResourceId); err == nil {
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

	return pluginsdk.SetResourceIdentityData(d, id)
}
