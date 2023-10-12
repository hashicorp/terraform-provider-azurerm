package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizedAcceleratorProperties struct {
	AcceleratorTags   *[]string                               `json:"acceleratorTags,omitempty"`
	Description       *string                                 `json:"description,omitempty"`
	DisplayName       *string                                 `json:"displayName,omitempty"`
	GitRepository     AcceleratorGitRepository                `json:"gitRepository"`
	IconUrl           *string                                 `json:"iconUrl,omitempty"`
	ProvisioningState *CustomizedAcceleratorProvisioningState `json:"provisioningState,omitempty"`
}
