package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevToolPortalFeatureDetail struct {
	Route *string                    `json:"route,omitempty"`
	State *DevToolPortalFeatureState `json:"state,omitempty"`
}
