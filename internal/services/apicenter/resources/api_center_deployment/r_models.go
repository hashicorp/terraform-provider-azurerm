// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

type ApiCenterDeploymentModel struct {
	Name             string   `tfschema:"name"`
	CustomProperties string   `tfschema:"custom_properties"`
	DefinitionId     string   `tfschema:"api_definition_id"`
	Description      string   `tfschema:"description"`
	EnvironmentId    string   `tfschema:"environment_id"`
	ServerRuntimeUri []string `tfschema:"server_runtime_uri"`
	State            string   `tfschema:"state"`
	Title            string   `tfschema:"display_name"`
}
