// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_managed_certificate

type ContainerAppEnvironmentManagedCertificateModel struct {
	Name                      string         `tfschema:"name"`
	ContainerAppEnvironmentId string         `tfschema:"container_app_environment_id"`
	SubjectName               string         `tfschema:"subject_name"`
	DomainControlValidation   string         `tfschema:"domain_control_validation"`
	Tags                      map[string]any `tfschema:"tags"`
	ValidationToken           string         `tfschema:"validation_token"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentManagedCertificateModel{}
}
