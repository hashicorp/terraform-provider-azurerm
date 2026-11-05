// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_txt_record

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_txt_record/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "dns_zone_name:zone_name,resource_group_name,name" -known-values "record_type:TXT"

const azurermDnsTxtRecordResourceName = "azurerm_dns_txt_record"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDnsTxtRecordCreateUpdate,
		Read:   resourceDnsTxtRecordRead,
		Update: resourceDnsTxtRecordCreateUpdate,
		Delete: resourceDnsTxtRecordDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&recordsets.RecordTypeId{}),
		},
		Importer: pluginsdk.ImporterValidatingIdentityThen(&recordsets.RecordTypeId{}, helper.ResourceDnsRecordImporter(recordsets.RecordTypeTXT)),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TXTRecordV0ToV1{},
		}),

		Schema: dnsTxtRecordSchema(),
	}
}
