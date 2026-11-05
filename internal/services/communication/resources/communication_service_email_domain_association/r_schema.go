// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service_email_domain_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"communication_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: communicationservices.ValidateCommunicationServiceID,
		},
		"email_service_domain_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: domains.ValidateDomainID,
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
