// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"
)

func AgentLifetime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	days, hours, minutes, seconds, err := parseTimeSpan(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid timespan (expected format: dd.hh:mm:ss): %s", k, err))
		return
	}

	if hours >= 24 {
		errors = append(errors, fmt.Errorf("%q has invalid hours value %d: must be between 0 and 23", k, hours))
	}
	if minutes >= 60 {
		errors = append(errors, fmt.Errorf("%q has invalid minutes value %d: must be between 0 and 59", k, minutes))
	}
	if seconds >= 60 {
		errors = append(errors, fmt.Errorf("%q has invalid seconds value %d: must be between 0 and 59", k, seconds))
	}

	totalSeconds := days*86400 + hours*3600 + minutes*60 + seconds
	if totalSeconds > 7*86400 {
		errors = append(errors, fmt.Errorf("%q must not exceed 7 days (7.00:00:00), got %q", k, v))
	}

	return
}

var (
	timeSpanWithDays    = regexp.MustCompile(`^(\d+)\.(\d{1,2}):(\d{2}):(\d{2})$`)
	timeSpanWithoutDays = regexp.MustCompile(`^(\d{1,2}):(\d{2}):(\d{2})$`)
)

func parseTimeSpan(s string) (days, hours, minutes, seconds int, err error) {
	if m := timeSpanWithDays.FindStringSubmatch(s); m != nil {
		days, err = strconv.Atoi(m[1])
		if err != nil {
			return
		}
		hours, err = strconv.Atoi(m[2])
		if err != nil {
			return
		}
		minutes, err = strconv.Atoi(m[3])
		if err != nil {
			return
		}
		seconds, err = strconv.Atoi(m[4])
		if err != nil {
			return
		}
	} else if m := timeSpanWithoutDays.FindStringSubmatch(s); m != nil {
		hours, err = strconv.Atoi(m[1])
		if err != nil {
			return
		}
		minutes, err = strconv.Atoi(m[2])
		if err != nil {
			return
		}
		seconds, err = strconv.Atoi(m[3])
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("value %q does not match expected format dd.hh:mm:ss", s)
	}

	return
}
