// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func confidentialLedgerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Required
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConfidentialLedgerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"ledger_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(confidentialledger.LedgerTypePrivate),
				string(confidentialledger.LedgerTypePublic),
			}, false),
		},

		"azuread_based_service_principal": {
			// this is Required since if none are specified then the calling SP gets added
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ledger_role_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(confidentialledger.PossibleValuesForLedgerRoleName(), false),
					},
					"principal_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		// Optional
		"certificate_based_security_principal": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ledger_role_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(confidentialledger.PossibleValuesForLedgerRoleName(), false),
					},
					"pem_public_key": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"identity_service_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"ledger_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
