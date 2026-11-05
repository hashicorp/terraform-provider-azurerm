// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"email_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: domains.ValidateEmailServiceID,
		},

		"domain_management": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(domains.PossibleValuesForDomainManagement(), false),
		},

		"user_engagement_tracking_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"from_sender_domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"mail_from_sender_domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"verification_records": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain": helper.DomainVerificationRecordsCommonSchema(),
					"dkim":   helper.DomainVerificationRecordsCommonSchema(),
					"dkim2":  helper.DomainVerificationRecordsCommonSchema(),
					"dmarc":  helper.DomainVerificationRecordsCommonSchema(),
					"spf":    helper.DomainVerificationRecordsCommonSchema(),
				},
			},
		},
	}
}
