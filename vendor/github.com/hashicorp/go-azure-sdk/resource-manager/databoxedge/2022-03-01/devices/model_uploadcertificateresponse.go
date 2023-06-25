package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UploadCertificateResponse struct {
	AadAudience                     *string             `json:"aadAudience,omitempty"`
	AadAuthority                    *string             `json:"aadAuthority,omitempty"`
	AadTenantId                     *string             `json:"aadTenantId,omitempty"`
	AuthType                        *AuthenticationType `json:"authType,omitempty"`
	AzureManagementEndpointAudience *string             `json:"azureManagementEndpointAudience,omitempty"`
	ResourceId                      *string             `json:"resourceId,omitempty"`
	ServicePrincipalClientId        *string             `json:"servicePrincipalClientId,omitempty"`
	ServicePrincipalObjectId        *string             `json:"servicePrincipalObjectId,omitempty"`
}
