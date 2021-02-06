package consumption

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/shopspring/decimal"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"time"
)

// expand and flatten
func ExpandConsumptionBudgetTimePeriod(i []interface{}) (*consumption.BudgetTimePeriod, error) {
	if len(i) == 0 || i[0] == nil {
		return nil, nil
	}

	input := i[0].(map[string]interface{})
	timePeriod := consumption.BudgetTimePeriod{}

	if startDateInput, ok := input["start_date"].(string); ok {
		startDate, err := date.ParseTime(time.RFC3339, startDateInput)
		if err != nil {
			return nil, fmt.Errorf("start_date '%s' was not in the correct format: %+v", startDateInput, err)
		}

		timePeriod.StartDate = &date.Time{
			Time: startDate,
		}
	}

	if endDateInput, ok := input["end_date"].(string); ok {
		if endDateInput != "" {
			endDate, err := date.ParseTime(time.RFC3339, endDateInput)
			if err != nil {
				return nil, fmt.Errorf("end_date '%s' was not in the correct format: %+v", endDateInput, err)
			}

			timePeriod.EndDate = &date.Time{
				Time: endDate,
			}
		}
	}

	return &timePeriod, nil
}

func FlattenConsumptionBudgetTimePeriod(input *consumption.BudgetTimePeriod) []interface{} {
	if input == nil {
		return nil
	}

	timePeriod := make([]interface{}, 0)
	timePeriodBlock := make(map[string]interface{})

	timePeriodBlock["start_date"] = input.StartDate.String()
	timePeriodBlock["end_date"] = input.EndDate.String()

	return append(timePeriod, timePeriodBlock)
}

func ExpandConsumptionBudgetNotifications(input []interface{}) map[string]*consumption.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]*consumption.Notification)

	for _, v := range input {
		notificationRaw := v.(map[string]interface{})
		notification := consumption.Notification{}

		notification.Enabled = utils.Bool(notificationRaw["enabled"].(bool))
		notification.Operator = consumption.OperatorType(notificationRaw["operator"].(string))

		thresholdDecimal := decimal.NewFromInt(int64(notificationRaw["threshold"].(int)))
		notification.Threshold = &thresholdDecimal

		notification.ContactEmails = utils.ExpandStringSlice(notificationRaw["contact_emails"].([]interface{}))
		notification.ContactRoles = utils.ExpandStringSlice(notificationRaw["contact_roles"].([]interface{}))
		notification.ContactGroups = utils.ExpandStringSlice(notificationRaw["contact_groups"].([]interface{}))

		notificationKey := fmt.Sprintf("actual_%s_%s_Percent", string(notification.Operator), notification.Threshold.StringFixed(0))
		notifications[notificationKey] = &notification
	}

	return notifications
}

func FlattenConsumptionBudgetNotifications(input map[string]*consumption.Notification) []interface{} {
	if input == nil {
		return nil
	}

	notifications := make([]interface{}, 0)

	for _, v := range input {
		notificationBlock := make(map[string]interface{})

		notificationBlock["enabled"] = *v.Enabled
		notificationBlock["operator"] = string(v.Operator)
		threshold, _ := v.Threshold.Float64()
		notificationBlock["threshold"] = int(threshold)
		notificationBlock["contact_emails"] = utils.FlattenStringSlice(v.ContactEmails)
		notificationBlock["contact_roles"] = utils.FlattenStringSlice(v.ContactRoles)
		notificationBlock["contact_groups"] = utils.FlattenStringSlice(v.ContactGroups)

		notifications = append(notifications, notificationBlock)
	}

	return notifications
}

func ExpandConsumptionBudgetFilterTags(input []interface{}) map[string][]string {
	output := make(map[string][]string, len(input))

	for _, v := range input {
		tagInput := v.(map[string]interface{})

		values := utils.ExpandStringSlice(tagInput["values"].([]interface{}))
		output[tagInput["name"].(string)] = *values
	}

	return output
}

func FlattenConsumptionBudgetFilterTags(input map[string][]string) []interface{} {
	output := make([]interface{}, 0)

	for i, v := range input {
		if v == nil {
			continue
		}
		value := utils.FlattenStringSlice(&v)

		tagBlock := make(map[string]interface{})
		tagBlock["name"] = i
		tagBlock["values"] = value

		output = append(output, tagBlock)
	}

	return output
}

func ExpandConsumptionBudgetFilter(i []interface{}) *consumption.Filters {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	filters := consumption.Filters{}

	input := i[0].(map[string]interface{})

	filters.ResourceGroups = utils.ExpandStringSlice(input["resource_groups"].([]interface{}))
	filters.Resources = utils.ExpandStringSlice(input["resources"].([]interface{}))
	filters.Tags = ExpandConsumptionBudgetFilterTags(input["tag"].(*schema.Set).List())
	filters.Meters = utils.ExpandUUIDSlice(input["meters"].([]interface{}))

	return &filters
}

func FlattenConsumptionBudgetFilter(input *consumption.Filters) []interface{} {
	filters := make([]interface{}, 0)

	if input == nil {
		return filters
	}

	filterBlock := make(map[string]interface{})

	if input.ResourceGroups != nil && len(*input.ResourceGroups) != 0 {
		filterBlock["resource_groups"] = utils.FlattenStringSlice(input.ResourceGroups)
	}

	if input.Resources != nil && len(*input.Resources) != 0 {
		filterBlock["resources"] = utils.FlattenStringSlice(input.Resources)
	}

	if input.Meters != nil && len(*input.Meters) != 0 {
		filterBlock["meters"] = utils.FlattenUUIDSlice(input.Meters)
	}

	if input.Tags != nil && len(input.Tags) != 0 {
		filterBlock["tag"] = schema.NewSet(schema.HashResource(SchemaAzureConsumptionBudgetFilterTagElement()), FlattenConsumptionBudgetFilterTags(input.Tags))
	}

	if len(filterBlock) != 0 {
		filters = append(filters, filterBlock)
	}

	return filters
}

// DiffSuppressFuncs
func DiffSuppressFuncConsumptionBudgetTimePeriodEndDate(k, old, new string, d *schema.ResourceData) bool {
	// If the end date is not set by the user, Azure defaults this to 10 years
	// from the start date. Therefore, the diff is suppressed if the user didn't
	// set the value previously or in the current config
	return new == ""
}
