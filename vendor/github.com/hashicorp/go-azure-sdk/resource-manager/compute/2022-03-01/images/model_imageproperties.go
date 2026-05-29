package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageProperties struct {
	HyperVGeneration     *HyperVGenerationTypes `json:"hyperVGeneration,omitempty"`
	ProvisioningState    *string                `json:"provisioningState,omitempty"`
	SourceVirtualMachine *SubResource           `json:"sourceVirtualMachine,omitempty"`
	StorageProfile       *ImageStorageProfile   `json:"storageProfile,omitempty"`
}
