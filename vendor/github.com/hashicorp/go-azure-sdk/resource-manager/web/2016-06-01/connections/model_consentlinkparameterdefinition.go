package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsentLinkParameterDefinition struct {
	ObjectId      *string `json:"objectId,omitempty"`
	ParameterName *string `json:"parameterName,omitempty"`
	RedirectUrl   *string `json:"redirectUrl,omitempty"`
	TenantId      *string `json:"tenantId,omitempty"`
}
