// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

type ApiCenterDeploymentModel struct {
	Name             string   `tfschema:"name"`
	CustomProperties string   `tfschema:"custom_properties_json"`
	DefinitionId     string   `tfschema:"api_center_api_definition_id"`
	Description      string   `tfschema:"description"`
	EnvironmentId    string   `tfschema:"api_center_environment_id"`
	ServerRuntimeUrl []string `tfschema:"server_runtime_url"`
	State            string   `tfschema:"state"`
	Title            string   `tfschema:"display_name"`
}
