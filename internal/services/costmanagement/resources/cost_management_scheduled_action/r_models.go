// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_management_scheduled_action

type CostManagementScheduledActionModel struct {
	Name               string   `tfschema:"name"`
	DisplayName        string   `tfschema:"display_name"`
	ViewId             string   `tfschema:"view_id"`
	EmailSubject       string   `tfschema:"email_subject"`
	EmailAddresses     []string `tfschema:"email_addresses"`
	Message            string   `tfschema:"message"`
	EmailAddressSender string   `tfschema:"email_address_sender"`
	Frequency          string   `tfschema:"frequency"`
	DaysOfWeek         []string `tfschema:"days_of_week"`
	WeeksOfMonth       []string `tfschema:"weeks_of_month"`
	HourOfDay          int64    `tfschema:"hour_of_day"`
	DayOfMonth         int64    `tfschema:"day_of_month"`
	StartDate          string   `tfschema:"start_date"`
	EndDate            string   `tfschema:"end_date"`
}

func (r Resource) ModelObject() any {
	return &CostManagementScheduledActionModel{}
}
