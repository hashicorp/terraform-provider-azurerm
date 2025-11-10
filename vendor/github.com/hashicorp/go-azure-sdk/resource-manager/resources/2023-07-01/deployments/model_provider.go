package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Provider struct {
	Id                                *string                            `json:"id,omitempty"`
	Namespace                         *string                            `json:"namespace,omitempty"`
	ProviderAuthorizationConsentState *ProviderAuthorizationConsentState `json:"providerAuthorizationConsentState,omitempty"`
	RegistrationPolicy                *string                            `json:"registrationPolicy,omitempty"`
	RegistrationState                 *string                            `json:"registrationState,omitempty"`
	ResourceTypes                     *[]ProviderResourceType            `json:"resourceTypes,omitempty"`
}
