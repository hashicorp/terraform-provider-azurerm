// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func AgentLifetime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	days, hours, minutes, seconds, err := parseTimeSpan(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid timespan (expected format: dd.hh:mm:ss or hh:mm:ss): %s", k, err))
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

func parseTimeSpan(s string) (days, hours, minutes, seconds int, err error) {
	timePart := s

	// Check for "dd.hh:mm:ss" format
	if dotIdx := strings.IndexByte(s, '.'); dotIdx >= 0 {
		days, err = strconv.Atoi(s[:dotIdx])
		if err != nil || days < 0 {
			err = fmt.Errorf("value %q has invalid days component", s)
			return
		}
		timePart = s[dotIdx+1:]
	}

	parts := strings.Split(timePart, ":")
	if len(parts) != 3 {
		err = fmt.Errorf("value %q does not match expected format dd.hh:mm:ss or hh:mm:ss", s)
		return
	}

	hours, err = strconv.Atoi(parts[0])
	if err != nil || hours < 0 {
		err = fmt.Errorf("value %q has invalid hours component", s)
		return
	}

	minutes, err = strconv.Atoi(parts[1])
	if err != nil || minutes < 0 {
		err = fmt.Errorf("value %q has invalid minutes component", s)
		return
	}

	seconds, err = strconv.Atoi(parts[2])
	if err != nil || seconds < 0 {
		err = fmt.Errorf("value %q has invalid seconds component", s)
		return
	}

	return
}
