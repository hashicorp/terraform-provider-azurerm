// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_batch_configuration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountBatchConfigurationCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountBatchConfigurationRead,
		Update: resourceLogicAppIntegrationAccountBatchConfigurationCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountBatchConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(id)
			return err
		}),

		Schema: logicAppIntegrationAccountBatchConfigurationSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			frequency := strings.ToLower(diff.Get("release_criteria.0.recurrence.0.frequency").(string))

			_, hasWeekDays := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.week_days")
			if hasWeekDays && frequency != "week" {
				return fmt.Errorf("`week_days` can only be set when frequency is `Week`")
			}

			_, hasMonthDays := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.month_days")
			if hasMonthDays && frequency != "month" {
				return fmt.Errorf("`month_days` can only be set when frequency is `Month`")
			}

			_, hasMonthlyOccurrences := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.monthly")
			if hasMonthlyOccurrences && frequency != "month" {
				return fmt.Errorf("`monthly` can only be set when frequency is `Month`")
			}

			return nil
		}),
	}
}
