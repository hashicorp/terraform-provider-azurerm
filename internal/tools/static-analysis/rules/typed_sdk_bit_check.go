// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

var _ Rule = TypedSDKBitCheck{}

type TypedSDKBitCheck struct{}

func (r TypedSDKBitCheck) Run() (errors []error) {
	for _, s := range provider.SupportedTypedServices() {
		for _, resource := range s.Resources() {
			modelType := reflect.TypeOf(resource.ModelObject())
			switch {
			case modelType != nil && modelType.Kind() == reflect.Ptr:
				model := modelType.Elem()
				if model.Kind() != reflect.Struct {
					errors = append(errors, fmt.Errorf("%s is not a pointer to a struct\n", modelType.Name()))
					continue
				}

				errors = append(errors, checkForBits(model)...)

			case modelType == nil:
				continue

			default:
				errors = append(errors, fmt.Errorf("%q cannot be bit checked, ModelObject did not return a pointer\n", resource.ResourceType()))
			}

		}
		for _, datasource := range s.DataSources() {
			modelType := reflect.TypeOf(datasource.ModelObject())
			if modelType != nil && modelType.Kind() == reflect.Ptr { // Have to nil-check here due to base types not having a model. e.g. roleAssignmentBaseResource
				model := modelType.Elem()
				if model.Kind() != reflect.Struct {
					errors = append(errors, fmt.Errorf("%s is not a pointer to a struct", modelType.Name()))
					continue
				}

				errors = append(errors, checkForBits(model)...)
			} else {
				errors = append(errors, fmt.Errorf("%q cannot be bit checked, ModelObject did not return a pointer", datasource.ResourceType()))
			}
		}
	}

	return
}

func (r TypedSDKBitCheck) Name() string {
	return "checkBittiness"
}

func (r TypedSDKBitCheck) Description() string {
	return fmt.Sprintf(`
The '%s' check function is used to check the correct Go types are used in TypedSDK e.g. 'int64' rather than 'int'.
`, r.Name())
}

func checkForBits(model reflect.Type) (errors []error) {
	for i := 0; i < model.NumField(); i++ {
		switch t := model.Field(i).Type; t.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32:
			errors = append(errors, fmt.Errorf("property %s in model %s should be type int64, got `%s`\n", model.Field(i).Name, model.Name(), t.String()))
		case reflect.Float32:
			errors = append(errors, fmt.Errorf("property %s in model %s should be type float64, got `%s`\n", model.Field(i).Name, model.Name(), t.String()))
		case reflect.Slice, reflect.Array:
			switch t.Elem().Kind() {
			case reflect.Struct:
				errors = append(errors, checkForBits(t.Elem())...)

			case reflect.Int, reflect.Int16, reflect.Int32:
				errors = append(errors, fmt.Errorf("property %s in model %s should be type []int64, got `%s`\n", model.Field(i).Name, model.Name(), t.String()))

			case reflect.Float32:
				errors = append(errors, fmt.Errorf("property %s in model %s should be type []float64, got `%s`\n", model.Field(i).Name, model.Name(), t.String()))
			default:
			}

		default:
		}
	}
	return
}
