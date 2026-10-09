package virtualmachineimagetemplate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunOutputProperties struct {
	ArtifactId        *string            `json:"artifactId,omitempty"`
	ArtifactUri       *string            `json:"artifactUri,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
