// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

type ContainerAppCustomDomainResourceModel struct {
	Name                 string `tfschema:"name"`
	ContainerAppId       string `tfschema:"container_app_id"`
	CertificateId        string `tfschema:"container_app_environment_certificate_id"`
	BindingType          string `tfschema:"certificate_binding_type"`
	ManagedCertificateId string `tfschema:"container_app_environment_managed_certificate_id"`
}

func (a Resource) ModelObject() any {
	return &ContainerAppCustomDomainResourceModel{}
}
