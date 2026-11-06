package newrelicmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LatestLinkedSaaSResponse struct {
	IsHiddenSaaS   *bool   `json:"isHiddenSaaS,omitempty"`
	SaaSResourceId *string `json:"saaSResourceId,omitempty"`
}
