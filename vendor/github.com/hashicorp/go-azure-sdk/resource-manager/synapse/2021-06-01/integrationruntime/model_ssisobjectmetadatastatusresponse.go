package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisObjectMetadataStatusResponse struct {
	Error      *string `json:"error,omitempty"`
	Name       *string `json:"name,omitempty"`
	Properties *string `json:"properties,omitempty"`
	Status     *string `json:"status,omitempty"`
}
