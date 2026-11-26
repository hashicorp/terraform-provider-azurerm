// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DomainVerificationRecords struct {
	Type  string `tfschema:"type"`
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
	TTL   int64  `tfschema:"ttl"`
}

func DomainVerificationRecordsCommonSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Description: "Type of the DNS record. Example: TXT.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"name": {
					Description: "Name of the DNS record.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"value": {
					Description: "Value of the DNS record.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"ttl": {
					Description: "Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},
			},
		},
	}
}

func DomainVerificationRecordsToModel(record *domains.DnsRecord) []DomainVerificationRecords {
	return []DomainVerificationRecords{{
		Name:  pointer.From(record.Name),
		Type:  pointer.From(record.Type),
		Value: pointer.From(record.Value),
		TTL:   pointer.From(record.Ttl),
	}}
}
