package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListCredentialResponse struct {
	Error      *ListCredentialResponseError      `json:"error,omitempty"`
	Id         *string                           `json:"id,omitempty"`
	Name       *string                           `json:"name,omitempty"`
	Properties *ListCredentialResponseProperties `json:"properties,omitempty"`
	ResourceId *string                           `json:"resourceId,omitempty"`
	Status     *ResourceProvisioningState        `json:"status,omitempty"`
}
