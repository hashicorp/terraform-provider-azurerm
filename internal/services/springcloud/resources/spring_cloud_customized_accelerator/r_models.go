// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_customized_accelerator

type SpringCloudCustomizedAcceleratorModel struct {
	Name                     string               `tfschema:"name"`
	SpringCloudAcceleratorId string               `tfschema:"spring_cloud_accelerator_id"`
	AcceleratorTags          []string             `tfschema:"accelerator_tags"`
	AcceleratorType          string               `tfschema:"accelerator_type"`
	Description              string               `tfschema:"description"`
	DisplayName              string               `tfschema:"display_name"`
	GitRepository            []GitRepositoryModel `tfschema:"git_repository"`
	IconURL                  string               `tfschema:"icon_url"`
}

type GitRepositoryModel struct {
	BasicAuth         []BasicAuthModel `tfschema:"basic_auth"`
	SshAuth           []SshAuthModel   `tfschema:"ssh_auth"`
	Branch            string           `tfschema:"branch"`
	CaCertificateId   string           `tfschema:"ca_certificate_id"`
	Commit            string           `tfschema:"commit"`
	GitTag            string           `tfschema:"git_tag"`
	IntervalInSeconds int64            `tfschema:"interval_in_seconds"`
	Url               string           `tfschema:"url"`
	Path              string           `tfschema:"path"`
}

type BasicAuthModel struct {
	Username string `tfschema:"username"`
	Password string `tfschema:"password"`
}

type SshAuthModel struct {
	PrivateKey          string `tfschema:"private_key"`
	HostKey             string `tfschema:"host_key"`
	PrivateKeyAlgorithm string `tfschema:"host_key_algorithm"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudCustomizedAcceleratorModel{}
}
