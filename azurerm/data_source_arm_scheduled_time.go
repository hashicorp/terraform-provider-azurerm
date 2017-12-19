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
				Computed: true,
			},
			"minute": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"day_of_week": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"day_of_month"},
			},
			"day_of_month": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"day_of_week"},
			},
			"minimum_delay_from_now_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timezone": {
				Type:         schema.TypeString,
				Default:      time.UTC.String(),
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateTimezone,
			},
			"next_run_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateTimezone(v interface{}, k string) (ws []string, errors []error) {
	locationString := v.(string)

	if _, err := time.LoadLocation(locationString); err != nil {
		errors = append(errors, fmt.Errorf("Cannot parse timezone: %q", locationString))
	}

	return
}

func dataSourceArmScheduledTimeRead(d *schema.ResourceData, meta interface{}) error {

	var shiftTime int
	var location *time.Location

	if v, exists := d.GetOkExists("minimum_delay_from_now_in_minutes"); exists {
		shiftTime = v.(int)
	}

	closestValidStartTime := time.Now().In(location).Add(time.Duration(shiftTime) * time.Minute)

	var firstRunSec, firstRunMinute, firstRunHour, firstRunDayOfWeek, firstRunDayOfMonth int

	if v, exists := d.GetOkExists("second"); exists {
		firstRunSec = v.(int)
	} else {
		firstRunSec = closestValidStartTime.Second()
	}

	if v, exists := d.GetOkExists("minute"); exists {
		firstRunMinute = v.(int)
	} else {
		firstRunMinute = closestValidStartTime.Minute()
	}

	if v, exists := d.GetOkExists("hour"); exists {
		firstRunHour = v.(int)
	} else {
		firstRunHour = closestValidStartTime.Hour()
	}

	if v, exists := d.GetOkExists("day_of_week"); exists {
		firstRunDayOfWeek = v.(int)
	} else {
		firstRunDayOfWeek = int(closestValidStartTime.Weekday())
	}

	if v, exists := d.GetOkExists("day_of_month"); exists {
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
		} else {
			dayadd := (firstRunDayOfWeek - int(closestValidStartTime.Weekday()))
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
