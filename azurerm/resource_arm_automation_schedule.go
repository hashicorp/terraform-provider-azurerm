package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/automation"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expiry_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"frequency": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.Day),
					string(automation.Hour),
					string(automation.Month),
					string(automation.OneTime),
					string(automation.Week),
				}, true),
			},
			"first_run": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"second": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"minute": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"hour": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"day_of_week": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"day_of_month": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
				Set: resourceAzureRMAutomationScheduleFreqConstraint,
			},
		},
	}
}

func resourceAzureRMAutomationScheduleFreqConstraint(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%d-%d-%d-%d-%d-%d", m["second"], m["minute"], m["hour"], m["day_of_week"], m["day_of_month"]))
	return hashcode.String(buf.String())
}

func computeValidStartTime(firstRunSet *schema.Set, freq automation.ScheduleFrequency) time.Time {

	firstRunSetList := firstRunSet.List()

	closestValidStartTime := time.Now().UTC().Add(time.Duration(7) * time.Minute)
	if len(firstRunSetList) == 0 {
		return closestValidStartTime
	}

	firstRun := firstRunSetList[0].(map[string]interface{})

	firstRunSec := int(firstRun["second"].(int))
	if firstRunSec == -1 {
		firstRunSec = closestValidStartTime.Second()
	}

	firstRunMinute := int(firstRun["minute"].(int))
	if firstRunMinute == -1 {
		firstRunMinute = closestValidStartTime.Minute()
	}

	firstRunHour := int(firstRun["hour"].(int))
	if firstRunHour == -1 {
		firstRunHour = closestValidStartTime.Hour()
	}

	firstRunDayOfWeek := int(firstRun["day_of_week"].(int))
	if firstRunDayOfWeek == -1 {
		firstRunDayOfWeek = int(closestValidStartTime.Weekday())
	}

	firstRunDayOfMonth := int(firstRun["day_of_month"].(int))
	if firstRunDayOfMonth == -1 {
		firstRunDayOfMonth = closestValidStartTime.Day()
	}

	switch freq {
	case automation.Hour:
		validStartTime := time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), closestValidStartTime.Hour(), firstRunMinute, firstRunSec, 0, time.UTC)
		if firstRunMinute <= closestValidStartTime.Minute() {
			validStartTime.Add(time.Duration(1) * time.Hour)
		}

		return validStartTime

	case automation.Day:
		validStartTime := time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), firstRunHour, firstRunMinute, firstRunSec, 0, time.UTC)
		if firstRunHour <= closestValidStartTime.Hour() {
			validStartTime.AddDate(0, 0, 1)
		}
		return validStartTime

	case automation.Week:
		validStartTime := time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), firstRunHour, firstRunMinute, firstRunSec, 0, time.UTC)
		if firstRunDayOfWeek <= int(closestValidStartTime.Weekday()) {
			dayadd := 7 - (int(closestValidStartTime.Weekday()) - firstRunDayOfWeek)
			validStartTime.AddDate(0, 0, dayadd)
		}
		return validStartTime

	case automation.Month:
		validStartTime := time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), firstRunDayOfMonth, firstRunHour, firstRunMinute, firstRunSec, 0, time.UTC)
		if firstRunDayOfMonth <= closestValidStartTime.Day() {
			validStartTime.AddDate(0, 1, 0)
		}
		return validStartTime

	case automation.OneTime:
		validStartTime := time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), firstRunDayOfMonth, firstRunHour, firstRunMinute, firstRunSec, 0, time.UTC)
		return validStartTime

	}

	fmt.Errorf("Error compute first valid run time")
	return closestValidStartTime
}

func resourceArmAutomationScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient
	log.Printf("[INFO] preparing arguments for AzureRM Automation Schedule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	accName := d.Get("account_name").(string)
	freqstr := d.Get("frequency").(string)
	freq := automation.ScheduleFrequency(freqstr)

	cst := d.Get("start_time").(string)
	var starttime time.Time

	if cst != "" {
		var tperr error
		starttime, tperr = time.Parse(time.RFC3339, cst)
		if tperr != nil {
			return fmt.Errorf("Cannot parse start_time: %s", cst)
		}
	} else {
		starttime = computeValidStartTime(d.Get("first_run").(*schema.Set), freq)
	}

	ardt := date.Time{Time: starttime}

	//TODO Interval handling
	//interval :=
	description := d.Get("description").(string)
 	timezone := "UTC"
	
	parameters := automation.ScheduleCreateOrUpdateParameters{
		Name: &name,
		ScheduleCreateOrUpdateProperties: &automation.ScheduleCreateOrUpdateProperties{
			Description: &description,
			//Interval:    &interval,
			Frequency: freq,
			StartTime: &ardt,
			TimeZone: &timezone,
		},
	}

	_, err := client.CreateOrUpdate(resGroup, accName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, accName, name)
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
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["schedules"]

	resp, err := client.Get(resGroup, accName, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Schedule '%s': %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("account_name", accName)
	d.Set("frequency", resp.Frequency)
	d.Set("interval", resp.Interval)
	d.Set("description", resp.Description)
	d.Set("start_time", string(resp.StartTime.Format(time.RFC3339)))
	d.Set("first_run", nil)
	return nil
}

func resourceArmAutomationScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["schedules"]

	resp, err := client.Delete(resGroup, accName, name)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Schedule '%s': %+v", name, err)
	}

	return nil
}
