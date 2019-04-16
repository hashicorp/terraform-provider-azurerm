package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateAzureRMDataFactoryLinkedServiceDatasetName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if regexp.MustCompile(`^[-.+?/<>*%&:\\]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("any of '-' '.', '+', '?', '/', '<', '>', '*', '%%', '&', ':', '\\', are not allowed in %q: %q", k, value))
	}

	return warnings, errors
}

func expandDataFactoryLinkedServiceIntegrationRuntime(integrationRuntimeName string) *datafactory.IntegrationRuntimeReference {
	typeString := "IntegrationRuntimeReference"

	return &datafactory.IntegrationRuntimeReference{
		ReferenceName: &integrationRuntimeName,
		Type:          &typeString,
	}
}

// Because the password isn't returned from the api in the connection string, we'll check all
// but the password string and return true if they match.
func azureRmDataFactoryLinkedServiceConnectionStringDiff(k, old string, new string, d *schema.ResourceData) bool {
	oldSplit := strings.Split(strings.ToLower(old), ";")
	newSplit := strings.Split(strings.ToLower(new), ";")

	sort.Strings(oldSplit)
	sort.Strings(newSplit)

	// We need to remove the password from the new string since it isn't returned from the api
	for i, v := range newSplit {
		if strings.HasPrefix(v, "password") {
			newSplit = append(newSplit[:i], newSplit[i+1:]...)
		}
	}

	if len(oldSplit) != len(newSplit) {
		return false
	}

	// We'll error out if we find any differences between the old and the new connection strings
	for i := range oldSplit {
		if !strings.EqualFold(oldSplit[i], newSplit[i]) {
			return false
		}
	}

	return true
}

func expandDataFactoryParameters(input map[string]interface{}) map[string]*datafactory.ParameterSpecification {
	output := make(map[string]*datafactory.ParameterSpecification)

	for k, v := range input {
		output[k] = &datafactory.ParameterSpecification{
			Type:         datafactory.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func flattenDataFactoryParameters(input map[string]*datafactory.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func flattenDataFactoryAnnotations(input *[]interface{}) []string {
	annotations := make([]string, 0)
	if input == nil {
		return annotations
	}

	for _, annotation := range *input {
		val, ok := annotation.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping annotation %q since it's not a string", val)
		}
		annotations = append(annotations, val)
	}
	return annotations
}

func expandDataFactoryVariables(input map[string]interface{}) map[string]*datafactory.VariableSpecification {
	output := make(map[string]*datafactory.VariableSpecification)

	for k, v := range input {
		output[k] = &datafactory.VariableSpecification{
			Type:         datafactory.VariableTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func flattenDataFactoryVariables(input map[string]*datafactory.VariableSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping variable %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}
