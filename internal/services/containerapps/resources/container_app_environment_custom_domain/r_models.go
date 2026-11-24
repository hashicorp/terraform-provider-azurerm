// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

type ContainerAppEnvironmentCustomDomainModel struct {
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	CertificatePassword string `tfschema:"certificate_password"`
	CertificateValue    string `tfschema:"certificate_blob_base64"`
	DnsSuffix           string `tfschema:"dns_suffix"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentCustomDomainModel{}
}
