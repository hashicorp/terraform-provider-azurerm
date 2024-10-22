// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ApiManagementChildName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// from the portal: `The field may contain only numbers, letters, underscore (-), and dash (-) sign when preceded and followed by number or a letter.`
	if matched := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-_]{0,78}[a-zA-Z0-9])?$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, underscores and dashes up to 80 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementServiceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 50 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementUserName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q may not start or end with '-' character", k))
	}
	// from the portal: `The field may contain only numbers, letters, and dash (-) sign when preceded and followed by number or a letter.`
	if matched := regexp.MustCompile(`(^([a-zA-Z0-9-]{1,80})$)`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 80 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementServicePublisherName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^.{1,100}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 100 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementServicePublisherEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[\S*]{1,100}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 100 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementApiName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[^*#&+]{1,256}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only be up to 256 characters in length and not include the characters `*`, `#`, `&` or `+`", k))
	}
	return ws, es
}

func ApiManagementApiPath(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^(?:|[\w]|[\w.][\w-/.:]{0,398}[\w-])$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only be up to 400 characters in length, not start or end with `/` and only contain valid url characters", k))
	}
	return ws, es
}

func ApiManagementBackendName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// From https://learn.microsoft.com/en-us/rest/api/apimanagement/backend/create-or-update#uri-parameters
	if matched := regexp.MustCompile(`(^[\w]+$)|(^[\w][\w\-]+[\w]$)`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 50 characters in length", k))
	}

	return warnings, errors
}

func ApiManagementNamedValueDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// From the portal: `Name may contain only letters, digits, periods, dash, and underscore.`
	// `The value must have a length of at most 256.`
	if matched := regexp.MustCompile(`^[0-9a-zA-Z_.-]{1,256}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, periods, underscores and dashes", k))
	}

	return warnings, errors
}
