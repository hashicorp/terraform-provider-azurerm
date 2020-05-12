package validate

import (
	"fmt"
	"regexp"
	"time"
)

const (
	scheduleDateLayout = "01/02/2006"
	scheduleTimeLayout = "15:04:05"
)

func ActionRuleName(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if !regexp.MustCompile(`^([a-zA-Z\d])[a-zA-Z\d-_]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s should begin with a letter or number, contain only letters, numbers, underscores and hyphens.", k))
	}

	return
}

func ActionRuleScheduleDate(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	_, err := time.Parse(scheduleDateLayout, v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s format does not align to: MM/DD/YYYY.", k))
	}

	return
}

func ActionRuleScheduleTime(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	_, err := time.Parse(scheduleTimeLayout, v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s format does not align to: HH:MM:SS.", k))
	}

	return
}
