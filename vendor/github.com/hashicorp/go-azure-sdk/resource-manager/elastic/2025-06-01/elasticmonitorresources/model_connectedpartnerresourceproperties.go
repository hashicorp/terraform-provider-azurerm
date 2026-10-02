package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedPartnerResourceProperties struct {
	AzureResourceId       *string `json:"azureResourceId,omitempty"`
	Location              *string `json:"location,omitempty"`
	PartnerDeploymentName *string `json:"partnerDeploymentName,omitempty"`
	PartnerDeploymentUri  *string `json:"partnerDeploymentUri,omitempty"`
	Type                  *string `json:"type,omitempty"`
}
