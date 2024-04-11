// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AppServiceCertificateOrderResourceV0ToV1{}

type AppServiceCertificateOrderResourceV0ToV1 struct{}

func (AppServiceCertificateOrderResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"auto_renew": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"certificates": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key_vault_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key_vault_secret_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"provisioning_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"csr": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"distinguished_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"key_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"product_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"validity_in_years": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"domain_verification_token": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"is_private_key_external": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"app_service_certificate_not_renewable_reasons": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"signed_certificate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"root_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"intermediate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (AppServiceCertificateOrderResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := parse.CertificateOrderOldID(oldIdRaw)
		if err != nil {
			return rawState, fmt.Errorf("parsing ID %q to upgrade: %+v", oldIdRaw, err)
		}

		appServiceCertOrderId := parse.NewCertificateOrderID(oldId.SubscriptionId, oldId.ResourceGroup, oldId.CertificateOrderName)
		newId := appServiceCertOrderId.ID()

		rawState["id"] = newId
		return rawState, nil
	}
}
