package groupquotalimits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupQuotaLimitListProperties struct {
	NextLink          *string            `json:"nextLink,omitempty"`
	ProvisioningState *RequestState      `json:"provisioningState,omitempty"`
	Value             *[]GroupQuotaLimit `json:"value,omitempty"`
}
