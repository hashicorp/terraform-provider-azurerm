package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizedAcceleratorProperties struct {
	AcceleratorTags   *[]string                               `json:"acceleratorTags,omitempty"`
	AcceleratorType   *CustomizedAcceleratorType              `json:"acceleratorType,omitempty"`
	Description       *string                                 `json:"description,omitempty"`
	DisplayName       *string                                 `json:"displayName,omitempty"`
	GitRepository     AcceleratorGitRepository                `json:"gitRepository"`
	IconURL           *string                                 `json:"iconUrl,omitempty"`
	Imports           *[]string                               `json:"imports,omitempty"`
	ProvisioningState *CustomizedAcceleratorProvisioningState `json:"provisioningState,omitempty"`
}
