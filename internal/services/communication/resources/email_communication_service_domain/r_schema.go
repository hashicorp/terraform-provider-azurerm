// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Email Communication Service resource. If `domain_management` is `AzureManaged`, the name must be `AzureManagedDomain`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"email_service_id": {
			Description:  "The resource ID of the Email Communication Service where the Domain belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: domains.ValidateEmailServiceID,
		},

		"domain_management": {
			Description:  "Describes how a Domains resource is being managed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(domains.PossibleValuesForDomainManagement(), false),
		},

		"user_engagement_tracking_enabled": {
			Description: "Whether user engagement tracking is enabled or disabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"from_sender_domain": {
			Description: "P2 sender domain that is displayed to the email recipients [RFC 5322].",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"mail_from_sender_domain": {
			Description: "P1 sender domain that is present on the email envelope [RFC 5321].",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"verification_records": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain": helpers.DomainVerificationRecordsCommonSchema(),
					"dkim":   helpers.DomainVerificationRecordsCommonSchema(),
					"dkim2":  helpers.DomainVerificationRecordsCommonSchema(),
					"dmarc":  helpers.DomainVerificationRecordsCommonSchema(),
					"spf":    helpers.DomainVerificationRecordsCommonSchema(),
				},
			},
		},
	}
}
