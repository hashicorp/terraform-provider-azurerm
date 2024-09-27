package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevToolPortalFeatureSettings struct {
	ApplicationAccelerator *DevToolPortalFeatureDetail `json:"applicationAccelerator,omitempty"`
	ApplicationLiveView    *DevToolPortalFeatureDetail `json:"applicationLiveView,omitempty"`
}
