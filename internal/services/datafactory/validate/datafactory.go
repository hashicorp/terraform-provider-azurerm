// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func DataFactoryPipelineAndTriggerName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		value := i.(string)
		if !regexp.MustCompile(`^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("invalid name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules %q: %q", k, value))
		}

		return warnings, errors
	}
}

func DataFactoryName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		value := i.(string)
		if !regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("invalid data_factory_name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules %q: %q", k, value))
		}

		return warnings, errors
	}
}

func DataFactoryManagedPrivateEndpointName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", k))
			return
		}

		if !regexp.MustCompile(`^([[:alnum:]][-._[:alnum:]]{0,78}[_[:alnum:]])$`).MatchString(v) {
			errors = append(errors, fmt.Errorf("invalid Data Factory Managed Private Endpoint name, must match the regular expression ^^([[:alnum:]][-._[:alnum:]]{0,78}[_[:alnum:]])$"))
		}

		return warnings, errors
	}
}

func CMKIdentityIdRequiredAtCreation(ctx context.Context, d *pluginsdk.ResourceDiff, meta interface{}) error {
	if d.Id() == "" &&
		d.Get("customer_managed_key_id").(string) != "" &&
		d.Get("customer_managed_key_identity_id").(string) == "" {
		return fmt.Errorf("`customer_managed_key_identity_id` is required when creating a new Data Factory with `customer_managed_key_id`")
	}
	return nil
}
