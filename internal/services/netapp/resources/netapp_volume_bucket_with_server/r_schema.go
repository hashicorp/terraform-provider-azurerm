// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket_with_server

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	arguments := helpers.NetAppBucketResourceCommonArguments()

	// When the certificate is sourced from Key Vault the inline `server.0.certificate_pem`
	// must not be supplied, so wire up the mutual exclusivity that only applies when the
	// `server` block exists.
	arguments["key_vault"].ConflictsWith = []string{"server.0.certificate_pem"}

	arguments["server"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"fqdn": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"certificate_pem": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					Sensitive:     true,
					ValidateFunc:  validation.StringIsBase64,
					ConflictsWith: []string{"key_vault"},
				},
				"on_certificate_conflict_action": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(buckets.OnCertificateConflictActionFail),
					ValidateFunc: validation.StringInSlice(buckets.PossibleValuesForOnCertificateConflictAction(), false),
				},
			},
		},
	}

	return arguments
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return helpers.NetAppBucketResourceCommonAttributes()
}
