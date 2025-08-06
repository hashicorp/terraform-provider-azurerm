package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiPortalProperties struct {
	ApiTryOutEnabledState *ApiPortalApiTryOutEnabledState `json:"apiTryOutEnabledState,omitempty"`
	GatewayIds            *[]string                       `json:"gatewayIds,omitempty"`
	HTTPSOnly             *bool                           `json:"httpsOnly,omitempty"`
	Instances             *[]ApiPortalInstance            `json:"instances,omitempty"`
	ProvisioningState     *ApiPortalProvisioningState     `json:"provisioningState,omitempty"`
	Public                *bool                           `json:"public,omitempty"`
	ResourceRequests      *ApiPortalResourceRequests      `json:"resourceRequests,omitempty"`
	SourceURLs            *[]string                       `json:"sourceUrls,omitempty"`
	SsoProperties         *SsoProperties                  `json:"ssoProperties,omitempty"`
	Url                   *string                         `json:"url,omitempty"`
}
