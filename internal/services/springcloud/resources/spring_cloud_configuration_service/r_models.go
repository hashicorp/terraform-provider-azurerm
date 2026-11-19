// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_configuration_service

type SpringCloudConfigurationServiceModel struct {
	Name                 string                       `tfschema:"name"`
	SpringCloudServiceId string                       `tfschema:"spring_cloud_service_id"`
	Generation           string                       `tfschema:"generation"`
	RefreshInterval      int64                        `tfschema:"refresh_interval_in_seconds"`
	Repository           []SpringCloudRepositoryModel `tfschema:"repository"`
}

type SpringCloudRepositoryModel struct {
	Name                  string   `tfschema:"name"`
	Label                 string   `tfschema:"label"`
	Patterns              []string `tfschema:"patterns"`
	Uri                   string   `tfschema:"uri"`
	CaCertificateId       string   `tfschema:"ca_certificate_id"`
	HostKey               string   `tfschema:"host_key"`
	HostKeyAlgorithm      string   `tfschema:"host_key_algorithm"`
	Password              string   `tfschema:"password"`
	PrivateKey            string   `tfschema:"private_key"`
	SearchPaths           []string `tfschema:"search_paths"`
	StrictHostKeyChecking bool     `tfschema:"strict_host_key_checking"`
	Username              string   `tfschema:"username"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudConfigurationServiceModel{}
}
