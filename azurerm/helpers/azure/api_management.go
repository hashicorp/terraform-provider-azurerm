package azure

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func ValidateApiManagementName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,50}$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 50 characters in length", k))
	}

	return
}

func ValidateApiManagementApiName(v interface{}, k string) (ws []string, es []error) {
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

func ResourceArmApiManagementCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	hostnameManagement, hasHostnameManagement := diff.GetOk("hostname_configurations.0.management")
	hostnamePortal, hasHostnamePortal := diff.GetOk("hostname_configurations.0.portal")
	hostnameProxy, hasHostnameProxy := diff.GetOk("hostname_configurations.0.proxy")
	hostnameScm, hasHostnameScm := diff.GetOk("hostname_configurations.0.scm")

	if hasHostnameManagement {
		if err := validateHostnameConfig(hostnameManagement.([]interface{})); err != nil {
			return err
		}
	}

	if hasHostnamePortal {
		if err := validateHostnameConfig(hostnamePortal.([]interface{})); err != nil {
			return err
		}
	}

	if hasHostnameProxy {
		if err := validateHostnameConfig(hostnameProxy.([]interface{})); err != nil {
			return err
		}
	}

	if hasHostnameScm {
		if err := validateHostnameConfig(hostnameScm.([]interface{})); err != nil {
			return err
		}
	}

	return nil
}

func validateHostnameConfig(hostnames []interface{}) error {
	if len(hostnames) > 0 {
		for _, v := range hostnames {
			hostnameConfig := v.(map[string]interface{})

			keyVaultId := hostnameConfig["key_vault_id"].(string)
			cert := hostnameConfig["certificate"].(string)
			certPass := hostnameConfig["certificate_password"].(string)

			if keyVaultId != "" && (cert != "" || certPass != "") {
				return fmt.Errorf("`keyVaultId` cannot be used together with `certificate ` or `certificate_password`.")
			}
		}
	}
	return nil
}
