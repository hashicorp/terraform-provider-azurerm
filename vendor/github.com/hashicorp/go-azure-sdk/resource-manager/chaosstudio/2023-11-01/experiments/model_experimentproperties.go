package experiments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Selectors         []Selector         `json:"selectors"`
	Steps             []Step             `json:"steps"`
}
