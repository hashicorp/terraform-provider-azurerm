package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KubernetesVersionProfileProperties struct {
	ProvisioningState *ResourceProvisioningState     `json:"provisioningState,omitempty"`
	Values            *[]KubernetesVersionProperties `json:"values,omitempty"`
}
