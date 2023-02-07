package workflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenAuthenticationAccessPolicy struct {
	Claims *[]OpenAuthenticationPolicyClaim `json:"claims,omitempty"`
	Type   *OpenAuthenticationProviderType  `json:"type,omitempty"`
}
