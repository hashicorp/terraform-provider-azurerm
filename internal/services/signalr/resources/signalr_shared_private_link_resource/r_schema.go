// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_shared_private_link_resource

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func signalrSharedPrivateLinkResourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Signalr Shared Private Link Resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_service_id": {
			Description:  "The ID of the Signalr Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateSignalRID,
		},

		"sub_resource_name": {
			Description:  "The sub resource name which the Signalr Private Endpoint can connect to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkSubResourceName,
		},

		"target_resource_id": {
			Description:  "The ID of the Shared Private Link Enabled Remote Resource which this Signalr Private Endpoint should be connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"request_message": {
			Description:  "The request message for requesting approval of the Shared Private Link Enabled Remote Resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"status": {
			Description: "The status of a Private Endpoint connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
