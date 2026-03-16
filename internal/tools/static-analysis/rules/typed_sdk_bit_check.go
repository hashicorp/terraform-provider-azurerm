// Copyright IBM Corp. 2014, 2025
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
			if errs := checkModelObject(resource.ResourceType(), resource.ModelObject()); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
		for _, datasource := range s.DataSources() {
			if errs := checkModelObject(datasource.ResourceType(), datasource.ModelObject()); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
	}

	return
}

func checkModelObject(resourceType string, modelObj interface{}) (errors []error) {
	modelType := reflect.TypeOf(modelObj)
	if modelType == nil {
		// TODO: uncomment the following error once all typed resources return a proper model from ModelObject().
		// Currently 16 resources use metadata.ResourceData directly and return nil here. These need to be
		// migrated to use metadata.Decode/Encode with a typed model struct before this check can be enforced.
		// return []error{fmt.Errorf("%q cannot be bit checked, ModelObject returned nil", resourceType)}
		return nil
	}

	if modelType.Kind() != reflect.Ptr {
		return []error{fmt.Errorf("%q cannot be bit checked, ModelObject did not return a pointer", resourceType)}
	}

	model := modelType.Elem()
	if model.Kind() != reflect.Struct {
		return []error{fmt.Errorf("%q ModelObject is not a pointer to a struct", resourceType)}
	}

	return checkForBits(model)
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
			errors = append(errors, fmt.Errorf("property %s in model %s should be type int64, got `%s`", model.Field(i).Name, model.Name(), t.String()))
		case reflect.Float32:
			errors = append(errors, fmt.Errorf("property %s in model %s should be type float64, got `%s`", model.Field(i).Name, model.Name(), t.String()))
		case reflect.Slice, reflect.Array:
			switch t.Elem().Kind() {
			case reflect.Struct:
				errors = append(errors, checkForBits(t.Elem())...)

			case reflect.Int, reflect.Int16, reflect.Int32:
				errors = append(errors, fmt.Errorf("property %s in model %s should be type []int64, got `%s`", model.Field(i).Name, model.Name(), t.String()))

			case reflect.Float32:
				errors = append(errors, fmt.Errorf("property %s in model %s should be type []float64, got `%s`", model.Field(i).Name, model.Name(), t.String()))
			default:
			}

		default:
		}
	}
	return
}
