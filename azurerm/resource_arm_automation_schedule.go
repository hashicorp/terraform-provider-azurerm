package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationScheduleCreateUpdate,
		Read:   resourceArmAutomationScheduleRead,
		Update: resourceArmAutomationScheduleCreateUpdate,
		Delete: resourceArmAutomationScheduleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[^<>*%&:\\?.+/]{0,127}[^<>*%&:\\?.+/\s]$`),
					`The name length must be from 1 to 128 characters. The name cannot contain special characters < > * % & : \ ? . + / and cannot end with a whitespace character.`,
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			//this is AutomationAccountName in the SDK
			"account_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "account_name has been renamed to automation_account_name for clarity and to match the azure API",
				ConflictsWith: []string{"automation_account_name"},
			},

			"automation_account_name": {
				Type:     schema.TypeString,
				Optional: true, //todo change to required once account_name has been removed
				Computed: true,
				//ForceNew:      true, //todo this needs to come back once account_name has been removed
				ConflictsWith: []string{"account_name"},
			},

			"frequency": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.Day),
					string(automation.Hour),
					string(automation.Month),
					string(automation.OneTime),
					string(automation.Week),
				}, true),
			},

			//ignored when frequency is `OneTime`
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true, //defaults to 1 if frequency is not OneTime
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"start_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validate.RFC3339DateInFutureBy(time.Duration(5) * time.Minute),
				//defaults to now + 7 minutes in create function if not set
			},

			"expiry_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true, //same as start time when OneTime, ridiculous value when recurring: "9999-12-31T15:59:00-08:00"
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.RFC3339Time,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
				//todo figure out how to validate this properly
			},

			"week_days": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(automation.Monday),
						string(automation.Tuesday),
						string(automation.Wednesday),
						string(automation.Thursday),
						string(automation.Friday),
						string(automation.Saturday),
						string(automation.Sunday),
					}, true),
				},
				Set:           set.HashStringIgnoreCase,
				ConflictsWith: []string{"month_days", "monthly_occurrence"},
			},

			"month_days": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validate.IntBetweenAndNot(-1, 31, 0),
				},
				Set:           set.HashInt,
				ConflictsWith: []string{"week_days", "monthly_occurrence"},
			},

			"monthly_occurrence": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(automation.Monday),
								string(automation.Tuesday),
								string(automation.Wednesday),
								string(automation.Thursday),
								string(automation.Friday),
								string(automation.Saturday),
								string(automation.Sunday),
							}, true),
						},
						"occurrence": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetweenAndNot(-1, 5, 0),
						},
					},
				},
				ConflictsWith: []string{"week_days", "month_days"},
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			frequency := strings.ToLower(diff.Get("frequency").(string))
			interval, _ := diff.GetOk("interval")
			if frequency == "onetime" && interval.(int) > 0 {
				return fmt.Errorf("`interval` cannot be set when frequency is `OneTime`")
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

			_, hasAccount := diff.GetOk("automation_account_name")
			_, hasAutomationAccountWeb := diff.GetOk("account_name")
			if !hasAccount && !hasAutomationAccountWeb {
				return fmt.Errorf("`automation_account_name` must be set")
			}

			//if automation_account_name changed or account_name changed to or from nil force a new resource
			//remove once we remove the deprecated property
			oAan, nAan := diff.GetChange("automation_account_name")
			if oAan != "" && nAan != "" {
				diff.ForceNew("automation_account_name")
			}

			oAn, nAn := diff.GetChange("account_name")
			if oAn != "" && nAn != "" {
				diff.ForceNew("account_name")
			}

			return nil
		},
	}
}

func resourceArmAutomationScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.ScheduleClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Schedule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	//CustomizeDiff should ensure one of these two is set
	//todo remove this once `account_name` is removed
	accountName := ""
	if v, ok := d.GetOk("automation_account_name"); ok {
		accountName = v.(string)
	} else if v, ok := d.GetOk("account_name"); ok {
		accountName = v.(string)
	}

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Schedule %q (Account %q / Resource Group %q): %s", name, accountName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_schedule", *existing.ID)
		}
	}

	frequency := d.Get("frequency").(string)
	timeZone := d.Get("timezone").(string)
	description := d.Get("description").(string)

	parameters := automation.ScheduleCreateOrUpdateParameters{
		Name: &name,
		ScheduleCreateOrUpdateProperties: &automation.ScheduleCreateOrUpdateProperties{
			Description: &description,
			Frequency:   automation.ScheduleFrequency(frequency),
			TimeZone:    &timeZone,
		},
	}
	properties := parameters.ScheduleCreateOrUpdateProperties

	//start time can default to now + 7 (5 could be invalid by the time the API is called)
	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) //should be validated by the schema
		properties.StartTime = &date.Time{Time: t}
	} else {
		properties.StartTime = &date.Time{Time: time.Now().Add(time.Duration(7) * time.Minute)}
	}

	if v, ok := d.GetOk("expiry_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) //should be validated by the schema
		properties.ExpiryTime = &date.Time{Time: t}
	}

	//only pay attention to interval if frequency is not OneTime, and default it to 1 if not set
	if properties.Frequency != automation.OneTime {
		if v, ok := d.GetOk("interval"); ok {
			properties.Interval = utils.Int32(int32(v.(int)))
		} else {
			properties.Interval = 1
		}
	}

	//only pay attention to the advanced schedule fields if frequency is either Week or Month
	if properties.Frequency == automation.Week || properties.Frequency == automation.Month {
		advancedRef, err := expandArmAutomationScheduleAdvanced(d, d.Id() != "")
		if err != nil {
			return err
		}
		properties.AdvancedSchedule = advancedRef
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, accountName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accountName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Schedule '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationScheduleRead(d, meta)
}

func resourceArmAutomationScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.ScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["schedules"]
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Get(ctx, resGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Schedule '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accountName)
	d.Set("account_name", accountName) //todo remove once `account_name` is removed
	d.Set("frequency", string(resp.Frequency))

	if v := resp.StartTime; v != nil {
		d.Set("start_time", v.Format(time.RFC3339))
	}
	if v := resp.ExpiryTime; v != nil {
		d.Set("expiry_time", v.Format(time.RFC3339))
	}
	if v := resp.Interval; v != nil {
		//seems to me missing its type in swagger, leading to it being a interface{} float64
		d.Set("interval", int(v.(float64)))
	}
	if v := resp.Description; v != nil {
		d.Set("description", v)
	}
	if v := resp.TimeZone; v != nil {
		d.Set("timezone", v)
	}

	if v := resp.AdvancedSchedule; v != nil {
		if err := d.Set("week_days", flattenArmAutomationScheduleAdvancedWeekDays(v)); err != nil {
			return fmt.Errorf("Error setting `week_days`: %+v", err)
		}
		if err := d.Set("month_days", flattenArmAutomationScheduleAdvancedMonthDays(v)); err != nil {
			return fmt.Errorf("Error setting `month_days`: %+v", err)
		}
		if err := d.Set("monthly_occurrence", flattenArmAutomationScheduleAdvancedMonthlyOccurrences(v)); err != nil {
			return fmt.Errorf("Error setting `monthly_occurrence`: %+v", err)
		}
	}
	return nil
}

func resourceArmAutomationScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.ScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["schedules"]
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Delete(ctx, resGroup, accountName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing AzureRM delete request for Automation Schedule '%s': %+v", name, err)
		}
	}

	return nil
}

func expandArmAutomationScheduleAdvanced(d *schema.ResourceData, isUpdate bool) (*automation.AdvancedSchedule, error) {

	expandedAdvancedSchedule := automation.AdvancedSchedule{}

	// If frequency is set to `Month` the `week_days` array cannot be set (even empty), otherwise the API returns an error.
	// During update it can be set and it will not return an error. Workaround for the APIs behavior
	if v, ok := d.GetOk("week_days"); ok {
		weekDays := v.(*schema.Set).List()
		expandedWeekDays := make([]string, len(weekDays))
		for i := range weekDays {
			expandedWeekDays[i] = weekDays[i].(string)
		}
		expandedAdvancedSchedule.WeekDays = &expandedWeekDays
	} else if isUpdate {
		expandedAdvancedSchedule.WeekDays = &[]string{}
	}

	// Same as above with `week_days`
	if v, ok := d.GetOk("month_days"); ok {
		monthDays := v.(*schema.Set).List()
		expandedMonthDays := make([]int32, len(monthDays))
		for i := range monthDays {
			expandedMonthDays[i] = int32(monthDays[i].(int))
		}
		expandedAdvancedSchedule.MonthDays = &expandedMonthDays
	} else if isUpdate {
		expandedAdvancedSchedule.MonthDays = &[]int32{}
	}

	monthlyOccurrences := d.Get("monthly_occurrence").([]interface{})
	expandedMonthlyOccurrences := make([]automation.AdvancedScheduleMonthlyOccurrence, len(monthlyOccurrences))
	for i := range monthlyOccurrences {
		m := monthlyOccurrences[i].(map[string]interface{})
		occurrence := int32(m["occurrence"].(int))

		expandedMonthlyOccurrences[i] = automation.AdvancedScheduleMonthlyOccurrence{
			Occurrence: &occurrence,
			Day:        automation.ScheduleDay(m["day"].(string)),
		}
	}
	expandedAdvancedSchedule.MonthlyOccurrences = &expandedMonthlyOccurrences

	return &expandedAdvancedSchedule, nil
}

func flattenArmAutomationScheduleAdvancedWeekDays(s *automation.AdvancedSchedule) *schema.Set {
	flattenedWeekDays := schema.NewSet(set.HashStringIgnoreCase, []interface{}{})
	if weekDays := s.WeekDays; weekDays != nil {
		for _, v := range *weekDays {
			flattenedWeekDays.Add(v)
		}
	}
	return flattenedWeekDays
}

func flattenArmAutomationScheduleAdvancedMonthDays(s *automation.AdvancedSchedule) *schema.Set {
	flattenedMonthDays := schema.NewSet(set.HashInt, []interface{}{})
	if monthDays := s.MonthDays; monthDays != nil {
		for _, v := range *monthDays {
			flattenedMonthDays.Add(int(v))
		}
	}
	return flattenedMonthDays
}

func flattenArmAutomationScheduleAdvancedMonthlyOccurrences(s *automation.AdvancedSchedule) []map[string]interface{} {
	flattenedMonthlyOccurrences := make([]map[string]interface{}, 0)
	if monthlyOccurrences := s.MonthlyOccurrences; monthlyOccurrences != nil {
		for _, v := range *monthlyOccurrences {
			f := make(map[string]interface{})
			f["day"] = v.Day
			f["occurrence"] = int(*v.Occurrence)
			flattenedMonthlyOccurrences = append(flattenedMonthlyOccurrences, f)
		}
	}
	return flattenedMonthlyOccurrences
}
