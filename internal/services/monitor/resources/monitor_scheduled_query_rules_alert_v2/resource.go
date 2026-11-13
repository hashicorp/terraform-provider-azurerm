// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_monitor_scheduled_query_rules_alert_v2"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &scheduledqueryrules.ScheduledQueryRuleId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledqueryrules.ValidateScheduledQueryRuleID
}
