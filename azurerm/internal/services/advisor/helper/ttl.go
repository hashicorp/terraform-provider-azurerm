package helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AdvisorTtl struct {
	days  *int32
	times *AdvisorTtlTime
}
type AdvisorTtlTime struct {
	hours   *int32
	minutes *int32
	seconds *int32
}

// TTL valid form includes: dd.hh:mm:ss, hh:mm:ss, dd, dd.hh, dd.hh:mm
func ParseAdvisorSuppresionTTL(input string) (*AdvisorTtl, error) {
	// dd.hh:mm:ss form
	if strings.Contains(input, ".") {
		daySplit := strings.Split(input, ".")
		days, err := parseAdvisorSuppresionTTLDayshelper(daySplit[0])
		if err != nil {
			return nil, err
		}
		times, err := parseAdvisorSuppresionTTLTimeshelper(daySplit[1])
		if err != nil {
			return nil, err
		}
		return &AdvisorTtl{
			days:  days,
			times: times,
		}, nil

		// hh:mm:ss form
	} else if strings.Contains(input, ":") {
		times, err := parseAdvisorSuppresionTTLTimeshelper(input)
		if err != nil {
			return nil, err
		}
		return &AdvisorTtl{
			days:  utils.Int32(0),
			times: times,
		}, nil

		// dd form
	} else {
		days, err := parseAdvisorSuppresionTTLDayshelper(input)
		if err != nil {
			return nil, err
		}
		return &AdvisorTtl{
			days: days,
			times: &AdvisorTtlTime{
				hours:   utils.Int32(0),
				minutes: utils.Int32(0),
				seconds: utils.Int32(0),
			},
		}, nil
	}
}

func parseAdvisorSuppresionTTLDayshelper(input string) (*int32, error) {
	v, err := strconv.Atoi(input)
	days := int32(v)
	if err != nil {
		return nil, fmt.Errorf("days must be numeric: %v", err)
	}
	if days < 0 || days > 24855 {
		return nil, fmt.Errorf("days must between 0 and 24855")
	}
	return &days, nil
}

func parseAdvisorSuppresionTTLTimeshelper(input string) (*AdvisorTtlTime, error) {
	times := strings.Split(input, ":")
	var hrs, mins, secs int32
	v, err := strconv.Atoi(times[0])
	if err != nil {
		return nil, fmt.Errorf("hours must be numeric: %v", err)
	}
	hrs = int32(v)
	if hrs < 0 || hrs >= 24 {
		return nil, fmt.Errorf("hours must between 0 and 24")
	}

	if len(times) >= 2 {
		v, err = strconv.Atoi(times[1])
		if err != nil {
			return nil, fmt.Errorf("minutes must be numeric: %v", err)
		}
		mins = int32(v)
		if mins < 0 || mins >= 60 {
			return nil, fmt.Errorf("minutes must between 0 and 60")
		}
	}
	if len(times) >= 3 {
		v, err = strconv.Atoi(times[2])
		if err != nil {
			return nil, fmt.Errorf("seconds must be numeric: %v", err)
		}
		secs = int32(v)
		if secs < 0 || secs >= 60 {
			return nil, fmt.Errorf("seconds must between 0 and 60")
		}
	}

	return &AdvisorTtlTime{
		hours:   utils.Int32(hrs),
		minutes: utils.Int32(mins),
		seconds: utils.Int32(secs),
	}, nil
}

func (t1 *AdvisorTtl) Equal(t2 *AdvisorTtl) bool {
	return *t1.days == *t2.days && t1.times.equal(t2.times)
}

func (t1 *AdvisorTtlTime) equal(t2 *AdvisorTtlTime) bool {
	return *t1.hours == *t2.hours && *t1.minutes == *t2.minutes && *t1.seconds == *t2.seconds
}

func (t1 *AdvisorTtl) toString() string {
	return fmt.Sprintf("days:%d, hours:%d,minutes:%d,seconds:%d", *t1.days, *t1.times.hours, *t1.times.minutes, *t1.times.seconds)
}

func (t1 *AdvisorTtl) IsZero() bool {
	return *t1.days == 0 && t1.times.isZero()
}

func (t1 *AdvisorTtlTime) isZero() bool {
	return *t1.hours == 0 && *t1.minutes == 0 && *t1.seconds == 0
}
