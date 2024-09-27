package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildpackBindingProperties struct {
	BindingType       *BindingType                       `json:"bindingType,omitempty"`
	LaunchProperties  *BuildpackBindingLaunchProperties  `json:"launchProperties,omitempty"`
	ProvisioningState *BuildpackBindingProvisioningState `json:"provisioningState,omitempty"`
}
