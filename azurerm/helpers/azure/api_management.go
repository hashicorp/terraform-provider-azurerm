package azure

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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

func SetCustomPropertyFrom(input map[string]*string, path string, output map[string]interface{}, key string) error {
	log.Printf("input to custom prop = %v", input)
	if v := input[path]; v != nil {
		val, err := strconv.ParseBool(*v)
		if err != nil {
			return fmt.Errorf("Error parsing `%s` %q: %+v", key, *v, err)
		}

		if val {
			output[key] = val
		}
	}

	return nil
}
