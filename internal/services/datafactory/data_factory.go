// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// @tombuildsstuff: these have been ported over from the Azure SDK for Go since the service team has removed them
// but the casing differs in the API, so we need to ensure these are normalized on our side.
const (
	TypeBasicDatasetCompressionTypeBZip2      string = "BZip2"
	TypeBasicDatasetCompressionTypeDeflate    string = "Deflate"
	TypeBasicDatasetCompressionTypeGZip       string = "GZip"
	TypeBasicDatasetCompressionTypeTar        string = "Tar"
	TypeBasicDatasetCompressionTypeTarGZip    string = "TarGZip"
	TypeBasicDatasetCompressionTypeZipDeflate string = "ZipDeflate"
)

func expandDataFactoryLinkedServiceIntegrationRuntime(integrationRuntimeName string) *datafactory.IntegrationRuntimeReference {
	typeString := "IntegrationRuntimeReference"

	return &datafactory.IntegrationRuntimeReference{
		ReferenceName: &integrationRuntimeName,
		Type:          &typeString,
	}
}

// Because the password isn't returned from the api in the connection string, we'll check all
// but the password string and return true if they match.
func azureRmDataFactoryLinkedServiceConnectionStringDiff(_, old string, new string, _ *pluginsdk.ResourceData) bool {
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

// DatasetSnowflakeSchemaColumn describes the attributes needed to specify a Snowflake schema column for a dataset
type DatasetSnowflakeSchemaColumn struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Precision int    `json:"precision,omitempty"`
	Scale     int    `json:"scale,omitempty"`
}

func expandDataFactoryDatasetSnowflakeSchema(input []interface{}) interface{} {
	columns := make([]DatasetSnowflakeSchemaColumn, 0)
	for _, column := range input {
		attrs := column.(map[string]interface{})

		datasetSnowflakeSchemaColumn := DatasetSnowflakeSchemaColumn{
			Name: attrs["name"].(string),
		}
		if attrs["type"] != nil {
			datasetSnowflakeSchemaColumn.Type = attrs["type"].(string)
		}

		if attrs["precision"] != nil {
			datasetSnowflakeSchemaColumn.Precision = attrs["precision"].(int)
		}

		if attrs["scale"] != nil {
			datasetSnowflakeSchemaColumn.Scale = attrs["scale"].(int)
		}

		columns = append(columns, datasetSnowflakeSchemaColumn)
	}
	return columns
}

func flattenDataFactorySnowflakeSchemaColumns(input interface{}) []interface{} {
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
		if column["precision"] != nil {
			result["precision"] = column["precision"]
		}
		if column["scale"] != nil {
			result["scale"] = column["scale"]
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

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}

func expandAzureKeyVaultSecretReference(input []interface{}) *datafactory.AzureKeyVaultSecretReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})

	return &datafactory.AzureKeyVaultSecretReference{
		SecretName: config["secret_name"].(string),
		Store: &datafactory.LinkedServiceReference{
			Type:          utils.String("LinkedServiceReference"),
			ReferenceName: utils.String(config["linked_service_name"].(string)),
		},
	}
}

func flattenAzureKeyVaultConnectionString(input map[string]interface{}) []interface{} {
	if input == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if v, ok := input["store"].(map[string]interface{}); ok {
		if v != nil {
			parameters["linked_service_name"] = v["referenceName"].(string)
		}
	}

	parameters["secret_name"] = input["secretName"]

	return []interface{}{parameters}
}

func flattenAzureKeyVaultSecretReference(secretReference *datafactory.AzureKeyVaultSecretReference) []interface{} {
	if secretReference == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if store := secretReference.Store; store != nil {
		if store.ReferenceName != nil {
			parameters["linked_service_name"] = *store.ReferenceName
		}
	}

	parameters["secret_name"] = secretReference.SecretName

	return []interface{}{parameters}
}

