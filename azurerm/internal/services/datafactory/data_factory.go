package datafactory

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
func azureRmDataFactoryLinkedServiceConnectionStringDiff(_, old string, new string, _ *schema.ResourceData) bool {
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

// DatasetColumn describes the attributes needed to specify a structure column for a dataset
type DatasetColumn struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

func expandDataFactoryDatasetStructure(input []interface{}) interface{} {
	columns := make([]DatasetColumn, 0)
	for _, column := range input {
		attrs := column.(map[string]interface{})

		datasetColumn := DatasetColumn{
			Name: attrs["name"].(string),
		}
		if attrs["description"] != nil {
			datasetColumn.Description = attrs["description"].(string)
		}
		if attrs["type"] != nil {
			datasetColumn.Type = attrs["type"].(string)
		}
		columns = append(columns, datasetColumn)
	}
	return columns
}

func flattenDataFactoryStructureColumns(input interface{}) []interface{} {
	output := make([]interface{}, 0)

	columns, ok := input.([]interface{})
	if !ok {
		return columns
	}

	for _, v := range columns {
		column, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result := make(map[string]interface{})
		if column["name"] != nil {
			result["name"] = column["name"]
		}
		if column["type"] != nil {
			result["type"] = column["type"]
		}
		if column["description"] != nil {
			result["description"] = column["description"]
		}
		output = append(output, result)
	}
	return output
}

func deserializeDataFactoryPipelineActivities(jsonData string) (*[]datafactory.BasicActivity, error) {
	jsonData = fmt.Sprintf(`{ "activities": %s }`, jsonData)
	pipeline := &datafactory.Pipeline{}
	err := pipeline.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		return nil, err
	}
	return pipeline.Activities, nil
}

func serializeDataFactoryPipelineActivities(activities *[]datafactory.BasicActivity) (string, error) {
	pipeline := &datafactory.Pipeline{Activities: activities}
	result, err := pipeline.MarshalJSON()
	if err != nil {
		return "nil", err
	}

	var m map[string]*json.RawMessage
	err = json.Unmarshal(result, &m)
	if err != nil {
		return "", err
	}

	activitiesJson, err := json.Marshal(m["activities"])
	if err != nil {
		return "", err
	}

	return string(activitiesJson), nil
}

func suppressJsonOrderingDifference(_, old, new string, _ *schema.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}

func expandAzureKeyVaultPasswordReference(input []interface{}) *datafactory.AzureKeyVaultSecretReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})

	keyVaultLinkedServiceName := config["linked_service_name"].(string)
	keyVaultPasswordSecretName := config["password_secret_name"].(string)
	linkedServiceType := "LinkedServiceReference"

	return &datafactory.AzureKeyVaultSecretReference{
		SecretName: keyVaultPasswordSecretName,
		Store: &datafactory.LinkedServiceReference{
			Type:          &linkedServiceType,
			ReferenceName: &keyVaultLinkedServiceName,
		},
	}
}

func flattenAzureKeyVaultPasswordReference(secretReference *datafactory.AzureKeyVaultSecretReference) []interface{} {
	if secretReference == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if store := secretReference.Store; store != nil {
		if store.ReferenceName != nil {
			parameters["linked_service_name"] = *store.ReferenceName
		}
	}

	parameters["password_secret_name"] = secretReference.SecretName

	return []interface{}{parameters}
}
