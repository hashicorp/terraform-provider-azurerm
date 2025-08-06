package api

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiContractUpdateProperties struct {
	ApiRevision                   *string                                `json:"apiRevision,omitempty"`
	ApiRevisionDescription        *string                                `json:"apiRevisionDescription,omitempty"`
	ApiVersion                    *string                                `json:"apiVersion,omitempty"`
	ApiVersionDescription         *string                                `json:"apiVersionDescription,omitempty"`
	ApiVersionSetId               *string                                `json:"apiVersionSetId,omitempty"`
	AuthenticationSettings        *AuthenticationSettingsContract        `json:"authenticationSettings,omitempty"`
	Contact                       *ApiContactInformation                 `json:"contact,omitempty"`
	Description                   *string                                `json:"description,omitempty"`
	DisplayName                   *string                                `json:"displayName,omitempty"`
	IsCurrent                     *bool                                  `json:"isCurrent,omitempty"`
	IsOnline                      *bool                                  `json:"isOnline,omitempty"`
	License                       *ApiLicenseInformation                 `json:"license,omitempty"`
	Path                          *string                                `json:"path,omitempty"`
	Protocols                     *[]Protocol                            `json:"protocols,omitempty"`
	ServiceURL                    *string                                `json:"serviceUrl,omitempty"`
	SubscriptionKeyParameterNames *SubscriptionKeyParameterNamesContract `json:"subscriptionKeyParameterNames,omitempty"`
	SubscriptionRequired          *bool                                  `json:"subscriptionRequired,omitempty"`
	TermsOfServiceURL             *string                                `json:"termsOfServiceUrl,omitempty"`
	Type                          *ApiType                               `json:"type,omitempty"`
}