func expandDataFactoryDatasetLocation(d *pluginsdk.ResourceData) datafactory.BasicDatasetLocation {
	if _, ok := d.GetOk("http_server_location"); ok {
		return expandDataFactoryDatasetHttpServerLocation(d)
	}

	if _, ok := d.GetOk("azure_blob_storage_location"); ok {
		return expandDataFactoryDatasetAzureBlobStorageLocation(d)
	}

	if _, ok := d.GetOk("azure_blob_fs_location"); ok {
		return expandDataFactoryDatasetAzureBlobFSLocation(d)
	}

	if _, ok := d.GetOk("sftp_server_location"); ok {
		return expandDataFactoryDatasetSFTPServerLocation(d)
	}

	return nil
}

func expandDataFactoryDatasetSFTPServerLocation(d *pluginsdk.ResourceData) datafactory.BasicDatasetLocation {
	sftpServerLocations := d.Get("sftp_server_location").([]interface{})
	if len(sftpServerLocations) == 0 || sftpServerLocations[0] == nil {
		return nil
	}

	props := sftpServerLocations[0].(map[string]interface{})

	sftpServerLocation := datafactory.SftpLocation{
		FolderPath: expandDataFactoryExpressionResultType(props["path"].(string), props["dynamic_path_enabled"].(bool)),
		FileName:   expandDataFactoryExpressionResultType(props["filename"].(string), props["dynamic_filename_enabled"].(bool)),
	}
	return sftpServerLocation
}

func expandDataFactoryDatasetHttpServerLocation(d *pluginsdk.ResourceData) datafactory.BasicDatasetLocation {
	httpServerLocations := d.Get("http_server_location").([]interface{})
	if len(httpServerLocations) == 0 || httpServerLocations[0] == nil {
		return nil
	}

	props := httpServerLocations[0].(map[string]interface{})

	httpServerLocation := datafactory.HTTPServerLocation{
		RelativeURL: props["relative_url"].(string),
		FolderPath:  expandDataFactoryExpressionResultType(props["path"].(string), props["dynamic_path_enabled"].(bool)),
		FileName:    expandDataFactoryExpressionResultType(props["filename"].(string), props["dynamic_filename_enabled"].(bool)),
	}
	return httpServerLocation
}

func expandDataFactoryDatasetAzureBlobStorageLocation(d *pluginsdk.ResourceData) datafactory.BasicDatasetLocation {
	azureBlobStorageLocations := d.Get("azure_blob_storage_location").([]interface{})
	if len(azureBlobStorageLocations) == 0 || azureBlobStorageLocations[0] == nil {
		return nil
	}

	props := azureBlobStorageLocations[0].(map[string]interface{})

	blobStorageLocation := datafactory.AzureBlobStorageLocation{
		Container:  expandDataFactoryExpressionResultType(props["container"].(string), props["dynamic_container_enabled"].(bool)),
		FolderPath: expandDataFactoryExpressionResultType(props["path"].(string), props["dynamic_path_enabled"].(bool)),
		FileName:   expandDataFactoryExpressionResultType(props["filename"].(string), props["dynamic_filename_enabled"].(bool)),
	}

	return blobStorageLocation
}

func expandDataFactoryDatasetAzureBlobFSLocation(d *pluginsdk.ResourceData) datafactory.BasicDatasetLocation {
	azureBlobFsLocations := d.Get("azure_blob_fs_location").([]interface{})
	if len(azureBlobFsLocations) == 0 || azureBlobFsLocations[0] == nil {
		return nil
	}

	props := azureBlobFsLocations[0].(map[string]interface{})

	blobStorageLocation := datafactory.AzureBlobFSLocation{
		FileSystem: props["file_system"].(string),
		Type:       datafactory.TypeBasicDatasetLocationTypeAzureBlobFSLocation,
	}
	if path := props["path"].(string); len(path) > 0 {
		blobStorageLocation.FolderPath = path
	}
	if filename := props["filename"].(string); len(filename) > 0 {
		blobStorageLocation.FileName = filename
	}

	return blobStorageLocation
}

