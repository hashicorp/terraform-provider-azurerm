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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
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

			"resource_group_name": resourceGroupNameSchema(),

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
				DiffSuppressFunc: suppress.Rfc3339Time,
				ValidateFunc:     validate.Rfc3339DateInFutureBy(time.Duration(5) * time.Minute),
				//defaults to now + 7 minutes in create function if not set
			},

			"expiry_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true, //same as start time when OneTime, ridiculous value when recurring: "9999-12-31T15:59:00-08:00"
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.Rfc3339Time,
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

			//todo missing properties: week_days, month_days, month_week_day from advanced automation section
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			interval, _ := diff.GetOk("interval")
			if strings.ToLower(diff.Get("frequency").(string)) == "onetime" && interval.(int) > 0 {
				return fmt.Errorf("interval canot be set when frequency is not OneTime")
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
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Schedule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	frequency := d.Get("frequency").(string)

	timeZone := d.Get("timezone").(string)
	description := d.Get("description").(string)

	//CustomizeDiff should ensure one of these two is set
	//todo remove this once `account_name` is removed
	accountName := ""
	if v, ok := d.GetOk("automation_account_name"); ok {
		accountName = v.(string)
	} else if v, ok := d.GetOk("account_name"); ok {
		accountName = v.(string)
	}

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
	_, err := client.CreateOrUpdate(ctx, resGroup, accountName, name, parameters)
	if err != nil {
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
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("start_time", string(v.Format(time.RFC3339)))
	}
	if v := resp.ExpiryTime; v != nil {
		d.Set("expiry_time", string(v.Format(time.RFC3339)))
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
	return nil
}

func resourceArmAutomationScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
