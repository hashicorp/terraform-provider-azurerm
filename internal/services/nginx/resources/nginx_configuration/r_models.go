// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_configuration

type ConfigurationModel struct {
	NginxDeploymentId string          `tfschema:"nginx_deployment_id"`
	ConfigFile        []ConfigFile    `tfschema:"config_file"`
	ProtectedFile     []ProtectedFile `tfschema:"protected_file"`
	PackageData       string          `tfschema:"package_data"`
	RootFile          string          `tfschema:"root_file"`
}

func (m Resource) ModelObject() any {
	return &ConfigurationModel{}
}
