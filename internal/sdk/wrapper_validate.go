package sdk

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// ValidateModelObject validates that the object contains the specified `tfschema` tags
// required to be used with the Encode and Decode functions
func ValidateModelObject(input interface{}, fields map[string]*schema.Schema) error {
	if input == nil {
		// model not used for this resource
		return nil
	}

	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return fmt.Errorf("need a pointer to the model object")
	}

	// TODO: could we also validate that each `tfschema` tag exists in the schema?

	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()

	if objVal.Kind() == reflect.Interface {
		return fmt.Errorf("cannot resolve pointer to interface")
	}

	return validateModelObjectRecursively("", objType, objVal, fields)
}

func validateModelObjectRecursively(prefix string, objType reflect.Type, objVal reflect.Value, schemafields map[string]*schema.Schema) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			out, ok := r.(error)
			if !ok {
				return
			}

			errOut = out
		}
	}()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)

		if field.Type.Kind() == reflect.Slice {
			sv := fieldVal.Slice(0, fieldVal.Len())
			innerType := sv.Type().Elem()
			innerVal := reflect.Indirect(reflect.New(innerType))
			fieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, field.Name), ".")

			if tag, exists := field.Tag.Lookup("tfschema"); !exists {
				return fmt.Errorf("field %q is missing an `tfschema` label", fieldName)
			} else {
				schemaField, ok := schemafields[tag]
				if !ok {
					return fmt.Errorf("field %q in model is missing from schema", tag)
				}
				switch schemaField.Elem.(type) {
				case *pluginsdk.Resource:
					innerFieldSchema, _ := schemaField.Elem.(*pluginsdk.Resource)
					if err := validateModelObjectRecursively(fieldName, innerType, innerVal, innerFieldSchema.Schema); err != nil {
						return err
					}
				case *schema.Schema:
					switch schemaField.Elem.(*schema.Schema).Type {
					case schema.TypeString:
						if field.Type.String() != "[]string" {
							return fmt.Errorf("expected field %q in model to be of type `[]string` but instead got %q", fieldName, field.Type.String())
						}
					case schema.TypeInt:
						if field.Type.String() != "[]int" {
							return fmt.Errorf("expected field %q in model to be of type `[]int` but instead got %q", fieldName, field.Type.String())
						}
					case schema.TypeList:
						if !strings.HasPrefix(field.Type.String(), "[]") {
							return fmt.Errorf("expected field %q in model to be slice but instead got %q", fieldName, field.Type.String())
						}
						switch schemaField.Elem.(*schema.Schema).Elem.(type) {
						case *schema.Resource:
							if err := validateModelObjectRecursively(fieldName, innerType, innerVal, schemaField.Elem.(*schema.Schema).Elem.(*schema.Resource).Schema); err != nil {
								return err
							}
						default:
							return fmt.Errorf("wtf is this")
						}
					default:
						return fmt.Errorf("unexpected List Type %q for field %q", schemaField.Elem.(*schema.Schema).Type, fieldName)
					}

				default:
					return fmt.Errorf("unexpected type %q for field %q", reflect.TypeOf(schemaField.Elem), tag)
				}
			}
		}

		if tag, exists := field.Tag.Lookup("tfschema"); !exists {
			fieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, field.Name), ".")
			return fmt.Errorf("field %q is missing an `tfschema` label", fieldName)
		} else {
			_, ok := schemafields[tag]
			if !ok {
				return fmt.Errorf(" field %q is missing from schema", tag)
			}
		}
	}

	return nil
}
