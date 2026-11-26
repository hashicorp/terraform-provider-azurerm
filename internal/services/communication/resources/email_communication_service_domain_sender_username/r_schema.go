// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain_sender_username

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/senderusernames"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Email Communication Service Domain Sender Username resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 253),
		},

		"email_service_domain_id": commonschema.ResourceIDReferenceRequiredForceNew(&senderusernames.DomainId{}, "Email Communication Service Domain"),

		"display_name": {
			Description:  "The display name for the Email Communication Service Domain Sender Username resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 253),
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
