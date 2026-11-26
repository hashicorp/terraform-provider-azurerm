// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Confidential Ledger.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConfidentialLedgerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"ledger_type": {
			Description: "The type of Confidential Ledger.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
						Description:  "The Ledger Role to grant this AzureAD Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(confidentialledger.PossibleValuesForLedgerRoleName(), false),
					},
					"principal_id": {
						Description:  "The Principal ID of the AzureAD Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"tenant_id": {
						Description:  "The Tenant ID for this AzureAD Service Principal.",
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
						Description:  "The Ledger Role to grant this Certificate Security Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(confidentialledger.PossibleValuesForLedgerRoleName(), false),
					},
					"pem_public_key": {
						Description:  "The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"identity_service_endpoint": {
			Description: "The Identity Service Endpoint for this Confidential Ledger.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"ledger_endpoint": {
			Description: "The Endpoint for this Confidential Ledger.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
