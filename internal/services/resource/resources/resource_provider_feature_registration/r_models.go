// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_feature_registration

type ResourceProviderFeatureRegistrationModel struct {
	Name         string `tfschema:"name"`
	ProviderName string `tfschema:"provider_name"`
}

func (r Resource) ModelObject() any {
	return &ResourceProviderFeatureRegistrationModel{}
}
