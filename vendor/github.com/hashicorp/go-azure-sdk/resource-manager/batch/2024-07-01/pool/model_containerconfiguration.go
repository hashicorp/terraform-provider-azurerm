package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerConfiguration struct {
	ContainerImageNames *[]string            `json:"containerImageNames,omitempty"`
	ContainerRegistries *[]ContainerRegistry `json:"containerRegistries,omitempty"`
	Type                ContainerType        `json:"type"`
}
