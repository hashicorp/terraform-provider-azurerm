// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_schedule

import (
	"context"
	"fmt"
	"strings"
	"time"

	// import time/tzdata to embed timezone information in the program
	// add this to resolve https://github.com/hashicorp/terraform-provider-azurerm/issues/20690
	_ "time/tzdata"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationScheduleCreate,
		Read:   resourceAutomationScheduleRead,
		Update: resourceAutomationScheduleUpdate,
		Delete: resourceAutomationScheduleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schedule.ParseScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: automationScheduleSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			frequency := strings.ToLower(diff.Get("frequency").(string))
			interval, _ := diff.GetOk("interval")
			if frequency == "onetime" && interval.(int) > 0 {
				// because `interval` is optional and computed, so interval value can exist even it removed from configuration
				// have to check it in raw config
				intervalVal := diff.GetRawConfig().GetAttr("interval")
				if !intervalVal.IsNull() {
					return fmt.Errorf("`interval` cannot be set when frequency is `OneTime`")
				}
			}

			_, hasWeekDays := diff.GetOk("week_days")
			if hasWeekDays && frequency != "week" {
				return fmt.Errorf("`week_days` can only be set when frequency is `Week`")
			}

			_, hasMonthDays := diff.GetOk("month_days")
			if hasMonthDays && frequency != "month" {
				return fmt.Errorf("`month_days` can only be set when frequency is `Month`")
			}

			_, hasMonthlyOccurrences := diff.GetOk("monthly_occurrence")
			if hasMonthlyOccurrences && frequency != "month" {
				return fmt.Errorf("`monthly_occurrence` can only be set when frequency is `Month`")
			}

			return nil
		}),
	}
}
