package provider

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func TestDataSourcesHaveEnabledFieldsMarkedAsBooleans(t *testing.T) {
	// This test validates that Data Sources do not contain a field suffixed with `_enabled` that isn't a Boolean.
	//
	// If this test is failing due to a new Data Source/new field within an existing Data Source, it'd be worth validating
	// the schema, since fields matching `{some_name}_enabled` should be Booleans. Should a Tri-State Boolean exist,
	// this field likely wants the `_enabled` suffix removing, to make the example `{some_name}` instead (with
	// validation for the possible values).
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	dataSourceNames := make([]string, 0)
	for dataSourceName := range provider.DataSourcesMap {
		dataSourceNames = append(dataSourceNames, dataSourceName)
	}
	sort.Strings(dataSourceNames)

	for _, dataSourceName := range dataSourceNames {
		dataSource := provider.DataSourcesMap[dataSourceName]
		if err := schemaContainsEnabledFieldsNotDefinedAsABoolean(dataSource.Schema, map[string]struct{}{}); err != nil {
			t.Fatalf("the Data Source %q contains an `_enabled` field which isn't defined as a boolean: %+v", dataSourceName, err)
		}
	}
}

func TestResourcesHaveEnabledFieldsMarkedAsBooleans(t *testing.T) {
	// This test validates that Resources do not contain a field suffixed with `_enabled` that isn't a Boolean.
	//
	// If this test is failing due to a new Resource/new field within an existing Resource, it'd be worth validating
	// the schema, since fields matching `{some_name}_enabled` should be Booleans. Should a Tri-State Boolean exist,
	// this field likely wants the `_enabled` suffix removing, to make the example `{some_name}` instead (with
	// validation for the possible values).
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	resourceNames := make([]string, 0)
	for resourceName := range provider.ResourcesMap {
		resourceNames = append(resourceNames, resourceName)
	}
	sort.Strings(resourceNames)

	// TODO: 4.0 - work through this list
	resourceFieldsWhichNeedToBeAddressed := map[string]map[string]struct{}{
		// 1: Fields which require renaming etc
		"azurerm_datadog_monitor_sso_configuration": {
			// should be fixed in 4.0, presumably ditching `_enabled` and adding Enum validation
			"single_sign_on_enabled": {},
		},
		"azurerm_netapp_volume": {
			// should be fixed in 4.0, presumably ditching `_enabled` and making this `protocols_to_use` or something?
			"protocols_enabled": {},
		},
		"azurerm_kubernetes_cluster": {
			// this either wants `enabled` removing, or to be marked as a false-positive
			"transparent_huge_page_enabled": {},
		},
		"azurerm_kubernetes_cluster_node_pool": {
			// this either wants `enabled` removing, or to be marked as a false-positive
			"transparent_huge_page_enabled": {},
		},

		// 2: False Positives
		"azurerm_iot_security_solution": {
			// this is a list of recommendations
			"recommendations_enabled": {},
		},
	}

	for _, resourceName := range resourceNames {
		resource := provider.ResourcesMap[resourceName]
		fieldsToBeAddressed := resourceFieldsWhichNeedToBeAddressed[resourceName]

		if err := schemaContainsEnabledFieldsNotDefinedAsABoolean(resource.Schema, fieldsToBeAddressed); err != nil {
			t.Fatalf("the Resource %q contains an `_enabled` field which isn't defined as a boolean: %+v", resourceName, err)
		}
	}
}

func schemaContainsEnabledFieldsNotDefinedAsABoolean(input map[string]*schema.Schema, fieldsToBeAddressed map[string]struct{}) error {
	// intentionally sorting these so the output is consistent
	fieldNames := make([]string, 0)
	for fieldName := range input {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		key := strings.ToLower(fieldName)
		field := input[fieldName]

		if strings.HasSuffix(key, "_enabled") {
			// @tombuildsstuff: we have some Resources which will need to be addressed in the next major version (v4.0)
			// if this field name matches one we're intentionally ignoring, let's ignore it for now
			if _, shouldIgnore := fieldsToBeAddressed[key]; shouldIgnore {
				continue
			}
			if field.Type != pluginsdk.TypeBool {
				return fmt.Errorf("field %q is an `_enabled` field so should be defined as a Boolean but got %+v", fieldName, field.Type)
			}
		}

		if field.Type == pluginsdk.TypeList && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsEnabledFieldsNotDefinedAsABoolean(val.Schema, fieldsToBeAddressed); err != nil {
					return fmt.Errorf("the field %q is a List: %+v", fieldName, err)
				}
			}
		}

		if field.Type == pluginsdk.TypeSet && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsEnabledFieldsNotDefinedAsABoolean(val.Schema, fieldsToBeAddressed); err != nil {
					return fmt.Errorf("the field %q is a Set: %+v", fieldName, err)
				}
			}
		}
	}

	return nil
}

func TestDataSourcesDoNotContainANameFieldWithADefaultOfDefault(t *testing.T) {
	// This test validates that Data Sources do not contain a field `name` with a default value of `default`, which
	// would signify that only a single instance of this resource can be created and is related to the parent resource.
	//
	// If a new Data Sources is failing because of this test, rather than adding a new Data Sources, you likely want to
	// embed the relevant structure (for example `sso_configuration {}`) within the parent Data Sources this is related to.
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	dataSourceNames := make([]string, 0)
	for dataSourceName := range provider.DataSourcesMap {
		dataSourceNames = append(dataSourceNames, dataSourceName)
	}
	sort.Strings(dataSourceNames)

	for _, dataSourceName := range dataSourceNames {
		dataSource := provider.DataSourcesMap[dataSourceName]
		if err := schemaContainsANameFieldWithADefaultValueOfDefault(dataSource.Schema, map[string]struct{}{}); err != nil {
			t.Fatalf("the Data Source %q contains a `name` field with a default value of `default` - this Data Source should be exposed as part of the parent Data Source it's located within: %+v", dataSourceName, err)
		}
	}
}

