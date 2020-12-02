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

func ValidateAzureConsumptionBudgetTimePeriodStartDate(i interface{}, k string) (warnings []string, errors []error) {
	validateRFC3339TimeWarnings, validateRFC3339TimeErrors := validation.IsRFC3339Time(i, k)
	errors = append(errors, validateRFC3339TimeErrors...)
	warnings = append(warnings, validateRFC3339TimeWarnings...)

	if len(errors) != 0 || len(warnings) != 0 {
		return warnings, errors
	}

	// Errors were already checked by validation.IsRFC3339Time
	startDate, _ := date.ParseTime(time.RFC3339, i.(string))

	// The start date must be first of the month
	if startDate.Day() != 1 {
		errors = append(errors, fmt.Errorf("%q must be first of the month, got day %d", k, startDate.Day()))
		return warnings, errors
	}

	// Budget start date must be on or after June 1, 2017.
	earliestPossibleStartDateString := "2017-06-01T00:00:00Z"
	earliestPossibleStartDate, _ := date.ParseTime(time.RFC3339, earliestPossibleStartDateString)
	if startDate.Before(earliestPossibleStartDate) {
		errors = append(errors, fmt.Errorf("%q must be on or after June 1, 2017, got %q", k, i.(string)))
		return warnings, errors
	}

	// Future start date should not be more than twelve months.
	if startDate.After(time.Now().AddDate(0, 12, 0)) {
		warnings = append(warnings, fmt.Sprintf("%q should not be more than twelve months in the future", k))
	}

	return warnings, errors
}

// parse
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

func FlattenAzureConsumptionBudgetNotifications(input map[string]*consumption.Notification) []interface{} {
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

func ExpandAzureConsumptionBudgetFilterTags(input []interface{}) map[string][]string {
	output := make(map[string][]string, len(input))

	for _, v := range input {
		tagInput := v.(map[string]interface{})

		values := utils.ExpandStringSlice(tagInput["values"].([]interface{}))
		output[tagInput["name"].(string)] = *values
	}

	return output
}

func FlattenAzureConsumptionBudgetFilterTags(input map[string][]string) []interface{} {
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

func ExpandAzureConsumptionBudgetFilter(i []interface{}) *consumption.Filters {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	filters := consumption.Filters{}

	input := i[0].(map[string]interface{})

	filters.ResourceGroups = utils.ExpandStringSlice(input["resource_groups"].([]interface{}))
	filters.Resources = utils.ExpandStringSlice(input["resources"].([]interface{}))
	filters.Tags = ExpandAzureConsumptionBudgetFilterTags(input["tag"].(*schema.Set).List())
	filters.Meters = utils.ExpandUUIDSlice(input["meters"].([]interface{}))

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
		filterBlock["tag"] = schema.NewSet(schema.HashResource(SchemaAzureConsumptionBudgetFilterTagElement()), FlattenAzureConsumptionBudgetFilterTags(input.Tags))
	}

	if len(filterBlock) != 0 {
		filters = append(filters, filterBlock)
	}

	return filters
}

// DiffSuppressFuncs
func DiffSuppressFuncAzureConsumptionBudgetTimePeriodEndDate(k, old, new string, d *schema.ResourceData) bool {
	// If the end date is not set by the user, Azure defaults this to 10 years
	// from the start date. Therefore, the diff is suppressed if the user didn't
	// set the value previously or in the current config
	return new == ""
}

// schema
func SchemaAzureConsumptionBudgetResourceGroupResource() map[string]*schema.Schema {
	resourceGroupNameSchema := map[string]*schema.Schema{
		"resource_group_name": SchemaResourceGroupName(),
	}

	return MergeSchema(SchemaAzureConsumptionBudgetSubscriptionResource(), resourceGroupNameSchema)
}

func SchemaAzureConsumptionBudgetSubscriptionResource() map[string]*schema.Schema {
	subscriptionIDSchema := map[string]*schema.Schema{
		"subscription_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}

	return MergeSchema(SchemaAzureConsumptionBudgetCommonResource(), subscriptionIDSchema)
}

func SchemaAzureConsumptionBudgetFilterTagElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaAzureConsumptionBudgetNotificationElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.EqualTo),
					string(consumption.GreaterThan),
					string(consumption.GreaterThanOrEqualTo),
				}, false),
			},

			"contact_emails": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaAzureConsumptionBudgetCommonResource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateAzureConsumptionBudgetName(),
		},

		"amount": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(1.0),
		},

		"category": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(consumption.Cost),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(consumption.Cost),
				string(consumption.Usage),
			}, false),
		},

		"filter": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource_groups": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"resources": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: ValidateResourceID,
						},
					},
					"meters": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
					"tag": {
						Type:     schema.TypeSet,
						Optional: true,
						Set:      schema.HashResource(SchemaAzureConsumptionBudgetFilterTagElement()),
						Elem:     SchemaAzureConsumptionBudgetFilterTagElement(),
					},
				},
			},
		},

		"notification": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Set:      schema.HashResource(SchemaAzureConsumptionBudgetNotificationElement()),
			Elem:     SchemaAzureConsumptionBudgetNotificationElement(),
		},

		"time_grain": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(consumption.TimeGrainTypeMonthly),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(consumption.TimeGrainTypeBillingAnnual),
				string(consumption.TimeGrainTypeBillingMonth),
				string(consumption.TimeGrainTypeBillingQuarter),
				string(consumption.TimeGrainTypeAnnually),
				string(consumption.TimeGrainTypeMonthly),
				string(consumption.TimeGrainTypeQuarterly),
			}, false),
		},

		"time_period": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_date": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: ValidateAzureConsumptionBudgetTimePeriodStartDate,
						ForceNew:     true,
					},
					"end_date": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: DiffSuppressFuncAzureConsumptionBudgetTimePeriodEndDate,
					},
				},
			},
		},
	}
}
