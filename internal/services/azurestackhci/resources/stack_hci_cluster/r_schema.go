// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func stackHciClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Stack HCI Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 260),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"client_id": {
			Description:  "The Client ID of the Azure Active Directory Application which is used by the Azure Stack HCI Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Description:  "The Tenant ID of the Azure Active Directory which is used by the Azure Stack HCI Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"automanage_configuration_id": {
			Description:  "The ID of the Automanage Configuration assigned to the Azure Stack HCI Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
		},

		"cloud_id": {
			Description: "An immutable UUID for the Azure Stack HCI Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"service_endpoint": {
			Description: "The region specific Data Path Endpoint of the Azure Stack HCI Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_provider_object_id": {
			Description: "The object ID of the Resource Provider Service Principal.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}
