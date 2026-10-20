// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type AIFoundryProjectModel struct {
	Name                        string                                     `tfschema:"name"`
	Location                    string                                     `tfschema:"location"`
	AIServicesHubId             string                                     `tfschema:"ai_services_hub_id"`
	Identity                    []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	HighBusinessImpactEnabled   bool                                       `tfschema:"high_business_impact_enabled"`
	Description                 string                                     `tfschema:"description"`
	PrimaryUserAssignedIdentity string                                     `tfschema:"primary_user_assigned_identity"`
	FriendlyName                string                                     `tfschema:"friendly_name"`
	ProjectId                   string                                     `tfschema:"project_id"`
	Tags                        map[string]any                             `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AIFoundryProjectModel{}
}
