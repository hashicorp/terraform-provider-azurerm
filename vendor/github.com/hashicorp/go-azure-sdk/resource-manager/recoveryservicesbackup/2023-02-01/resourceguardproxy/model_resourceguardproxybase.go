package resourceguardproxy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardProxyBase struct {
	Description                   *string                         `json:"description,omitempty"`
	LastUpdatedTime               *string                         `json:"lastUpdatedTime,omitempty"`
	ResourceGuardOperationDetails *[]ResourceGuardOperationDetail `json:"resourceGuardOperationDetails,omitempty"`
	ResourceGuardResourceId       *string                         `json:"resourceGuardResourceId,omitempty"`
}
