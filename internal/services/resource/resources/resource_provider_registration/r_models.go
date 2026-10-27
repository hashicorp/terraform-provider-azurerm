// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

type ResourceProviderRegistrationModel struct {
	Name     string                                     `tfschema:"name"`
	Features []ResourceProviderRegistrationFeatureModel `tfschema:"feature"`
}

type ResourceProviderRegistrationFeatureModel struct {
	Name       string `tfschema:"name"`
	Registered bool   `tfschema:"registered"`
}

func (r Resource) ModelObject() any {
	return &ResourceProviderRegistrationModel{}
}
