// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func healthcareMedtechServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MedTechServiceName(),
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"eventhub_namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubNamespaceName(),
		},

		"eventhub_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_consumer_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
		},

		"device_mapping_json": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
		},

		"tags": commonschema.Tags(),
	}
}
