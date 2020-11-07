package azure

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/shopspring/decimal"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	consumptionBudgetIDSeparator = "/providers/Microsoft.Consumption/budgets/"
)

type ConsumptionBudgetID struct {
	Name  string
	Scope string
}

// validation
func ValidateAzureConsumptionBudgetName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[-_a-zA-Z0-9]{1,63}$"),
		"The consumption budget name can contain only letters, numbers, underscores, and hyphens. The consumption budget name be between 6 and 63 characters long.",
	)
}

func ParseAzureConsumptionBudgetID(id string) (*ConsumptionBudgetID, error) {
	if !strings.Contains(id, consumptionBudgetIDSeparator) {
		return nil, fmt.Errorf("the provided ID %q does not contain the expected resource provider %s", id, consumptionBudgetIDSeparator)
	}

	components := strings.Split(id, consumptionBudgetIDSeparator)

	if len(components) != 2 {
		return nil, fmt.Errorf("the provided ID %q is not in the expected format", id)
	}

	name := components[1]
	scope := components[0]

	if name == "" || scope == "" {
		return nil, fmt.Errorf("name or scope cannot be empty strings. Name: '%s', Scope: '%s'", name, scope)
	}

	return &ConsumptionBudgetID{
		Name:  components[1],
		Scope: components[0],
	}, nil
}

// expand and flatten
func ExpandAzureConsumptionBudgetTimePeriod(i []interface{}) (*consumption.BudgetTimePeriod, error) {
	if len(i) == 0 || i[0] == nil {
		return nil, nil
	}

	input := i[0].(map[string]interface{})
	timePeriod := consumption.BudgetTimePeriod{}

	if startDateRaw, ok := input["start_date"]; ok {
		startDate, err := date.ParseTime(time.RFC3339, startDateRaw.(string))
		if err != nil {
			return nil, err
		}

		timePeriod.StartDate = &date.Time{
			Time: startDate,
		}
	}

	if endDateRaw, ok := input["end_date"]; ok {
		if endDateRaw.(string) != "" {
			endDate, err := date.ParseTime(time.RFC3339, endDateRaw.(string))
			if err != nil {
				return nil, err
			}

			timePeriod.EndDate = &date.Time{
				Time: endDate,
			}
		}
	}

	return &timePeriod, nil
}

func FlattenAzureConsumptionBudgetTimePeriod(input *consumption.BudgetTimePeriod) []interface{} {
	if input == nil {
		return nil
	}

	timePeriod := make([]interface{}, 0)
	timePeriodBlock := make(map[string]interface{})

	timePeriodBlock["start_date"] = input.StartDate.String()
	timePeriodBlock["end_date"] = input.EndDate.String()

	return append(timePeriod, timePeriodBlock)
}

func ExpandAzureConsumptionBudgetNotifications(input []interface{}) map[string]*consumption.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]*consumption.Notification)

	for _, v := range input {
		notificationRaw := v.(map[string]interface{})
		notification := consumption.Notification{}

		if enabled, ok := notificationRaw["enabled"].(bool); ok {
			notification.Enabled = utils.Bool(enabled)
		}

		if operator, ok := notificationRaw["operator"].(string); ok {
			notification.Operator = consumption.OperatorType(operator)
		}

		if threshold, ok := notificationRaw["threshold"].(int); ok {
			thresholdDecimal := decimal.NewFromInt(int64(threshold))
			notification.Threshold = &thresholdDecimal
		}

		if contactEmails, ok := notificationRaw["contact_emails"].([]interface{}); ok {
			notification.ContactEmails = utils.ExpandStringSlice(contactEmails)
		}

		if contactRoles, ok := notificationRaw["contact_roles"].([]interface{}); ok {
			notification.ContactRoles = utils.ExpandStringSlice(contactRoles)
		}

		if contactGroups, ok := notificationRaw["contact_groups"].([]interface{}); ok {
			notification.ContactGroups = utils.ExpandStringSlice(contactGroups)
		}

		notificationKey := fmt.Sprintf("actual_%s_%s_Percent", string(notification.Operator), notification.Threshold.StringFixed(0))
		notifications[notificationKey] = &notification
	}

	return notifications
}

func FlattenAzureConsumptionBudgetNotifications(input map[string]*consumption.Notification) []interface{} {
	if input == nil {
		return nil
	}

	notifications := make([]interface{}, 0)

	for _, v := range input {
		notificationBlock := make(map[string]interface{})

		notificationBlock["enabled"] = v.Enabled
		notificationBlock["operator"] = utils.String(string(v.Operator))
		threshold, _ := v.Threshold.Float64()
		notificationBlock["threshold"] = utils.Float(threshold)
		notificationBlock["contact_emails"] = v.ContactEmails
		notificationBlock["contact_roles"] = v.ContactRoles
		notificationBlock["contact_groups"] = v.ContactGroups

		notifications = append(notifications, notificationBlock)
	}

	return notifications
}

func ExpandAzureConsumptionBudgetFilterTags(input map[string]interface{}) map[string][]string {
	output := make(map[string][]string, len(input))

	for i, v := range input {
		value := utils.ExpandStringSlice(v.([]interface{}))
		output[i] = *value
	}

	return output
}

func FlattenAzureConsumptionBudgetFilterTags(input map[string][]string) map[string]interface{} {
	output := make(map[string]interface{}, len(input))

	for i, v := range input {
		if v == nil {
			continue
		}
		value := utils.FlattenStringSlice(&v)
		output[i] = value
	}

	return output
}

func ExpandAzureConsumptionBudgetFilter(i []interface{}) *consumption.Filters {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	filters := consumption.Filters{}

	input := i[0].(map[string]interface{})

	if resourceGroups, ok := input["resource_groups"].([]interface{}); ok {
		filters.ResourceGroups = utils.ExpandStringSlice(resourceGroups)
	}

	if resources, ok := input["resources"].([]interface{}); ok {
		filters.Resources = utils.ExpandStringSlice(resources)
	}

	if tags, ok := input["tags"].(map[string]interface{}); ok {
		filters.Tags = ExpandAzureConsumptionBudgetFilterTags(tags)
	}

	if meters, ok := input["meters"].([]interface{}); ok {
		filters.Meters = utils.ExpandUUIDSlice(meters)
	}

	return &filters
}

func FlattenAzureConsumptionBudgetFilter(input *consumption.Filters) []interface{} {
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
		filterBlock["tags"] = FlattenAzureConsumptionBudgetFilterTags(input.Tags)
	}

	if len(filterBlock) != 0 {
		filters = append(filters, filterBlock)
	}

	return filters
}
