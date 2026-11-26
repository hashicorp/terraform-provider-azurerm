// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cost_management_scheduled_action

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Cost Management Scheduled Action.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"display_name": {
			Description:  "User visible input name of the Cost Management Scheduled Action.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 25),
		},

		"view_id": {
			Description:  "The ID of the Cost Management View that is used by the Scheduled Action.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_subject": {
			Description:  "Subject of the email.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"email_addresses": {
			Description: "A list of email addresses that will receive the Scheduled Action.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    20,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"message": {
			Description:  "Message to be added in the email.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 250),
		},

		"email_address_sender": {
			Description: "Email address of the point of contact that should get the unsubscribe requests of Scheduled Action notification emails.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because when using User principal to run Terraform (e.g. `az login`), this attribute is optional in the request and will be set as the user's email address in return if omitted.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frequency": {
			Description:  "Frequency of the schedule. Value `Monthly` requires either `weeks_of_month` and `days_of_week` or `day_of_month` to be specified. Value `Weekly` requires `days_of_week` to be specified.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(scheduledactionoperationgroup.PossibleValuesForScheduleFrequency(), false),
		},

		"days_of_week": {
			Description: "A list of day names on which cost analysis data will be emailed. This property is applicable when frequency is `Weekly` or `Monthly`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(scheduledactionoperationgroup.PossibleValuesForDaysOfWeek(), false),
			},
		},

		"weeks_of_month": {
			Description: "A list of weeks in which cost analysis data will be emailed. This property is applicable when `frequency` is `Monthly` and used in combination with `days_of_week`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(scheduledactionoperationgroup.PossibleValuesForWeeksOfMonth(), false),
			},
		},

		"hour_of_day": {
			Description:  "UTC time at which cost analysis data will be emailed.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 23),
		},

		"day_of_month": {
			Description:  "UTC day on which cost analysis data will be emailed. This property is applicable when `frequency` is `Monthly`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 31),
		},

		"start_date": {
			Description:  "The start date and time of the Scheduled Action (UTC).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"end_date": {
			Description:  "The end date and time of the Scheduled Action (UTC).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.Client.Account.AuthenticatedAsAServicePrincipal && metadata.ResourceDiff.Get("email_address_sender") == "" {
				return errors.New("`email_address_sender` is required when authenticated as a service principal")
			}
			return nil
		},
	}
}