func TestResourcesDoNotContainANameFieldWithADefaultOfDefault(t *testing.T) {
	// This test validates that Resources do not contain a field `name` with a default value of `default`, which
	// would signify that only a single instance of this resource can be created and is related to the parent resource.
	//
	// If a new Resource is failing because of this test, rather than adding a new Resource, you likely want to
	// embed the relevant structure (for example `sso_configuration {}`) within the parent Resource this is related to.
	provider := TestAzureProvider()

	// intentionally sorting these so the output is consistent
	resourceNames := make([]string, 0)
	for resourceName := range provider.ResourcesMap {
		resourceNames = append(resourceNames, resourceName)
	}
	sort.Strings(resourceNames)

	// TODO: 4.0 - work through this list
	resourceFieldsWhichNeedToBeAddressed := map[string]map[string]struct{}{
		// 1: to be addressed in 4.0
		"azurerm_datadog_monitor_sso_configuration": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_datadog_monitor`
			// which'll also need the Monitor resource to have Create call Update
			"name": {},
		},
		"azurerm_datadog_monitor_tag_rule": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_datadog_monitor`
			// which'll also need the Monitor resource to have Create call Update
			"name": {},
		},
		"azurerm_spring_cloud_accelerator": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},
		"azurerm_spring_cloud_api_portal": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},
		"azurerm_spring_cloud_application_live_view": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},
		"azurerm_spring_cloud_configuration_service": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},
		"azurerm_spring_cloud_dev_tool_portal": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},
		"azurerm_spring_cloud_gateway": {
			// TODO: in 4.0 this resource probably wants embedding within `azurerm_spring_cloud_service`
			"name": {},
		},

		// 2: False Positives?
		"azurerm_redis_enterprise_database": {
			"name": {},
		},

		// 3: Deprecated / to be removed in 4.0
		"azurerm_cosmosdb_notebook_workspace": {
			"name": {},
		},
	}

	for _, resourceName := range resourceNames {
		resource := provider.ResourcesMap[resourceName]
		fieldsToBeAddressed := resourceFieldsWhichNeedToBeAddressed[resourceName]

		if err := schemaContainsANameFieldWithADefaultValueOfDefault(resource.Schema, fieldsToBeAddressed); err != nil {
			t.Fatalf("the Resource %q contains a `name` field with a default value of `default` - this Resource should be exposed as part of the parent Resource it's located within: %+v", resourceName, err)
		}
	}
}

func schemaContainsANameFieldWithADefaultValueOfDefault(input map[string]*schema.Schema, fieldsToBeAddressed map[string]struct{}) error {
	// intentionally sorting these so the output is consistent
	fieldNames := make([]string, 0)
	for fieldName := range input {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		key := strings.ToLower(fieldName)
		field := input[fieldName]

		// @tombuildsstuff: we have some Resources which will need to be addressed in the next major version (v4.0)
		// if this field name matches one we're intentionally ignoring, let's ignore it for now
		if _, shouldIgnore := fieldsToBeAddressed[key]; shouldIgnore {
			continue
		}

		if strings.EqualFold(key, "name") {
			defaultValue, err := field.DefaultValue()
			if err != nil {
				return fmt.Errorf("obtaining default value for %q: %+v", fieldName, err)
			}

			if v, ok := defaultValue.(string); ok {
				if strings.EqualFold(v, "default") {
					return fmt.Errorf("field %q is a `name` field which contains a default value of `default`", fieldName)
				}
			}

			// should the ValidateFunc allow `default`, we can assume this too
			if field.ValidateFunc != nil {
				allowsEmptyString := runInputForValidateFunction(field.ValidateFunc, "")
				allowsWhitespaceString := runInputForValidateFunction(field.ValidateFunc, " ")
				allowsPlaceholderValue := runInputForValidateFunction(field.ValidateFunc, "placeholder")
				allowsDefaultValue := runInputForValidateFunction(field.ValidateFunc, "default")
				if allowsDefaultValue && !allowsPlaceholderValue && !allowsWhitespaceString && !allowsEmptyString {
					return fmt.Errorf("field %q is a `name` field where the ValidateFunc explicitly allows a default value of `default`", fieldName)
				}
			}
		}

		if field.Type == pluginsdk.TypeList && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsANameFieldWithADefaultValueOfDefault(val.Schema, fieldsToBeAddressed); err != nil {
					return fmt.Errorf("the field %q is a List: %+v", fieldName, err)
				}
			}
		}

		if field.Type == pluginsdk.TypeSet && field.Elem != nil {
			if val, ok := field.Elem.(*pluginsdk.Resource); ok && val.Schema != nil {
				if err := schemaContainsANameFieldWithADefaultValueOfDefault(val.Schema, fieldsToBeAddressed); err != nil {
					return fmt.Errorf("the field %q is a Set: %+v", fieldName, err)
				}
			}
		}
	}

	return nil
}

func runInputForValidateFunction(validateFunc pluginsdk.SchemaValidateFunc, input string) bool {
	warnings, errs := validateFunc(input, input)
	return len(warnings) == 0 && len(errs) == 0
}
