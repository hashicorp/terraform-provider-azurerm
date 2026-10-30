// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_cname_record

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name private_dns_cname_record -properties "name" -compare-values "subscription_id:private_dns_zone_id,resource_group_name:private_dns_zone_id,private_dns_zone_name:private_dns_zone_id,record_type:id"

const azurePrivateDnsCNameRecordResourceName = "azurerm_private_dns_cname_record"

func RegisterResource() *pluginsdk.Resource {
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

		Schema: privateDnsCnameRecordSchema(),
	}
}
