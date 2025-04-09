package daprcomponents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprMetadata struct {
	Name      *string `json:"name,omitempty"`
	SecretRef *string `json:"secretRef,omitempty"`
	Value     *string `json:"value,omitempty"`
}
