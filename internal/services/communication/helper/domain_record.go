// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"value": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ttl": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
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