func flattenDataFactoryDatasetHTTPServerLocation(input *datafactory.HTTPServerLocation) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.RelativeURL != nil {
		result["relative_url"] = input.RelativeURL
	}
	if input.FolderPath != nil {
		path, dynamicPathEnabled := flattenDataFactoryExpressionResultType(input.FolderPath)
		result["path"] = path
		result["dynamic_path_enabled"] = dynamicPathEnabled
	}
	if input.FileName != nil {
		filename, dynamicFilenameEnabled := flattenDataFactoryExpressionResultType(input.FileName)
		result["filename"] = filename
		result["dynamic_filename_enabled"] = dynamicFilenameEnabled
	}

	return []interface{}{result}
}

func flattenDataFactoryDatasetAzureBlobStorageLocation(input *datafactory.AzureBlobStorageLocation) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.Container != nil {
		container, dynamicContainerEnabled := flattenDataFactoryExpressionResultType(input.Container)
		result["container"] = container
		result["dynamic_container_enabled"] = dynamicContainerEnabled
	}
	if input.FolderPath != nil {
		path, dynamicPathEnabled := flattenDataFactoryExpressionResultType(input.FolderPath)
		result["path"] = path
		result["dynamic_path_enabled"] = dynamicPathEnabled
	}
	if input.FileName != nil {
		filename, dynamicFilenameEnabled := flattenDataFactoryExpressionResultType(input.FileName)
		result["filename"] = filename
		result["dynamic_filename_enabled"] = dynamicFilenameEnabled
	}

	return []interface{}{result}
}

func flattenDataFactoryDatasetAzureBlobFSLocation(input *datafactory.AzureBlobFSLocation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	fileSystem, path, fileName := "", "", ""
	if input.FileSystem != nil {
		if v, ok := input.FileSystem.(string); ok {
			fileSystem = v
		}
	}
	if input.FolderPath != nil {
		if v, ok := input.FolderPath.(string); ok {
			path = v
		}
	}
	if input.FileName != nil {
		if v, ok := input.FileName.(string); ok {
			fileName = v
		}
	}

	return []interface{}{
		map[string]interface{}{
			"file_system": fileSystem,
			"path":        path,
			"filename":    fileName,
		},
	}
}

func flattenDataFactoryDatasetSFTPLocation(input *datafactory.SftpLocation) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.FolderPath != nil {
		path, dynamicPathEnabled := flattenDataFactoryExpressionResultType(input.FolderPath)
		result["path"] = path
		result["dynamic_path_enabled"] = dynamicPathEnabled
	}
	if input.FileName != nil {
		filename, dynamicFilenameEnabled := flattenDataFactoryExpressionResultType(input.FileName)
		result["filename"] = filename
		result["dynamic_filename_enabled"] = dynamicFilenameEnabled
	}

	return []interface{}{result}
}

func flattenDataFactoryDatasetCompression(input *datafactory.DatasetCompression) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// the Azure API returns these in a different case to what we're expecting, so we need to convert these
	compressionTypes := []string{
		TypeBasicDatasetCompressionTypeBZip2,
		TypeBasicDatasetCompressionTypeDeflate,
		TypeBasicDatasetCompressionTypeGZip,
		TypeBasicDatasetCompressionTypeTar,
		TypeBasicDatasetCompressionTypeTarGZip,
		TypeBasicDatasetCompressionTypeZipDeflate,
	}

	compressionType := ""
	if t, ok := input.Type.(string); ok {
		for _, v := range compressionTypes {
			if strings.EqualFold(v, t) {
				compressionType = v
			}
		}
	}

	level := ""
	if v, ok := input.Level.(string); ok {
		level = v
	}

	return []interface{}{
		map[string]interface{}{
			"type":  compressionType,
			"level": level,
		},
	}
}

func expandDataFactoryDatasetCompression(d *pluginsdk.ResourceData) *datafactory.DatasetCompression {
	compression := d.Get("compression").([]interface{})
	if len(compression) == 0 || compression[0] == nil {
		return nil
	}

	props := compression[0].(map[string]interface{})
	return &datafactory.DatasetCompression{
		Type:  props["type"].(string),
		Level: props["level"].(string),
	}
}
