// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type MonitorsResourceModel struct {
	Name                          string                         `tfschema:"name"`
	ResourceGroup                 string                         `tfschema:"resource_group_name"`
	Location                      string                         `tfschema:"location"`
	MonitoringStatus              bool                           `tfschema:"monitoring_enabled"`
	MarketplaceSubscriptionStatus string                         `tfschema:"marketplace_subscription"`
	Identity                      []identity.ModelSystemAssigned `tfschema:"identity"`
	EnvironmentProperties         []EnvironmentProperties        `tfschema:"environment_properties"`
	PlanData                      []PlanData                     `tfschema:"plan"`
	UserInfo                      []UserInfo                     `tfschema:"user"`
	Tags                          map[string]string              `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MonitorsResourceModel{}
}
