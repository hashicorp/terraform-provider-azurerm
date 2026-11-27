// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_managed_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func appServiceManagedCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_hostname_binding_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webapps.ValidateHostNameBindingID,
		},

		"canonical_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subject_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_names": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"issuer": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"issue_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
