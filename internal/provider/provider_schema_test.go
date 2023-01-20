package provider

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestDataSourcesHaveSensitiveFieldsMarkedAsSensitive(t *testing.T) {
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	dataSourceNames := make([]string, 0)
	for dataSourceName := range provider.DataSourcesMap {
		dataSourceNames = append(dataSourceNames, dataSourceName)
	}
	sort.Strings(dataSourceNames)

	for _, dataSourceName := range dataSourceNames {
		dataSource := provider.DataSourcesMap[dataSourceName]
		if err := schemaContainsSensitiveFieldsNotMarkedAsSensitive(dataSource.Schema); err != nil {
			t.Fatalf("the Data Source %q contains a sensitive field which isn't marked as sensitive: %+v", dataSourceName, err)
		}
	}
}

func TestResourcesHaveSensitiveFieldsMarkedAsSensitive(t *testing.T) {
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	resourceNames := make([]string, 0)
	for resourceName := range provider.ResourcesMap {
		resourceNames = append(resourceNames, resourceName)
	}
	sort.Strings(resourceNames)

	for _, resourceName := range resourceNames {
		resource := provider.ResourcesMap[resourceName]
		if err := schemaContainsSensitiveFieldsNotMarkedAsSensitive(resource.Schema); err != nil {
			t.Fatalf("the Resource %q contains a sensitive field which isn't marked as sensitive: %+v", resourceName, err)
		}
	}
}

func schemaContainsSensitiveFieldsNotMarkedAsSensitive(input map[string]*pluginsdk.Schema) error {
	exactMatchFieldNames := []string{
		"api_key",
		"api_secret_key",
		"password",
		"private_key",
		"ssh_private_key",
	}

	// intentionally sorting these so the output is consistent
	fieldNames := make([]string, 0)
	for fieldName := range input {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		key := strings.ToLower(fieldName)
		field := input[fieldName]

		for _, val := range exactMatchFieldNames {
			if strings.EqualFold(key, val) && !field.Sensitive {
				return fmt.Errorf("field %q is a sensitive value and should be marked as Sensitive", fieldName)
			}
		}

		if strings.HasSuffix(key, "_api_key") && field.Type == pluginsdk.TypeString && !field.Sensitive {
			return fmt.Errorf("field %q is a sensitive value and should be marked as Sensitive", fieldName)
		}

		if field.Type == pluginsdk.TypeList && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsSensitiveFieldsNotMarkedAsSensitive(val.Schema); err != nil {
					return fmt.Errorf("the field %q is a List: %+v", fieldName, err)
				}
			}
		}

		if field.Type == pluginsdk.TypeSet && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsSensitiveFieldsNotMarkedAsSensitive(val.Schema); err != nil {
					return fmt.Errorf("the field %q is a Set: %+v", fieldName, err)
				}
			}
		}
	}

	return nil
}
