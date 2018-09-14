package azure

import (
	"fmt"
	"regexp"
)

func ValidateApiManagementName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,50}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 50 characters in length", k))
	}

	return
}

func ValidateApiManagementPublisherName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[\S*]{1,100}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only be up to 100 characters in length", k))
	}

	return
}

func ValidateApiManagementPublisherEmail(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[\S*]{1,100}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only be up to 100 characters in length", k))
	}

	return
}
