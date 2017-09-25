package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type Frequency string

const (
	Day     Frequency = "Day"
	Hour    Frequency = "Hour"
	Month   Frequency = "Month"
	OneTime Frequency = "OneTime"
	Week    Frequency = "Week"
)

func dataSourceArmScheduledTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmScheduledTimeRead,
		Schema: map[string]*schema.Schema{
			"frequency": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(Day),
					string(Hour),
					string(Month),
					string(OneTime),
					string(Week),
				}, true),
			},
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
			"minimum_delay_from_now_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_run_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmScheduledTimeRead(d *schema.ResourceData, meta interface{}) error {

	var shiftTime int
	var location *time.Location

	if v, ok := d.GetOk("minimum_delay_from_now_in_minutes"); ok {
		shiftTime = v.(int)
	} else {
		shiftTime = 0
	}

	if v, ok := d.GetOk("timezone"); ok {
		var err error
		if location, err = time.LoadLocation(v.(string)); err != nil {
			return fmt.Errorf("Cannot parse timezone: %s", v.(string))
		}
	} else {
		location = time.UTC
	}

	closestValidStartTime := time.Now().In(location).Add(time.Duration(shiftTime) * time.Minute)

	var firstRunSec, firstRunMinute, firstRunHour, firstRunDayOfWeek, firstRunDayOfMonth int

	//TODO: GetOk is not suitable here because it returns ok=false if the second value is 0. The GetOkExists looks good but it hasn't merged yet
	if v := d.Get("second"); v.(int) > -1 {
		firstRunSec = v.(int)
	} else {
		firstRunSec = closestValidStartTime.Second()
	}

	if v := d.Get("minute"); v.(int) > -1 {
		firstRunMinute = v.(int)
	} else {
		firstRunMinute = closestValidStartTime.Minute()
	}

	if v := d.Get("hour"); v.(int) > -1 {
		firstRunHour = v.(int)
	} else {
		firstRunHour = closestValidStartTime.Hour()
	}

	if v := d.Get("day_of_week"); v.(int) > -1 {
		firstRunDayOfWeek = v.(int)
	} else {
		firstRunDayOfWeek = int(closestValidStartTime.Weekday())
	}

	if v := d.Get("day_of_month"); v.(int) > -1 {
		firstRunDayOfMonth = v.(int)
	} else {
		firstRunDayOfMonth = closestValidStartTime.Day()
	}

	freq := Frequency(d.Get("frequency").(string))

	var validStartTime time.Time
	switch freq {
	case Hour:
		validStartTime = time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), closestValidStartTime.Hour(), firstRunMinute, firstRunSec, 0, location)
		if firstRunMinute <= closestValidStartTime.Minute() {
			validStartTime = validStartTime.Add(time.Duration(1) * time.Hour)
		}

	case Day:
		validStartTime = time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), firstRunHour, firstRunMinute, firstRunSec, 0, location)
		if firstRunHour <= closestValidStartTime.Hour() {
			validStartTime = validStartTime.AddDate(0, 0, 1)
		}

	case Week:
		validStartTime = time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), closestValidStartTime.Day(), firstRunHour, firstRunMinute, firstRunSec, 0, location)
		if firstRunDayOfWeek <= int(closestValidStartTime.Weekday()) {
			dayadd := 7 - (int(closestValidStartTime.Weekday()) - firstRunDayOfWeek)
			validStartTime = validStartTime.AddDate(0, 0, dayadd)
		}

	case Month:
		validStartTime = time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), firstRunDayOfMonth, firstRunHour, firstRunMinute, firstRunSec, 0, location)
		if firstRunDayOfMonth <= closestValidStartTime.Day() {
			validStartTime = validStartTime.AddDate(0, 1, 0)
		}

	case OneTime:
		validStartTime = time.Date(closestValidStartTime.Year(), closestValidStartTime.Month(), firstRunDayOfMonth, firstRunHour, firstRunMinute, firstRunSec, 0, location)

	}

	d.SetId(time.Now().UTC().String())
	d.Set("next_run_time", validStartTime.Format(time.RFC3339))

	return nil
}
