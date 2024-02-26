package subaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoredResourceListResponse struct {
	NextLink *string              `json:"nextLink,omitempty"`
	Value    *[]MonitoredResource `json:"value,omitempty"`
}
