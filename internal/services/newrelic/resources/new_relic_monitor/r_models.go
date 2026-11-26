// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_monitor

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
)

type NewRelicMonitorModel struct {
	Name                  string                         `tfschema:"name"`
	ResourceGroupName     string                         `tfschema:"resource_group_name"`
	AccountCreationSource monitors.AccountCreationSource `tfschema:"account_creation_source"`
	AccountId             string                         `tfschema:"account_id"`
	IngestionKey          string                         `tfschema:"ingestion_key"`
	Location              string                         `tfschema:"location"`
	OrganizationId        string                         `tfschema:"organization_id"`
	OrgCreationSource     monitors.OrgCreationSource     `tfschema:"org_creation_source"`
	PlanData              []PlanDataModel                `tfschema:"plan"`
	UserId                string                         `tfschema:"user_id"`
	UserInfo              []UserInfoModel                `tfschema:"user"`
}

type PlanDataModel struct {
	EffectiveDate string             `tfschema:"effective_date"`
	BillingCycle  string             `tfschema:"billing_cycle"`
	PlanDetails   string             `tfschema:"plan_id"`
	UsageType     monitors.UsageType `tfschema:"usage_type"`
}

type UserInfoModel struct {
	EmailAddress string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	PhoneNumber  string `tfschema:"phone_number"`
}

func (r Resource) ModelObject() any {
	return &NewRelicMonitorModel{}
}
